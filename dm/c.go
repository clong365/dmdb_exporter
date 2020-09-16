/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"bytes"
	"compress/zlib"
	"dmdb_exporter/dm/util"
	"github.com/golang/snappy"
)

func Compress(srcBuffer *util.Dm_build_912, offset int, length int, compressID int) ([]byte, error) {
	if compressID == Dm_build_83 {
		return snappy.Encode(nil, srcBuffer.Dm_build_1237(offset, make([]byte, length))), nil
	}
	return GzlibCompress(srcBuffer, offset, length)
}

func UnCompress(srcBytes []byte, compressID int) ([]byte, error) {
	if compressID == Dm_build_83 {
		return snappy.Decode(nil, srcBytes)
	}
	return GzlibUncompress(srcBytes)
}

func GzlibCompress(srcBuffer *util.Dm_build_912, offset int, length int) ([]byte, error) {
	var ret bytes.Buffer
	var w = zlib.NewWriter(&ret)
	w.Write(srcBuffer.Dm_build_1237(offset, make([]byte, length)))
	w.Close()
	return ret.Bytes(), nil
}

func GzlibUncompress(srcBytes []byte) ([]byte, error) {
	var bytesBuf = new(bytes.Buffer)
	r, err := zlib.NewReader(bytes.NewReader(srcBytes))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	_, err = bytesBuf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return bytesBuf.Bytes(), nil
}
