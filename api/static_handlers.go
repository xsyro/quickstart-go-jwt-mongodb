package api

import (
	"net/http"
)

func index() HandleRequest {
	return HandleRequest{
		Uri:    "/",
		Method: GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Hello!"))
			return
		},
	}
}

// healthCheck return HTTP_STATUS_CODE of 200 to inform container service of it activeness
func healthCheck() HandleRequest {
	return HandleRequest{
		Uri:    "/status",
		Method: GET,
		Callback: func(w http.ResponseWriter, req *http.Request) {
			_, _ = w.Write([]byte("Good!"))
			return
		},
	}
}

func StaticHandlers(resources *WithResource) {
	resources.HttpRequest.HandleRequest(healthCheck())
	resources.HttpRequest.HandleRequest(index())
}
