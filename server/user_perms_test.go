package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestServerGetUserPerms(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/userperms", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/userperms").HandlerFunc(svr.GetUserPerms)
	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusOK,
			expBody: `[{"id":1,"user_id":1,"perm_id":1}]` + "\n",
		},
	}

	for _, c := range cases {
		rtr.ServeHTTP(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}

		bs := string(c.w.Body.Bytes())
		if bs != c.expBody {
			t.Errorf("Body expected: %v, got: %v", c.expBody, bs)
		}
	}
}

func TestGetUserPermByID(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/userperms/1", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/userperms/{id}").HandlerFunc(svr.GetUserPermByID)
	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusOK,
			expBody: `{"id":1,"user_id":1,"perm_id":1}` + "\n",
		},
	}

	for _, c := range cases {
		rtr.ServeHTTP(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}

		bs := string(c.w.Body.Bytes())
		if bs != c.expBody {
			t.Errorf("Body expected: %v, got: %v", c.expBody, bs)
		}
	}
}

func TestPostUserPerms(t *testing.T) {
	jbs := []byte("[{\"id\":1}]\n")
	fr, err := http.NewRequest("POST", "/dauth/userperms", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("POST").Path("/dauth/userperms").HandlerFunc(svr.PostUserPerms)
	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusOK,
			expBody: `{"number":1,"message":"User permissions saved","data":[{"id":1,"user_id":1,"perm_id":1}]}` + "\n",
		},
	}

	for _, c := range cases {
		rtr.ServeHTTP(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}

		bs := string(c.w.Body.Bytes())
		if bs != c.expBody {
			t.Errorf("Body expected: %v, got: %v", c.expBody, bs)
		}
	}
}

func TestPutUserPermByID(t *testing.T) {
	jbs := []byte("{\"id\":1}\n")
	fr, err := http.NewRequest("PUT", "/dauth/userperms/1", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("PUT").Path("/dauth/userperms/{id}").HandlerFunc(svr.PutUserPermByID)
	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusOK,
			expBody: `{"value":{"id":1,"user_id":1,"perm_id":1},"number":1,"message":"User permission saved"}` + "\n",
		},
	}

	for _, c := range cases {
		rtr.ServeHTTP(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}

		bs := string(c.w.Body.Bytes())
		if bs != c.expBody {
			t.Errorf("Body expected: %v, got: %v", c.expBody, bs)
		}
	}
}

func TestDeleteUserPerms(t *testing.T) {
	jbs := []byte("{\"id\":1}\n")
	fr, err := http.NewRequest("DELETE", "/dauth/userperms", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/userperms").HandlerFunc(svr.DeleteUserPerms)
	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusOK,
			expBody: `{"number":1,"message":"User permissions deleted"}` + "\n",
		},
	}

	for _, c := range cases {
		rtr.ServeHTTP(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}

		bs := string(c.w.Body.Bytes())
		if bs != c.expBody {
			t.Errorf("Body expected: %v, got: %v", c.expBody, bs)
		}
	}
}

func TestDeleteUserPermByID(t *testing.T) {
	fr, err := http.NewRequest("DELETE", "/dauth/userperms/1", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/userperms/{id}").HandlerFunc(svr.DeleteUserPermByID)
	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusOK,
			expBody: `{"number":1,"message":"User permission deleted"}` + "\n",
		},
	}

	for _, c := range cases {
		rtr.ServeHTTP(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}

		bs := string(c.w.Body.Bytes())
		if bs != c.expBody {
			t.Errorf("Body expected: %v, got: %v", c.expBody, bs)
		}
	}
}
