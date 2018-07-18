package controllers

import (
	"net/http"
	"github.com/gorilla/schema"
	"strconv"
	"strings"
)

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}

func paramToInt(param string) int {
	num, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return num
}

func paramToBool(param string) bool {
	lowerParam := strings.ToLower(param)
	if lowerParam == "true" {
		return true
	}
	return false
}