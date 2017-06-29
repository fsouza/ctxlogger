package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/web"
	"github.com/fsouza/ctxlogger"
	"github.com/sirupsen/logrus"
)

func (s *SimpleService) getMostPopular(r *http.Request) (int, interface{}, error) {
	logger := r.Context().Value(ctxlogger.ContextKey).(*logrus.Logger)
	logger.WithField("handler", "most-popular").Error("where are the popular cats?!")
	resourceType := web.Vars(r)["resourceType"]
	section := web.Vars(r)["section"]
	timeframe := web.GetUInt64Var(r, "timeframe")
	res, err := s.client.GetMostPopular(resourceType, section, uint(timeframe))
	if err != nil {
		return http.StatusInternalServerError, nil, &jsonErr{err.Error()}
	}
	return http.StatusOK, res, nil
}

type jsonErr struct {
	Err string `json:"error"`
}

func (e *jsonErr) Error() string {
	return e.Err
}
