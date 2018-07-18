package views

import (
	"net/http"
	"encoding/json"
	)

func Render(w http.ResponseWriter, r *http.Request, data interface{})  {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(js)
}
