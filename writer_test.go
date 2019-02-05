// Copyright 2019 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxlogger

import (
	"io"
	"sync"
)

type safeWriter struct {
	w   io.Writer
	mtx sync.Mutex
}

func (w *safeWriter) Write(p []byte) (int, error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	return w.w.Write(p)
}
