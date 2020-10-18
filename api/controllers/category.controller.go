package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rizalreza/golang-restful/api/models"
	"github.com/rizalreza/golang-restful/api/responses"
	"github.com/rizalreza/golang-restful/api/utils/formaterror"
)

func (server *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}

	category.Name = r.FormValue("name")

	category.Prepare()
	err := category.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	categoryCreated, err := category.SaveCategory(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, categoryCreated.ID))
	responses.JSON(w, http.StatusCreated, categoryCreated)
}

func (server *Server) GetCategories(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}

	categories, err := category.GetAllCategory(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, categories)
}

func (server *Server) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	category := models.Category{}
	categoryRecieved, err := category.FindCategoryById(server.DB, uint32(cid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, categoryRecieved)
}

func (server *Server) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if category id is valid
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	category := models.Category{}
	err = server.DB.Debug().Model(models.Category{}).Where("id = ?", cid).Take(&category).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Category not found"))
		return
	}

	categoryUpdate := models.Category{}
	categoryUpdate.Name = r.FormValue("name")

	categoryUpdate.Prepare()
	err = categoryUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	categoryUpdated, err := categoryUpdate.UpdateCategory(server.DB, uint32(cid))
	categoryUpdated.ID = uint32(cid)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, categoryUpdated)

}

func (server *Server) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// Check if the post exist
	category := models.Category{}
	err = server.DB.Debug().Model(models.Category{}).Where("id = ?", cid).Take(&category).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Data not found"))
		return
	}

	_, err = category.DeleteCategory(server.DB, uint32(cid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", cid))
	responses.JSON(w, http.StatusNoContent, "")
}
