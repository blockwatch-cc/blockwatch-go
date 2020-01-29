// Copyright (c) 2020 Blockwatch Data Inc.
// Author: alex@blockwatch.cc

package blockwatch

import (
	"bytes"
)

func indexByteColumnN(b []byte, sep byte, n int) (int, int) {
	var start int
	// find the n-th start offset
	for ; start != -1 && n > 0; n-- {
		n := bytes.IndexByte(b[start:], sep)
		if n == -1 {
			return -1, 0
		}
		start += n + 1
	}
	// find the next offset
	end := bytes.IndexByte(b[start:], sep)
	if end < 0 {
		return start, len(b) + 1
	}
	return start, end + start
}
