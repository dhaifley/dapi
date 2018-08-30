package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/gorilla/mux"
)

// GetUsers is the get handler function for users.
func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	var data []dauth.User
	q := dauth.User{}
	if len(r.URL.Query()) != 0 {
		if err := q.FromQueryValues(r.URL.Query()); err != nil {
			s.RespondWithError(err, w, r)
			return
		}
	}

	req := q.ToRequest()
	stream, err := s.Auth.GetUsers(context.Background(), &req)
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

		v := dauth.User{}
		err = v.FromResponse(res)
		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}

		v.Pass = ""
		data = append(data, v)
	}
}

// GetUserByID is the get by id handler function for users.
func (s *Server) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	var v *dauth.User
	req := ptypes.UserRequest{ID: id}
	stream, err := s.Auth.GetUsers(context.Background(), &req)
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

		v = new(dauth.User)
		err = v.FromResponse(res)
		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}

		v.Pass = ""
	}
}

// PostUsers is the post handler function for users.
func (s *Server) PostUsers(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	vals := []dauth.User{}
	err := dec.Decode(&vals)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	stream, err := s.Auth.SaveUsers(context.Background())
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	data := []dauth.User{}
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

			v := dauth.User{}
			v.FromResponse(res)
			count++
			v.Pass = ""
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
		Msg:  "Users saved",
		Num:  count,
		Data: data,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// PutUserByID is the put handler function for users.
func (s *Server) PutUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	dec := json.NewDecoder(r.Body)
	var v dauth.User
	err = dec.Decode(&v)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	stream, err := s.Auth.SaveUsers(context.Background())
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	var val dauth.User
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
			val.Pass = ""
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
		Msg: "User saved",
		Num: 1,
		Val: &val,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeleteUsers is the delete handler function for users.
func (s *Server) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var q dauth.User
	err := dec.Decode(&q)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	req := q.ToRequest()
	dres, err := s.Auth.DeleteUsers(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "Users deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeleteUserByID is the delete by id handler function for users.
func (s *Server) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	req := ptypes.UserRequest{ID: id}
	dres, err := s.Auth.DeleteUsers(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "User deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}
