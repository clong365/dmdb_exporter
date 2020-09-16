/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"bytes"
	"dmdb_exporter/dm/i18n"
	"fmt"
	"runtime"
)

// 驱动级错误
var (
	DSN_INVALID_SCHEMA       = newDmError(9001, "error.dsn.invalidSchema")
	UNSUPPORTED_SCAN         = newDmError(9002, "error.unsupported.scan")
	INVALID_PARAMETER_NUMBER = newDmError(9003, "error.invalidParameterNumber")
	THIRD_PART_CIPHER_INIT_FAILED = newDmError(9004, "error.initThirdPartCipherFailed")
)

var (
	ECGO_CONNECTION_SWITCH_FAILED    = newDmError(3001, "error.connectionSwitchFailed")
	ECGO_CONNECTION_SWITCHED         = newDmError(3002, "error.connectionSwitched")
	ECGO_INVALID_SERVER_MODE         = newDmError(4001, "error.invalidServerMode")
	ECGO_OSAUTH_ERROR                = newDmError(5001, "error.osauthError")
	ECGO_NOT_QUERY_SQL               = newDmError(5002, "error.notQuerySQL")
	ECGO_NOT_EXEC_SQL                = newDmError(5003, "error.notExecSQL")
	ECGO_INVALID_TRAN_ISOLATION      = newDmError(5004, "error.invalidTranIsolation")
	ECGO_COMMIT_IN_AUTOCOMMIT_MODE   = newDmError(5005, "errorCommitInAutoCommitMode")
	ECGO_ROLLBACK_IN_AUTOCOMMIT_MODE = newDmError(5006, "errorRollbackInAutoCommitMode")
	ECGO_STATEMENT_HANDLE_CLOSED     = newDmError(5007, "errorStatementHandleClosed")
	ECGO_RESULTSET_CLOSED            = newDmError(5008, "errorResultSetColsed")
	ECGO_COMMUNITION_ERROR           = newDmError(6001, "error.communicationError")
	ECGO_MSG_CHECK_ERROR             = newDmError(6002, "error.msgCheckError")
	ECGO_UNKOWN_NETWORK              = newDmError(6003, "error.unkownNetWork")
	ECGO_ERROR_SERVER_VERSION        = newDmError(6004, "error.serverVersion")
	ECGO_USERNAME_TOO_LONG           = newDmError(6005, "error.usernameTooLong")
	ECGO_PASSWORD_TOO_LONG           = newDmError(6006, "error.passwordTooLong")
	ECGO_DATA_TOO_LONG               = newDmError(6007, "error.dataTooLong")
	ECGO_INVALID_COLUMN_TYPE         = newDmError(6008, "error.invalidColumnType")
	ECGO_DATA_CONVERTION_ERROR       = newDmError(6009, "error.dataConvertionError")
	ECGO_INVALID_CONN                = newDmError(6010, "error.invalidConn")
	ECGO_INVALID_HEX                 = newDmError(6011, "error.invalidHex")
	ECGO_DATA_OVERFLOW               = newDmError(6012, "error.dataOverflow")
	ECGO_INVALID_DATETIME_FORMAT     = newDmError(6013, "error.invalidDateTimeFormat")
	ECGO_DATETIME_OVERFLOW           = newDmError(6014, "error.datetimeOverflow")
	ECGO_INVALID_TIME_INTERVAL       = newDmError(6015, "error.invalidTimeInterval")
	ECGO_UNSUPPORTED_INPARAM_TYPE    = newDmError(6016, "error.unsupportedInparamType")
	ECGO_UNSUPPORTED_OUTPARAM_TYPE   = newDmError(6017, "error.unsupportedOutparamType")
	ECGO_UNSUPPORTED_TYPE            = newDmError(6018, "error.unsupportedType")
	ECGO_INVALID_OBJ_BLOB            = newDmError(6019, "error.invalidObjBlob")
	ECGO_STRUCT_MEM_NOT_MATCH        = newDmError(6020, "error.structMemNotMatch")
	ECGO_INVALID_COMPLEX_TYPE_NAME   = newDmError(6021, "error.invalidComplexTypeName")
	ECGO_INVALID_PARAMETER_VALUE     = newDmError(6022, "error.invalidParamterValue")
	ECGO_INVALID_ARRAY_LEN           = newDmError(6023, "error.invalidArrayLen")
	//rows
	ECGO_INVALID_SEQUENCE_NUMBER = newDmError(6024, "error.invalidSequenceNumber")
	//lob
	ECGO_RESULTSET_IS_READ_ONLY   = newDmError(6029, "error.resultsetInReadOnlyStatus")
	ECGO_INIT_SSL_FAILED          = newDmError(20002, "error.SSLInitFailed")
	ECGO_LOB_FREED                = newDmError(20003, "error.LobDataHasFreed")
	ECGO_FATAL_ERROR              = newDmError(20004, "error.fatalError")
	ECGO_INVALID_LENGTH_OR_OFFSET = newDmError(6057, "error.invalidLenOrOffset")
	ECGO_INTERVAL_OVERFLOW        = newDmError(6066, "error.intervalValueOverflow")
	ECGO_INVALID_CIPHER           = newDmError(6069, "error.invalidCipher")
)

//Svr Msg Err
var (
	EC_RN_EXCEED_ROWSET_SIZE = &DmError{-7036, "", nil}
	EC_BP_WITH_ERROR         = &DmError{121, "", nil}
)

type DmError struct {
	ErrCode int32
	ErrText string
	stack   []uintptr
}

func newDmError(errCode int32, errText string) *DmError {
	de := new(DmError)
	de.ErrCode = errCode
	de.ErrText = i18n.Get(errText, Locale)
	de.stack = nil
	return de
}

func (dmError *DmError) throw() *DmError {
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:])
	dmError.stack = pcs[0:n]
	return dmError
}

func (dmError *DmError) FormatStack() string {
	if dmError == nil || dmError.stack == nil {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	index := 1
	space := "  "
	for _, p := range dmError.stack {
		if fn := runtime.FuncForPC(p - 1); fn != nil {
			file, line := fn.FileLine(p - 1)
			buffer.WriteString(fmt.Sprintf("   %d).%s%s\n    \t%s:%d\n", index, space, fn.Name(), file, line))
			index++
		}
	}
	return buffer.String()
}

func (dmError *DmError) Error() string {
	return fmt.Sprintf("Error %d: %s", dmError.ErrCode, dmError.ErrText) + "\n" + "stack info:\n" + dmError.FormatStack()
}

// 扩充ErrText
func (dmError *DmError) addDetail(detail string) *DmError {
	return &DmError{dmError.ErrCode, dmError.ErrText + detail, nil}
}
func (dmError *DmError) addDetailln(detail string) *DmError {
	return dmError.addDetail("\n" + detail)
}
