// Copyright 2016 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxlogger

import (
	"bytes"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

type goroutineHook struct {
	test.Hook
	t        *testing.T
	nentries int32
}

func (h *goroutineHook) Fire(e *logrus.Entry) error {
	go func() {
		for k, v := range e.Data {
			h.t.Logf("%s: %v", k, v)
		}
		atomic.AddInt32(&h.nentries, 1)
	}()
	return nil
}

func TestVarsLoggerIsSafe(t *testing.T) {
	var fakeHook test.Hook
	const N = 32
	var b bytes.Buffer
	logger := logrus.New()
	logger.Out = &safeWriter{w: &b}
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Hooks.Add(&fakeHook)
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			innerLogger := varsLogger(map[string]string{"name": "gopher"}, logger)
			innerLogger.WithField("some", "thing").Info("be advised")
			wg.Done()
		}(i)
	}
	wg.Wait()
	logLines := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(logLines) != N {
		t.Errorf("wrong log lines returned, wanted %d log lines, got %d:\n%s", N, len(logLines), b.String())
	}
	if len(fakeHook.Entries) != N {
		t.Errorf("wrong number of entries in the hook. want %d, got %d", N, len(fakeHook.AllEntries()))
	}
}

func TestAlwaysFirstInTheListOfLoggers(t *testing.T) {
	fakeHook := goroutineHook{t: t}
	const N = 32
	var b bytes.Buffer
	logger := logrus.New()
	logger.Out = &safeWriter{w: &b}
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Hooks.Add(&fakeHook)
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			innerLogger := varsLogger(map[string]string{"name": "gopher"}, logger)
			innerLogger.WithField("some", "thing").Info("be advised")
			time.Sleep(500 * time.Millisecond)
			wg.Done()
		}(i)
	}
	wg.Wait()
	logLines := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(logLines) != N {
		t.Errorf("wrong log lines returned, wanted %d log lines, got %d:\n%s", N, len(logLines), b.String())
	}
	nentries := atomic.LoadInt32(&fakeHook.nentries)
	if nentries != N {
		t.Errorf("wrong number of entries in the hook. want %d, got %d", N, nentries)
	}
}
