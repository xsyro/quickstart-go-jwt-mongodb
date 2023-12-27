package api

import (
	"net/http"
)

// Provide personal preferences to roles
const (
	Role1 = "SUPERVISOR"
	Role2 = "SALES_PERSON"
)

func securedRole1Only() HttpRequest {
	return HttpRequest{
		Uri:         "/secured/role-1",
		Method:      GET,
		Secure:      true,
		PermitRoles: []string{Role1},
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Access granted!"))
		},
	}
}

func securedRole2Only() HttpRequest {
	return HttpRequest{
		Uri:         "/secured/role-2",
		Method:      GET,
		Secure:      true,
		PermitRoles: []string{Role2},
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Access granted!"))
		},
	}
}

func securedRole1And2Only() HttpRequest {
	return HttpRequest{
		Uri:         "/secured/role-1-and-2",
		Method:      GET,
		Secure:      true,
		PermitRoles: []string{Role1, Role2},
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Access granted!"))
		},
	}
}

func UserHandlers(resources *WithResource) {

	resources.HttpRequest.RequestRegistry(securedRole1Only())
	resources.HttpRequest.RequestRegistry(securedRole2Only())
	resources.HttpRequest.RequestRegistry(securedRole1And2Only())

}
