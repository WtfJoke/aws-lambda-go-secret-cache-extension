package webserver

import (
	"aws-lambda-go-secret-cache-extension/extension"
	"net/http"

	"github.com/gorilla/mux"
)

func Start(port string) {
	go startHTTPServer(port)
}

func startHTTPServer(port string) {
	router := mux.NewRouter()
	router.Path("/secrets").Queries("name", "{name}").HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			variables := mux.Vars(request)
			secret := extension.GetSecretFromCache(variables["name"])

			if len(secret.SecretString) != 0 {
				_, _ = writer.Write([]byte(secret.SecretString))
			} else {
				_, _ = writer.Write([]byte("No secret found"))
			}
		})

	println("Starting http server on port: ", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}
