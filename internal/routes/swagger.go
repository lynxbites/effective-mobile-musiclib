package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

func NewSwaggerRouter() *chi.Mux {

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./docs/"))
	r.Handle("/docs/*", http.StripPrefix("/docs/", fs))
	r.Get("/doc/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8001/docs/swagger.json"), //The url pointing to API definition
	))

	return r
}
