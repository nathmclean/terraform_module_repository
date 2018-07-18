package controllers

import (
	"net/http"
	"github.com/nathmclean/terraform_module_repository/views"
)

type WellKnown struct {
	ModulesV1 string `json:"modules.v1"`
}

func NewWellKnown(registryUrl string) *WellKnown {
	return &WellKnown{
		ModulesV1: registryUrl,
	}
}

func (wk *WellKnown) Get(w http.ResponseWriter, r *http.Request) {
	views.Render(w, r, wk)
}
