package models

import (
	"task4/database"
	"time"
)

type ReservationStatus string

const (
	ReservationActive    ReservationStatus = "Active"
	ReservationFulfilled ReservationStatus = "Fulfilled"
	ReservationCanceled  ReservationStatus = "Canceled"
)

type ReservationRecord struct {
	ID              uint              `gorm:"primary_key" json:"id"`
	UserID          uint              `json:"user_id"`
	BookID          uint              `json:"book_id"`
	ReservationDate time.Time         `json:"reservation_date"`
	Status          ReservationStatus `gorm:"type:varchar(20)" json:"status"`
	User            User              `gorm:"foreignkey:UserID"`
	Book            Book              `gorm:"foreignkey:BookID"`
}

// 创建预定记录
func CreateReservation(record *ReservationRecord) error {
	return database.DB.Create(record).Error
}

// 查找某本书的下一个有效预定
func GetNextReservation(bookID uint) (*ReservationRecord, error) {
	var record ReservationRecord
	err := database.DB.Where("book_id = ? AND status = ?", bookID, ReservationActive).Order("reservation_date asc").First(&record).Error
	return &record, err
}

// 检查用户是否已预定某本书
func HasUserReserved(userID, bookID uint) bool {
	var count int
	database.DB.Model(&ReservationRecord{}).Where("user_id = ? AND book_id = ? AND status = ?", userID, bookID, ReservationActive).Count(&count)
	return count > 0
}

// 更新预定记录
func UpdateReservation(record *ReservationRecord) error {
	return database.DB.Save(record).Error
}
