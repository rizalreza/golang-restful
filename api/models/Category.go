package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Category struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `grom:"size:100;no null;unique" json:"name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Category) Prepare() {
	c.ID = 0
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("Required Name")
	}
	return nil
}

func (c *Category) SaveCategory(db *gorm.DB) (*Category, error) {
	var err error
	err = db.Debug().Model(&Category{}).Create(&c).Error
	if err != nil {
		return &Category{}, err
	}

	return c, nil
}

func (c *Category) GetAllCategory(db *gorm.DB) (*[]Category, error) {
	var err error
	categories := []Category{}
	err = db.Debug().Model(&Category{}).Limit(100).Find(&categories).Error
	if err != nil {
		return &[]Category{}, err
	}

	return &categories, err
}

func (c *Category) FindCategoryById(db *gorm.DB, cid uint32) (*Category, error) {
	var err error
	err = db.Debug().Model(Category{}).Where("id = ?", cid).Take(&c).Error
	if err != nil {
		return &Category{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &Category{}, errors.New("Category no found")
	}

	return c, err

}

func (c *Category) UpdateCategory(db *gorm.DB, cid uint32) (*Category, error) {
	var err error

	err = db.Debug().Model(&Category{}).Where("id = ?", cid).Updates(Category{Name: c.Name, UpdatedAt: time.Now()}).Error

	if err != nil {
		return &Category{}, err
	}

	return c, nil
}

func (c *Category) DeleteCategory(db *gorm.DB, cid uint32) (int64, error) {
	db = db.Debug().Model(&Category{}).Where("id = ?", cid).Take(&Category{}).Delete(&Category{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
