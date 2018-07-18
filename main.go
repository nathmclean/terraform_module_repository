package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nathmclean/terraform_module_repository/controllers"
	"github.com/nathmclean/terraform_module_repository/models"
	"net/http"
)

const version = "v1"
const host = "localhost"
const port = 4000

func main() {
	// services
	services, err := models.NewServices(fmt.Sprintf("host=%s port=%d user=%s dbname=%s "+
		"sslmode=disable", "localhost", 5432, "nathan", "terraform_repo"))
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	// controllers
	wellKnownC := controllers.NewWellKnown(fmt.Sprintf("https://localhost:%d/%s/", port, version))
	modulesC := controllers.NewModules(services.Models)

	r := mux.NewRouter()

	// well known
	r.HandleFunc("/.well-known/terraform.json", wellKnownC.Get).Methods("GET")

	// list
	// list all modules
	r.HandleFunc(fmt.Sprintf("/%s/", version), modulesC.List).Methods("GET")
	// list all modules in a namespace
	r.HandleFunc(fmt.Sprintf("/%s/{namespace}", version), modulesC.ListWithParams).Methods("GET")
	// latest version of each provider for a module
	r.HandleFunc(fmt.Sprintf("/%s/{namespace}/{name}", version), modulesC.ListWithParams).Methods("GET")
	// list all versions of a module
	r.HandleFunc(fmt.Sprintf("/%s/{namespace}/{name}/{provider}/versions", version), modulesC.ListWithParams).Methods("GET")

	// get modules
	// get the latest version for a specific module provider
	r.HandleFunc(fmt.Sprintf("/%s/{namespace}/{name}/{provider}", version), nil).Methods("GET")
	// get a specific version of a module for a single provider
	r.HandleFunc(fmt.Sprintf("/%s/{namespace}/{name}/{provider}/{version}", version), nil).Methods("GET")

	// download
	// download specific version
	r.HandleFunc(fmt.Sprintf("/%s/{namespace}/{name}/{provider}/{version}/download", version), nil).Methods("GET")
	// download latest version
	r.HandleFunc(fmt.Sprintf("/%s/{namespace}/{name}/{provider}/download", version), nil).Methods("GET")

	// search
	r.HandleFunc(fmt.Sprintf("/%s/search", version), nil).Methods("GET")

	// create
	r.HandleFunc(fmt.Sprintf("/%s/", version), modulesC.Create).Methods("Post")

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r)
	if err != nil {
		panic(err)
	}
}
