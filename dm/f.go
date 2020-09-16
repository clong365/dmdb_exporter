/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"context"
	"strconv"
	"strings"
	"time"
)

const (
	STATUS_VALID_TIME = 20 * time.Second // ms

	// sort 值
	SORT_SERVER_MODE_INVALID = -1 // 不允许连接的模式

	SORT_SERVER_NOT_ALIVE = -2 // 站点无法连接

	SORT_UNKNOWN = 0 // 站点还未连接过，模式未知

	SORT_PRIMARY = 10

	SORT_STANDBY = 20

	SORT_NORMAL = 30

	// OPEN>MOUNT>SUSPEND
	SORT_OPEN = 3

	SORT_MOUNT = 2

	SORT_SUSPEND = 1
)

type DB struct {
	host            string
	port            int
	alive           bool
	statusRefreshTs int // 状态更新的时间点
	serverMode      int
	serverStatus    int
	sort            int
}

func newDB(host string, port int) *DB {
	db := new(DB)
	db.host = host
	db.port = port
	db.serverMode = -1
	db.serverStatus = -1
	db.sort = SORT_UNKNOWN
	return db
}

func (db *DB) getSort(checkTime bool) int {
	if checkTime {
		if time.Now().Nanosecond()-db.statusRefreshTs < int(STATUS_VALID_TIME) {
			return db.sort
		} else {
			return SORT_UNKNOWN
		}
	}

	return db.sort
}

func (db *DB) calcSort(loginMode int) int {
	sort := 0
	switch loginMode {
	case LOGIN_MODE_ALL_AND_PRIMARY_FIRST:
		{
			// 主机优先：PRIMARY>NORMAL>STANDBY
			switch db.serverMode {
			case SERVER_MODE_NORMAL:
				sort += SORT_NORMAL * 10
			case SERVER_MODE_PRIMARY:
				sort += SORT_PRIMARY * 100
			case SERVER_MODE_STANDBY:
				sort += SORT_STANDBY
			}
		}
	case LOGIN_MODE_ALL_AND_STANDBY_FIRST:
		{
			// STANDBY优先: STANDBY>PRIMARY>NORMAL
			switch db.serverMode {
			case SERVER_MODE_NORMAL:
				sort += SORT_NORMAL
			case SERVER_MODE_PRIMARY:
				sort += SORT_PRIMARY * 10
			case SERVER_MODE_STANDBY:
				sort += SORT_STANDBY * 100
			}
		}
	case LOGIN_MODE_ALL_AND_NORMAL_FIRST:
		{
			// NORMAL优先：NORMAL>PRIMARY>STANDBY
			switch db.serverMode {
			case SERVER_MODE_NORMAL:
				sort += SORT_NORMAL * 100
			case SERVER_MODE_PRIMARY:
				sort += SORT_PRIMARY * 10
			case SERVER_MODE_STANDBY:
				sort += SORT_STANDBY
			}
		}
	case LOGIN_MODE_PRIMARY:
		if db.serverMode != SERVER_MODE_PRIMARY {
			return SORT_SERVER_MODE_INVALID
		}
		sort += SORT_PRIMARY
	case LOGIN_MODE_STANDBY:
		if db.serverMode != SERVER_MODE_STANDBY {
			return SORT_SERVER_MODE_INVALID
		}
		sort += SORT_STANDBY
	}

	switch db.serverStatus {
	case SERVER_STATUS_MOUNT:
		sort += SORT_MOUNT
	case SERVER_STATUS_OPEN:
		sort += SORT_OPEN
	case SERVER_STATUS_SUSPEND:
		sort += SORT_SUSPEND
	}
	return sort
}

func (db *DB) refreshStatus(alive bool, conn *DmConnection) {
	db.alive = alive
	db.statusRefreshTs = time.Now().Nanosecond()
	if alive {
		db.serverMode = int(conn.SvrMode)
		db.serverStatus = int(conn.SvrStat)
		db.sort = db.calcSort(conn.dmConnector.loginMode)
	} else {
		db.serverMode = -1
		db.serverStatus = -1
		db.sort = SORT_SERVER_NOT_ALIVE
	}
}

func (db *DB) connect(connector *DmConnector, first bool) (*DmConnection, error) {
	connector.host = db.host
	connector.port = db.port
	conn, err := connector.connectSingle(context.Background())
	if err != nil {
		db.refreshStatus(false, conn)
		return nil, err
	}
	db.refreshStatus(true, conn)

	// 模式不匹配, 这里使用的是连接之前的sort，连接之后server的状态可能发生改变sort也可能改变
	if conn.dmConnector.loginStatus == SERVER_STATUS_OPEN && int(conn.SvrStat) != SERVER_STATUS_OPEN {
		conn.close()
		return nil, ECGO_INVALID_SERVER_MODE.throw()
	}

	if !db.checkServerMode(first, conn) {
		conn.close()
		return nil, ECGO_INVALID_SERVER_MODE.throw()
	}

	return conn, nil
}

func (db *DB) checkServerMode(first bool, conn *DmConnection) bool {
	if !first {
		switch conn.dmConnector.loginMode {
		case LOGIN_MODE_PRIMARY:
			return int(conn.SvrMode) == SERVER_MODE_PRIMARY

		case LOGIN_MODE_STANDBY:
			return int(conn.SvrMode) == SERVER_MODE_STANDBY

		default:
			return true
		}
	}

	switch conn.dmConnector.loginMode {
	case LOGIN_MODE_ALL_AND_PRIMARY_FIRST, LOGIN_MODE_PRIMARY:
		return int(conn.SvrMode) == SERVER_MODE_PRIMARY

	case LOGIN_MODE_ALL_AND_STANDBY_FIRST, LOGIN_MODE_STANDBY:
		return int(conn.SvrMode) == SERVER_MODE_STANDBY

	case LOGIN_MODE_ALL_AND_NORMAL_FIRST:
		return int(conn.SvrMode) == SERVER_MODE_NORMAL

	}

	return false
}

func (db *DB) getServerStatusDesc(serverStatus int) string {
	ret := ""
	switch db.serverStatus {
	case SERVER_STATUS_OPEN:
		ret = "OPEN"
	case SERVER_STATUS_MOUNT:
		ret = "MOUNT"
	case SERVER_STATUS_SUSPEND:
		ret = "SUSPEND"
	default:
		ret = "UNKNOW"
	}
	return ret
}

func (db *DB) getServerModeDesc(serverMode int) string {
	ret := ""
	switch db.serverMode {
	case SERVER_MODE_NORMAL:
		ret = "NORMAL"
	case SERVER_MODE_PRIMARY:
		ret = "PRIMARY"
	case SERVER_MODE_STANDBY:
		ret = "STANDBY"
	default:
		ret = "UNKNOW"
	}
	return ret
}

func (db *DB) String() string {
	return strings.TrimSpace(db.host) + ":" + strconv.Itoa(db.port) +
		" (" + db.getServerModeDesc(db.serverMode) + ", " + db.getServerStatusDesc(db.serverStatus) + ")"
}
