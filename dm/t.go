/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"container/list"
	"context"
	"database/sql"
	"database/sql/driver"
	"dmdb_exporter/dm/util"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	BIND_IN byte = 0x01

	BIND_OUT byte = 0x10
)

var rp = newRsPool()

type DmStatement struct {
	filterable

	dmConn    *DmConnection
	rsMap     map[int16]*innerRows
	inUse     bool
	innerUsed bool

	id int32

	cursorName string

	readBaseColName bool

	execInfo *execInfo

	resultSetType int

	resultSetConcurrency int

	resultSetHoldability int

	nativeSql string

	maxFieldSize int

	maxRows int64

	escapeProcessing bool

	queryTimeout int32

	fetchDirection int

	fetchSize int

	cursorUpdateRow int64

	closeOnCompletion bool

	closed bool

	columns []column

	params []parameter

	paramCount int16

	curRowBindIndicator []byte
}

type stmtPoolInfo struct {
	id int32

	cursorName string

	readBaseColName bool
}

type rsPoolKey struct {
	dbGuid        string
	currentSchema string
	sql           string
	paramCount    int
}

func newRsPoolKey(stmt *DmStatement, sql string) rsPoolKey {
	rpk := new(rsPoolKey)
	rpk.dbGuid = stmt.dmConn.Guid
	rpk.currentSchema = stmt.dmConn.Schema
	rpk.paramCount = int(stmt.paramCount)

	rpk.sql = sql
	return *rpk
}

func (key rsPoolKey) equals(destKey rsPoolKey) bool {
	return key.dbGuid == destKey.dbGuid &&
		key.currentSchema == destKey.currentSchema &&
		key.sql == destKey.sql &&
		key.paramCount == destKey.paramCount

}

type rsPoolValue struct {
	m_lastChkTime int
	m_TbIds       []int32
	m_TbTss       []int64
	execInfo      *execInfo
}

func newRsPoolValue(execInfo *execInfo) rsPoolValue {
	rpv := new(rsPoolValue)
	rpv.execInfo = execInfo
	rpv.m_lastChkTime = time.Now().Nanosecond()
	copy(rpv.m_TbIds, execInfo.tbIds)
	copy(rpv.m_TbTss, execInfo.tbTss)
	return *rpv
}

func (rpv rsPoolValue) refreshed(conn *DmConnection) (bool, error) {

	if conn.dmConnector.rsRefreshFreq == 0 {
		return false, nil
	}

	if rpv.m_lastChkTime+conn.dmConnector.rsRefreshFreq*int(time.Second) > time.Now().Nanosecond() {
		return false, nil
	}

	tss, err := conn.Access.Dm_build_1408(interface{}(rpv.m_TbIds).([]uint32))
	if err != nil {
		return false, err
	}
	rpv.m_lastChkTime = time.Now().Nanosecond()

	var tbCount int
	if tss != nil {
		tbCount = len(tss)
	}

	if tbCount != len(rpv.m_TbTss) {
		return true, nil
	}

	for i := 0; i < tbCount; i++ {
		if rpv.m_TbTss[i] != tss[i] {
			return true, nil
		}

	}
	return false, nil
}

func (rpv rsPoolValue) getResultSet(stmt *DmStatement) *innerRows {
	destDatas := rpv.execInfo.rsDatas
	var totalRows int
	if rpv.execInfo.rsDatas != nil {
		totalRows = len(rpv.execInfo.rsDatas)
	}

	if stmt.maxRows > 0 && stmt.maxRows < int64(totalRows) {
		destDatas = make([][][]byte, stmt.maxRows)
		copy(destDatas[:len(destDatas)], rpv.execInfo.rsDatas[:len(destDatas)])
	}

	rs := newLocalInnerRows(stmt, stmt.columns, destDatas)
	rs.id = 1
	return rs
}

func (rpv rsPoolValue) getDataLen() int {
	return rpv.execInfo.rsSizeof
}

type rsPool struct {
	rsMap        map[rsPoolKey]rsPoolValue
	rsList       *list.List
	totalDataLen int
}

func newRsPool() *rsPool {
	rp := new(rsPool)
	rp.rsMap = make(map[rsPoolKey]rsPoolValue, 100)
	rp.rsList = list.New()
	return rp
}

func (rp *rsPool) removeInList(key rsPoolKey) {
	for e := rp.rsList.Front(); e != nil && e.Value.(rsPoolKey).equals(key); e = e.Next() {
		rp.rsList.Remove(e)
	}
}

func (rp *rsPool) put(stmt *DmStatement, sql string, execInfo *execInfo) {
	var dataLen int
	if execInfo != nil {
		dataLen = execInfo.rsSizeof
	}

	cacheSize := stmt.dmConn.dmConnector.rsCacheSize * 1024 * 1024

	for rp.totalDataLen+dataLen > cacheSize {
		if rp.totalDataLen == 0 {
			return
		}

		lk := rp.rsList.Back().Value.(rsPoolKey)
		rp.totalDataLen -= rp.rsMap[lk].getDataLen()
		rp.rsList.Remove(rp.rsList.Back())
		delete(rp.rsMap, rp.rsList.Back().Value.(rsPoolKey))
	}

	key := newRsPoolKey(stmt, sql)
	value := newRsPoolValue(execInfo)

	if _, ok := rp.rsMap[key]; !ok {
		rp.rsList.PushFront(key)
	} else {
		rp.removeInList(key)
		rp.rsList.PushFront(key)
	}

	rp.rsMap[key] = value
	rp.totalDataLen += dataLen
}

func (rp *rsPool) get(stmt *DmStatement, sql string) (*rsPoolValue, error) {
	key := newRsPoolKey(stmt, sql)

	v, ok := rp.rsMap[key]
	if ok {
		b, err := v.refreshed(stmt.dmConn)
		if err != nil {
			return nil, err
		}

		if b {
			rp.removeInList(key)
			delete(rp.rsMap, key)
			return nil, nil
		}

		rp.removeInList(key)
		rp.rsList.PushFront(key)
		return &v, nil
	} else {
		return nil, nil
	}
}

func (s *DmStatement) Close() error {
	if s.closed {
		return nil
	}
	if len(s.filterChain.filters) == 0 {
		return s.close()
	}
	return s.filterChain.reset().DmStatementClose(s)
}

func (s *DmStatement) NumInput() int {
	if err := s.checkClosed(); err != nil {
		return 0
	}
	if len(s.filterChain.filters) == 0 {
		return s.numInput()
	}
	return s.filterChain.reset().DmStatementNumInput(s)
}

func (s *DmStatement) Exec(args []driver.Value) (driver.Result, error) {
	if err := s.checkClosed(); err != nil {
		return nil, err
	}
	if len(s.filterChain.filters) == 0 {
		return s.exec(args)
	}
	return s.filterChain.reset().DmStatementExec(s, args)
}

func (s *DmStatement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if err := s.checkClosed(); err != nil {
		return nil, err
	}
	if len(s.filterChain.filters) == 0 {
		return s.execContext(ctx, args)
	}
	return s.filterChain.reset().DmStatementExecContext(s, ctx, args)
}

func (s *DmStatement) Query(args []driver.Value) (driver.Rows, error) {
	if err := s.checkClosed(); err != nil {
		return nil, err
	}
	if len(s.filterChain.filters) == 0 {
		return s.query(args)
	}
	return s.filterChain.reset().DmStatementQuery(s, args)
}

func (s *DmStatement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if err := s.checkClosed(); err != nil {
		return nil, err
	}
	if len(s.filterChain.filters) == 0 {
		s.queryContext(ctx, args)
	}
	return s.filterChain.reset().DmStatementQueryContext(s, ctx, args)
}

func (s *DmStatement) CheckNamedValue(nv *driver.NamedValue) error {
	if len(s.filterChain.filters) == 0 {
		s.checkNamedValue(nv)
	}
	return s.filterChain.reset().DmStatementCheckNamedValue(s, nv)
}

func (st *DmStatement) prepare() error {
	var err error
	if st.dmConn.dmConnector.escapeProcess {
		st.nativeSql, err = st.dmConn.escape(st.nativeSql, st.dmConn.dmConnector.keyWords)
		if err != nil {
			return err
		}
	}

	st.execInfo, err = st.dmConn.Access.Dm_build_1351(st, Dm_build_85)
	if err != nil {
		return err
	}
	st.curRowBindIndicator = make([]byte, st.paramCount)
	return nil
}

func (stmt *DmStatement) close() error {
	stmt.inUse = true
	if stmt.dmConn.stmtPool != nil && len(stmt.dmConn.stmtPool) < stmt.dmConn.dmConnector.stmtPoolMaxSize {
		stmt.pool()
		return nil
	} else {
		return stmt.free()
	}
}

func (stmt *DmStatement) numInput() int {
	return int(stmt.paramCount)
}

func (stmt *DmStatement) checkNamedValue(nv *driver.NamedValue) error {
	var err error
	nv.Value, err = converter{stmt.dmConn}.ConvertValue(nv.Value)
	return err
}

func (stmt *DmStatement) exec(args []driver.Value) (*DmResult, error) {
	var err error

	stmt.inUse = true
	err = stmt.executeInner(args, Dm_build_87)
	if err != nil {
		return nil, err
	}
	return newDmResult(stmt, stmt.execInfo), nil
}

func (stmt *DmStatement) execContext(ctx context.Context, args []driver.NamedValue) (*DmResult, error) {
	stmt.inUse = true
	dargs, err := namedValueToValue(stmt, args)
	if err != nil {
		return nil, err
	}

	if err := stmt.dmConn.watchCancel(ctx); err != nil {
		return nil, err
	}
	defer stmt.dmConn.finish()

	return stmt.exec(dargs)
}

func (stmt *DmStatement) query(args []driver.Value) (*DmRows, error) {
	var err error
	stmt.inUse = true
	err = stmt.executeInner(args, Dm_build_86)
	if err != nil {
		return nil, err
	}

	return newDmRows(newInnerRows(0, stmt, stmt.execInfo)), nil
}

func (stmt *DmStatement) queryContext(ctx context.Context, args []driver.NamedValue) (*DmRows, error) {
	stmt.inUse = true
	dargs, err := namedValueToValue(stmt, args)
	if err != nil {
		return nil, err
	}

	if err := stmt.dmConn.watchCancel(ctx); err != nil {
		return nil, err
	}

	rows, err := stmt.query(dargs)
	if err != nil {
		stmt.dmConn.finish()
		return nil, err
	}
	rows.finish = stmt.dmConn.finish
	return rows, err
}

func newDmStmt(conn *DmConnection, sql string) (*DmStatement, error) {
	var s *DmStatement

	if conn.stmtMap != nil && len(conn.stmtMap) > 0 {
		for _, sv := range conn.stmtMap {
			if !sv.inUse {
				sv.inUse = true
				sv.nativeSql = sql
				s = sv
				break
			}
		}
	}

	if s == nil {
		s = new(DmStatement)
		s.resetFilterable(&conn.filterable)
		s.objId = -1
		s.idGenerator = dmStmtIDGenerator
		s.dmConn = conn
		s.maxRows = int64(conn.dmConnector.maxRows)
		s.nativeSql = sql
		s.rsMap = make(map[int16]*innerRows)
		s.inUse = true

		if conn.stmtPool != nil && len(conn.stmtPool) > 0 {
			len := len(conn.stmtPool)
			spi := conn.stmtPool[0]
			copy(conn.stmtPool, conn.stmtPool[1:])
			conn.stmtPool = conn.stmtPool[:len-1]
			s.id = spi.id
			s.cursorName = spi.cursorName
			s.readBaseColName = spi.readBaseColName
		} else {
			err := conn.Access.Dm_build_1333(s)
			if err != nil {
				return nil, err
			}
		}
		conn.stmtMap[s.id] = s
	}

	return s, nil

}

func (stmt *DmStatement) checkClosed() error {
	if stmt.dmConn.closed.IsSet() {
		return driver.ErrBadConn
	} else if stmt.closed {
		return ECGO_STATEMENT_HANDLE_CLOSED.throw()
	}

	return nil
}

func (stmt *DmStatement) pool() {
	for _, rs := range stmt.rsMap {
		rs.Close()
	}

	stmt.dmConn.stmtPool = append(stmt.dmConn.stmtPool, stmtPoolInfo{stmt.id, stmt.cursorName, stmt.readBaseColName})
	delete(stmt.dmConn.stmtMap, stmt.id)
	stmt.inUse = false
	stmt.closed = true
}

func (stmt *DmStatement) free() error {
	for _, rs := range stmt.rsMap {
		rs.Close()
	}

	err := stmt.dmConn.Access.Dm_build_1338(int32(stmt.id))
	if err != nil {
		return err
	}
	delete(stmt.dmConn.stmtMap, stmt.id)
	stmt.inUse = false
	stmt.closed = true
	return nil
}

func encodeArgs(stmt *DmStatement, args []driver.Value) ([][]interface{}, error) {
	bytes := [][]interface{}{make([]interface{}, len(args), len(args))}

	var err error

	for i, arg := range args {
	nextSwitch:
		if arg == nil {
			if resetColType(stmt, i, NULL) {
				bytes[0][i] = ParamDataEnum_Null
			} else if stmt.params[i].colType == CURSOR && stmt.params[i].cursorStmt == nil {
				stmt.curRowBindIndicator[i] |= BIND_OUT
				stmt.params[i].cursorStmt = &DmStatement{dmConn: stmt.dmConn}
				err = stmt.params[i].cursorStmt.dmConn.Access.Dm_build_1333(stmt.params[i].cursorStmt)
			}

			continue
		}

		switch v := arg.(type) {
		case bool:
			if resetColType(stmt, i, BIT) {
				bytes[0][i], err = G2DB.fromBool(v, stmt.params[i], stmt.dmConn)
			}
		case int8:
			if resetColType(stmt, i, TINYINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(v), stmt.params[i], stmt.dmConn)
			}
		case int16:
			if resetColType(stmt, i, SMALLINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(v), stmt.params[i], stmt.dmConn)
			}
		case int32:
			if resetColType(stmt, i, INT) {
				bytes[0][i], err = G2DB.fromInt64(int64(v), stmt.params[i], stmt.dmConn)
			}
		case int64:
			if resetColType(stmt, i, BIGINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(v), stmt.params[i], stmt.dmConn)
			}
		case int:
			if resetColType(stmt, i, BIGINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(v), stmt.params[i], stmt.dmConn)
			}
		case uint8:
			if resetColType(stmt, i, SMALLINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(v), stmt.params[i], stmt.dmConn)
			}
		case uint16:
			if resetColType(stmt, i, INT) {
				bytes[0][i], err = G2DB.fromInt64(int64(v), stmt.params[i], stmt.dmConn)
			}
		case uint32:
			if resetColType(stmt, i, BIGINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(v), stmt.params[i], stmt.dmConn)
			}

		case float32:
			if resetColType(stmt, i, REAL) {
				bytes[0][i], err = G2DB.fromFloat32(v, stmt.params[i], stmt.dmConn)
			}
		case float64:
			if resetColType(stmt, i, DOUBLE) {
				bytes[0][i], err = G2DB.fromFloat64(float64(v), stmt.params[i], stmt.dmConn)
			}
		case []byte:
			if resetColType(stmt, i, VARBINARY) {
				bytes[0][i], err = G2DB.fromBytes(v, stmt.params[i], stmt.dmConn)
			}
		case string:

			if v == "" && emptyStringToNil(stmt.params[i].colType) {
				arg = nil
				goto nextSwitch
			}
			if resetColType(stmt, i, VARCHAR) {
				bytes[0][i], err = G2DB.fromString(v, stmt.params[i], stmt.dmConn)
			}
		case time.Time:
			if resetColType(stmt, i, DATETIME_TZ) {
				bytes[0][i], err = G2DB.fromTime(v, stmt.params[i], stmt.dmConn)
			}
		case DmTimestamp:
			if resetColType(stmt, i, DATETIME_TZ) {
				bytes[0][i], err = G2DB.fromTime(v.ToTime(), stmt.params[i], stmt.dmConn)
			}
		case DmIntervalDT:
			if resetColType(stmt, i, INTERVAL_DT) {
				bytes[0][i], err = G2DB.fromDmIntervalDT(v, stmt.params[i], stmt.dmConn)
			}
		case DmIntervalYM:
			if resetColType(stmt, i, INTERVAL_YM) {
				bytes[0][i], err = G2DB.fromDmdbIntervalYM(v, stmt.params[i], stmt.dmConn)
			}
		case DmDecimal:
			if resetColType(stmt, i, DECIMAL) {
				bytes[0][i], err = G2DB.fromDecimal(v, stmt.params[i], stmt.dmConn)
			}

		case DmBlob:
			if resetColType(stmt, i, BLOB) {
				bytes[0][i], err = G2DB.fromBlob(DmBlob(v), stmt.params[i], stmt.dmConn)
				if err != nil {
					return nil, err
				}
			}
		case DmClob:
			bytes[0][i], err = G2DB.fromClob(DmClob(v), stmt.params[i], stmt.dmConn)
			if err != nil {
				return nil, err
			}
		case DmArray:
			if resetColType(stmt, i, ARRAY) {
				da := &v
				da, err = da.create(stmt.dmConn)
				if err != nil {
					return nil, err
				}

				bytes[0][i], err = G2DB.fromArray(da, stmt.params[i], stmt.dmConn)
			}
		case DmStruct:
			if resetColType(stmt, i, RECORD) {
				ds := &v
				ds, err = ds.create(stmt.dmConn)
				if err != nil {
					return nil, err
				}

				bytes[0][i], err = G2DB.fromStruct(ds, stmt.params[i], stmt.dmConn)
			}
		case sql.Out:
			arg = v.Dest
			goto nextSwitch

		case *bool:
			if resetColType(stmt, i, BIT) {
				if *v {
					bytes[0][i], err = G2DB.fromInt64(int64(1), stmt.params[i], stmt.dmConn)
				} else {
					bytes[0][i], err = G2DB.fromInt64(int64(0), stmt.params[i], stmt.dmConn)
				}
			}
		case *int8:
			if resetColType(stmt, i, TINYINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(*v), stmt.params[i], stmt.dmConn)
			}
		case *int16:
			if resetColType(stmt, i, SMALLINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(*v), stmt.params[i], stmt.dmConn)
			}
		case *int32:
			if resetColType(stmt, i, INT) {
				bytes[0][i], err = G2DB.fromInt64(int64(*v), stmt.params[i], stmt.dmConn)
			}
		case *int64:
			if resetColType(stmt, i, BIGINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(*v), stmt.params[i], stmt.dmConn)
			}
		case *int:
			if resetColType(stmt, i, BIGINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(*v), stmt.params[i], stmt.dmConn)
			}
		case *uint8:
			if resetColType(stmt, i, SMALLINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(*v), stmt.params[i], stmt.dmConn)
			}
		case *uint16:
			if resetColType(stmt, i, INT) {
				bytes[0][i], err = G2DB.fromInt64(int64(*v), stmt.params[i], stmt.dmConn)
			}
		case *uint32:
			if resetColType(stmt, i, BIGINT) {
				bytes[0][i], err = G2DB.fromInt64(int64(*v), stmt.params[i], stmt.dmConn)
			}

		case *float32:
			if resetColType(stmt, i, REAL) {
				bytes[0][i], err = G2DB.fromFloat64(float64(*v), stmt.params[i], stmt.dmConn)
			}
		case *float64:
			if resetColType(stmt, i, DOUBLE) {
				bytes[0][i], err = G2DB.fromFloat64(float64(*v), stmt.params[i], stmt.dmConn)
			}
		case *[]byte:
			if resetColType(stmt, i, VARBINARY) {
				bytes[0][i], err = G2DB.fromBytes(*v, stmt.params[i], stmt.dmConn)
			}
		case *string:

			if *v == "" && emptyStringToNil(stmt.params[i].colType) {
				arg = nil
				goto nextSwitch
			}
			if resetColType(stmt, i, VARCHAR) {
				bytes[0][i], err = G2DB.fromString(*v, stmt.params[i], stmt.dmConn)
			}
		case *time.Time:
			if resetColType(stmt, i, DATETIME_TZ) {
				bytes[0][i], err = G2DB.fromTime(*v, stmt.params[i], stmt.dmConn)
			}
		case *DmTimestamp:
			if resetColType(stmt, i, DATETIME_TZ) {
				bytes[0][i], err = G2DB.fromTime(v.ToTime(), stmt.params[i], stmt.dmConn)
			}
		case *DmIntervalDT:
			if resetColType(stmt, i, INTERVAL_DT) {
				bytes[0][i], err = G2DB.fromDmIntervalDT(*v, stmt.params[i], stmt.dmConn)
			}
		case *DmIntervalYM:
			if resetColType(stmt, i, INTERVAL_YM) {
				bytes[0][i], err = G2DB.fromDmdbIntervalYM(*v, stmt.params[i], stmt.dmConn)
			}
		case *DmDecimal:
			if resetColType(stmt, i, DECIMAL) {
				bytes[0][i], err = G2DB.fromDecimal(*v, stmt.params[i], stmt.dmConn)
			}
		case *DmBlob:
			if resetColType(stmt, i, BLOB) {
				bytes[0][i], err = G2DB.fromBlob(DmBlob(*v), stmt.params[i], stmt.dmConn)
			}
		case *DmArray:
			if resetColType(stmt, i, ARRAY) {
				v, err = v.create(stmt.dmConn)
				if err != nil {
					return nil, err
				}

				bytes[0][i], err = G2DB.fromArray(v, stmt.params[i], stmt.dmConn)
			}
		case *DmStruct:
			if resetColType(stmt, i, RECORD) {
				v, err = v.create(stmt.dmConn)
				if err != nil {
					return nil, err
				}

				bytes[0][i], err = G2DB.fromStruct(v, stmt.params[i], stmt.dmConn)
			}
		case *driver.Rows:
			if stmt.params[i].colType == CURSOR && !resetColType(stmt, i, CURSOR) && stmt.params[i].cursorStmt == nil {
				stmt.params[i].cursorStmt = &DmStatement{dmConn: stmt.dmConn}
				err = stmt.params[i].cursorStmt.dmConn.Access.Dm_build_1333(stmt.params[i].cursorStmt)
			}
		case io.Reader:
			bytes[0][1], err = G2DB.fromReader(io.Reader(v), stmt.params[i], stmt.dmConn)
			if err != nil {
				return nil, err
			}
		default:
			err = ECGO_UNSUPPORTED_INPARAM_TYPE.throw()
		}

		if err != nil {
			return nil, err
		}

	}

	return bytes, nil
}

type converter struct {
	conn *DmConnection
}
type decimalDecompose interface {
	Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32)
}

func (c converter) ConvertValue(v interface{}) (driver.Value, error) {
	if driver.IsValue(v) {
		return v, nil
	}

	switch vr := v.(type) {
	case driver.Valuer:
		sv, err := callValuerValue(vr)
		if err != nil {
			return nil, err
		}
		if !driver.IsValue(sv) {
			return nil, fmt.Errorf("non-Value type %T returned from Value", sv)
		}
		return sv, nil

	case decimalDecompose, big.Int, big.Float, DmDecimal, *DmDecimal, DmTimestamp, *DmTimestamp, DmIntervalDT, *DmIntervalDT,
		DmIntervalYM, *DmIntervalYM, driver.Rows, *driver.Rows, DmArray, *DmArray, DmStruct, *DmStruct, sql.Out:
		return vr, nil
	case DmClob:

		if vr.connection == nil {
			vr.connection = c.conn
		}
		return vr, nil
	case DmBlob:

		if vr.connection == nil {
			vr.connection = c.conn
		}
		return vr, nil
	case *DmBlob:

		if vr.connection == nil {
			vr.connection = c.conn
		}
		return vr, nil
	case io.Reader:
		return vr, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		} else {
			return converter{c.conn}.ConvertValue(rv.Elem().Interface())
		}
	case reflect.Int:
		return rv.Int(), nil
	case reflect.Int8:
		return int8(rv.Int()), nil
	case reflect.Int16:
		return int16(rv.Int()), nil
	case reflect.Int32:
		return int32(rv.Int()), nil
	case reflect.Int64:
		return int64(rv.Int()), nil
	case reflect.Uint8:
		return uint8(rv.Uint()), nil
	case reflect.Uint16:
		return uint16(rv.Uint()), nil
	case reflect.Uint32:
		return uint32(rv.Uint()), nil
	case reflect.Uint64, reflect.Uint:
		u64 := rv.Uint()
		if u64 >= 1<<63 {
			bigInt := &big.Int{}
			bigInt.SetString(strconv.FormatUint(u64, 10), 10)
			return FromBigInt(bigInt)
		}
		return int64(u64), nil
	case reflect.Float32:
		return float32(rv.Float()), nil
	case reflect.Float64:
		return float64(rv.Float()), nil
	case reflect.Bool:
		return rv.Bool(), nil
	case reflect.Slice:
		ek := rv.Type().Elem().Kind()
		if ek == reflect.Uint8 {
			return rv.Bytes(), nil
		}
		return nil, fmt.Errorf("unsupported type %T, a slice of %s", v, ek)
	case reflect.String:
		return rv.String(), nil
	}
	return nil, fmt.Errorf("unsupported type %T, a %s", v, rv.Kind())
}

var valuerReflectType = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

func callValuerValue(vr driver.Valuer) (v driver.Value, err error) {
	if rv := reflect.ValueOf(vr); rv.Kind() == reflect.Ptr &&
		rv.IsNil() &&
		rv.Type().Elem().Implements(valuerReflectType) {
		return nil, nil
	}
	return vr.Value()
}

func namedValueToValue(stmt *DmStatement, named []driver.NamedValue) ([]driver.Value, error) {
	if int(stmt.paramCount) != len(named) {
		return nil, INVALID_PARAMETER_NUMBER
	}
	dargs := make([]driver.Value, stmt.paramCount)
	for i, _ := range dargs {
		found := false
		for _, nv := range named {
			if nv.Name != "" && strings.ToUpper(nv.Name) == strings.ToUpper(stmt.params[i].name) {
				dargs[i] = nv.Value
				found = true
				break
			}
		}

		if !found {
			dargs[i] = named[i].Value
		}

	}
	return dargs, nil
}

func (stmt *DmStatement) executeInner(args []driver.Value, executeType int16) (err error) {

	var bytes [][]interface{}

	if stmt.paramCount > 0 {
		bytes, err = encodeArgs(stmt, args)
		if err != nil {
			return err
		}
	}
	stmt.execInfo, err = stmt.dmConn.Access.Dm_build_1365(stmt, bytes)
	if err != nil {
		return err
	}
	if stmt.execInfo.outParamDatas != nil {
		for i, outParamData := range stmt.execInfo.outParamDatas {

			if stmt.curRowBindIndicator[i]&BIND_OUT == BIND_OUT {

				if args[i] == nil {
					switch stmt.params[i].colType {
					case BOOLEAN:
						args[i], err = DB2G.toBool(outParamData, &stmt.params[i].column, stmt.dmConn)
					case BIT:
						if strings.ToLower(stmt.params[i].typeName) == "boolean" {
							args[i], err = DB2G.toBool(outParamData, &stmt.params[i].column, stmt.dmConn)
						}

						args[i], err = DB2G.toInt8(outParamData, &stmt.params[i].column, stmt.dmConn)
					case TINYINT:
						args[i], err = DB2G.toInt8(outParamData, &stmt.params[i].column, stmt.dmConn)
					case SMALLINT:
						args[i], err = DB2G.toInt16(outParamData, &stmt.params[i].column, stmt.dmConn)
					case INT:
						args[i], err = DB2G.toInt32(outParamData, &stmt.params[i].column, stmt.dmConn)
					case BIGINT:
						args[i], err = DB2G.toInt64(outParamData, &stmt.params[i].column, stmt.dmConn)
					case REAL:
						args[i], err = DB2G.toFloat32(outParamData, &stmt.params[i].column, stmt.dmConn)
					case DOUBLE:
						args[i], err = DB2G.toFloat64(outParamData, &stmt.params[i].column, stmt.dmConn)
					case DATE, TIME, DATETIME, TIME_TZ, DATETIME_TZ:
						args[i], err = DB2G.toTime(outParamData, &stmt.params[i].column, stmt.dmConn)
					case INTERVAL_DT:
						args[i] = newDmIntervalDTByBytes(outParamData)
					case INTERVAL_YM:
						args[i] = newDmIntervalYMByBytes(outParamData)
					case DECIMAL:
						args[i], err = DB2G.toDmDecimal(outParamData, &stmt.params[i].column, stmt.dmConn)
					case BINARY, VARBINARY:
						args[i] = util.StringUtil.BytesToHexString(outParamData, false)
					case BLOB:
						args[i] = DB2G.toDmBlob(outParamData, &stmt.params[i].column, stmt.dmConn)
					case CHAR, VARCHAR2, VARCHAR:
						args[i] = DB2G.toString(outParamData, &stmt.params[i].column, stmt.dmConn)
					case CLOB:
						args[i] = DB2G.toDmClob(outParamData, stmt.dmConn, &stmt.params[i].column)
					default:
						err = ECGO_UNSUPPORTED_OUTPARAM_TYPE.throw()
					}
				} else {
				nextSwitch:
					switch v := args[i].(type) {
					case sql.Out:
						args[i] = v.Dest
						goto nextSwitch
					case *string:
						args[i] = DB2G.toString(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *[]byte:
						args[i], err = DB2G.toBytes(outParamData, &stmt.params[i].column, stmt.dmConn)

					case *bool:
						args[i], err = DB2G.toBool(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *int8:
						args[i], err = DB2G.toInt8(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *int16:
						args[i], err = DB2G.toInt16(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *int32:
						args[i], err = DB2G.toInt32(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *int64:
						args[i], err = DB2G.toInt64(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *uint8:
						args[i], err = DB2G.toByte(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *uint16:
						args[i], err = DB2G.toUInt16(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *uint32:
						args[i], err = DB2G.toUInt32(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *uint64:
						args[i], err = DB2G.toUInt64(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *int:
						args[i], err = DB2G.toInt(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *uint:
						args[i], err = DB2G.toUInt(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *float32:
						args[i], err = DB2G.toFloat32(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *float64:
						args[i], err = DB2G.toFloat64(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *time.Time:
						args[i], err = DB2G.toTime(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *DmIntervalDT:
						args[i] = newDmIntervalDTByBytes(outParamData)
					case *DmIntervalYM:
						args[i] = newDmIntervalYMByBytes(outParamData)
					case *DmDecimal:
						args[i], err = DB2G.toDmDecimal(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *DmBlob:
						args[i] = DB2G.toDmBlob(outParamData, &stmt.params[i].column, stmt.dmConn)
					case *driver.Rows:
						if stmt.params[i].colType == CURSOR {
							var tmpExecInfo *execInfo
							tmpExecInfo, err = stmt.dmConn.Access.Dm_build_1375(stmt.params[i].cursorStmt, 1)
							if err != nil {
								return err
							}

							if tmpExecInfo.hasResultSet {
								*v = newDmRows(newInnerRows(0, stmt.params[i].cursorStmt, tmpExecInfo))
							} else {
								*v = nil
							}
						}
					case *DmArray:
						args[i], err = TypeDataSV.bytesToArray(outParamData, nil, stmt.params[i].typeDescriptor)
					case *DmStruct:
						args[i], err = TypeDataSV.bytesToRecord(outParamData, nil, stmt.params[i].typeDescriptor)
					default:
						err = ECGO_UNSUPPORTED_OUTPARAM_TYPE.throw()
					}
				}
			}
		}
	}

	return err
}
