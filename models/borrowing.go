package models

import (
	"task4/database"
	"time"
)

type BorrowingStatus string

const (
	StatusBorrowed BorrowingStatus = "Borrowed"
	StatusReturned BorrowingStatus = "Returned"
	StatusOverdue  BorrowingStatus = "Overdue"
	// 可以添加 StatusReserved, StatusRenewed 等
)

type BorrowingRecord struct {
	ID         uint            `gorm:"primary_key" json:"id"`
	UserID     uint            `json:"user_id"`
	BookID     uint            `json:"book_id"`
	BorrowDate time.Time       `json:"borrow_date"`
	DueDate    time.Time       `json:"due_date"`
	ReturnDate *time.Time      `json:"return_date"`
	Status     BorrowingStatus `gorm:"type:varchar(20)" json:"status"`
	RenewCount int             `gorm:"default:0" json:"renew_count"` // 新增：续借次数
	User       User            `gorm:"foreignkey:UserID"`
	Book       Book            `gorm:"foreignkey:BookID"`
}

// 创建借阅记录
func CreateBorrowingRecord(record *BorrowingRecord) error {
	return database.DB.Create(record).Error
}

// 获取用户的借阅记录
func GetBorrowingRecordsByUser(userID uint) ([]BorrowingRecord, error) {
	var records []BorrowingRecord
	err := database.DB.Preload("Book").Where("user_id = ?", userID).Order("borrow_date desc").Find(&records).Error
	return records, err
}

// 查找某本书当前有效的借阅记录
func GetActiveBorrowingRecordByBook(bookID uint) (*BorrowingRecord, error) {
	var record BorrowingRecord
	err := database.DB.Where("book_id = ? AND status = ?", bookID, StatusBorrowed).First(&record).Error
	return &record, err
}

// 更新借阅记录
func UpdateBorrowingRecord(record *BorrowingRecord) error {
	return database.DB.Save(record).Error
}
