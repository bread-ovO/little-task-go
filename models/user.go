package models

import "task4/database"

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Username string `gorm:"type:varchar(50);unique" json:"username"`
	Password string `gorm:"type:varchar(255)" json:"-"`
	Nickname string `gorm:"type:varchar(50)" json:"nickname"`
	Gender   string `gorm:"type:enum('Male','Female','Other')" json:"gender"`
	Books    []Book `json:"books"`
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据用户 ID 获取用户信息
func GetUserByID(userID uint) (*User, error) {
	var user User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	return database.DB.Save(user).Error
}
