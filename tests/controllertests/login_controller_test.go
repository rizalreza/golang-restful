package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("This is the error %v\n", err)
	}

	samples := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     "password", //Note the password has to be this, not the hashed one from the database
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "Wrong email",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {

		token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

func TestLogin(t *testing.T) {

	refreshUserTable()

	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("This is the error %v\n", err)
	}
	samples := []struct {
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        "john@gmail.com",
			password:     "password",
			statusCode:   200,
			errorMessage: "",
		},
		{
			email:        "john@gmail.com",
			password:     "wrong password",
			statusCode:   422,
			errorMessage: "Incorrect Password",
		},
		{
			email:        "wow@gmail.com",
			password:     "password",
			statusCode:   404,
			errorMessage: "Incorrect Email",
		},
		{
			email:        "wow.com",
			password:     "password",
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			email:        "",
			password:     "password",
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			email:        "wow.com",
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
		err := writer.Close()
		if err != nil {
			fmt.Println(err)
		}

		req, err := http.NewRequest("POST", "/login", payload)
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}

		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
