package util

import (
	"path"
	"runtime"
)

func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(b), "..")
	return dir
}
