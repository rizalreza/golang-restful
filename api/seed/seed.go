package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/rizalreza/golang-restful/api/models"
)

var users = []models.User{
	models.User{
		Username: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	models.User{
		Username: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var categories = []models.Category{
	models.Category{
		Name: "Category 1",
	},
	models.Category{
		Name: "Category 2",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}, &models.Category{}).Error
	if err != nil {
		log.Fatalf("Cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Category{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("Cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("Attaching User foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("category_id", "categories(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("Attaching Category foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed users table: %v", err)
		}

		err = db.Debug().Model(&models.Category{}).Create(&categories[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed categories table: %v", err)
		}

		posts[i].AuthorID = users[i].ID
		posts[i].CategoryID = categories[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed posts table: %v", err)
		}
	}
}
