package api

import (
	"net/http"
)

func index() HttpRequest {
	return HttpRequest{
		Uri:    "/",
		Secure: false,
		Method: GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Hello!"))
			return
		},
	}
}

// healthCheck return HTTP_STATUS_CODE of 200 to inform container service of it activeness
func healthCheck() HttpRequest {
	return HttpRequest{
		Uri:    "/status",
		Secure: false,
		Method: GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Good!"))
			return
		},
	}
}

func StaticHandlers(resources *WithResource) {
	resources.HttpRequest.RequestRegistry(healthCheck())
	resources.HttpRequest.RequestRegistry(index())
}
