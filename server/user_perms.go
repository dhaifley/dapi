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

// GetUserPerms is the get handler function for user_perms.
func (s *Server) GetUserPerms(w http.ResponseWriter, r *http.Request) {
	var data []dauth.UserPerm
	q := dauth.UserPerm{}
	if len(r.URL.Query()) != 0 {
		if err := q.FromQueryValues(r.URL.Query()); err != nil {
			s.RespondWithError(err, w, r)
			return
		}
	}

	req := q.ToRequest()
	stream, err := s.Auth.GetUserPerms(context.Background(), &req)
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

		v := dauth.UserPerm{}
		err = v.FromResponse(res)
		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}

		data = append(data, v)
	}
}

// GetUserPermByID is the get by id handler function for user_perms.
func (s *Server) GetUserPermByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	var v *dauth.UserPerm
	req := ptypes.UserPermRequest{ID: id}
	stream, err := s.Auth.GetUserPerms(context.Background(), &req)
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

		v = new(dauth.UserPerm)
		err = v.FromResponse(res)
		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}
	}
}

// PostUserPerms is the post handler function for user_perms.
func (s *Server) PostUserPerms(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	vals := []dauth.UserPerm{}
	err := dec.Decode(&vals)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	stream, err := s.Auth.SaveUserPerms(context.Background())
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	data := []dauth.UserPerm{}
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

			v := dauth.UserPerm{}
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
		Msg:  "User permissions saved",
		Num:  count,
		Data: data,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// PutUserPermByID is the put handler function for user_perms.
func (s *Server) PutUserPermByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	dec := json.NewDecoder(r.Body)
	var v dauth.UserPerm
	err = dec.Decode(&v)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	stream, err := s.Auth.SaveUserPerms(context.Background())
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	var val dauth.UserPerm
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
		Msg: "User permission saved",
		Num: 1,
		Val: &val,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeleteUserPerms is the delete handler function for user_perms.
func (s *Server) DeleteUserPerms(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var q dauth.UserPerm
	err := dec.Decode(&q)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	req := q.ToRequest()
	dres, err := s.Auth.DeleteUserPerms(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "User permissions deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeleteUserPermByID is the delete by id handler function for user_perms.
func (s *Server) DeleteUserPermByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	req := ptypes.UserPermRequest{ID: id}
	dres, err := s.Auth.DeleteUserPerms(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "User permission deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}
