package controllers

import (
	"net/http"

	"github.com/rizalreza/golang-restful/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to this API")
}
