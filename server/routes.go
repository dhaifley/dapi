package server

import (
	"net/http"

	"github.com/dhaifley/dlib/dauth"
	"github.com/gorilla/mux"
)

// Route type defines an api route for use by the router.
type Route struct {
	Service     string
	Name        string
	Path        string
	Method      string
	Auth        bool
	HandlerFunc http.HandlerFunc
}

// InitRouter initializes the server router.
// It configures and attaches all required middleware and attaches the routes
// specified in the routes.go file.
func (s *Server) InitRouter() {
	s.Router = mux.NewRouter().StrictSlash(true)
	for _, route := range s.GetRoutes() {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = s.Header(handler)
		if route.Auth {
			handler = s.AuthHandler(handler, &dauth.Perm{
				Service: route.Service,
				Name:    route.Name,
			})
		}

		s.Router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)
	}

	s.Router.NotFoundHandler = http.HandlerFunc(s.NotFoundHandler)
}

// GetRoutes returns all routes for the server.
func (s *Server) GetRoutes() []Route {
	return []Route{
		Route{
			Service:     "dapi",
			Name:        "index",
			Path:        "/",
			Method:      "GET",
			Auth:        false,
			HandlerFunc: s.GetIndex,
		},
		Route{
			Service:     "dapi",
			Name:        "docs",
			Path:        "/docs",
			Method:      "GET",
			Auth:        false,
			HandlerFunc: s.GetDocs,
		},
		Route{
			Service:     "dapi",
			Name:        "icon",
			Path:        "/favicon.ico",
			Method:      "GET",
			Auth:        false,
			HandlerFunc: s.GetIcon,
		},
		Route{
			Service:     "dauth",
			Name:        "GetTokens",
			Path:        "/dauth/tokens",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.GetTokens,
		},
		Route{
			Service:     "dauth",
			Name:        "GetTokens",
			Path:        "/dauth/tokens/{id}",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.GetTokenByID,
		},
		Route{
			Service:     "dauth",
			Name:        "SaveTokens",
			Path:        "/dauth/tokens",
			Method:      "POST",
			Auth:        true,
			HandlerFunc: s.PostTokens,
		},
		Route{
			Service:     "dauth",
			Name:        "SaveTokens",
			Path:        "/dauth/tokens/{id}",
			Method:      "PUT",
			Auth:        true,
			HandlerFunc: s.PutTokenByID,
		},
		Route{
			Service:     "dauth",
			Name:        "DeleteTokens",
			Path:        "/dauth/tokens/old",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeleteTokensOld,
		},
		Route{
			Service:     "dauth",
			Name:        "DeleteTokens",
			Path:        "/dauth/tokens/old/{age}",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeleteTokensOld,
		},
		Route{
			Service:     "dauth",
			Name:        "DeleteTokens",
			Path:        "/dauth/tokens/{id}",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeleteTokenByID,
		},
		Route{
			Service:     "dauth",
			Name:        "DeleteTokens",
			Path:        "/dauth/tokens}",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeleteTokens,
		},
		Route{
			Service:     "dauth",
			Name:        "GetUsers",
			Path:        "/dauth/users",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.GetUsers,
		},
		Route{
			Service:     "dauth",
			Name:        "GetUsers",
			Path:        "/dauth/users/{id}",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.GetUserByID,
		},
		Route{
			Service:     "dauth",
			Name:        "SaveUsers",
			Path:        "/dauth/users",
			Method:      "POST",
			Auth:        true,
			HandlerFunc: s.PostUsers,
		},
		Route{
			Service:     "dauth",
			Name:        "SaveUsers",
			Path:        "/dauth/users/{id}",
			Method:      "PUT",
			Auth:        true,
			HandlerFunc: s.PutUserByID,
		},
		Route{
			Service:     "dauth",
			Name:        "DeleteUsers",
			Path:        "/dauth/users/{id}",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeleteUserByID,
		},
		Route{
			Service:     "dauth",
			Name:        "DeleteUsers",
			Path:        "/dauth/users",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeleteUsers,
		},
		Route{
			Service:     "dauth",
			Name:        "GetPerms",
			Path:        "/dauth/perms",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.GetPerms,
		},
		Route{
			Service:     "dauth",
			Name:        "GetPerms",
			Path:        "/dauth/perms/{id}",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.GetPermByID,
		},
		Route{
			Service:     "dauth",
			Name:        "SavePerms",
			Path:        "/dauth/perms",
			Method:      "POST",
			Auth:        true,
			HandlerFunc: s.PostPerms,
		},
		Route{
			Service:     "dauth",
			Name:        "SavePerms",
			Path:        "/dauth/perms/{id}",
			Method:      "PUT",
			Auth:        true,
			HandlerFunc: s.PutPermByID,
		},
		Route{
			Service:     "dauth",
			Name:        "DeletePerms",
			Path:        "/dauth/perms/{id}",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeletePermByID,
		},
		Route{
			Service:     "dauth",
			Name:        "DeletePerms",
			Path:        "/dauth/perms",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeletePerms,
		},
		Route{
			Service:     "dauth",
			Name:        "GetUserPerms",
			Path:        "/dauth/userperms",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.GetUserPerms,
		},
		Route{
			Service:     "dauth",
			Name:        "GetUserPerms",
			Path:        "/dauth/userperms/{id}",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.GetUserPermByID,
		},
		Route{
			Service:     "dauth",
			Name:        "SaveUserPerms",
			Path:        "/dauth/userperms",
			Method:      "POST",
			Auth:        true,
			HandlerFunc: s.PostUserPerms,
		},
		Route{
			Service:     "dauth",
			Name:        "SaveUserPerms",
			Path:        "/dauth/userperms/{id}",
			Method:      "PUT",
			Auth:        true,
			HandlerFunc: s.PutUserPermByID,
		},
		Route{
			Service:     "dauth",
			Name:        "DeleteUserPerms",
			Path:        "/dauth/userperms/{id}",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeleteUserPermByID,
		},
		Route{
			Service:     "dauth",
			Name:        "DeleteUserPerms",
			Path:        "/dauth/userperms",
			Method:      "DELETE",
			Auth:        true,
			HandlerFunc: s.DeleteUserPerms,
		},
		Route{
			Service:     "dauth",
			Name:        "Auth",
			Path:        "/dauth/auth",
			Method:      "GET",
			Auth:        true,
			HandlerFunc: s.Authenticate,
		},
		Route{
			Service:     "dauth",
			Name:        "Login",
			Path:        "/dauth/login",
			Method:      "POST",
			Auth:        false,
			HandlerFunc: s.Login,
		},
		Route{
			Service:     "dauth",
			Name:        "Logout",
			Path:        "/dauth/logout",
			Method:      "POST",
			Auth:        false,
			HandlerFunc: s.Logout,
		},
	}
}
