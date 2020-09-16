/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"dmdb_exporter/dm/i18n"
)

func init() {
	load()
	if GlobalProperties == nil {
		GlobalProperties = NewProperties()
	}
	if GlobalProperties != nil && GlobalProperties.Len() > 0 {
		setDriverAttributes(GlobalProperties)
	}

	switch Locale {
	//case 0:
	//	i18n.InitConfig(util.FileUtil.Search("/dm/i18n/messages_zh-CN.json"))
	//case 1:
	//	i18n.InitConfig(util.FileUtil.Search("/dm/i18n/messages_en-US.json"))
	//case 2:
	//	i18n.InitConfig(util.FileUtil.Search("/dm/i18n/messages_zh-TW.json"))

	case 0:
		i18n.InitConfig(`{
  "language": "zh-Hans",
  "messages": [
    {
      "id": "error.dsn.invalidSchema",
      "translation": "DSN串必须以dm://开头"
    },
    {
      "id": "error.communicationError",
      "translation": "网络通信异常"
    },
    {
      "id": "error.msgCheckError",
      "translation": "消息校验异常"
    },
    {
      "id": "error.unkownNetWork",
      "translation": "未知的网络"
    },
    {
      "id": "error.invalidConn",
      "translation": "连接失效"
    }
  ]
}`)
	case 1:
	i18n.InitConfig(`{
  "language": "en-US",
  "messages": [
    {
      "id": "error.dsn.invalidSchema",
      "translation": "DSN must start with dm://"
    },
    {
      "id": "error.communicationError",
      "translation": "Communication  error"
    },
    {
      "id": "error.msgCheckError",
      "translation": "Message check error"
    },
    {
      "id": "error.unkownNetWork",
      "translation": "Unkown net work"
    },
    {
      "id": "error.invalidConn",
      "translation": "Invalid connection"
    }
  ]
}`)
	case 2:
		i18n.InitConfig(`{
  "language": "zh-Hant",
  "messages": [
    {
      "id": "error.dsn.invalidSchema",
      "translation": "DSN串必須以dm://開頭"
    },
    {
      "id": "error.communicationError",
      "translation": "網絡通信異常"
    },
    {
      "id": "error.msgCheckError",
      "translation": "消息校驗異常"
    },
    {
      "id": "error.unkownNetWork",
      "translation": "未知的網絡"
    }
  ]
}`)

	}

	sql.Register("dm", newDmDriver())
}

type DmDriver struct {
	filterable
}

func newDmDriver() *DmDriver {
	d := new(DmDriver)
	d.createFilterChain(nil, GlobalProperties)
	d.idGenerator = dmDriverIDGenerator
	return d
}

/*************************************************************
 ** PUBLIC METHODS AND FUNCTIONS
 *************************************************************/
func (d *DmDriver) Open(dsn string) (driver.Conn, error) {
	if len(d.filterChain.filters) == 0 {
		return d.open(dsn)
	} else {
		return d.filterChain.reset().DmDriverOpen(d, dsn)
	}
}

func (d *DmDriver) OpenConnector(dsn string) (driver.Connector, error) {
	if len(d.filterChain.filters) == 0 {
		return d.openConnector(dsn)
	} else {
		return d.filterChain.reset().DmDriverOpenConnector(d, dsn)
	}
}

func (d *DmDriver) open(dsn string) (*DmConnection, error) {
	c, err := d.openConnector(dsn)
	if err != nil {
		return nil, err
	}
	return c.connect(context.Background())
}

func (d *DmDriver) openConnector(dsn string) (*DmConnector, error) {
	connector := new(DmConnector).init()
	connector.url = dsn
	connector.dmDriver = d
	err := connector.mergeConfigs(dsn)
	if err != nil {
		return nil, err
	}
	connector.createFilterChain(connector, nil)
	return connector, nil
}
