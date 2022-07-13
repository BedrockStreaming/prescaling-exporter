package handlers

import (
	"fmt"
	"net/http"
)

type IStatusHandlers interface {
	Index(w http.ResponseWriter, r *http.Request)
}

type StatusHandlers struct{}

func NewStatusHandlers() IStatusHandlers {
	return &StatusHandlers{}
}

func (s StatusHandlers) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
