package server

import "testing"

func TestServerGetRoutes(t *testing.T) {
	svr := Server{}
	cases := []struct {
		server   *Server
		expected string
	}{
		{
			server:   &svr,
			expected: "index",
		},
	}

	for _, c := range cases {
		actual := c.server.GetRoutes()
		if actual[0].Name != c.expected {
			t.Errorf("Name expected: %v, got: %v", c.expected, actual[0].Name)
		}
	}
}

func TestServerInitRouter(t *testing.T) {
	cases := []struct {
		name   string
		params []string
		server Server
		exp    string
	}{
		{
			name:   "index",
			params: []string{},
			server: Server{},
			exp:    "/",
		},
	}

	for _, c := range cases {
		c.server.InitRouter()
		url, err := c.server.Router.Get(c.name).URL(c.params...)
		if err != nil {
			t.Fatal("Unable to parse URL", c)
		}

		if url.Path != c.exp {
			t.Errorf("Path expected: %v, got: %v", c.exp, url.Path)
		}
	}
}
