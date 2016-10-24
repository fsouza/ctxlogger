// Copyright 2016 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctxlogger

import (
	"context"
	"net/http"

	"github.com/NYTimes/gizmo/web"
	"github.com/Sirupsen/logrus"
)

// ContextKey is the key used by the middleware to set the logger.
var ContextKey = &struct{ key string }{key: "ctxlogger"}

// ContextLogger takes the logger and returns the middleware that will always
// add the logger to the request context. It also expands the logger with any
// path variable on the given request.
func ContextLogger(baseLogger *logrus.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := baseLogger
			if vars := web.Vars(r); len(vars) > 0 {
				logger = varsLogger(vars, baseLogger)
			}
			h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ContextKey, logger)))
		})
	}
}

func varsLogger(vars map[string]string, logger *logrus.Logger) *logrus.Logger {
	newLogger := logrus.Logger{
		Out:       logger.Out,
		Formatter: logger.Formatter,
		Hooks:     logger.Hooks,
		Level:     logger.Level,
	}
	newLogger.Hooks.Add(&varsLogHook{vars})
	return &newLogger
}

type varsLogHook struct {
	vars map[string]string
}

func (h *varsLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *varsLogHook) Fire(e *logrus.Entry) error {
	for k, v := range h.vars {
		if _, ok := e.Data[k]; !ok {
			e.Data[k] = v
		}
	}
	return nil
}
