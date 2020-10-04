package controllers

import (
	"errors"
	"net/http"

	"github.com/rizalreza/golang-restful/api/auth"
	"github.com/rizalreza/golang-restful/api/models"
	"github.com/rizalreza/golang-restful/api/responses"
	"github.com/rizalreza/golang-restful/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {

	// Initial data from request body form-data
	email := r.FormValue("email")
	password := r.FormValue("password")

	user := models.User{}
	user.Email = email
	user.Password = password

	user.Prepare()

	err := user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Incorrect Email"))
		return
	}

	token, err := server.SignIn(email, password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.AUTH_JSON(w, http.StatusOK, user, map[string]string{
		"token": token,
	})
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	user.Prepare()
	err := user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	token, err := auth.CreateToken(userCreated.ID)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.AUTH_JSON(w, http.StatusCreated, userCreated, map[string]string{
		"token": token,
	})
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
