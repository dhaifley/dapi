package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/dauth"
	"github.com/dhaifley/dlib/ptypes"
	"google.golang.org/grpc"
)

// RESTClient values are used to communicate with REST APIs.
type RESTClient struct {
	URL     string
	AuthURL string
	Token   *dauth.Token
	User    *dauth.User
	Client  *http.Client
}

// NewRESTClient creates and returns a pointer to a RESTClient value.
func NewRESTClient(url, authURL, cert string) (*RESTClient, error) {
	rc := RESTClient{
		URL:     url,
		AuthURL: authURL,
	}

	cli, err := dlib.GetHTTPSClient(cert)
	if err != nil {
		return nil, err
	}

	rc.Client = cli
	return &rc, nil
}

// Do executes a request to the REST API server.
func (rc *RESTClient) Do(req *http.Request) (*http.Response, error) {
	return rc.Client.Do(req)
}

// Login obtains an authorization token for the REST client.
func (rc *RESTClient) Login(user *dauth.User) <-chan *dlib.Result {
	rc.User = user
	ch := make(chan *dlib.Result)
	go func() {
		defer close(ch)
		if rc.User == nil {
			ch <- dlib.NewErrorResult(dlib.NewError(http.StatusBadRequest,
				"no user provided for client login"))
			return
		}

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(rc.User)
		req, err := http.NewRequest("POST", rc.AuthURL+"/login", b)
		if err != nil {
			ch <- dlib.NewErrorResult(err)
			return
		}

		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		res, err := rc.Client.Do(req)
		if err != nil {
			ch <- dlib.NewErrorResult(err)
			return
		}

		if res.StatusCode != http.StatusOK {
			var e dlib.Error
			err = json.NewDecoder(res.Body).Decode(&e)
			if err != nil {
				ch <- dlib.NewErrorResult(err)
				return
			}

			ch <- dlib.NewErrorResult(err)
			return
		}

		err = json.NewDecoder(res.Body).Decode(&rc.Token)
		if err != nil {
			ch <- dlib.NewErrorResult(err)
			return
		}

		ch <- dlib.NewResult(nil, nil, "msg", 0, "login completed", nil, nil)
	}()

	return ch
}

// Logout detroys an authorization token for the REST client.
func (rc *RESTClient) Logout() <-chan *dlib.Result {
	ch := make(chan *dlib.Result)
	go func() {
		defer close(ch)
		if rc.Token != nil {
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(rc.Token)
			req, err := http.NewRequest("POST", rc.AuthURL+"/logout", b)
			if err != nil {
				ch <- dlib.NewErrorResult(err)
				return
			}

			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			res, err := rc.Client.Do(req)
			if err != nil {
				ch <- dlib.NewErrorResult(err)
				return
			}

			if res.StatusCode != http.StatusOK {
				var e dlib.Error
				err = json.NewDecoder(res.Body).Decode(&e)
				if err != nil {
					ch <- dlib.NewErrorResult(err)
					return
				}

				ch <- dlib.NewErrorResult(err)
				return
			}
		}

		rc.Token = nil
		ch <- dlib.NewResult(nil, nil, "msg", 0, "logout completed", nil, nil)
	}()

	return ch
}

// RPCClient values are used to communicate with gRPC APIs.
type RPCClient struct {
	URL     string
	AuthURL string
	Token   *dauth.Token
	User    *dauth.User
	Cert    string
	Opts    []grpc.DialOption
	Conn    *grpc.ClientConn
	Context context.Context
}

// NewRPCClient creates and returns a pointer to a RPCClient value.
func NewRPCClient(url, authURL, cert string) (*RPCClient, error) {
	rpc := RPCClient{
		URL:     url,
		AuthURL: authURL,
		Context: context.Background(),
		Cert:    cert,
		Opts:    []grpc.DialOption{},
	}

	creds, err := dlib.GetGRPCClientCredentials(rpc.Cert)
	if err != nil {
		return nil, err
	}

	rpc.Opts = append(rpc.Opts, grpc.WithTransportCredentials(creds))
	rpc.Conn, err = grpc.Dial(rpc.URL, rpc.Opts...)
	if err != nil {
		return nil, err
	}

	return &rpc, nil
}

// Close shuts down the gRPC client connection.
func (rpc *RPCClient) Close() error {
	return rpc.Conn.Close()
}

// Login obtains an authorization token for the RPC client.
func (rpc *RPCClient) Login(user *dauth.User) <-chan *dlib.Result {
	rpc.User = user
	ch := make(chan *dlib.Result)
	go func() {
		defer close(ch)
		cli, err := NewRPCClient(rpc.AuthURL, "", rpc.Cert)
		if err != nil {
			ch <- dlib.NewErrorResult(err)
			return
		}

		acli := ptypes.NewAuthClient(cli.Conn)
		usr := rpc.User.ToRequest()
		if rpc.AuthURL != "test" {
			res, err := acli.Login(rpc.Context, &usr)
			if err != nil {
				ch <- dlib.NewErrorResult(err)
				return
			}

			rpc.Token.FromResponse(res)
		}

		ch <- dlib.NewResult(nil, nil, "msg", 0, "login completed", nil, nil)
	}()

	return ch
}

// Logout obtains an authorization token for the RPC client.
func (rpc *RPCClient) Logout() <-chan *dlib.Result {
	ch := make(chan *dlib.Result)
	go func() {
		defer close(ch)
		if rpc.Token != nil {
			token := rpc.Token.ToRequest()
			cli, err := NewRPCClient(rpc.AuthURL, "", rpc.Cert)
			if err != nil {
				ch <- dlib.NewErrorResult(err)
				return
			}

			acli := ptypes.NewAuthClient(cli.Conn)
			_, err = acli.Logout(rpc.Context, &token)
			if err != nil {
				ch <- dlib.NewErrorResult(err)
				return
			}
		}

		rpc.Token = nil
		ch <- dlib.NewResult(nil, nil, "msg", 0, "logout completed", nil, nil)
	}()

	return ch
}
