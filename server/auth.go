package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
)

// Authenticate is the get handler for validating tokens.
func (s *Server) Authenticate(w http.ResponseWriter, r *http.Request) {
	req := ptypes.AuthRequest{}
	req.Token = &ptypes.TokenRequest{Token: r.URL.Query().Get("token")}
	req.Perm = &ptypes.PermRequest{Service: r.URL.Query().Get("service"), Name: r.URL.Query().Get("name")}
	u := dauth.User{}
	ch := s.CheckAuth(&req)
	for ar := range ch {
		if ar.Err != nil {
			s.RespondWithError(ar.Err, w, r)
			return
		}

		switch v := ar.Val.(type) {
		case *ptypes.AuthResponse:
			err := u.FromResponse(v.User)
			if err != nil {
				s.RespondWithError(err, w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(u); err != nil {
				s.Log.Error(err.Error())
			}

			return
		default:
			continue
		}
	}

	s.RespondWithError(dlib.NewError(http.StatusInternalServerError,
		"invalid authentication resposnse"), w, r)
}

// Login is the post handler for authorizing new tokens.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var u dauth.User
	err := dec.Decode(&u)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	req := u.ToRequest()
	res, err := s.Auth.Login(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	t := dauth.Token{}
	err = t.FromResponse(res)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		s.Log.Error(err.Error())
	}
}

// Logout is the post handeler for destroying tokens.
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var t dauth.Token
	err := dec.Decode(&t)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	req := t.ToRequest()
	res, err := s.Auth.Logout(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	tk := dauth.Token{}
	err = tk.FromResponse(res)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tk); err != nil {
		s.Log.Error(err.Error())
	}
}
