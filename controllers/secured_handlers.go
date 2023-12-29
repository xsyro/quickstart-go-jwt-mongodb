package controllers

import (
	"context"
	"net/http"
	"quickstart-go-jwt-mongodb/internal"
	"quickstart-go-jwt-mongodb/server"
)

// Provide personal preferences to roles
const (
	Role1 = "SUPERVISOR"
	Role2 = "SALES_PERSON"
)

func SecuredRole1Only(database internal.MongoDatabase, ctx context.Context) server.Controller {
	return server.Controller{
		Uri:         "/secured/role-1",
		Method:      server.GET,
		Secure:      true,
		PermitRoles: []string{Role1},
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Access granted!"))
		},
	}
}

func SecuredRole2Only(database internal.MongoDatabase, ctx context.Context) server.Controller {
	return server.Controller{
		Uri:         "/secured/role-2",
		Method:      server.GET,
		Secure:      true,
		PermitRoles: []string{Role2},
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Access granted!"))
		},
	}
}

func SecuredRole1And2Only(database internal.MongoDatabase, ctx context.Context) server.Controller {
	return server.Controller{
		Uri:         "/secured/role-1-and-2",
		Method:      server.GET,
		Secure:      true,
		PermitRoles: []string{Role1, Role2},
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Access granted!"))
		},
	}
}
