package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/rizalreza/golang-restful/api/controllers"
	"github.com/rizalreza/golang-restful/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var categortInstance = models.Category{}
var postInstance = models.Post{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())

}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	server.DB.Exec("SET foreign_key_checks=0")
	err := server.DB.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	server.DB.Exec("SET foreign_key_checks=1")
	err = server.DB.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func refreshCategoryTable() error {
	server.DB.Exec("SET foreign_key_checks=0")
	err := server.DB.DropTableIfExists(&models.Category{}).Error
	if err != nil {
		return err
	}

	server.DB.Exec("SET foreign_key_checks=1")
	err = server.DB.AutoMigrate(&models.Category{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed categories table")
	log.Printf("refreshCategoryTable routine OK !!!")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Username: "john",
		Email:    "john@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedOneCategory() (models.Category, error) {
	_ = refreshCategoryTable()

	category := models.Category{
		Name: "Category 1",
	}

	err := server.DB.Model(&models.Category{}).Create(&category).Error
	if err != nil {
		log.Fatalf("Cannot seed categories table: %v", err)
	}

	log.Printf("seedOneCategory routine OK !!!")
	return category, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			Username: "john",
			Email:    "john@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "doe",
			Email:    "doe@gmail.com",
			Password: "password",
		},
	}
	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserCategoryAndPostTable() error {
	server.DB.Exec("SET foreign_key_checks=0")
	err := server.DB.Debug().DropTableIfExists(&models.Post{}, &models.Category{}, &models.User{}).Error
	if err != nil {
		return err
	}

	server.DB.Exec("SET foreign_key_checks=1")
	err = server.DB.Debug().AutoMigrate(&models.User{}, &models.Category{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	log.Printf("refreshUserCategoryAndPostTable routine OK !!!")
	return nil
}

func seedOneUserOneCategoryAndOnePost() (models.Post, error) {

	err := refreshUserCategoryAndPostTable()
	if err != nil {
		return models.Post{}, err
	}
	user := models.User{
		Username: "john",
		Email:    "john@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Post{}, err
	}

	category := models.Category{
		Name: "Category Test 1",
	}

	err = server.DB.Model(&models.Category{}).Create(&category).Error
	if err != nil {
		return models.Post{}, err
	}

	post := models.Post{
		Title:      "This is the title sam",
		Content:    "This is the content sam",
		AuthorID:   user.ID,
		CategoryID: category.ID,
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func SeedUsersCategoriesAndPosts() ([]models.User, []models.Category, []models.Post, error) {
	var err error

	if err != nil {
		return []models.User{}, []models.Category{}, []models.Post{}, err
	}

	var users = []models.User{
		models.User{
			Username: "mike",
			Email:    "mike@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "shinoda",
			Email:    "shinoda@gmail.com",
			Password: "password",
		},
	}
	var categories = []models.Category{
		models.Category{
			Name: "Category Test 1",
		},
		models.Category{
			Name: "Category Test 2",
		},
	}
	var posts = []models.Post{
		models.Post{
			Title:   "First Title",
			Content: "First Content",
		},
		models.Post{
			Title:   "Second Title",
			Content: "Second Content",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

		err = server.DB.Model(&models.Category{}).Create(&categories[i]).Error
		if err != nil {
			log.Fatalf("cannot seed categories table: %v", err)
		}

		posts[i].AuthorID = users[i].ID
		posts[i].CategoryID = categories[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
	return users, categories, posts, nil
}
