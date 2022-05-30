// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package file

import (
	"io"
	"strconv"

	logsconfig "flashcat.cloud/categraf/config/logs"
	"flashcat.cloud/categraf/pkg/logs/auditor"
)

// Position returns the position from where logs should be collected.
func Position(registry auditor.Registry, identifier string, mode logsconfig.TailingMode) (int64, int, error) {
	var offset int64
	var whence int
	var err error

	value := registry.GetOffset(identifier)

	switch {
	case mode == logsconfig.ForceBeginning:
		offset, whence = 0, io.SeekStart
	case mode == logsconfig.ForceEnd:
		offset, whence = 0, io.SeekEnd
	case value != "":
		// an offset was registered, tailing mode is not forced, tail from the offset
		whence = io.SeekStart
		offset, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			offset = 0
			if mode == logsconfig.End {
				whence = io.SeekEnd
			} else if mode == logsconfig.Beginning {
				whence = io.SeekStart
			}
		}
	case mode == logsconfig.Beginning:
		offset, whence = 0, io.SeekStart
	case mode == logsconfig.End:
		fallthrough
	default:
		offset, whence = 0, io.SeekEnd
	}
	return offset, whence, err
}
