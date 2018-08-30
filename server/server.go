package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dhaifley/dapi/lib"
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server values implement API server functionality.
type Server struct {
	Log    logrus.FieldLogger
	Router *mux.Router
	Auth   ptypes.AuthClient
}

// CheckAuth authenticates the provided token using the dauth service.
func (s *Server) CheckAuth(req *ptypes.AuthRequest) <-chan *dlib.Result {
	ch := make(chan *dlib.Result)
	go func() {
		defer close(ch)
		res, err := s.Auth.Auth(context.Background(), req)
		retry := 0
		for retry < 10 && err != nil {
			if err.Error() != "rpc error: code = Unavailable desc = transport is closing" {
				ch <- dlib.NewErrorResult(err)
				return
			}

			res, err = s.Auth.Auth(context.Background(), req)
			retry++
		}

		if retry >= 10 || res == nil {
			ch <- dlib.NewErrorResult(err)
			return
		}

		if !res.Ok {
			ch <- dlib.NewErrorResult(dlib.NewError(
				http.StatusUnauthorized, "unauthorized user"))
		}

		ch <- dlib.NewResult(req, res, "result", 0, "authentication successful", nil, nil)
	}()

	return ch
}

// AuthHandler wraps an http handler function with authentication verification.
func (s *Server) AuthHandler(handler http.Handler, perm *dauth.Perm) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathparts := strings.Split(r.URL.Path, "/")
		if r.Header.Get("Token") == "" || len(pathparts) < 2 {
			s.RespondWithError(dlib.NewError(http.StatusUnauthorized, "unauthorized request"), w, r)
			return
		}

		preq := perm.ToRequest()
		areq := ptypes.AuthRequest{
			Token: &ptypes.TokenRequest{Token: r.Header.Get("Token")},
			Perm:  &preq,
		}

		ac := s.CheckAuth(&areq)
		for ar := range ac {
			if ar.Err != nil {
				s.RespondWithError(ar.Err, w, r)
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}

// Header wraps a handler function to set default header values.
func (s *Server) Header(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		handler.ServeHTTP(w, r)
	})
}

// Logger wraps a handler function with logging functionality.
func (s *Server) Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		s.Log.WithFields(logrus.Fields{
			"method":  r.Method,
			"uri":     r.RequestURI,
			"remote":  r.RemoteAddr,
			"elapsed": time.Since(start).String(),
		}).Info("Request processed")
	})
}

// RespondWithError responds to the current request with a standard error response.
func (s *Server) RespondWithError(err error, w http.ResponseWriter, r *http.Request) {
	s.Log.WithFields(logrus.Fields{
		"method": r.Method,
		"uri":    r.RequestURI,
		"remote": r.RemoteAddr,
	}).Error(err)

	switch v := err.(type) {
	case *dlib.Error:
		w.WriteHeader(v.Code)
		if err := json.NewEncoder(w).Encode(v); err != nil {
			s.Log.Error(err)
		}

		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(dlib.Error{
			Code: http.StatusInternalServerError,
			Msg:  v.Error(),
		}); err != nil {
			s.Log.Error(err)
		}
	}
}

// RespondNotFound responds to the current request with a 404 not found error.
func (s *Server) RespondNotFound(w http.ResponseWriter, r *http.Request) {
	s.RespondWithError(dlib.NewError(http.StatusNotFound,
		"resource not found"), w, r)
}

// NotFoundHandler is the handler function for 404 errors.
func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.RespondNotFound(w, r)
}

// GetDocs is the handler function for documentation requests.
func (s *Server) GetDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if _, err := os.Stat("docs/docs.html"); err == nil {
		http.ServeFile(w, r, "docs/docs.html")
		return
	}

	if _, err := os.Stat("../docs/docs.html"); err == nil {
		http.ServeFile(w, r, "../docs/docs.html")
		return
	}

	s.RespondNotFound(w, r)
}

// GetIndex is the handler function for the root path.
func (s *Server) GetIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(lib.ServiceInfo); err != nil {
		s.Log.Error(err)
	}
}

// GetIcon is the handler function for the application icon.
func (s *Server) GetIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	if _, err := os.Stat("docs/favicon.ico"); err == nil {
		http.ServeFile(w, r, "docs/favicon.ico")
		return
	}

	if _, err := os.Stat("../docs/favicon.ico"); err == nil {
		http.ServeFile(w, r, "../docs/favicon.ico")
		return
	}

	s.RespondNotFound(w, r)
}
