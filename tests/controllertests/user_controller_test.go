package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rizalreza/golang-restful/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		statusCode   int
		username     string
		email        string
		password     string
		errorMessage string
	}{
		{
			statusCode:   201,
			username:     "john",
			email:        "john@gmail.com",
			password:     "password",
			errorMessage: "",
		},
		{
			statusCode:   500,
			username:     "cena",
			email:        "john@gmail.com",
			password:     "password",
			errorMessage: "Email Already Taken",
		},
		{
			statusCode:   500,
			username:     "john",
			email:        "cena@gmail.com",
			password:     "password",
			errorMessage: "Username Already Taken",
		},
		{
			statusCode:   422,
			username:     "doe",
			email:        "cena.com",
			password:     "password",
			errorMessage: "Invalid Email",
		},
		{
			statusCode:   422,
			username:     "",
			email:        "cena@gmail.com",
			password:     "password",
			errorMessage: "Required Username",
		},
		{
			statusCode:   422,
			username:     "cena",
			email:        "",
			password:     "password",
			errorMessage: "Required Email",
		},
		{
			username:     "doe",
			email:        "cena@gmail.com",
			password:     "",
			statusCode:   422,
			errorMessage: "Required Password",
		},
	}

	for _, v := range samples {

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("email", v.email)
		_ = writer.WriteField("password", v.password)
		_ = writer.WriteField("username", v.username)
		err := writer.Close()
		if err != nil {
			fmt.Println(err)
		}

		req, err := http.NewRequest("POST", "/users", payload)
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["username"], v.username)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserByID(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id           string
		statusCode   int
		username     string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			username:   user.Username,
			email:      user.Email,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUserById)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Username, responseMap["username"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := seedUsers() //we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updateUsename  string
		updateEmail    string
		updatePassword string
		tokenGiven     string
		errorMessage   string
	}{
		{
			// Convert int32 to int first before converting to string
			id: strconv.Itoa(int(AuthID)),
			// updateJSON:    `{"username":"Grand", "email": "grand@gmail.com", "password": "password"}`,
			statusCode:     200,
			updateUsename:  "Grand",
			updateEmail:    "grand@gmail.com",
			updatePassword: "password",
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			// When password field is empty
			id: strconv.Itoa(int(AuthID)),
			// updateJSON:   `{"username":"Woman", "email": "woman@gmail.com", "password": ""}`,
			updateUsename:  "Woman",
			updateEmail:    "woman@gmail.com",
			updatePassword: "",
			statusCode:     422,
			tokenGiven:     tokenString,
			errorMessage:   "Required Password",
		},
		{
			// When no token was passed
			id: strconv.Itoa(int(AuthID)),
			// updateJSON:   `{"username":"Man", "email": "man@gmail.com", "password": "password"}`,
			updateUsename:  "Man",
			updateEmail:    "man@gmail.com",
			updatePassword: "password",
			statusCode:     401,
			tokenGiven:     "",
			errorMessage:   "Unauthorized",
		},
		{
			// When incorrect token was passed
			id: strconv.Itoa(int(AuthID)),
			// updateJSON:   `{"username":"Woman", "email": "woman@gmail.com", "password": "password"}`,
			updateUsename:  "Woman",
			updateEmail:    "woman@gmail.com",
			updatePassword: "password",
			statusCode:     401,
			tokenGiven:     "This is incorrect token",
			errorMessage:   "Unauthorized",
		},
		{
			// Remember "kenny@gmail.com" belongs to user 2
			id: strconv.Itoa(int(AuthID)),
			// updateJSON:   `{"username":"Frank", "email": "kenny@gmail.com", "password": "password"}`,
			updateUsename:  "Frank",
			updateEmail:    "doe@gmail.com",
			updatePassword: "password",
			statusCode:     500,
			tokenGiven:     tokenString,
			errorMessage:   "Email Already Taken",
		},
		{
			// Remember "Kenny Morris" belongs to user 2
			id: strconv.Itoa(int(AuthID)),
			// updateJSON:   `{"username":"Kenny Morris", "email": "grand@gmail.com", "password": "password"}`,
			updateUsename:  "doe",
			updateEmail:    "kenny@gmail.com",
			updatePassword: "password",
			statusCode:     500,
			tokenGiven:     tokenString,
			errorMessage:   "Username Already Taken",
		},
		{
			id: strconv.Itoa(int(AuthID)),
			// updateJSON:     `{"username":"Kan", "email": "kangmail.com", "password": "password"}`,
			updateUsename:  "kan",
			updateEmail:    "kan.com",
			updatePassword: "password",
			statusCode:     422,
			tokenGiven:     tokenString,
			errorMessage:   "Invalid Email format",
		},
		{
			id:             strconv.Itoa(int(AuthID)),
			updateUsename:  "",
			updateEmail:    "kan@gmail.com",
			updatePassword: "password",
			statusCode:     422,
			tokenGiven:     tokenString,
			errorMessage:   "Required Username",
		},
		{
			id:             strconv.Itoa(int(AuthID)),
			updateUsename:  "kan",
			updateEmail:    "",
			updatePassword: "password",
			statusCode:     422,
			tokenGiven:     tokenString,
			errorMessage:   "Required Email",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When user 2 is using user 1 token
			id:             strconv.Itoa(int(2)),
			updateUsename:  "mike",
			updateEmail:    "mike@gmail.com",
			updatePassword: "password",
			tokenGiven:     tokenString,
			statusCode:     401,
			errorMessage:   "Unauthorized",
		},
	}

	for _, v := range samples {

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("email", v.updateEmail)
		_ = writer.WriteField("password", v.updatePassword)
		_ = writer.WriteField("username", v.updateUsename)
		err := writer.Close()
		if err != nil {
			fmt.Println(err)
		}

		req, err := http.NewRequest("PUT", "/users", payload)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["username"], v.updateUsename)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {

	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() //we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Get only the first and log him in
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		AuthID = user.ID
		AuthEmail = user.Email
		AuthPassword = "password" ////Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When no token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is given
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// User 2 trying to use User 1 token
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
