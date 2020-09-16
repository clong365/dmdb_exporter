/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"dmdb_exporter/dm/util"
	"io"
)

type DmBlob struct {
	lob
	data   []byte
	offset int64
}

func newDmBlob() *DmBlob {
	return &DmBlob{
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
		offset: 1,
	}
}

func newBlobFromDB(value []byte, conn *DmConnection, column *column, fetchAll bool) *DmBlob {
	var blob = newDmBlob()
	blob.connection = conn
	blob.lobFlag = LOB_FLAG_BYTE
	blob.compatibleOracle = conn.CompatibleOracle()
	blob.local = false
	blob.updateable = !column.readonly
	blob.tabId = column.lobTabId
	blob.colId = column.lobColId

	blob.inRow = util.Dm_build_586.Dm_build_677(value, NBLOB_HEAD_IN_ROW_FLAG) == LOB_IN_ROW
	blob.blobId = util.Dm_build_586.Dm_build_691(value, NBLOB_HEAD_BLOBID)
	if !blob.inRow {
		blob.groupId = util.Dm_build_586.Dm_build_681(value, NBLOB_HEAD_OUTROW_GROUPID)
		blob.fileId = util.Dm_build_586.Dm_build_681(value, NBLOB_HEAD_OUTROW_FILEID)
		blob.pageNo = util.Dm_build_586.Dm_build_686(value, NBLOB_HEAD_OUTROW_PAGENO)
	}
	if conn.NewLobFlag {
		blob.tabId = util.Dm_build_586.Dm_build_686(value, NBLOB_EX_HEAD_TABLE_ID)
		blob.colId = util.Dm_build_586.Dm_build_681(value, NBLOB_EX_HEAD_COL_ID)
		blob.rowId = util.Dm_build_586.Dm_build_691(value, NBLOB_EX_HEAD_ROW_ID)
		blob.exGroupId = util.Dm_build_586.Dm_build_681(value, NBLOB_EX_HEAD_FPA_GRPID)
		blob.exFileId = util.Dm_build_586.Dm_build_681(value, NBLOB_EX_HEAD_FPA_FILEID)
		blob.exPageNo = util.Dm_build_586.Dm_build_686(value, NBLOB_EX_HEAD_FPA_PAGENO)
	}
	blob.resetCurrentInfo()

	blob.length = blob.getLengthFromHead(value)
	if blob.inRow {
		blob.data = make([]byte, blob.length)
		if conn.NewLobFlag {
			util.Dm_build_586.Dm_build_642(blob.data, 0, value, NBLOB_EX_HEAD_SIZE, len(blob.data))
		} else {
			util.Dm_build_586.Dm_build_642(blob.data, 0, value, NBLOB_INROW_HEAD_SIZE, len(blob.data))
		}
	} else if fetchAll {
		blob.loadAllData()
	}
	return blob
}

func newBlobOfLocal(value []byte, conn *DmConnection) *DmBlob {
	var blob = newDmBlob()
	blob.connection = conn
	blob.lobFlag = LOB_FLAG_BYTE
	blob.data = value
	blob.length = int64(len(blob.data))
	return blob
}

func NewBlob(value []byte) *DmBlob {
	var blob = newDmBlob()

	blob.lobFlag = LOB_FLAG_BYTE
	blob.data = value
	blob.length = int64(len(blob.data))
	return blob
}

func (blob *DmBlob) Read(dest []byte) (n int, err error) {
	result, err := blob.getBytes(blob.offset, int32(len(dest)))
	if err != nil {
		return 0, err
	}
	blob.offset += int64(len(result))
	copy(dest, result)
	if len(result) == 0 {
		return 0, io.EOF
	}
	return len(result), nil
}

func (blob *DmBlob) ReadAt(pos int, dest []byte) (n int, err error) {
	result, err := blob.getBytes(int64(pos), int32(len(dest)))
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, io.EOF
	}
	copy(dest[0:len(result)], result)
	return len(result), nil
}

func (blob *DmBlob) Write(pos int, src []byte) (n int, err error) {
	if err = blob.checkFreed(); err != nil {
		return
	}
	if pos < 1 {
		err = ECGO_INVALID_LENGTH_OR_OFFSET.throw()
		return
	}
	if !blob.updateable {
		err = ECGO_RESULTSET_IS_READ_ONLY.throw()
		return
	}
	pos -= 1
	if blob.local || blob.fetchAll {
		if int64(pos) > blob.length {
			err = ECGO_INVALID_LENGTH_OR_OFFSET.throw()
			return
		}
		blob.setLocalData(pos, src)
		n = len(src)
	} else {

		var writeLen, err = blob.connection.Access.dm_build_1456(blob, pos, src)
		if err != nil {
			return -1, err
		}

		if blob.groupId == -1 {
			blob.setLocalData(pos, src)
		} else {
			blob.inRow = false
			blob.length = -1
		}
		n = writeLen

	}
	blob.modify = true
	return
}

func (dest *DmBlob) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch src := src.(type) {
	case []byte:
		*dest = *NewBlob(src)
		return nil
	case *DmBlob:
		*dest = *src
		return nil
	default:
		return UNSUPPORTED_SCAN
	}
}

func (blob *DmBlob) getBytes(pos int64, length int32) ([]byte, error) {
	var err error
	var leaveLength int64
	if err = blob.checkFreed(); err != nil {
		return nil, err
	}
	if pos < 1 || length < 0 {
		return nil, ECGO_INVALID_LENGTH_OR_OFFSET.throw()
	}
	pos = pos - 1
	if leaveLength, err = blob.GetLength(); err != nil {
		return nil, err
	}
	leaveLength -= pos
	if leaveLength < 0 {
		return nil, ECGO_INVALID_LENGTH_OR_OFFSET.throw()
	}
	if int64(length) > leaveLength {
		length = int32(leaveLength)
	}
	if blob.local || blob.inRow || blob.fetchAll {
		return blob.data[pos : pos+int64(length)], nil
	} else {

		return blob.connection.Access.dm_build_1423(blob, int32(pos), length)
	}
}

func (blob *DmBlob) loadAllData() {
	blob.checkFreed()
	if blob.local || blob.inRow || blob.fetchAll {
		return
	}
	len, _ := blob.GetLength()
	blob.data, _ = blob.getBytes(1, int32(len))
	blob.fetchAll = true
}

func (blob *DmBlob) setLocalData(pos int, p []byte) {
	if pos+len(p) >= int(blob.length) {
		var tmp = make([]byte, pos+len(p))
		util.Dm_build_586.Dm_build_642(tmp, 0, blob.data, 0, pos)
		util.Dm_build_586.Dm_build_642(tmp, pos, p, 0, len(p))
		blob.data = tmp
	} else {
		util.Dm_build_586.Dm_build_642(blob.data, pos, p, 0, len(p))
	}
	blob.length = int64(len(blob.data))
}
