package handler

import (
	"net/http"
	"github.com/wolfogre/qiniuauth/internal/log"
	"strings"
	"go.uber.org/zap"
	"github.com/wolfogre/qiniuauth/internal/judge"
	"fmt"
)

type Handler struct {

}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger.With(
		"ip", strings.Split(r.RemoteAddr, ":")[0],
		"method", r.Method,
		"host", r.Host,
		"url", r.URL.String(),
		"user_agent", r.UserAgent(),
	)

	switch r.Method {
	case "HEAD":
		switch r.URL.Path {
		case "/auth":
			if ok, msg := judge.Status(); ok {
				response(logger, w, http.StatusOK, msg)
			} else {
				response(logger, w, http.StatusInternalServerError, msg)
			}
		default:
			abort(logger, w, http.StatusNotFound)
		}
	case "GET":
		switch r.URL.Path {
		case "/_status":
			if ok, msg := judge.Status(); ok {
				response(logger, w, http.StatusOK, msg)
			} else {
				response(logger, w, http.StatusInternalServerError, msg)
			}
		default:
			abort(logger, w, http.StatusNotFound)
		}
	default:
		abort(logger, w, http.StatusMethodNotAllowed)
	}
}

func abort(logger *zap.SugaredLogger, w http.ResponseWriter, code int) {
	response(logger, w, code, http.StatusText(code))
}

func response(logger *zap.SugaredLogger, w http.ResponseWriter, code int, content string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, content)
	logger.With(
		"status_code", code,
	)
	if code >= 400 {
		logger.Error(content)
	} else {
		logger.Info(content)
	}
}