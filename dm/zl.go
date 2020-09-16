/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"dmdb_exporter/dm/util"
	"io"
)

const (
	READ_LEN = Dm_build_117
)

type iOffRowBinder interface {
	read(buf *util.Dm_build_834)
	isReadOver() bool
	getObj() interface{}
}

type offRowBinder struct {
	obj          interface{}
	encoding     string
	readOver     bool
	buffer       *util.Dm_build_834
	position     int32
	offRow       bool
	targetLength int64
}

func newOffRowBinder(obj interface{}, encoding string, targetLength int64) *offRowBinder {
	return &offRowBinder{
		obj:          obj,
		encoding:     encoding,
		targetLength: targetLength,
		readOver:     false,
		buffer:       util.Dm_build_838(),
		position:     0,
	}
}

type offRowBytesBinder struct {
	*offRowBinder
}

func newOffRowBytesBinder(obj []byte, encoding string) *offRowBytesBinder {
	var binder = &offRowBytesBinder{
		newOffRowBinder(obj, encoding, int64(IGNORE_TARGET_LENGTH)),
	}
	binder.read(binder.buffer)
	binder.offRow = binder.buffer.Dm_build_839() > Dm_build_114
	return binder
}

func (b *offRowBytesBinder) read(buf *util.Dm_build_834) {
	if b.buffer.Dm_build_839() > 0 {
		buf.Dm_build_871(b.buffer)
	} else if !b.readOver {
		var obj = b.obj.([]byte)
		buf.Dm_build_860(obj, 0, len(obj))
		b.readOver = true
	}
}

func (b *offRowBytesBinder) isReadOver() bool {
	return b.readOver
}

func (b *offRowBytesBinder) getObj() interface{} {
	return b.obj
}

type offRowBlobBinder struct {
	*offRowBinder
}

func newOffRowBlobBinder(blob DmBlob, encoding string) *offRowBlobBinder {
	var binder = &offRowBlobBinder{
		newOffRowBinder(blob, encoding, int64(IGNORE_TARGET_LENGTH)),
	}
	binder.read(binder.buffer)
	binder.offRow = binder.buffer.Dm_build_839() > Dm_build_114
	return binder
}

func (b *offRowBlobBinder) read(buf *util.Dm_build_834) {
	if b.buffer.Dm_build_839() > 0 {
		buf.Dm_build_871(b.buffer)
	} else if !b.readOver {
		var obj = b.obj.(DmBlob)
		var totalLen, _ = obj.GetLength()
		var leaveLen = totalLen - int64(b.position)
		var readLen = int32(leaveLen)
		if leaveLen > READ_LEN {
			readLen = READ_LEN
		}
		var bytes, _ = obj.getBytes(int64(b.position)+1, readLen)
		b.position += readLen
		if b.position == int32(totalLen) {
			b.readOver = true
		}
		buf.Dm_build_860(bytes, 0, len(bytes))
	}
}

func (b *offRowBlobBinder) isReadOver() bool {
	return b.readOver
}

func (b *offRowBlobBinder) getObj() interface{} {
	return b.obj
}

type offRowClobBinder struct {
	*offRowBinder
}

func newOffRowClobBinder(clob DmClob, encoding string) *offRowClobBinder {
	var binder = &offRowClobBinder{
		newOffRowBinder(clob, encoding, int64(IGNORE_TARGET_LENGTH)),
	}
	binder.read(binder.buffer)
	binder.offRow = binder.buffer.Dm_build_839() > Dm_build_114
	return binder
}

func (b *offRowClobBinder) read(buf *util.Dm_build_834) {
	if b.buffer.Dm_build_839() > 0 {
		buf.Dm_build_871(b.buffer)
	} else if !b.readOver {
		var obj = b.obj.(DmClob)
		var totalLen, _ = obj.GetLength()
		var leaveLen = totalLen - int64(b.position)
		var readLen = int32(leaveLen)
		if leaveLen > READ_LEN {
			readLen = READ_LEN
		}
		var str, _ = obj.getSubString(int64(b.position)+1, readLen)
		var bytes = util.Dm_build_586.Dm_build_793(str, b.encoding)
		b.position += readLen
		if b.position == int32(totalLen) {
			b.readOver = true
		}
		buf.Dm_build_860(bytes, 0, len(bytes))
	}
}

func (b *offRowClobBinder) isReadOver() bool {
	return b.readOver
}

func (b *offRowClobBinder) getObj() interface{} {
	return b.obj
}

type offRowReaderBinder struct {
	*offRowBinder
}

func newOffRowReaderBinder(reader io.Reader, encoding string) *offRowReaderBinder {
	var binder = &offRowReaderBinder{
		newOffRowBinder(reader, encoding, int64(IGNORE_TARGET_LENGTH)),
	}
	binder.read(binder.buffer)
	binder.offRow = binder.buffer.Dm_build_839() > Dm_build_114
	return binder
}

func (b *offRowReaderBinder) read(buf *util.Dm_build_834) {
	if b.buffer.Dm_build_839() > 0 {
		buf.Dm_build_871(b.buffer)
	} else if !b.readOver {
		var err error
		var readLen = READ_LEN
		var reader = b.obj.(io.Reader)
		var bytes = make([]byte, readLen)
		readLen, err = reader.Read(bytes)
		if err == io.EOF {
			b.readOver = true
			return
		}
		b.position += int32(readLen)
		if readLen < len(bytes) || b.targetLength != int64(IGNORE_TARGET_LENGTH) && int64(b.position) == b.targetLength {
			b.readOver = true
		}
		buf.Dm_build_860(bytes, 0, len(bytes))
	}
}

func (b *offRowReaderBinder) readAll() []byte {
	var byteArray = util.Dm_build_838()
	b.read(byteArray)
	for !b.readOver {
		b.read(byteArray)
	}
	return byteArray.Dm_build_881()
}

func (b *offRowReaderBinder) isReadOver() bool {
	return b.readOver
}

func (b *offRowReaderBinder) getObj() interface{} {
	return b.obj
}
