package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dhaifley/dapi/lib"
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/grpc"
)

type FakeAuthGetTokensClient struct {
	grpc.ClientStream
	count int
}

func (x *FakeAuthGetTokensClient) Recv() (*ptypes.TokenResponse, error) {
	x.count++
	if x.count > 1 {
		return nil, io.EOF
	}

	m := new(ptypes.TokenResponse)
	m.ID = 1
	m.Token = "test"
	return m, nil
}

func (x *FakeAuthGetTokensClient) CloseSend() error {
	return nil
}

type FakeAuthSaveTokensClient struct {
	grpc.ClientStream
	count int
}

func (x *FakeAuthSaveTokensClient) Send(m *ptypes.TokenRequest) error {
	return nil
}

func (x *FakeAuthSaveTokensClient) Recv() (*ptypes.TokenResponse, error) {
	x.count++
	if x.count > 1 {
		return nil, io.EOF
	}

	m := new(ptypes.TokenResponse)
	m.ID = 1
	m.Token = "test"
	return m, nil
}

func (x *FakeAuthSaveTokensClient) CloseSend() error {
	return nil
}

type FakeAuthGetUsersClient struct {
	grpc.ClientStream
	count int
}

func (x *FakeAuthGetUsersClient) Recv() (*ptypes.UserResponse, error) {
	x.count++
	if x.count > 1 {
		return nil, io.EOF
	}

	m := new(ptypes.UserResponse)
	m.ID = 1
	m.User = "test"
	return m, nil
}

func (x *FakeAuthGetUsersClient) CloseSend() error {
	return nil
}

type FakeAuthSaveUsersClient struct {
	grpc.ClientStream
	count int
}

func (x *FakeAuthSaveUsersClient) Send(m *ptypes.UserRequest) error {
	return nil
}

func (x *FakeAuthSaveUsersClient) Recv() (*ptypes.UserResponse, error) {
	x.count++
	if x.count > 1 {
		return nil, io.EOF
	}

	m := new(ptypes.UserResponse)
	m.ID = 1
	m.User = "test"
	return m, nil
}

func (x *FakeAuthSaveUsersClient) CloseSend() error {
	return nil
}

type FakeAuthGetPermsClient struct {
	grpc.ClientStream
	count int
}

func (x *FakeAuthGetPermsClient) Recv() (*ptypes.PermResponse, error) {
	x.count++
	if x.count > 1 {
		return nil, io.EOF
	}

	m := new(ptypes.PermResponse)
	m.ID = 1
	m.Service = "test"
	m.Name = "test"
	return m, nil
}

func (x *FakeAuthGetPermsClient) CloseSend() error {
	return nil
}

type FakeAuthSavePermsClient struct {
	grpc.ClientStream
	count int
}

func (x *FakeAuthSavePermsClient) Send(m *ptypes.PermRequest) error {
	return nil
}

func (x *FakeAuthSavePermsClient) Recv() (*ptypes.PermResponse, error) {
	x.count++
	if x.count > 1 {
		return nil, io.EOF
	}

	m := new(ptypes.PermResponse)
	m.ID = 1
	m.Service = "test"
	m.Name = "test"
	return m, nil
}

func (x *FakeAuthSavePermsClient) CloseSend() error {
	return nil
}

type FakeAuthGetUserPermsClient struct {
	grpc.ClientStream
	count int
}

func (x *FakeAuthGetUserPermsClient) Recv() (*ptypes.UserPermResponse, error) {
	x.count++
	if x.count > 1 {
		return nil, io.EOF
	}

	m := new(ptypes.UserPermResponse)
	m.ID = 1
	m.UserID = 1
	m.PermID = 1
	return m, nil
}

func (x *FakeAuthGetUserPermsClient) CloseSend() error {
	return nil
}

type FakeAuthSaveUserPermsClient struct {
	grpc.ClientStream
	count int
}

func (x *FakeAuthSaveUserPermsClient) Send(m *ptypes.UserPermRequest) error {
	return nil
}

func (x *FakeAuthSaveUserPermsClient) Recv() (*ptypes.UserPermResponse, error) {
	x.count++
	if x.count > 1 {
		return nil, io.EOF
	}

	m := new(ptypes.UserPermResponse)
	m.ID = 1
	m.UserID = 1
	m.PermID = 1
	return m, nil
}

func (x *FakeAuthSaveUserPermsClient) CloseSend() error {
	return nil
}

type FakeAuthClient struct{}

func (fc *FakeAuthClient) GetTokens(ctx context.Context, in *ptypes.TokenRequest, opts ...grpc.CallOption) (ptypes.Auth_GetTokensClient, error) {
	fgc := FakeAuthGetTokensClient{}
	return &fgc, nil
}

func (fc *FakeAuthClient) SaveTokens(ctx context.Context, opts ...grpc.CallOption) (ptypes.Auth_SaveTokensClient, error) {
	fsc := FakeAuthSaveTokensClient{}
	return &fsc, nil
}

func (fc *FakeAuthClient) DeleteTokens(ctx context.Context, in *ptypes.TokenRequest, opts ...grpc.CallOption) (*ptypes.DeleteResponse, error) {
	res := ptypes.DeleteResponse{Num: 1}
	return &res, nil
}

func (fc *FakeAuthClient) GetPerms(ctx context.Context, in *ptypes.PermRequest, opts ...grpc.CallOption) (ptypes.Auth_GetPermsClient, error) {
	fgc := FakeAuthGetPermsClient{}
	return &fgc, nil
}

func (fc *FakeAuthClient) SavePerms(ctx context.Context, opts ...grpc.CallOption) (ptypes.Auth_SavePermsClient, error) {
	fsc := FakeAuthSavePermsClient{}
	return &fsc, nil
}

func (fc *FakeAuthClient) DeletePerms(ctx context.Context, in *ptypes.PermRequest, opts ...grpc.CallOption) (*ptypes.DeleteResponse, error) {
	res := ptypes.DeleteResponse{Num: 1}
	return &res, nil
}

func (fc *FakeAuthClient) GetUsers(ctx context.Context, in *ptypes.UserRequest, opts ...grpc.CallOption) (ptypes.Auth_GetUsersClient, error) {
	fgc := FakeAuthGetUsersClient{}
	return &fgc, nil
}

func (fc *FakeAuthClient) SaveUsers(ctx context.Context, opts ...grpc.CallOption) (ptypes.Auth_SaveUsersClient, error) {
	fsc := FakeAuthSaveUsersClient{}
	return &fsc, nil
}

func (fc *FakeAuthClient) DeleteUsers(ctx context.Context, in *ptypes.UserRequest, opts ...grpc.CallOption) (*ptypes.DeleteResponse, error) {
	res := ptypes.DeleteResponse{Num: 1}
	return &res, nil
}

func (fc *FakeAuthClient) GetUserPerms(ctx context.Context, in *ptypes.UserPermRequest, opts ...grpc.CallOption) (ptypes.Auth_GetUserPermsClient, error) {
	fgc := FakeAuthGetUserPermsClient{}
	return &fgc, nil
}

func (fc *FakeAuthClient) SaveUserPerms(ctx context.Context, opts ...grpc.CallOption) (ptypes.Auth_SaveUserPermsClient, error) {
	fsc := FakeAuthSaveUserPermsClient{}
	return &fsc, nil
}

func (fc *FakeAuthClient) DeleteUserPerms(ctx context.Context, in *ptypes.UserPermRequest, opts ...grpc.CallOption) (*ptypes.DeleteResponse, error) {
	res := ptypes.DeleteResponse{Num: 1}
	return &res, nil
}

func (fc *FakeAuthClient) Login(ctx context.Context, in *ptypes.UserRequest, opts ...grpc.CallOption) (*ptypes.TokenResponse, error) {
	res := ptypes.TokenResponse{
		ID:     1,
		Token:  "test",
		UserID: 1,
		Created: &timestamp.Timestamp{
			Seconds: time.Date(1983, 2, 2, 0, 0, 0, 0, time.Local).Unix(),
		},
		Expires: &timestamp.Timestamp{
			Seconds: time.Date(2083, 2, 2, 0, 0, 0, 0, time.Local).Unix(),
		},
	}

	return &res, nil
}

func (fc *FakeAuthClient) Logout(ctx context.Context, in *ptypes.TokenRequest, opts ...grpc.CallOption) (*ptypes.TokenResponse, error) {
	res := ptypes.TokenResponse{
		Token:  "logout",
		UserID: 0,
	}

	return &res, nil
}

func (fc *FakeAuthClient) Auth(ctx context.Context, in *ptypes.AuthRequest, opts ...grpc.CallOption) (*ptypes.AuthResponse, error) {
	res := ptypes.AuthResponse{
		Ok: true,
		User: &ptypes.UserResponse{
			ID:   1,
			User: "test",
		},
		Perm: &ptypes.PermResponse{
			ID:      1,
			Service: "test",
			Name:    "test",
		},
	}

	return &res, nil
}

type FakeHandler struct {
	Served bool
}

func (f FakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.Served = true
}

func TestServerAuthHandler(t *testing.T) {
	fr, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	fc := FakeAuthClient{}
	lm, _ := test.NewNullLogger()
	svr := Server{Auth: &fc, Log: lm}
	fh := FakeHandler{false}
	authFunc, ok := svr.AuthHandler(fh, &dauth.Perm{Service: "test", Name: "test"}).(http.HandlerFunc)
	if !ok {
		t.Fatal("Auth did not return expected handler function")
	}

	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		token   string
		expCode int
		expBody []byte
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			token:   "test",
			expCode: http.StatusOK,
		},
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			token:   "",
			expCode: http.StatusUnauthorized,
		},
	}

	for _, c := range cases {
		c.r.Header.Set("Token", c.token)
		authFunc(c.w, c.r)
		if c.expCode != c.w.Code {
			t.Errorf("Status Code expected: %v, got: %v", c.expCode, c.w.Code)
		}
	}
}

func TestServerHeader(t *testing.T) {
	fr, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	s := Server{}
	fh := FakeHandler{false}
	headFunc, ok := s.Header(fh).(http.HandlerFunc)
	if !ok {
		t.Fatal("Header did not return expected handler function")
	}

	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expType string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expType: "application/json; charset=UTF-8",
		},
	}

	for _, c := range cases {
		headFunc(c.w, c.r)
		gotContentType := c.w.Header().Get("Content-Type")
		if c.expType != gotContentType {
			t.Errorf("Content-Type expected: %v, got: %v", c.expType, gotContentType)
		}
	}
}

func TestServerLogger(t *testing.T) {
	lm, hook := test.NewNullLogger()
	fr, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	s := Server{Log: lm}
	fh := FakeHandler{false}
	logFunc, ok := s.Logger(fh).(http.HandlerFunc)
	if !ok {
		t.Fatal("Logger did not return expected handler function")
	}

	cases := []struct {
		w               *httptest.ResponseRecorder
		r               *http.Request
		expectedLevel   logrus.Level
		expectedMessage string
	}{
		{
			w:               httptest.NewRecorder(),
			r:               fr,
			expectedLevel:   logrus.InfoLevel,
			expectedMessage: "Request processed",
		},
	}

	for _, c := range cases {
		logFunc(c.w, c.r)
		if hook.LastEntry().Level != c.expectedLevel {
			t.Errorf("Log level expected: %v, got: %v", c.expectedLevel, hook.LastEntry().Level)
		}

		if hook.LastEntry().Message != c.expectedMessage {
			t.Errorf("Log level expected: %v, got: %v", c.expectedMessage, hook.LastEntry().Message)
		}
	}
}

func TestServerRespondWithError(t *testing.T) {
	fr, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	lm, _ := test.NewNullLogger()
	svr := Server{Log: lm}
	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		svr.RespondWithError(dlib.NewError(
			http.StatusInternalServerError, "testerror"), c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}
	}
}

func TestServerRespondNotFound(t *testing.T) {
	fr, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	lm, _ := test.NewNullLogger()
	svr := Server{Log: lm}
	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusNotFound,
		},
	}

	for _, c := range cases {
		svr.RespondNotFound(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}
	}
}

func TestServerNotFoundHandler(t *testing.T) {
	lm, _ := test.NewNullLogger()
	svr := Server{Log: lm}
	fr, err := http.NewRequest("GET", "/badurl", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
		expBody string
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusNotFound,
		},
	}

	for _, c := range cases {
		svr.NotFoundHandler(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}
	}
}

func TestServerGetDocs(t *testing.T) {
	lm, _ := test.NewNullLogger()
	svr := Server{Log: lm}
	fr, err := http.NewRequest("GET", "/docs", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	cases := []struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		expCode int
	}{
		{
			w:       httptest.NewRecorder(),
			r:       fr,
			expCode: http.StatusOK,
		},
	}

	for _, c := range cases {
		svr.GetDocs(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}
	}
}

func TestServerGetIndex(t *testing.T) {
	lm, _ := test.NewNullLogger()
	svr := Server{Log: lm}
	fr, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

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
			expBody: lib.ServiceInfo.String() + "\n",
		},
	}

	for _, c := range cases {
		svr.GetIndex(c.w, c.r)
		if c.w.Code != c.expCode {
			t.Errorf("Code expected: %v, got: %v", c.expCode, c.w.Code)
		}

		bs := string(c.w.Body.Bytes())
		if bs != c.expBody {
			t.Errorf("Body expected: %v, got: %v", c.expBody, bs)
		}
	}
}

func TestServerGetIcon(t *testing.T) {
	lm, _ := test.NewNullLogger()
	svr := Server{Log: lm}
	fr, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal("Failed to initialize request", err)
	}

	cases := []struct {
		w    *httptest.ResponseRecorder
		r    *http.Request
		code int
	}{
		{
			w:    httptest.NewRecorder(),
			r:    fr,
			code: http.StatusOK,
		},
	}

	for _, c := range cases {
		svr.GetIcon(c.w, c.r)
		if c.w.Code != c.code {
			t.Errorf("Code expected: %v, got: %v", c.code, c.w.Code)
		}
	}
}
