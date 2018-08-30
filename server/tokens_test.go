package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestServerGetTokens(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/tokens", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/tokens").HandlerFunc(svr.GetTokens)
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
			expBody: `[{"id":1,"token":"test"}]` + "\n",
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

func TestGetTokenByID(t *testing.T) {
	fr, err := http.NewRequest("GET", "/dauth/tokens/1", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("GET").Path("/dauth/tokens/{id}").HandlerFunc(svr.GetTokenByID)
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
			expBody: `{"id":1,"token":"test"}` + "\n",
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

func TestPostTokens(t *testing.T) {
	jbs := []byte("[{\"id\":1}]\n")
	fr, err := http.NewRequest("POST", "/dauth/tokens", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("POST").Path("/dauth/tokens").HandlerFunc(svr.PostTokens)
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
			expBody: `{"number":1,"message":"Tokens saved","data":[{"id":1,"token":"test"}]}` + "\n",
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

func TestPutTokenByID(t *testing.T) {
	jbs := []byte("{\"id\":1}\n")
	fr, err := http.NewRequest("PUT", "/dauth/tokens/1", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("PUT").Path("/dauth/tokens/{id}").HandlerFunc(svr.PutTokenByID)
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
			expBody: `{"value":{"id":1,"token":"test"},"number":1,"message":"Token saved"}` + "\n",
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

func TestDeleteTokens(t *testing.T) {
	jbs := []byte("{\"id\":1}\n")
	fr, err := http.NewRequest("DELETE", "/dauth/tokens", bytes.NewBuffer(jbs))
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/tokens").HandlerFunc(svr.DeleteTokens)
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
			expBody: `{"number":1,"message":"Tokens deleted"}` + "\n",
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

func TestDeleteTokenByID(t *testing.T) {
	fr, err := http.NewRequest("DELETE", "/dauth/tokens/1", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/tokens/{id}").HandlerFunc(svr.DeleteTokenByID)
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
			expBody: `{"number":1,"message":"Token deleted"}` + "\n",
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

func TestDeleteTokensOld(t *testing.T) {
	fr, err := http.NewRequest("DELETE", "/dauth/tokens/old/99", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	rtr := mux.NewRouter()
	rtr.Methods("DELETE").Path("/dauth/tokens/old/{age}").HandlerFunc(svr.DeleteTokensOld)
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
			expBody: `{"number":1,"message":"Old tokens deleted"}` + "\n",
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
