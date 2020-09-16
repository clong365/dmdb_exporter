/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"dmdb_exporter/dm/util"
	"io"
)

type DmClob struct {
	lob
	data           []rune
	serverEncoding string
}

func newDmClob() *DmClob {
	return &DmClob{
		lob: lob{
			inRow:            true,
			groupId:          -1,
			fileId:           -1,
			pageNo:           -1,
			readOver:         false,
			local:            true,
			updateable:       true,
			length:           -1,
			compatibleOracle: false,
			fetchAll:         false,
			freed:            false,
			modify:           false,
		},
	}
}

func newClobFromDB(value []byte, conn *DmConnection, column *column, fetchAll bool) *DmClob {
	var clob = newDmClob()
	clob.connection = conn
	clob.lobFlag = LOB_FLAG_CHAR
	clob.compatibleOracle = conn.CompatibleOracle()
	clob.local = false
	clob.updateable = !column.readonly
	clob.tabId = column.lobTabId
	clob.colId = column.lobColId

	clob.inRow = util.Dm_build_586.Dm_build_677(value, NBLOB_HEAD_IN_ROW_FLAG) == LOB_IN_ROW
	clob.blobId = util.Dm_build_586.Dm_build_691(value, NBLOB_HEAD_BLOBID)
	if !clob.inRow {
		clob.groupId = util.Dm_build_586.Dm_build_681(value, NBLOB_HEAD_OUTROW_GROUPID)
		clob.fileId = util.Dm_build_586.Dm_build_681(value, NBLOB_HEAD_OUTROW_FILEID)
		clob.pageNo = util.Dm_build_586.Dm_build_686(value, NBLOB_HEAD_OUTROW_PAGENO)
	}
	if conn.NewLobFlag {
		clob.tabId = util.Dm_build_586.Dm_build_686(value, NBLOB_EX_HEAD_TABLE_ID)
		clob.colId = util.Dm_build_586.Dm_build_681(value, NBLOB_EX_HEAD_COL_ID)
		clob.rowId = util.Dm_build_586.Dm_build_691(value, NBLOB_EX_HEAD_ROW_ID)
		clob.exGroupId = util.Dm_build_586.Dm_build_681(value, NBLOB_EX_HEAD_FPA_GRPID)
		clob.exFileId = util.Dm_build_586.Dm_build_681(value, NBLOB_EX_HEAD_FPA_FILEID)
		clob.exPageNo = util.Dm_build_586.Dm_build_686(value, NBLOB_EX_HEAD_FPA_PAGENO)
	}
	clob.resetCurrentInfo()

	clob.serverEncoding = conn.getServerEncoding()
	if clob.inRow {
		if conn.NewLobFlag {
			clob.data = []rune(util.Dm_build_586.Dm_build_741(value, NBLOB_EX_HEAD_SIZE, int(clob.getLengthFromHead(value)), clob.serverEncoding))
		} else {
			clob.data = []rune(util.Dm_build_586.Dm_build_741(value, NBLOB_INROW_HEAD_SIZE, int(clob.getLengthFromHead(value)), clob.serverEncoding))
		}
		clob.length = int64(len(clob.data))
	} else if fetchAll {
		clob.loadAllData()
	}
	return clob
}

func newClobOfLocal(value string, conn *DmConnection) *DmClob {
	var clob = newDmClob()
	clob.connection = conn
	clob.lobFlag = LOB_FLAG_CHAR
	clob.data = []rune(value)
	clob.length = int64(len(clob.data))
	return clob
}

func NewClob(value string) *DmClob {
	var clob = newDmClob()

	clob.lobFlag = LOB_FLAG_CHAR
	clob.data = []rune(value)
	clob.length = int64(len(clob.data))
	return clob
}

func (clob *DmClob) ReadString(pos int, length int) (result string, err error) {
	result, err = clob.getSubString(int64(pos), int32(length))
	if err != nil {
		return
	}
	if len(result) == 0 {
		err = io.EOF
		return
	}
	return
}

func (clob *DmClob) WriteString(pos int, s string) (n int, err error) {
	if err = clob.checkFreed(); err != nil {
		return
	}
	if pos < 1 {
		err = ECGO_INVALID_LENGTH_OR_OFFSET.throw()
		return
	}
	if !clob.updateable {
		err = ECGO_RESULTSET_IS_READ_ONLY.throw()
		return
	}
	pos -= 1
	if clob.local || clob.fetchAll {
		if int64(pos) > clob.length {
			err = ECGO_INVALID_LENGTH_OR_OFFSET.throw()
			return
		}
		clob.setLocalData(pos, s)
		n = len(s)
	} else {

		var writeLen, err = clob.connection.Access.dm_build_1440(clob, pos, s, clob.serverEncoding)
		if err != nil {
			return -1, err
		}

		if clob.groupId == -1 {
			clob.setLocalData(pos, s)
		} else {
			clob.inRow = false
			clob.length = -1
		}
		n = writeLen
	}
	clob.modify = true
	return
}

func (dest *DmClob) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch src := src.(type) {
	case string:
		*dest = *NewClob(src)
		return nil
	case *DmClob:
		*dest = *src
		return nil
	default:
		return UNSUPPORTED_SCAN
	}
}

func (clob *DmClob) getSubString(pos int64, len int32) (string, error) {
	var err error
	var leaveLength int64
	if err = clob.checkFreed(); err != nil {
		return "", err
	}
	if pos < 1 || len < 0 {
		return "", ECGO_INVALID_LENGTH_OR_OFFSET.throw()
	}
	pos = pos - 1
	if leaveLength, err = clob.GetLength(); err != nil {
		return "", err
	}
	if pos > leaveLength {
		pos = leaveLength
	}
	leaveLength -= pos
	if leaveLength < 0 {
		return "", ECGO_INVALID_LENGTH_OR_OFFSET.throw()
	}
	if int64(len) > leaveLength {
		len = int32(leaveLength)
	}
	if clob.local || clob.inRow || clob.fetchAll {
		if pos > clob.length {
			return "", ECGO_INVALID_LENGTH_OR_OFFSET.throw()
		}
		return string(clob.data[pos : pos+int64(len)]), nil
	} else {

		return clob.connection.Access.dm_build_1431(clob, int32(pos), len)
	}
}

func (clob *DmClob) loadAllData() {
	clob.checkFreed()
	if clob.local || clob.inRow || clob.fetchAll {
		return
	}
	len, _ := clob.GetLength()
	s, _ := clob.getSubString(1, int32(len))
	clob.data = []rune(s)
	clob.fetchAll = true
}

func (clob *DmClob) setLocalData(pos int, str string) {
	if pos+len(str) >= int(clob.length) {
		clob.data = []rune(string(clob.data[0:pos]) + str)
	} else {
		clob.data = []rune(string(clob.data[0:pos]) + str + string(clob.data[pos+len(str):len(clob.data)]))
	}
	clob.length = int64(len(clob.data))
}
