package models

import (
	"task4/database"
)

type BookStatus string

const (
	BookStatusAvailable BookStatus = "Available"
	BookStatusBorrowed  BookStatus = "Borrowed"
	BookStatusReserved  BookStatus = "Reserved" // 确保已存在
)

type Book struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	UserID      uint       `json:"user_id"`
	BookNumber  string     `gorm:"type:varchar(50)" json:"book_number"`
	Title       string     `gorm:"type:varchar(255)" json:"title"`
	Author      string     `gorm:"type:varchar(255)" json:"author"`
	Publisher   string     `gorm:"type:varchar(255)" json:"publisher"`
	Description string     `gorm:"type:text" json:"description"`
	CoverImage  string     `gorm:"type:varchar(255)" json:"cover_image"`
	Status      BookStatus `gorm:"type:varchar(20);default:'Available'" json:"status"` // 新增字段
}

func GetBooks(keyword string, page int, pageSize int) ([]Book, int, error) {
	var books []Book
	var total int

	query := database.DB.Model(&Book{}) // 获取所有书

	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR author LIKE ? OR publisher LIKE ?", likeKeyword, likeKeyword, likeKeyword, likeKeyword)
	}
	query.Count(&total)
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

func GetBookForUpdate(bookID uint) (*Book, error) {
	var book Book
	tx := database.DB.Begin()
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&book, bookID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// 注意: 这里需要返回 tx，并在控制器中 Commit 或 Rollback
	// 或者设计一个更完善的事务处理模式
	tx.Commit() // 简化示例，实际应用中应更谨慎
	return &book, nil
}

// 更新书籍状态
func UpdateBookStatus(bookID uint, status BookStatus) error {
	return database.DB.Model(&Book{}).Where("id = ?", bookID).Update("status", status).Error
}
