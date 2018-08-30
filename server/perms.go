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

// GetPerms is the get handler function for perms.
func (s *Server) GetPerms(w http.ResponseWriter, r *http.Request) {
	var data []dauth.Perm
	q := dauth.Perm{}
	if len(r.URL.Query()) != 0 {
		if err := q.FromQueryValues(r.URL.Query()); err != nil {
			s.RespondWithError(err, w, r)
			return
		}
	}

	req := q.ToRequest()
	stream, err := s.Auth.GetPerms(context.Background(), &req)
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

		v := dauth.Perm{}
		err = v.FromResponse(res)
		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}

		data = append(data, v)
	}
}

// GetPermByID is the get by id handler function for perms.
func (s *Server) GetPermByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	var v *dauth.Perm
	req := ptypes.PermRequest{ID: id}
	stream, err := s.Auth.GetPerms(context.Background(), &req)
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

		v = new(dauth.Perm)
		err = v.FromResponse(res)
		if err != nil {
			s.RespondWithError(err, w, r)
			return
		}
	}
}

// PostPerms is the post handler function for perms.
func (s *Server) PostPerms(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	vals := []dauth.Perm{}
	err := dec.Decode(&vals)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	stream, err := s.Auth.SavePerms(context.Background())
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	data := []dauth.Perm{}
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

			v := dauth.Perm{}
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
		Msg:  "Permissions saved",
		Num:  count,
		Data: data,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// PutPermByID is the put handler function for perms.
func (s *Server) PutPermByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	dec := json.NewDecoder(r.Body)
	var v dauth.Perm
	err = dec.Decode(&v)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	stream, err := s.Auth.SavePerms(context.Background())
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	var val dauth.Perm
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
		Msg: "Permission saved",
		Num: 1,
		Val: &val,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeletePerms is the delete handler function for perms.
func (s *Server) DeletePerms(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var q dauth.Perm
	err := dec.Decode(&q)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	defer r.Body.Close()
	req := q.ToRequest()
	dres, err := s.Auth.DeletePerms(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "Permissions deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}

// DeletePermByID is the delete by id handler function for perms.
func (s *Server) DeletePermByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.RespondWithError(dlib.NewError(http.StatusBadRequest,
			"invalid id value"), w, r)
		return
	}

	req := ptypes.PermRequest{ID: id}
	dres, err := s.Auth.DeletePerms(context.Background(), &req)
	if err != nil {
		s.RespondWithError(err, w, r)
		return
	}

	if dres.Num == 0 {
		s.RespondNotFound(w, r)
		return
	}

	res := dlib.Result{
		Msg: "Permission deleted",
		Num: int(dres.Num),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.Log.Error(err)
	}
}
