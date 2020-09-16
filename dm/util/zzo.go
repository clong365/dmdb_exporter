/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package util

import (
	"os"
	"runtime"
	"strings"
)

const (
	PathSeparator     = string(os.PathSeparator)
	PathListSeparator = string(os.PathListSeparator)
)

var (
	goRoot = os.Getenv("GOROOT")
	goPath = os.Getenv("GOPATH")
)

type fileUtil struct {
}

var FileUtil = &fileUtil{}

func (fileUtil *fileUtil) Exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func (fileUtil *fileUtil) Search(relativePath string) (path string) {

	if strings.Contains(runtime.GOOS, "windows") {
		relativePath = strings.ReplaceAll(relativePath, "/", "\\")
	}

	pathArray := make([]string, 0)

	if workDir, _ := os.Getwd(); fileUtil.Exists(workDir) {
		pathArray = append(pathArray, workDir)
	}

	if fileUtil.Exists(goPath) {
		for _, s := range strings.Split(goPath, PathListSeparator) {
			pathArray = append(pathArray, s)
		}
	}

	if fileUtil.Exists(goRoot) {
		pathArray = append(pathArray, goRoot)
	}

	for _, path = range pathArray {
		path = path  + relativePath

		if fileUtil.Exists(path) {
			return path
		}
	}

	return ""
}
