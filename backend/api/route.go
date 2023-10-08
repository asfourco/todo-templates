package api

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

type route mux.Route

func (r *route) String() string {
	rr := (*mux.Route)(r)
	name := rr.GetName()

	if rr.GetError() != nil {
		if name == "" {
			name = "<unnamed>"
		}
		zlog.Warn("routes", zap.String("route", name), zap.Error(rr.GetError()))
		return ""
	}

	if name == "" {
		return fmt.Sprintf("%s %s", r.formattedMethods(), r.formattedRegexp())
	}
	return fmt.Sprintf("%s %s (%s)", r.formattedMethods(), r.formattedRegexp(), name)
}

func (r *route) formattedMethods() string {
	methods, err := (*mux.Route)(r).GetMethods()
	if err != nil {
		return "[Any]"
	}
	return fmt.Sprintf("[%s]", strings.Join(methods, ", "))
}

func (r *route) formattedRegexp() string {
	pathRegexp, err := (*mux.Route)(r).GetPathRegexp()
	if err != nil {
		return "<None>"
	}
	return pathRegexp
}
