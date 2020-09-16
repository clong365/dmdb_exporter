/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"bufio"
	"dmdb_exporter/dm/util"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var LogDirDef, _ = os.Getwd()

var StatDirDef, _ = os.Getwd()

const (
	DEFAULT_PORT int = 5236

	//log level
	LOG_OFF int = 0

	LOG_ERROR int = 1

	LOG_WARN int = 2

	LOG_SQL int = 3

	LOG_INFO int = 4

	LOG_DEBUG int = 5

	LOG_ALL int = 9

	//stat
	STAT_SQL_REMOVE_LATEST int = 0

	STAT_SQL_REMOVE_OLDEST int = 1

	// 编码字符集
	ENCODING_UTF8 string = "UTF-8"

	ENCODING_EUCKR string = "EUC-KR"

	ENCODING_GB18030 string = "GB18030"

	LoadBalanceFreqDef = 60 * 1000

	DbAliveCheckFreqDef = 0

	LocaleDef = 0

	// log
	LogLevelDef = LOG_OFF // 日志级别：off, error, warn, sql, info, all

	LogFlushFreqDef = 10 // 日志刷盘时间s (>=0)

	LogFlushQueueSizeDef = 100 //日志队列大小

	LogBufferSizeDef = 32 * 1024 // 日志缓冲区大小 (>0)

	// stat
	StatEnableDef = false //

	StatFlushFreqDef = 3 // 日志刷盘时间s (>=0)

	StatSlowSqlCountDef = 100 // 慢sql top行数，(0-1000)

	StatHighFreqSqlCountDef = 100 // 高频sql top行数， (0-1000)

	StatSqlMaxCountDef = 100000 // sql 统计最大值(0-100000)

	StatSqlRemoveModeDef = STAT_SQL_REMOVE_LATEST // 记录sql数超过最大值时，sql淘汰方式
)

var (
	LoadBalanceFreq = LoadBalanceFreqDef

	DbAliveCheckFreq = DbAliveCheckFreqDef

	Locale = LocaleDef // 0:简体中文 1：英文 2:繁体中文

	// log
	LogLevel = LogLevelDef // 日志级别：off, error, warn, sql,
	// info, all

	LogDir = LogDirDef

	LogFlushFreq = LogFlushFreqDef // 日志刷盘时间s (>=0)

	LogFlushQueueSize = LogFlushQueueSizeDef

	LogBufferSize = LogBufferSizeDef // 日志缓冲区大小 (>0)

	// stat
	StatEnable = StatEnableDef //

	StatDir = StatDirDef // jdbc工作目录,所有生成的文件都在该目录下

	StatFlushFreq = StatFlushFreqDef // 日志刷盘时间s (>=0)

	StatSlowSqlCount = StatSlowSqlCountDef // 慢sql top行数，(0-1000)

	StatHighFreqSqlCount = StatHighFreqSqlCountDef // 高频sql top行数， (0-1000)

	StatSqlMaxCount = StatSqlMaxCountDef // sql 统计最大值(0-100000)

	StatSqlRemoveMode = StatSqlRemoveModeDef // 记录sql数超过最大值时，sql淘汰方式

	/*---------------------------------------------------------------*/
	ServerGroupMap = make(map[string]*DBGroup)

	GlobalProperties *Properties

	addressRemap = make(map[string]string)

	userRemap = make(map[string]string)
)

func AddressRemap(url string) string {
	if url != "" {
		startIdx := strings.Index(url, "//")
		endIdx := strings.Index(url, "?")

		if startIdx != -1 {
			startIdx += 2
		}

		if endIdx == -1 {
			endIdx = len(url)
		}

		if startIdx != -1 {
			newAddress, ok := addressRemap[strings.TrimSpace(url[startIdx:endIdx])]
			if ok && newAddress != "" {
				url = url[0:startIdx] + newAddress + url[endIdx:]
			}
		}
	}
	return url
}

func UserRemap(props *Properties) {
	user := props.GetString("user", "")

	if user == "" {
		return
	}

	tmp, ok := userRemap[user]
	if ok {
		user = tmp
	}
	props.Set("user", user)
}

func load() {
	var filePath string // dm_svc.conf 文件路径

	switch runtime.GOOS {
	case "windows":
		filePath = os.Getenv("SystemRoot") + "\\system32\\dm_svc.conf"
	case "linux":
		filePath = "/etc/dm_svc.conf"
	default:
		return
	}
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return
	}
	fileReader := bufio.NewReader(file)

	GlobalProperties = NewProperties()
	var groupProps *Properties
	var line string //dm_svc.conf读取到的一行

	for line, err = fileReader.ReadString('\n'); line != "" && (err == nil || err == io.EOF); line, err = fileReader.ReadString('\n') {
		// 去除#标记的注释
		if notesIndex := strings.IndexByte(line, '#'); notesIndex != -1 {
			line = line[:notesIndex]
		}
		// 去除前后多余的空格
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			groupName := strings.ToLower(line[1 : len(line)-1])
			dbGroup, ok := ServerGroupMap[groupName]
			if groupName == "" || !ok {
				continue
			}
			groupProps = dbGroup.Props
			if groupProps.IsNil() {
				groupProps = NewProperties()
				groupProps.SetProperties(GlobalProperties)
				dbGroup.Props = groupProps
			}

		} else {
			cfgInfo := strings.Split(line, "=")
			if len(cfgInfo) < 2 {
				continue
			}
			key := strings.TrimSpace(cfgInfo[0])
			value := strings.TrimSpace(cfgInfo[1])
			if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
				value = strings.TrimSpace(value[1 : len(value)-1])
			} else {
				continue
			}
			if key == "" || value == "" {
				continue
			}
			// 区分属性是全局的还是组的
			var success bool
			if groupProps.IsNil() {
				success = SetServerGroupProperties(GlobalProperties, key, value)
			} else {
				success = SetServerGroupProperties(groupProps, key, value)
			}
			if !success {
				var serverGroup = parseServerName(key, value)
				if serverGroup != nil {
					serverGroup.Props = NewProperties()
					serverGroup.Props.SetProperties(GlobalProperties)
					ServerGroupMap[strings.ToLower(key)] = serverGroup
				}
			}
		}
	}
}

func SetServerGroupProperties(props *Properties, key string, value string) bool {
	if util.StringUtil.EqualsIgnoreCase(key, "TIME_ZONE") {
		props.Set("localTimezone", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SESS_ENCODE") {
		if IsSupportedCharset(value) {
			props.Set("sessEncode", value)
		}
	} else if util.StringUtil.EqualsIgnoreCase(key, "ENABLE_RS_CACHE") {
		props.Set("enRsCache", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RS_CACHE_SIZE") {
		props.Set("rsCacheSize", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RS_REFRESH_FREQ") {
		props.Set("rsRefreshFreq", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOGIN_PRIMARY") || util.StringUtil.EqualsIgnoreCase(key, "LOGIN_MODE") {
		props.Set("loginMode", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOGIN_STATUS") {
		props.Set("loginStatus", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SWITCH_TIME") {
		props.Set("switchTimes", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "SWITCH_INTERVAL") {
		props.Set("switchInterval", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "PRIMARY_KEY") || util.StringUtil.EqualsIgnoreCase(key, "KEYWORDS") {
		props.Set("keyWords", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "COMPRESS_MSG") {
		props.Set("compress", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "COMPRESS_ID") {
		props.Set("compressID", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOGIN_ENCRYPT") || util.StringUtil.EqualsIgnoreCase(key, "COMMUNICATION_ENCRYPT") {
		props.Set("loginEncrypt", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "DIRECT") {
		props.Set("direct", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "DEC2DOUB") {
		props.Set("dec2Double", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_SEPARATE") {
		props.Set("rwSeparate", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_PERCENT") {
		props.Set("rwPercent", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_AUTO_DISTRIBUTE") {
		props.Set("rwAutoDistribute", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "COMPATIBLE_MODE") {
		props.Set("compatibleMode", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "CIPHER_PATH") {
		props.Set("cipherPath", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOAD_BALANCE") {
		props.Set("loadBalance", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOAD_BALANCE_PERCENT") {
		props.Set("loadBalancePercent", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "DO_SWITCH") {
		props.Set("doSwitch", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LANGUAGE") {
		props.Set("language", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOAD_BALANCE_FREQ") {
		props.Set("loadBalanceFreq", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "DB_ALIVE_CHECK_FREQ") {
		props.Set("dbAliveCheckFreq", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "RW_STANDBY_RECOVER_TIME") {
		props.Set("rwStandbyRecoverTime", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_LEVEL") {
		props.Set("logLevel", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_DIR") {
		props.Set("logDir", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_BUFFER_POOL_SIZE") {
		props.Set("logBufferPoolSize", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_BUF_SIZE") {
		props.Set("logBufferSize", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_FLUSHER_QUEUE_SIZE") {
		props.Set("logFlusherQueueSize", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOG_FLUSH_FREQ") {
		props.Set("logFlushFreq", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_ENABLE") {
		props.Set("statEnable", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_DIR") {
		props.Set("statDir", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_FLUSH_FREQ") {
		props.Set("statFlushFreq", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_HIGH_FREQ_SQL_COUNT") {
		props.Set("statHighFreqSqlCount", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_SLOW_SQL_COUNT") {
		props.Set("statSlowSqlCount", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_MAX_SQL_COUNT") {
		props.Set("statSqlMaxCount", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "STAT_SQL_REMOVE_LATEST") {
		props.Set("statSqlRemoveMode", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "ADDRESS_REMAP") {
		address := strings.Split(value, ",")
		if len(address) == 2 {
			addressRemap[address[0]] = address[1]
		}
	} else if util.StringUtil.EqualsIgnoreCase(key, "USER_REMAP") {
		user := strings.Split(value, ",")
		if len(user) == 2 {
			userRemap[user[0]] = user[1]
		}
	} else if util.StringUtil.EqualsIgnoreCase(key, "CONNECT_TIMEOUT") {
		props.Set("connectTimeout", value)
	} else if util.StringUtil.EqualsIgnoreCase(key, "LOGIN_CERTIFICATE") {
		props.Set("loginCertificate", value)
	} else {
		return false
	}
	return true
}

func parseServerName(name string, value string) *DBGroup {
	values := strings.Split(value, ",")

	var tmpVals []string
	var tmpName string
	var tmpPort int
	var svrList = make([]DB, 0, len(values))

	for _, v := range values {

		var tmp DB
		// 先查找IPV6,以[]包括
		begin := strings.IndexByte(v, '[')
		end := -1
		if begin != -1 {
			end = strings.IndexByte(v[begin:], ']')
		}
		if end != -1 {
			tmpName = v[begin+1 : end]

			// port
			if portIndex := strings.IndexByte(v[end:], ':'); portIndex != -1 {
				tmpPort, _ = strconv.Atoi(strings.TrimSpace(v[portIndex+1:]))
			} else {
				tmpPort = DEFAULT_PORT
			}
			tmp = DB{host: tmpName, port: tmpPort}
			svrList = append(svrList, tmp)
			continue
		}
		// IPV4
		tmpVals = strings.Split(v, ":")
		tmpName = strings.TrimSpace(tmpVals[0])
		if len(tmpVals) >= 2 {
			tmpPort, _ = strconv.Atoi(tmpVals[1])
		} else {
			tmpPort = DEFAULT_PORT
		}
		tmp = DB{host: tmpName, port: tmpPort}
		svrList = append(svrList, tmp)
	}

	if len(svrList) == 0 {
		return nil
	}
	return &DBGroup{
		Name:       name,
		ServerList: svrList,
	}
}

func setDriverAttributes(props *Properties) {
	if props == nil || props.Len() == 0 {
		return
	}

	parseLanguage(props.GetString(LanguageKey, "cn"))
	DbAliveCheckFreq = props.GetInt(DbAliveCheckFreqKey, DbAliveCheckFreqDef, 1, int(INT32_MAX))
	LoadBalanceFreq = props.GetInt(LoadBalanceFreqKey, LoadBalanceFreqDef, 1, int(INT32_MAX))

	// log
	LogLevel = ParseLogLevel(props)
	LogDir = util.StringUtil.FormatDir(props.GetTrimString(LogDirKey, LogDirDef))
	LogBufferSize = props.GetInt(LogBufferSizeKey, LogBufferSizeDef, 1, int(INT32_MAX))
	LogFlushFreq = props.GetInt(LogFlushFreqKey, LogFlushFreqDef, 1, int(INT32_MAX))
	LogFlushQueueSize = props.GetInt(LogFlusherQueueSizeKey, LogFlushQueueSizeDef, 1, int(INT32_MAX))

	// stat
	StatEnable = props.GetBool(StatEnableKey, StatEnableDef)
	StatDir = util.StringUtil.FormatDir(props.GetTrimString(StatDirKey, StatDirDef))
	StatFlushFreq = props.GetInt(StatFlushFreqKey, StatFlushFreqDef, 1, int(INT32_MAX))
	StatHighFreqSqlCount = props.GetInt(StatHighFreqSqlCountKey, StatHighFreqSqlCountDef, 0, 1000)
	StatSlowSqlCount = props.GetInt(StatSlowSqlCountKey, StatSlowSqlCountDef, 0, 1000)
	StatSqlMaxCount = props.GetInt(StatSqlMaxCountKey, StatSqlMaxCountDef, 0, 100000)
	parseStatSqlRemoveMode(props)
}

func parseLanguage(value string) {
	if util.StringUtil.EqualsIgnoreCase("cn", value) {
		Locale = 0
	} else if util.StringUtil.EqualsIgnoreCase("en", value) {
		Locale = 1
	}
}

func ParseLogLevel(props *Properties) int {
	logLevel := LOG_OFF
	value := props.GetString(LogLevelKey, "")
	if value != "" && !util.StringUtil.IsDigit(value) {
		if util.StringUtil.EqualsIgnoreCase("debug", value) {
			logLevel = LOG_DEBUG
		} else if util.StringUtil.EqualsIgnoreCase("info", value) {
			logLevel = LOG_INFO
		} else if util.StringUtil.EqualsIgnoreCase("sql", value) {
			logLevel = LOG_SQL
		} else if util.StringUtil.EqualsIgnoreCase("warn", value) {
			logLevel = LOG_WARN
		} else if util.StringUtil.EqualsIgnoreCase("error", value) {
			logLevel = LOG_ERROR
		} else if util.StringUtil.EqualsIgnoreCase("off", value) {
			logLevel = LOG_OFF
		} else if util.StringUtil.EqualsIgnoreCase("all", value) {
			logLevel = LOG_ALL
		}
	} else {
		logLevel = props.GetInt(LogLevelKey, logLevel, LOG_OFF, LOG_INFO)
	}

	return logLevel
}

func parseStatSqlRemoveMode(props *Properties) {
	value := props.GetString(StatSqlRemoveModeKey, "")
	if value != "" && !util.StringUtil.IsDigit(value) {
		if util.StringUtil.EqualsIgnoreCase("oldest", value) || util.StringUtil.EqualsIgnoreCase("eldest", value) {
			StatSqlRemoveMode = STAT_SQL_REMOVE_OLDEST
		} else if util.StringUtil.EqualsIgnoreCase("latest", value) {
			StatSqlRemoveMode = STAT_SQL_REMOVE_LATEST
		}
	} else {
		StatSqlRemoveMode = props.GetInt(StatSqlRemoveModeKey, StatSqlRemoveModeDef, 1, 2)
	}
}

func IsSupportedCharset(charset string) bool {
	if util.StringUtil.EqualsIgnoreCase(ENCODING_UTF8, charset) || util.StringUtil.EqualsIgnoreCase(ENCODING_GB18030, charset) || util.StringUtil.EqualsIgnoreCase(ENCODING_EUCKR, charset) {
		return true
	}
	return false
}
