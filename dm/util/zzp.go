/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package util

const (
	LINE_SEPARATOR = "\n"
)

func SliceEquals(src []byte, dest []byte) bool {
	if len(src) != len(dest) {
		return false
	}

	for i, _ := range src {
		if src[i] != dest[i] {
			return false
		}
	}

	return true
}
