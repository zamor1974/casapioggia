package main

import (
	"casapioggia/config"
	"casapioggia/controllers"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	dbsqlx := config.ConnectDBSqlx()
	hsqlx := controllers.NewBaseHandlerSqlx(dbsqlx)

	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "docs"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/docs", sh1)

	pioggia := r.PathPrefix("/").Subrouter()
	pioggia.HandleFunc("/rain", hsqlx.PostRainSqlx).Methods("POST")
	pioggia.HandleFunc("/rains", hsqlx.GetRainsSqlx).Methods("GET")
	pioggia.HandleFunc("/lasthour", hsqlx.GetLastHourSqlx).Methods("GET")

	http.Handle("/", r)
	s := &http.Server{
		Addr: fmt.Sprintf("%s:%s", "", "5555"),
	}
	s.ListenAndServe()
}
