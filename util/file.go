/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

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
