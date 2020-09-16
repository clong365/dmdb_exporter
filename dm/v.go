/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"time"
)

const (
	Seconds_1900_1970 = 2209017600

	OFFSET_YEAR = 0

	OFFSET_MONTH = 1

	OFFSET_DAY = 2

	OFFSET_HOUR = 3

	OFFSET_MINUTE = 4

	OFFSET_SECOND = 5

	OFFSET_MILLISECOND = 6

	OFFSET_TIMEZONE = 7

	DT_LEN = 8

	INVALID_VALUE = int(INT32_MIN)
)

type DmTimestamp struct {
	dt                  []int
	dtype               int
	scale               int
	oracleFormatPattern string
	oracleDateLanguage  int
}

func newDmTimestampFromDt(dt []int, dtype int, scale int) *DmTimestamp {
	dmts := new(DmTimestamp)
	dmts.dt = dt
	dmts.dtype = dtype
	dmts.scale = scale
	return dmts
}

func newDmTimestampFromBytes(bytes []byte, column column, conn *DmConnection) *DmTimestamp {
	dmts := new(DmTimestamp)
	dmts.dt = decode(bytes, column.isBdta, int(column.colType), int(column.scale), int(conn.dmConnector.localTimezone), int(conn.DbTimezone))

	if isLocalTimeZone(int(column.colType), int(column.scale)) {
		dmts.scale = getLocalTimeZoneScale(int(column.colType), int(column.scale))
	} else {
		dmts.scale = int(column.scale)
	}

	dmts.dtype = int(column.colType)
	dmts.scale = int(column.scale)
	dmts.oracleDateLanguage = int(conn.OracleDateLanguage)
	switch column.colType {
	case DATE:
		dmts.oracleFormatPattern = conn.OracleDateFormat
	case TIME:
		dmts.oracleFormatPattern = conn.OracleTimeFormat
	case TIME_TZ:
		dmts.oracleFormatPattern = conn.OracleTimeTZFormat
	case DATETIME:
		dmts.oracleFormatPattern = conn.OracleTimestampFormat
	case DATETIME_TZ:
		dmts.oracleFormatPattern = conn.OracleTimestampTZFormat
	}
	return dmts
}

func NewDmTimestampFromString(str string) (*DmTimestamp, error) {
	dt := make([]int, DT_LEN)
	dtype, err := toDTFromString(str, dt)
	if err != nil {
		return nil, err
	}

	if dtype == DATE {
		return newDmTimestampFromDt(dt, dtype, 0), nil
	}
	return newDmTimestampFromDt(dt, dtype, 6), nil
}

func NewDmTimestampFromTime(time time.Time) *DmTimestamp {
	dt := toDTFromTime(time)
	return newDmTimestampFromDt(dt, DATETIME, 6)
}

func (dmTimestamp *DmTimestamp) ToTime() time.Time {
	return toTimeFromDT(dmTimestamp.dt, 0)
}

func (dest *DmTimestamp) Scan(src interface{}) error {
	switch src := src.(type) {
	case *DmTimestamp:
		*dest = *src
		return nil
	case time.Time:
		ret := NewDmTimestampFromTime(src)
		*dest = *ret
		return nil
	default:
		return UNSUPPORTED_SCAN
	}
}

func (dmTimestamp *DmTimestamp) toBytes() ([]byte, error) {
	return encode(dmTimestamp.dt, dmTimestamp.dtype, dmTimestamp.scale, dmTimestamp.dt[OFFSET_TIMEZONE])
}

/**
 * 获取当前对象的年月日时分秒，如果原来没有decode会先decode;
 */
func (dmTimestamp *DmTimestamp) getDt() []int {
	return dmTimestamp.dt
}

func (dmTimestamp *DmTimestamp) getTime() int64 {
	sec := toTimeFromDT(dmTimestamp.dt, 0).Unix()
	return sec + int64(dmTimestamp.dt[OFFSET_MILLISECOND])
}

func (dmTimestamp *DmTimestamp) setTime(time int64) {
	timeInMillis := (time / 1000) * 1000
	nanos := (int64)((time % 1000) * 1000000)
	if nanos < 0 {
		nanos = 1000000000 + nanos
		timeInMillis = (((time / 1000) - 1) * 1000)
	}
	dmTimestamp.dt = toDTFromUnix(timeInMillis, nanos)
}

func (dmTimestamp *DmTimestamp) setTimezone(tz int) error {
	// DM中合法的时区取值范围为-12:59至+14:00
	if tz <= -13*60 || tz > 14*60 {
		return ECGO_INVALID_DATETIME_FORMAT.throw()
	}
	dmTimestamp.dt[OFFSET_TIMEZONE] = tz
	return nil
}

func (dmTimestamp *DmTimestamp) getNano() int64 {
	return int64(dmTimestamp.dt[OFFSET_MILLISECOND] * 1000)
}

func (dmTimestamp *DmTimestamp) setNano(nano int64) {
	dmTimestamp.dt[OFFSET_MILLISECOND] = (int)(nano / 1000)
}

func (dmTimestamp *DmTimestamp) string() string {
	if dmTimestamp.oracleFormatPattern != "" {
		return dtToStringByOracleFormat(dmTimestamp.dt, dmTimestamp.oracleFormatPattern, dmTimestamp.oracleDateLanguage)
	}
	return dtToString(dmTimestamp.dt, dmTimestamp.dtype, dmTimestamp.scale)
}
