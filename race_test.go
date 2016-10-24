// Copyright 2016 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxlogger

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestVarsLoggerIsSafe(t *testing.T) {
	const N = 32
	var b bytes.Buffer
	logger := logrus.New()
	logger.Out = &b
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}
	logger = varsLogger(map[string]string{"name": "gopher"}, logger)
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			logger.WithField("i", i).Info("be advised")
			wg.Done()
		}(i)
	}
	wg.Wait()
	logLines := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(logLines) != N {
		t.Errorf("wrong log lines returned, wanted %d log lines, got %d:\n%s", N, len(logLines), b.String())
	}
}
