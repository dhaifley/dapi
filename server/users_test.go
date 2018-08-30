package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestServerGetUsers(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/users", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/users").HandlerFunc(svr.GetUsers)
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
			expBody: `[{"id":1,"user":"test"}]` + "\n",
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

func TestGetUserByID(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/users/1", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/users/{id}").HandlerFunc(svr.GetUserByID)
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
			expBody: `{"id":1,"user":"test"}` + "\n",
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

func TestPostUsers(t *testing.T) {
	jbs := []byte("[{\"id\":1}]\n")
	fr, err := http.NewRequest("POST", "/dauth/users", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("POST").Path("/dauth/users").HandlerFunc(svr.PostUsers)
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
			expBody: `{"number":1,"message":"Users saved","data":[{"id":1,"user":"test"}]}` + "\n",
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

func TestPutUserByID(t *testing.T) {
	jbs := []byte("{\"id\":1}\n")
	fr, err := http.NewRequest("PUT", "/dauth/users/1", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("PUT").Path("/dauth/users/{id}").HandlerFunc(svr.PutUserByID)
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
			expBody: `{"value":{"id":1,"user":"test"},"number":1,"message":"User saved"}` + "\n",
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

func TestDeleteUsers(t *testing.T) {
	jbs := []byte("{\"id\":1}\n")
	fr, err := http.NewRequest("DELETE", "/dauth/users", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/users").HandlerFunc(svr.DeleteUsers)
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
			expBody: `{"number":1,"message":"Users deleted"}` + "\n",
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

func TestDeleteUserByID(t *testing.T) {
	fr, err := http.NewRequest("DELETE", "/dauth/users/1", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/users/{id}").HandlerFunc(svr.DeleteUserByID)
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
			expBody: `{"number":1,"message":"User deleted"}` + "\n",
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
