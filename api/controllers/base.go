package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rizalreza/golang-restful/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Driver, User, Password, Port, Host, Name string) {
	var err error

	if Driver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", User, Password, Host, Port, Name)
		server.DB, err = gorm.Open(Driver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Driver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Driver)
		}
	}

	if Driver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", Host, Port, User, Name, Password)
		server.DB, err = gorm.Open(Driver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Driver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Driver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{})
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) Run(address string) {
	fmt.Println("Listening to port 8090")
	log.Fatal(http.ListenAndServe(address, server.Router))
}
