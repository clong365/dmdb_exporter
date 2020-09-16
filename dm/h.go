/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"bytes"
	"dmdb_exporter/dm/util"
	"time"
)

/**
 * dm_svc.conf中配置的服务名对应的一组实例, 以及相关属性和状态信息
 */
type DBGroup struct {
	Name       string
	ServerList []DB
	Props      *Properties
}

var curServerPos = -1

func newDBGroup(name string, serverList []DB) *DBGroup {
	g := new(DBGroup)
	g.Name = name
	g.ServerList = serverList
	return g
}

func (g *DBGroup) connect(connector *DmConnector) (*DmConnection, error) {
	serverCount := len(g.ServerList)
	startPos := 0

	curServerPos = (curServerPos + 1) % serverCount
	startPos = curServerPos

	// 指定开始序号，是为了相同状态的站点上的连接能均匀分布
	count := len(g.ServerList)
	orgServers := make([]DB, count)
	for i := 0; i < count; i++ {
		orgServers[i] = g.ServerList[(i+startPos)%count]
	}
	return g.tryConnectServerList(connector, orgServers)
}

/**
* 按sort从大到小排序，相同sort值顺序不变
 */
func sortServer(orgServers []DB, checkTime bool) []DB {
	count := len(orgServers)
	sortServers := make([]DB, count)
	var max, tmp DB
	for i := 0; i < count; i++ {
		max = orgServers[i]
		for j := i + 1; j < count; j++ {
			if max.getSort(checkTime) < orgServers[j].getSort(checkTime) {
				tmp = max
				max = orgServers[j]
				orgServers[j] = tmp
			}
		}
		sortServers[i] = max
	}
	return sortServers
}

/**
* 遍历连接服务名列表中的各个站点
*
 */
func (g *DBGroup) tryConnectServerList(connector *DmConnector, servers []DB) (*DmConnection, error) {
	var sortServers []DB
	var ex error
	for i := 0; i < connector.switchTimes; i++ {
		// 循环了一遍，如果没有符合要求的, 重新排序, 再尝试连接
		sortServers = sortServer(servers, i == 0)
		conn, err := g.traverseServerList(connector, sortServers, i == 0)
		if err != nil {
			ex = err
			time.Sleep(time.Duration(connector.switchInterval) * time.Millisecond)
			continue
		}
		return conn, nil
	}

	return nil, ex
}

/**
* 从指定编号开始，遍历一遍服务名中的ip列表，只连接指定类型（主机或备机）的ip
* @param servers
* @param checkTime
*
* @exception
* DBError.ECJDBC_INVALID_SERVER_MODE 有站点的模式不匹配
* DBError.ECJDBC_COMMUNITION_ERROR 所有站点都连不上
 */
func (g *DBGroup) traverseServerList(connector *DmConnector, servers []DB, checkTime bool) (*DmConnection, error) {
	errorMsg := bytes.NewBufferString("")
	var invalidModeErr error
	for _, server := range servers {
		conn, err := server.connect(connector, checkTime)
		if err != nil {
			if err == ECGO_INVALID_SERVER_MODE {
				invalidModeErr = err
			}
			errorMsg.WriteString("[")
			errorMsg.WriteString(server.String())
			errorMsg.WriteString("]")
			errorMsg.WriteString(err.Error())
			errorMsg.WriteString(util.StringUtil.LineSeparator())
			continue
		}
		return conn, nil
	}

	if invalidModeErr != nil {
		return nil, ECGO_INVALID_SERVER_MODE.addDetail("(" + errorMsg.String() + ")")
	} else if errorMsg.Len() > 0 {
		return nil, ECGO_COMMUNITION_ERROR.addDetail("(" + errorMsg.String() + ")")
	}

	return nil, ECGO_COMMUNITION_ERROR.throw()
}
