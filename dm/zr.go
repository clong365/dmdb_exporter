/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"math/rand"
	"strconv"
	"time"
)

var rwMap = make(map[string]*rwCounter)

type rwCounter struct {
	ntrx_primary int64

	ntrx_total int64

	primaryPercent float64

	standbyPercent float64

	standbyNTrxMap map[string]int64

	standbyCount int
}

func newRWCounter(primaryPercent int, standbyCount int) *rwCounter {
	rwc := new(rwCounter)
	rwc.reset(primaryPercent, standbyCount)
	return rwc
}

func (rwc *rwCounter) reset(primaryPercent int, standbyCount int) {
	rwc.ntrx_primary = 0
	rwc.ntrx_total = 0
	rwc.standbyNTrxMap = make(map[string]int64)
	rwc.standbyCount = standbyCount
	if standbyCount > 0 {
		rwc.primaryPercent = float64(primaryPercent) / 100.0
		rwc.standbyPercent = float64(100-primaryPercent) / 100.0 / float64(standbyCount)
	} else {
		rwc.primaryPercent = 1
		rwc.standbyPercent = 0
	}
}

/**
* 连接创建成功后调用，需要服务器返回standbyCount
 */
func getRwCounterInstance(conn *DmConnection) *rwCounter {
	key := conn.dmConnector.host + "_" + strconv.Itoa(conn.dmConnector.port) + "_" + strconv.Itoa(conn.dmConnector.rwPercent)

	rwc, ok := rwMap[key]
	if !ok {
		rwc = newRWCounter(conn.dmConnector.rwPercent, int(conn.StandbyCount))
		rwMap[key] = rwc
	} else if rwc.standbyCount != int(conn.StandbyCount) {
		rwc.reset(conn.dmConnector.rwPercent, int(conn.StandbyCount))
	}
	return rwc
}

/**
* @return 主机;
 */
func (rwc *rwCounter) countPrimary() RWSiteEnum {
	rwc.adjustNtrx()
	rwc.ntrx_primary++
	rwc.ntrx_total++
	return PRIMARY
}

/**
* @param dest 主机; 备机; any;
* @return 主机; 备机
 */
func (rwc *rwCounter) count(dest RWSiteEnum, standby *DmConnection) RWSiteEnum {
	rwc.adjustNtrx()
	switch dest {
	case ANYOF:
		{
			if rwc.primaryPercent != 1 && (rwc.primaryPercent == 0 || float64(rwc.getStandbyNtrx(standby)) < float64(rwc.ntrx_total)*rwc.standbyPercent || float64(rwc.ntrx_primary) > float64(rwc.ntrx_total)*rwc.primaryPercent) {
				rwc.incrementStandbyNtrx(standby)
				dest = STANDBY
			} else {
				rwc.ntrx_primary++
				dest = PRIMARY
			}
		}
	case STANDBY:
		{
			rwc.incrementStandbyNtrx(standby)
		}
	case PRIMARY:
		{
			rwc.ntrx_primary++
		}
	}

	rwc.ntrx_total++
	return dest
}

/**
* 防止ntrx超出有效范围，等比调整
 */
func (rwc *rwCounter) adjustNtrx() {
	if rwc.ntrx_total < INT64_MAX {
		return
	}

	var min int64
	i := 0
	for _, value := range rwc.standbyNTrxMap {
		if i == 0 {
			min = value
		} else {
			if value < min {
				min = value
			}
		}
		i++
	}

	if min >= rwc.ntrx_primary {
		min = rwc.ntrx_primary
	}

	rwc.ntrx_primary = rwc.ntrx_primary / min
	rwc.ntrx_total = rwc.ntrx_total / min

	for key, value := range rwc.standbyNTrxMap {
		rwc.standbyNTrxMap[key] = value / min
	}

}

func (rwc *rwCounter) getStandbyNtrx(standby *DmConnection) int64 {
	key := standby.dmConnector.host + ":" + strconv.Itoa(standby.dmConnector.port)
	ret, ok := rwc.standbyNTrxMap[key]
	if !ok {
		ret = 0
	}

	return ret
}

func (rwc *rwCounter) incrementStandbyNtrx(standby *DmConnection) {
	key := standby.dmConnector.host + ":" + strconv.Itoa(standby.dmConnector.port)
	ret, ok := rwc.standbyNTrxMap[key]
	if ok {
		ret += 1
	} else {
		ret = 1
	}
	rwc.standbyNTrxMap[key] = ret
}

func (rwc *rwCounter) random(rowCount int) int {
	rand.Seed(time.Now().UnixNano())
	return int(rand.Int31n(int32(rowCount)))
}

func (rwc *rwCounter) String() string {
	return "PERCENT(P/S) : " + strconv.FormatFloat(rwc.primaryPercent, 'f', -1, 64) + "/" + strconv.FormatFloat(rwc.standbyPercent, 'f', -1, 64) + "\nNTRX_PRIMARY : " +
		strconv.FormatInt(rwc.ntrx_primary, 10) + "\nNTRX_TOTAL : " + strconv.FormatInt(rwc.ntrx_total, 10) + "\nNTRX_STANDBY : "
}
