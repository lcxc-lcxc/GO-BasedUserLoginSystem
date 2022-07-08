package v0_0_0

import (
	"net/http"
	"v0.0.0/internel/router"
)

func main() {
	router := router.NewRouter()
	s := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	s.ListenAndServe()

}
