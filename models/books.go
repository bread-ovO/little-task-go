package models

import (
	"task4/database"
)

type Book struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	UserID      uint   `json:"user_id"`
	BookNumber  string `gorm:"type:varchar(50)" json:"book_number"`
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Author      string `gorm:"type:varchar(255)" json:"author"`
	Publisher   string `gorm:"type:varchar(255)" json:"publisher"`
	Description string `gorm:"type:text" json:"description"`
	CoverImage  string `gorm:"type:varchar(255)" json:"cover_image"`
}

func GetBooksByUser(userID uint, keyword string, page int, pageSize int) ([]Book, int, error) {
	var books []Book
	var total int

	query := database.DB.Where("user_id = ?", userID)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", likeKeyword, likeKeyword)
	}
	query.Model(&Book{}).Count(&total)
	err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&books).Error
	return books, total, err
}

func CreateBook(book *Book) error {
	return database.DB.Create(book).Error
}

func UpdateBook(book *Book) error {
	return database.DB.Save(book).Error
}

func DeleteBook(bookID uint, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", bookID, userID).Delete(&Book{}).Error
}
