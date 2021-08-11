package repository

import (
	"errors"
	"gorm.io/gorm"
	"zhiyudong.cn/gin-test/common"
	"zhiyudong.cn/gin-test/model"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return CategoryRepository{
		DB: common.GetDB(),
	}
}

func (c CategoryRepository) Create(name string) (*model.Category, error) {
	category := model.Category{
		Name: name,
	}
	if err := c.DB.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) Update(category model.Category, name string) (*model.Category, error) {
	if err := c.DB.Model(&category).Update("name", name).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) SelectById(id int) (*model.Category, error) {
	var category model.Category
	if err := c.DB.First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) DeleteById(id int) error {
	result := c.DB.Delete(model.Category{}, id)
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("id is not exist")
	}

	return nil
}
