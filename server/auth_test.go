package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestAuth(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/auth?token=test&service=test&name=test", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/auth").HandlerFunc(svr.Authenticate)
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

func TestLogin(t *testing.T) {
	jbs := []byte("{\"user\":\"testuser\",\"pass\":\"testpass\"}\n")
	fr, err := http.NewRequest("POST", "/dauth/login", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("POST").Path("/dauth/login").HandlerFunc(svr.Login)
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
			expBody: `{"id":1,"token":"test","user_id":1,"created":"1983-02-02T00:00:00-05:00","expires":"2083-02-02T00:00:00-05:00"}` + "\n",
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

func TestLogout(t *testing.T) {
	jbs := []byte("{\"token\":\"testtoken\"}\n")
	fr, err := http.NewRequest("POST", "/dauth/logout", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("POST").Path("/dauth/logout").HandlerFunc(svr.Logout)
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
			expBody: `{"token":"logout"}` + "\n",
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
