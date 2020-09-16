/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"bytes"
	"context"
	"database/sql/driver"
	"dmdb_exporter/dm/util"
	"net"
	"net/url"
	"os"
	"runtime"
	"strconv"
)

const (
	TimeZoneKey              = "timeZone"
	EnRsCacheKey             = "enRsCache"
	RsCacheSizeKey           = "rsCacheSize"
	RsRefreshFreqKey         = "rsRefreshFreq"
	LoginPrimary             = "loginPrimary"
	LoginModeKey             = "loginMode"
	LoginStatusKey           = "loginStatus"
	SwitchTimesKey           = "switchTimes"
	SwitchIntervalKey        = "switchInterval"
	PrimaryKey               = "primaryKey"
	KeywordsKey              = "keywords"
	CompressKey              = "compress"
	CompressIdKey            = "compressId"
	LoginEncryptKey          = "loginEncrypt"
	CommunicationEncryptKey  = "communicationEncrypt"
	DirectKey                = "direct"
	Dec2DoubleKey            = "dec2double"
	RwSeparateKey            = "rwSeparate"
	RwPercentKey             = "rwPercent"
	RwAutoDistributeKey      = "rwAutoDistribute"
	CompatibleModeKey        = "compatibleMode"
	CompatibleOraKey         = "comOra"
	CipherPathKey            = "cipherPath"
	LoadBalanceKey           = "loadBalance"
	LoadBalancePercentKey    = "loadBalancePercent"
	DoSwitchKey              = "doSwitch"
	LanguageKey              = "language"
	LoadBalanceFreqKey       = "loadBalanceFreq"
	DbAliveCheckFreqKey      = "dbAliveCheckFreq"
	RwStandbyRecoverTimeKey  = "rwStandbyRecoverTime"
	LogLevelKey              = "logLevel"
	LogDirKey                = "logDir"
	LogBufferPoolSizeKey     = "logBufferPoolSize"
	LogBufferSizeKey         = "logBufferSize"
	LogFlusherQueueSizeKey   = "logFlusherQueueSize"
	LogFlushFreqKey          = "logFlushFreq"
	StatEnableKey            = "statEnable"
	StatDirKey               = "statDir"
	StatFlushFreqKey         = "statFlushFreq"
	StatHighFreqSqlCountKey  = "statHighFreqSqlCount"
	StatSlowSqlCountKey      = "statSlowSqlCount"
	StatSqlMaxCountKey       = "statSqlMaxCount"
	StatSqlRemoveModeKey     = "statSqlRemoveMode"
	AddressRemapKey          = "addressreMap"
	UserreMapKey             = "userreMap"
	ConnectTimeoutKey        = "connectTimeout"
	LoginCertificateKey      = "loginCertificate"
	UrlKey                   = "url"
	HostKey                  = "host"
	PortKey                  = "port"
	UserKey                  = "user"
	PasswordKey              = "password"
	RwStandbyKey             = "rwStandby"
	IsCompressKey            = "isCompress"
	RwHAKey                  = "rwHA"
	AppNameKey               = "appName"
	MppLocalKey              = "mppLocal"
	SocketTimeoutKey         = "socketTimeout"
	SessionTimeoutKey        = "sessionTimeout"
	ContinueBatchOnErrorKey  = "continueBatchOnError"
	EscapeProcessKey         = "escapeProcess"
	AutoCommitKey            = "autoCommit"
	MaxRowsKey               = "maxRows"
	RowPrefetchKey           = "rowPrefetch"
	BufPrefetchKey           = "bufPrefetch"
	LobModeKey               = "LobMode"
	StmtPoolSizeKey          = "StmtPoolSize"
	IgnoreCaseKey            = "ignoreCase"
	AlwayseAllowCommitKey    = "AlwayseAllowCommit"
	BatchTypeKey             = "batchType"
	IsBdtaRSKey              = "isBdtaRS"
	ClobAsStringKey          = "clobAsString"
	CallBatchNotKey          = "callBatchNot"
	SslCertPathKey           = "sslCertPath"
	SslKeyPathKey            = "sslKeyPath"
	KerberosLoginConfPathKey = "kerberosLoginConfPath"
	UKeyNameKey              = "uKeyName"
	UKeyPinKey               = "uKeyPin"
	ColumnNameUpperCaseKey   = "columnNameUpperCase"
	ColumnNameCaseKey        = "columnNameCase"
	DatabaseProductNameKey   = "databaseProductName"
	OsAuthTypeKey            = "osAuthType"
	SchemaKey                = "schema"

	TIME_ZONE_DEFAULT int16 = 480

	LOGIN_MODE_ALL_AND_PRIMARY_FIRST int = 0

	LOGIN_MODE_PRIMARY int = 1

	LOGIN_MODE_STANDBY int = 2

	LOGIN_MODE_ALL_AND_STANDBY_FIRST int = 3

	LOGIN_MODE_ALL_AND_NORMAL_FIRST int = 4

	LOGIN_MODE_DEFAULT int = LOGIN_MODE_ALL_AND_PRIMARY_FIRST

	SERVER_MODE_NORMAL int = 0

	SERVER_MODE_PRIMARY int = 1

	SERVER_MODE_STANDBY int = 2

	SERVER_STATUS_OPEN int = 4

	SERVER_STATUS_MOUNT int = 3

	SERVER_STATUS_SUSPEND int = 5

	COMPATIBLE_MODE_ORACLE int = 1

	COMPATIBLE_MODE_MYSQL int = 2

	LANGUAGE_CN int = 0

	LANGUAGE_EN int = 1

	COLUMN_NAME_UPPDER_CASE = 1

	COLUMN_NAME_LOWER_CASE = 2

	compressDef   = Dm_build_81
	compressIDDef = Dm_build_82

	charCodeDef = ""

	enRsCacheDef = false

	rsCacheSizeDef = 20

	rsRefreshFreqDef = 10

	loginModeDef = LOGIN_MODE_DEFAULT

	loginStatusDef = 0

	switchTimesDef = 3

	switchIntervalDef = 2000

	loginEncryptDef = true

	loginCertificateDef = ""

	dec2DoubleDef = false

	rwHADef = false

	rwStandbyDef = false

	rwSeparateDef = false

	rwPercentDef = 25

	rwAutoDistributeDef = true

	rwStandbyRecoverTimeDef = 1000

	doSwitchDef = false

	cipherPathDef = ""

	loadBalanceDef = false

	loadBalancePercentDef = 10

	urlDef = ""

	userDef = "SYSDBA"

	passwordDef = "SYSDBA"

	hostDef = "localhost"

	portDef = DEFAULT_PORT

	appNameDef = ""

	osNameDef = runtime.GOOS

	mppLocalDef = false

	socketTimeoutDef = 0

	connectTimeoutDef = 5000

	sessionTimeoutDef = 0

	osAuthTypeDef = Dm_build_64

	continueBatchOnErrorDef = false

	escapeProcessDef = false

	autoCommitDef = true

	maxRowsDef = 0

	rowPrefetchDef = Dm_build_65

	bufPrefetchDef = 0

	lobModeDef = 1

	stmtPoolMaxSizeDef = 15

	ignoreCaseDef = true

	alwayseAllowCommitDef = true

	batchTypeDef = 1

	isBdtaRSDef = false

	callBatchNotDef = false

	kerberosLoginConfPathDef = ""

	uKeyNameDef = ""

	uKeyPinDef = ""

	databaseProductNameDef = ""

	columnNameCaseDef = 0

	caseSensitiveDef = true

	compatibleModeDef = 0

	localTimezoneDef = TIME_ZONE_DEFAULT
)

type DmConnector struct {
	filterable

	dmDriver *DmDriver

	compress int

	compressID int8

	charCode string

	enRsCache bool

	rsCacheSize int

	rsRefreshFreq int

	loginMode int

	loginStatus int

	switchTimes int

	switchInterval int

	keyWords []string

	loginEncrypt bool

	loginCertificate string

	logLevel int

	dec2Double bool

	rwHA bool

	rwStandby bool

	rwSeparate bool

	rwPercent int

	rwAutoDistribute bool

	rwStandbyRecoverTime int

	doSwitch bool

	cipherPath string

	loadBalance bool

	loadBalancePercent int

	url string

	user string

	password string

	host string

	group *DBGroup

	port int

	appName string

	osName string

	mppLocal bool

	socketTimeout int

	connectTimeout int

	sessionTimeout int

	osAuthType byte

	continueBatchOnError bool

	escapeProcess bool

	autoCommit bool

	maxRows int

	rowPrefetch int

	bufPrefetch int

	lobMode int

	stmtPoolMaxSize int

	ignoreCase bool

	alwayseAllowCommit bool

	batchType int

	isBdtaRS bool

	callBatchNot bool

	sslCertPath string

	sslKeyPath string

	kerberosLoginConfPath string

	uKeyName string

	uKeyPin string

	databaseProductName string

	columnNameCase int

	caseSensitive bool

	compatibleMode int

	localTimezone int16

	schema string

	reConnection *DmConnection
}

func (c *DmConnector) init() *DmConnector {
	c.compress = compressDef
	c.compressID = compressIDDef
	c.charCode = charCodeDef
	c.enRsCache = enRsCacheDef
	c.rsCacheSize = rsCacheSizeDef
	c.rsRefreshFreq = rsRefreshFreqDef
	c.loginMode = loginModeDef
	c.loginStatus = loginStatusDef
	c.switchTimes = switchTimesDef
	c.switchInterval = switchIntervalDef
	c.keyWords = nil
	c.loginEncrypt = loginEncryptDef
	c.loginCertificate = loginCertificateDef
	c.dec2Double = dec2DoubleDef
	c.rwHA = rwHADef
	c.rwStandby = rwStandbyDef
	c.rwSeparate = rwSeparateDef
	c.rwPercent = rwPercentDef
	c.rwAutoDistribute = rwAutoDistributeDef
	c.rwStandbyRecoverTime = rwStandbyRecoverTimeDef
	c.doSwitch = doSwitchDef
	c.cipherPath = cipherPathDef
	c.loadBalance = loadBalanceDef
	c.loadBalancePercent = loadBalancePercentDef
	c.url = urlDef
	c.user = userDef
	c.password = passwordDef
	c.host = hostDef
	c.port = portDef
	c.appName = appNameDef
	c.osName = osNameDef
	c.mppLocal = mppLocalDef
	c.socketTimeout = socketTimeoutDef
	c.connectTimeout = connectTimeoutDef
	c.sessionTimeout = sessionTimeoutDef
	c.osAuthType = osAuthTypeDef
	c.continueBatchOnError = continueBatchOnErrorDef
	c.escapeProcess = escapeProcessDef
	c.autoCommit = autoCommitDef
	c.maxRows = maxRowsDef
	c.rowPrefetch = rowPrefetchDef
	c.bufPrefetch = bufPrefetchDef
	c.lobMode = lobModeDef
	c.stmtPoolMaxSize = stmtPoolMaxSizeDef
	c.ignoreCase = ignoreCaseDef
	c.alwayseAllowCommit = alwayseAllowCommitDef
	c.batchType = batchTypeDef
	c.isBdtaRS = isBdtaRSDef
	c.callBatchNot = callBatchNotDef
	c.kerberosLoginConfPath = kerberosLoginConfPathDef
	c.uKeyName = uKeyNameDef
	c.uKeyPin = uKeyPinDef
	c.databaseProductName = databaseProductNameDef
	c.columnNameCase = columnNameCaseDef
	c.caseSensitive = caseSensitiveDef
	c.compatibleMode = compatibleModeDef
	c.localTimezone = localTimezoneDef

	c.idGenerator = dmConntorIDGenerator
	return c
}

func (c *DmConnector) setAttributes(props *Properties) error {
	if props == nil || props.Len() == 0 {
		return nil
	}

	c.url = props.GetTrimString(UrlKey, c.url)
	c.host = props.GetTrimString(HostKey, c.host)
	c.port = props.GetInt(PortKey, c.port, 0, 65535)
	c.user = props.GetString(UserKey, c.user)
	c.password = props.GetString(PasswordKey, c.password)
	c.rwStandby = props.GetBool(RwStandbyKey, c.rwStandby)

	if b := props.GetBool(IsCompressKey, false); b {
		c.compress = Dm_build_80
	}

	c.compress = props.GetInt(CompressKey, c.compress, 0, 2)
	c.compressID = int8(props.GetInt(CompressIdKey, int(c.compressID), -1, 1))
	c.enRsCache = props.GetBool(EnRsCacheKey, c.enRsCache)
	c.localTimezone = int16(props.GetInt(TimeZoneKey, int(c.localTimezone), -720, 720))
	c.rsCacheSize = props.GetInt(RsCacheSizeKey, c.rsCacheSize, 0, int(INT32_MAX))
	c.rsRefreshFreq = props.GetInt(RsRefreshFreqKey, c.rsRefreshFreq, 0, int(INT32_MAX))
	c.loginMode = props.GetInt(LoginModeKey, c.loginMode, 0, 4)
	c.loginStatus = props.GetInt(LoginStatusKey, c.loginStatus, 0, int(INT32_MAX))
	c.switchTimes = props.GetInt(SwitchTimesKey, c.switchTimes, 0, int(INT32_MAX))
	c.switchInterval = props.GetInt(SwitchIntervalKey, c.switchInterval, 0, int(INT32_MAX))
	c.loginEncrypt = props.GetBool(LoginEncryptKey, c.loginEncrypt)
	c.loginCertificate = props.GetTrimString(LoginCertificateKey, c.loginCertificate)
	c.dec2Double = props.GetBool(Dec2DoubleKey, c.dec2Double)

	c.rwSeparate = props.GetBool(RwSeparateKey, c.rwSeparate)
	c.rwAutoDistribute = props.GetBool(RwAutoDistributeKey, c.rwAutoDistribute)
	c.rwPercent = props.GetInt(RwPercentKey, c.rwPercent, 0, 100)
	c.rwHA = props.GetBool(RwHAKey, c.rwHA)
	c.rwStandbyRecoverTime = props.GetInt(RwStandbyRecoverTimeKey, c.rwStandbyRecoverTime, 0, int(INT32_MAX))
	c.doSwitch = props.GetBool(DoSwitchKey, c.doSwitch)
	c.cipherPath = props.GetTrimString(CipherPathKey, c.cipherPath)

	if props.GetBool(CompatibleOraKey, false) {
		c.compatibleMode = int(COMPATIBLE_MODE_ORACLE)
	}
	c.parseCompatibleMode(props)
	c.loadBalance = props.GetBool(LoadBalanceKey, c.loadBalance)
	c.loadBalancePercent = props.GetInt(LoadBalancePercentKey, c.loadBalancePercent, 0, 100)
	c.keyWords = props.GetStringArray(KeywordsKey, c.keyWords)

	c.appName = props.GetTrimString(AppNameKey, c.appName)
	c.mppLocal = props.GetBool(MppLocalKey, c.mppLocal)
	c.socketTimeout = props.GetInt(SocketTimeoutKey, c.socketTimeout, 0, int(INT32_MAX))
	c.connectTimeout = props.GetInt(ConnectTimeoutKey, c.connectTimeout, 0, int(INT32_MAX))
	c.sessionTimeout = props.GetInt(SessionTimeoutKey, c.sessionTimeout, 0, int(INT32_MAX))

	err := c.parseOsAuthType(props)
	if err != nil {
		return err
	}
	c.continueBatchOnError = props.GetBool(ContinueBatchOnErrorKey,
		c.continueBatchOnError)
	c.escapeProcess = props.GetBool(EscapeProcessKey,
		c.escapeProcess)
	c.autoCommit = props.GetBool(AutoCommitKey, c.autoCommit)
	c.maxRows = props.GetInt(MaxRowsKey, c.maxRows, 0, int(INT32_MAX))
	c.rowPrefetch = props.GetInt(RowPrefetchKey, c.rowPrefetch, 0, int(INT32_MAX))
	c.bufPrefetch = props.GetInt(BufPrefetchKey, c.bufPrefetch, int(Dm_build_66), int(Dm_build_67))
	c.lobMode = props.GetInt(LobModeKey, c.lobMode, 1, 2)
	c.stmtPoolMaxSize = props.GetInt(StmtPoolSizeKey, c.stmtPoolMaxSize, 0, int(INT32_MAX))
	c.ignoreCase = props.GetBool(IgnoreCaseKey, c.ignoreCase)
	c.alwayseAllowCommit = props.GetBool(AlwayseAllowCommitKey, c.alwayseAllowCommit)
	c.batchType = props.GetInt(BatchTypeKey, c.batchType, 1, 2)
	c.isBdtaRS = props.GetBool(IsBdtaRSKey, c.isBdtaRS)

	c.callBatchNot = props.GetBool(CallBatchNotKey, c.callBatchNot)
	c.sslCertPath = props.GetTrimString(SslCertPathKey, c.sslCertPath)
	c.sslKeyPath = props.GetTrimString(SslKeyPathKey, c.sslKeyPath)
	c.kerberosLoginConfPath = props.GetTrimString(KerberosLoginConfPathKey, c.kerberosLoginConfPath)

	c.uKeyName = props.GetTrimString(UKeyNameKey, c.uKeyName)
	c.uKeyPin = props.GetTrimString(UKeyPinKey, c.uKeyPin)

	if b := props.GetBool(ColumnNameUpperCaseKey, false); b {
		c.columnNameCase = COLUMN_NAME_UPPDER_CASE
	} else {
		c.columnNameCase = 0
	}

	v := props.GetTrimString(ColumnNameCaseKey, "")
	if util.StringUtil.EqualsIgnoreCase(v, "upper") {
		c.columnNameCase = COLUMN_NAME_UPPDER_CASE
	} else if util.StringUtil.EqualsIgnoreCase(v, "lower") {
		c.columnNameCase = COLUMN_NAME_LOWER_CASE
	}

	c.databaseProductName = props.GetTrimString(DatabaseProductNameKey, c.databaseProductName)
	c.schema = props.GetTrimString(SchemaKey, c.schema)
	return nil
}

func (c *DmConnector) parseOsAuthType(props *Properties) error {
	value := props.GetString(OsAuthTypeKey, "")
	if value != "" && !util.StringUtil.IsDigit(value) {
		if util.StringUtil.EqualsIgnoreCase(value, "ON") {
			c.osAuthType = Dm_build_64
		} else if util.StringUtil.EqualsIgnoreCase(value, "SYSDBA") {
			c.osAuthType = Dm_build_60
		} else if util.StringUtil.EqualsIgnoreCase(value, "SYSAUDITOR") {
			c.osAuthType = Dm_build_62
		} else if util.StringUtil.EqualsIgnoreCase(value, "SYSSSO") {
			c.osAuthType = Dm_build_61
		} else if util.StringUtil.EqualsIgnoreCase(value, "AUTO") {
			c.osAuthType = Dm_build_63
		} else if util.StringUtil.EqualsIgnoreCase(value, "OFF") {
			c.osAuthType = Dm_build_59
		}
	} else {
		c.osAuthType = byte(props.GetInt(OsAuthTypeKey, int(c.osAuthType), 0, 4))
	}
	if c.user == "" && c.osAuthType == Dm_build_59 {
		c.user = "SYSDBA"
	} else if c.osAuthType != Dm_build_59 && c.user != "" {
		ECGO_OSAUTH_ERROR.throw()
	} else if c.osAuthType != Dm_build_59 {
		c.user = os.Getenv("user")
		c.password = ""
	}
	return nil
}

func (c *DmConnector) parseCompatibleMode(props *Properties) {
	value := props.GetString(CompatibleModeKey, "")
	if value != "" && !util.StringUtil.IsDigit(value) {
		if util.StringUtil.EqualsIgnoreCase(value, "oracle") {
			c.compatibleMode = COMPATIBLE_MODE_ORACLE
		} else if util.StringUtil.EqualsIgnoreCase(value, "mysql") {
			c.compatibleMode = COMPATIBLE_MODE_MYSQL
		}
	} else {
		c.compatibleMode = props.GetInt(CompatibleModeKey, c.compatibleMode, 0, 2)
	}
}

func (c *DmConnector) parseDSN(dsn string) (*Properties, error) {
	dsn = AddressRemap(dsn)

	var dsnProps = NewProperties()
	url, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "dm" {
		return nil, DSN_INVALID_SCHEMA
	}

	if url.User != nil {
		c.user = url.User.Username()
		c.password, _ = url.User.Password()
	}

	for name, group := range ServerGroupMap {
		if util.StringUtil.EqualsIgnoreCase(name, url.Host) {
			c.group = group
			break
		}
	}

	if c.group == nil {
		host, port, err := net.SplitHostPort(url.Host)

		if err != nil || net.ParseIP(host) == nil {
			c.host = hostDef
		} else {
			c.host = host
		}

		c.port, err = strconv.Atoi(port)
		if err != nil {
			c.port = portDef
		}
	}

	q := url.Query()
	for k := range q {
		dsnProps.Set(k, q.Get(k))
	}

	return dsnProps, nil
}

func (c *DmConnector) BuildDSN() string {
	var buf bytes.Buffer

	buf.WriteString("dm://")

	if len(c.user) > 0 {
		buf.WriteString(url.QueryEscape(c.user))
		if len(c.password) > 0 {
			buf.WriteByte(':')
			buf.WriteString(url.QueryEscape(c.password))
		}
		buf.WriteByte('@')
	}

	if len(c.host) > 0 {
		buf.WriteString(c.host)
		if c.port > 0 {
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(c.port))
		}
	}

	hasParam := false
	if c.connectTimeout > 0 {
		if hasParam {
			buf.WriteString("&timeout=")
		} else {
			buf.WriteString("?timeout=")
			hasParam = true
		}
		buf.WriteString(strconv.Itoa(c.connectTimeout))
	}
	return buf.String()
}

func (c *DmConnector) mergeConfigs(dsn string) error {
	props, err := c.parseDSN(dsn)
	if err != nil {
		return err
	}

	UserRemap(props)

	if c.group != nil {
		props.SetDiffProperties(c.group.Props)

		if props.GetBool(RwSeparateKey, false) {
			props.Set(LoginModeKey, strconv.Itoa(LOGIN_MODE_PRIMARY))
			props.Set(LoginStatusKey, strconv.Itoa(SERVER_STATUS_OPEN))
		}
	} else {
		props.SetDiffProperties(GlobalProperties)
	}

	err = c.setAttributes(props)
	if err != nil {
		return err
	}

	setDriverAttributes(props)
	return nil
}

func (c *DmConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return c.filterChain.reset().DmConnectorConnect(c, ctx)
}

func (c *DmConnector) Driver() driver.Driver {
	return c.filterChain.reset().DmConnectorDriver(c)
}

func (c *DmConnector) connect(ctx context.Context) (*DmConnection, error) {
	if c.group != nil {
		return c.group.connect(c)
	} else {
		return c.connectSingle(ctx)
	}
}

func (c *DmConnector) driver() *DmDriver {
	return c.dmDriver
}

func (c *DmConnector) connectSingle(ctx context.Context) (*DmConnection, error) {
	var err error
	var dc *DmConnection
	if c.reConnection == nil {
		dc = &DmConnection{
			closech: make(chan struct{}),
		}
		dc.dmConnector = c
		dc.createFilterChain(c, nil)

		dc.objId = -1
		dc.init()

		dc.startWatcher()
		if err = dc.watchCancel(ctx); err != nil {
			return nil, err
		}
		defer dc.finish()
	} else {
		dc = c.reConnection
		dc.reset()
	}

	dc.Access, err = dm_build_1288(dc)
	if err != nil {
		return nil, err
	}

	if err = dc.Access.dm_build_1329(); err != nil {

		if !dc.closed.IsSet() {
			close(dc.closech)
			if dc.Access != nil {
				dc.Access.Close()
			}
			dc.closed.Set(true)
		}
		return nil, err
	}

	if c.schema != "" {
		_, err = dc.exec("set schema "+c.schema, nil)
		if err != nil {
			return nil, err
		}
	}

	return dc, nil
}
