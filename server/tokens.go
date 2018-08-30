package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/gorilla/mux"
)

// GetTokens is the get handler function for tokens.
func (s *Server) GetTokens(w http.ResponseWriter, r *http.Request) {
	var data []dauth.Token
	q := dauth.Token{}
	if len(r.URL.Query()) != 0 {
		if err := q.FromQueryValues(r.URL.Query()); err != nil {
			s.RespondWithError(err, w, r)
			return
		}
	}

	req := q.ToRequest()
	stream, err := s.Auth.GetTokens(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			if len(data) == 0 {
				s.RespondNotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(data); err != nil {
				s.Log.Error(err)
			}

			return
		}

		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}

		v := dauth.Token{}
		err = v.FromResponse(res)
		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}

		data = append(data, v)
	}
}

// GetTokenByID is the get by id handler function for tokens.
func (s *Server) GetTokenByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	var v *dauth.Token
	req := ptypes.TokenRequest{ID: id}
	stream, err := s.Auth.GetTokens(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			if v == nil {
				s.RespondNotFound(w, r)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(v); err != nil {
				s.Log.Error(err)
			}

			return
		}

		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}

		v = new(dauth.Token)
		err = v.FromResponse(res)
		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}
	}
}

// PostTokens is the post handler function for tokens.
func (s *Server) PostTokens(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	vals := []dauth.Token{}
	err := dec.Decode(&vals)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	stream, err := s.Auth.SaveTokens(context.Background())
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	data := []dauth.Token{}
	count := 0
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					s.RespondWithError(err, w, r)
				}

				wg.Done()
				return
			}

			v := dauth.Token{}
			v.FromResponse(res)
			count++
			data = append(data, v)
		}
	}()

	for _, v := range vals {
		req := v.ToRequest()
		if err := stream.Send(&req); err != nil {
			s.RespondWithError(err, w, r)
			return
		}
	}

	err = stream.CloseSend()
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	wg.Wait()
	res := dlib.Result{
		Msg:  "Tokens saved",
		Num:  count,
		Data: data,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// PutTokenByID is the put handler function for tokens.
func (s *Server) PutTokenByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	dec := json.NewDecoder(r.Body)
	var v dauth.Token
	err = dec.Decode(&v)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	stream, err := s.Auth.SaveTokens(context.Background())
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	var val dauth.Token
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					s.RespondWithError(err, w, r)
				}

				wg.Done()
				return
			}

			val.FromResponse(res)
		}
	}()

	v.ID = id
	req := v.ToRequest()
	if err := stream.Send(&req); err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	err = stream.CloseSend()
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	wg.Wait()
	res := dlib.Result{
		Msg: "Token saved",
		Num: 1,
		Val: &val,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeleteTokens is the delete handler function for tokens.
func (s *Server) DeleteTokens(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var q dauth.Token
	err := dec.Decode(&q)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	req := q.ToRequest()
	dres, err := s.Auth.DeleteTokens(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "Tokens deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeleteTokenByID is the delete by id handler function for tokens.
func (s *Server) DeleteTokenByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	req := ptypes.TokenRequest{ID: id}
	dres, err := s.Auth.DeleteTokens(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "Token deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeleteTokensOld is the delete by id handler function for old tokens.
func (s *Server) DeleteTokensOld(w http.ResponseWriter, r *http.Request) {
	age, err := strconv.ParseInt(mux.Vars(r)["age"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	oh := age * -24
	old := timestamp.Timestamp{
		Seconds: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).
			Add(time.Hour * time.Duration(oh)).Unix(),
	}

	req := ptypes.TokenRequest{Old: &old}
	dres, err := s.Auth.DeleteTokens(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "Old tokens deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}
