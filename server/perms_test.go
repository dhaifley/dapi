package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestServerGetPerms(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/perms", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/perms").HandlerFunc(svr.GetPerms)
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
			expBody: `[{"id":1,"service":"test","name":"test"}]` + "\n",
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

func TestGetPermByID(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/perms/1", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/perms/{id}").HandlerFunc(svr.GetPermByID)
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
			expBody: `{"id":1,"service":"test","name":"test"}` + "\n",
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

func TestPostPerms(t *testing.T) {
	jbs := []byte("[{\"id\":1}]\n")
	fr, err := http.NewRequest("POST", "/dauth/perms", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("POST").Path("/dauth/perms").HandlerFunc(svr.PostPerms)
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
			expBody: `{"number":1,"message":"Permissions saved","data":[{"id":1,"service":"test","name":"test"}]}` + "\n",
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

func TestPutPermByID(t *testing.T) {
	jbs := []byte("{\"id\":1}\n")
	fr, err := http.NewRequest("PUT", "/dauth/perms/1", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("PUT").Path("/dauth/perms/{id}").HandlerFunc(svr.PutPermByID)
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
			expBody: `{"value":{"id":1,"service":"test","name":"test"},"number":1,"message":"Permission saved"}` + "\n",
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

func TestDeletePerms(t *testing.T) {
	jbs := []byte("{\"id\":1}\n")
	fr, err := http.NewRequest("DELETE", "/dauth/perms", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/perms").HandlerFunc(svr.DeletePerms)
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
			expBody: `{"number":1,"message":"Permissions deleted"}` + "\n",
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

func TestDeletePermByID(t *testing.T) {
	fr, err := http.NewRequest("DELETE", "/dauth/perms/1", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/perms/{id}").HandlerFunc(svr.DeletePermByID)
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
			expBody: `{"number":1,"message":"Permission deleted"}` + "\n",
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
