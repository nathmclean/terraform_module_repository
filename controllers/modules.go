package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nathmclean/terraform_module_repository/models"
	"github.com/nathmclean/terraform_module_repository/views"
	"net/http"
)

type Modules struct {
	ms models.ModuleService
}

type Meta struct {
	Limit         int `json:"limit"`
	CurrentOffset int `json:"current_offset"`
	NextOffset    int `json:"next_offset"`
	PrevOffset    int `json:"prev_offset"`
}

func NewModules(s models.ModuleService) Modules {
	return Modules{
		ms: s,
	}
}

type ModuleList struct {
	Meta
	Modules []models.Module `json:"modules"`
}

type ModuleGet struct {
	*models.Module
}

func (m *Modules) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controllers: Create")
	var module models.Module
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&module)
	if err != nil {
		fmt.Println("controllers, Create - ", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	err = m.ms.Create(&module)
	if err != nil {
		switch err {
		case models.ErrModuleExists:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte("Success"))
}

func (m *Modules) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controller List")
	v := r.URL.Query()
	limit := paramToInt(v.Get("limit"))
	offset := paramToInt(v.Get("offset"))

	modules, err := m.ms.List(&models.ListRequest{Limit: limit, Offset: offset})
	moduleList := ModuleList{
		Modules: modules,
	}
	if modules == nil {
		fmt.Println("It's nil")
	}
	if err != nil {
		fmt.Println("controllers", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	views.Render(w, r, moduleList)
	return
}

func (m *Modules) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get")
	namespace := mux.Vars(r)["namespace"]
	provider := mux.Vars(r)["provider"]
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	request := &models.GetRequest{
		Namespace: namespace,
		Name:      name,
		Provider:  provider,
		Version:   version,
	}

	module, err := m.ms.Get(request)
	if err != nil {
		fmt.Println("controllers", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	moduleResponse := ModuleGet{
		module,
	}
	views.Render(w, r, moduleResponse)
	return
}

func (m *Modules) ListWithParams(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controller ListWithParams")
	v := r.URL.Query()
	limit := paramToInt(v.Get("limit"))
	offset := paramToInt(v.Get("offset"))
	verified := paramToBool(v.Get("verified"))
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	provider := mux.Vars(r)["provider"]
	version := mux.Vars(r)["version"]

	meta := Meta{
		Limit:         0,
		CurrentOffset: offset,
		NextOffset:    0,
		PrevOffset:    offset,
	}

	req := &models.ListRequest{
		Limit:     limit,
		Offset:    offset,
		Namespace: namespace,
		Verified:  verified,
		Name:      name,
		Provider:  provider,
		Version:   version,
	}

	modules, err := m.ms.ListWithParams(req)
	moduleList := ModuleList{
		Modules: modules,
	}
	if modules == nil {
		fmt.Println("It's nil")
	}
	if err != nil {
		fmt.Println("controllers", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	meta.Limit = req.Limit
	// TODO only include next offset if there's more to get...
	meta.NextOffset = req.Offset + req.Limit
	moduleList.Meta = meta
	views.Render(w, r, moduleList)
	return
}
