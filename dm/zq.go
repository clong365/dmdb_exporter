/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"dmdb_exporter/dm/util"
	"io"
	"strings"
	"time"
)

const (
	SQL_SELECT_STANDBY = "select distinct mailIni.inst_name, mailIni.INST_IP, mailIni.INST_PORT, archIni.arch_status " +
		"from  v$arch_status archIni " +
		"left join (select * from V$DM_MAL_INI) mailIni on archIni.arch_dest = mailIni.inst_name " +
		"left join V$MAL_LINK_STATUS on CTL_LINK_STATUS  = 'CONNECTED' AND DATA_LINK_STATUS = 'CONNECTED' " +
		"where archIni.arch_type in ('TIMELY', 'REALTIME')"

	SQL_SELECT_STANDBY2 = "select distinct " +
		"mailIni.mal_inst_name, mailIni.mal_INST_HOST, mailIni.mal_INST_PORT, archIni.arch_status " +
		"from v$arch_status archIni " + "left join (select * from V$DM_MAL_INI) mailIni " +
		"on archIni.arch_dest = mailIni.mal_inst_name " + "left join V$MAL_LINK_STATUS " +
		"on CTL_LINK_STATUS  = 'CONNECTED' AND DATA_LINK_STATUS = 'CONNECTED' " +
		"where archIni.arch_type in ('TIMELY', 'REALTIME')"
)

type rwUtil struct {
}

var RWUtil = rwUtil{}

func (RWUtil rwUtil) reconnect(connection *DmConnection) error {
	if connection.rwInfo == nil {
		return nil
	}

	RWUtil.removeStandby(connection)

	err := connection.reconnect()
	if err != nil {
		return err
	}
	connection.rwInfo.cleanup()
	connection.rwInfo.rwCounter = getRwCounterInstance(connection)

	RWUtil.connectStandby(connection)

	return nil
}

func (RWUtil rwUtil) recoverStandby(connection *DmConnection) {
	if connection.closed.IsSet() || RWUtil.isStandbyAlive(connection) {
		return
	}

	ts := time.Now().Nanosecond()

	freq := connection.dmConnector.rwStandbyRecoverTime * int(time.Second)
	if freq <= 0 || ts-connection.rwInfo.tryRecoverTs < freq {
		return
	}

	RWUtil.connectStandby(connection)
	connection.rwInfo.tryRecoverTs = ts

}

func (RWUtil rwUtil) connectStandby(connection *DmConnection) {
	db := RWUtil.chooseValidStandby(connection)
	if db == nil {
		return
	}

	standbyConnector := *connection.dmConnector
	standbyConnector.host = db.host
	standbyConnector.port = db.port
	standbyConnector.rwStandby = true
	standbyConnectorPointer := &standbyConnector
	var err error
	connection.rwInfo.connStandby, err = standbyConnectorPointer.connectSingle(context.Background())
	if err != nil {
		return
	}

	if int(connection.rwInfo.connStandby.SvrMode) != SERVER_MODE_STANDBY || int(connection.rwInfo.connStandby.SvrStat) != SERVER_STATUS_OPEN {
		RWUtil.removeStandby(connection)
	}
}

func (RWUtil rwUtil) chooseValidStandby(connection *DmConnection) *DB {
	rs, err := connection.query(SQL_SELECT_STANDBY2, nil)
	if err != nil {
		rs, err = connection.query(SQL_SELECT_STANDBY, nil)
		if err != nil {
			return nil
		}
	}
	rowsCount := int(rs.CurrentRows.getRowCount())
	if rowsCount > 0 {
		i := 0
		rowIndex := connection.rwInfo.rwCounter.random(rowsCount)
		dest := make([]driver.Value, 3)
		for err := rs.next(dest); err != io.EOF; i++ {
			if i == rowIndex {
				db := newDB(dest[1].(string), int(dest[2].(int32)))
				rs.close()
				return db
			}
			rs.next(dest)
		}
	}

	rs.close()
	return nil
}

func (RWUtil rwUtil) afterExceptionOnStandby(connection *DmConnection, e error) {
	if e.(*DmError).ErrCode == ECGO_COMMUNITION_ERROR.ErrCode {
		RWUtil.removeStandby(connection)
	}
}

func (RWUtil rwUtil) removeStandby(connection *DmConnection) {
	if connection.rwInfo.connStandby != nil {
		connection.rwInfo.connStandby.close()
		connection.rwInfo.connStandby = nil
	}
}

func (RWUtil rwUtil) executeByConn(conn *DmConnection, execute1 func() (interface{}, error), execute2 func(otherConn *DmConnection) (interface{}, error)) (interface{}, error) {
	turnToPrimary := false
	var standbyException error

	RWUtil.recoverStandby(conn)

	ret, err := execute1()
	if err != nil {
		standbyException = err

		if conn.rwInfo.connCurrent == conn.rwInfo.connStandby {
			RWUtil.afterExceptionOnStandby(conn, err)
			turnToPrimary = true
		} else {
			return nil, err
		}
	}

	curConn := conn.rwInfo.connCurrent
	var otherConn *DmConnection
	if curConn != conn {
		otherConn = conn
	} else {
		otherConn = conn.rwInfo.connStandby
	}

	switch conn.lastExecInfo.retSqlType {
	case Dm_build_91, Dm_build_92, Dm_build_96, Dm_build_103, Dm_build_102, Dm_build_94:
		{

			execute2(otherConn)
		}
	case Dm_build_100:
		{

			if conn.dmConnector.rwHA && conn == conn.rwInfo.connStandby &&
				(conn.lastExecInfo.rsDatas == nil || len(conn.lastExecInfo.rsDatas) == 0) {
				turnToPrimary = true
			}
		}
	}

	if turnToPrimary {
		conn.rwInfo.distribute = conn.rwInfo.rwCounter.countPrimary()

		t, err := execute2(conn)
		if err != nil {
			if standbyException != nil {
				conn.rwInfo.stmtCurrent = conn.rwInfo.stmtStandby
			}
			return t, err
		}
		return t, nil

	}
	return ret, nil
}

func (RWUtil rwUtil) executeByStmt(stmt *DmStatement, execute1 func() (interface{}, error), execute2 func(otherStmt *DmStatement) (interface{}, error)) (interface{}, error) {
	turnToPrimary := false
	var standbyException error

	RWUtil.recoverStandby(stmt.dmConn)

	ret, err := execute1()
	if err != nil {
		standbyException = err

		if stmt.rwInfo.stmtCurrent == stmt.rwInfo.stmtStandby {
			RWUtil.afterExceptionOnStandby(stmt.dmConn, err)
			turnToPrimary = true
		} else {
			return nil, err
		}
	}

	curStmt := stmt.rwInfo.stmtCurrent
	var otherStmt *DmStatement
	if curStmt != stmt {
		otherStmt = stmt
	} else {
		otherStmt = stmt.rwInfo.stmtStandby
	}

	switch stmt.execInfo.retSqlType {
	case Dm_build_91, Dm_build_92, Dm_build_96, Dm_build_103, Dm_build_102, Dm_build_94:
		{

			RWUtil.copyStatement(curStmt, otherStmt)
			execute2(otherStmt)
		}
	case Dm_build_100:
		{

			if stmt.dmConn.dmConnector.rwHA && curStmt == stmt.rwInfo.stmtStandby &&
				(curStmt.execInfo.rsDatas == nil || len(curStmt.execInfo.rsDatas) == 0) {
				turnToPrimary = true
			}
		}
	}

	if turnToPrimary {
		stmt.dmConn.rwInfo.distribute = stmt.dmConn.rwInfo.rwCounter.countPrimary()
		stmt.rwInfo.stmtCurrent = stmt

		RWUtil.copyStatement(stmt.rwInfo.stmtStandby, stmt)

		t, err := execute2(stmt)
		if err != nil {
			if standbyException != nil {
				stmt.rwInfo.stmtCurrent = stmt.rwInfo.stmtStandby
			}
			return t, err
		}
		return t, nil

	}
	return ret, nil
}

func (RWUtil rwUtil) distributeSqlByConn(connection *DmConnection, query string) RWSiteEnum {
	if !RWUtil.isStandbyAlive(connection) {
		connection.rwInfo.connCurrent = connection
		return connection.rwInfo.rwCounter.countPrimary()
	}

	if (connection.rwInfo.distribute == PRIMARY && !connection.trxFinish) ||
		(connection.rwInfo.distribute == STANDBY && !connection.rwInfo.connStandby.trxFinish) {
		if connection.rwInfo.distribute == PRIMARY {
			connection.rwInfo.connCurrent = connection
		} else {
			connection.rwInfo.connCurrent = connection.rwInfo.connStandby
		}

		return connection.rwInfo.distribute
	}

	readonly := true

	if query != "" {
		tmpsql := strings.TrimSpace(query)
		sqlhead := strings.SplitN(tmpsql, " ", 2)[0]
		if util.StringUtil.EqualsIgnoreCase(sqlhead, "INSERT") ||
			util.StringUtil.EqualsIgnoreCase(sqlhead, "UPDATE") ||
			util.StringUtil.EqualsIgnoreCase(sqlhead, "DELETE") ||
			util.StringUtil.EqualsIgnoreCase(sqlhead, "CREATE") ||
			util.StringUtil.EqualsIgnoreCase(sqlhead, "TRUNCATE") ||
			util.StringUtil.EqualsIgnoreCase(sqlhead, "DROP") ||
			util.StringUtil.EqualsIgnoreCase(sqlhead, "ALTER") ||
			util.StringUtil.EqualsIgnoreCase(sqlhead, "SP_SET_SESSION_READONLY") {
			readonly = false
		} else {
			readonly = true
		}
	}

	if readonly && connection.IsoLevel != int32(sql.LevelSerializable) {
		rWSiteEnum := connection.rwInfo.rwCounter.count(ANYOF, connection.rwInfo.connStandby)
		if rWSiteEnum == PRIMARY {
			connection.rwInfo.connCurrent = connection
		} else if rWSiteEnum == STANDBY {
			connection.rwInfo.connCurrent = connection.rwInfo.connStandby
		}
		return rWSiteEnum
	}

	connection.rwInfo.connCurrent = connection
	return connection.rwInfo.rwCounter.countPrimary()
}

func (RWUtil rwUtil) distributeSqlByStmt(statement *DmStatement, sql string) {
	if !RWUtil.isStandbyAlive(statement.dmConn) || !RWUtil.isStandbyStatementValid(statement) {
		statement.dmConn.rwInfo.rwCounter.countPrimary()
		return
	}

	statement.dmConn.rwInfo.distribute = RWUtil.distributeSqlByConn(statement.dmConn, sql)

	if statement.dmConn.rwInfo.distribute == PRIMARY {
		statement.rwInfo.stmtCurrent = statement
	} else {
		statement.rwInfo.stmtCurrent = statement.rwInfo.stmtStandby
	}
}

func (RWUtil rwUtil) copyStatement(srcStmt *DmStatement, destStmt *DmStatement) {
	destStmt.nativeSql = srcStmt.nativeSql
	destStmt.params = srcStmt.params
	destStmt.paramCount = srcStmt.paramCount
	destStmt.curRowBindIndicator = srcStmt.curRowBindIndicator
}

func (RWUtil rwUtil) isStandbyAlive(connection *DmConnection) bool {
	return connection.rwInfo.connStandby != nil && !connection.rwInfo.connStandby.closed.IsSet()
}

func (RWUtil rwUtil) isStandbyStatementValid(statement *DmStatement) bool {
	return statement.rwInfo.stmtStandby != nil && !statement.rwInfo.stmtStandby.closed
}
