package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"task4/database"
	"task4/models"

	"time"
)

// 借阅书籍
func BorrowBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, "无效的书籍ID")
		return
	}

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		handleError(c, http.StatusNotFound, "书籍不存在")
		return
	}

	// 检查书籍是否属于该用户 (如果系统设计为用户管理自己的书库，则此检查可能需要调整)
	// 在借阅系统中，通常是检查书籍是否可借
	// if book.UserID != userID.(uint) {
	//     handleError(c, http.StatusForbidden, "您无权操作此书籍")
	//     return
	// }

	// 检查书籍状态
	if book.Status != models.BookStatusAvailable {
		handleError(c, http.StatusBadRequest, "该书籍当前不可借阅")
		return
	}

	// 开始事务 (推荐用于多步数据库操作)
	tx := database.DB.Begin()

	// 更新书籍状态
	if err := tx.Model(&models.Book{}).Where("id = ?", bookID).Update("status", models.BookStatusBorrowed).Error; err != nil {
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "更新书籍状态失败")
		return
	}

	// 创建借阅记录
	borrowingRecord := models.BorrowingRecord{
		UserID:     userID.(uint),
		BookID:     uint(bookID),
		BorrowDate: time.Now(),
		DueDate:    time.Now().AddDate(0, 1, 0), // 假设借期为1个月
		Status:     models.StatusBorrowed,
	}

	if err := tx.Create(&borrowingRecord).Error; err != nil {
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "创建借阅记录失败")
		return
	}

	// 提交事务
	tx.Commit()

	c.Redirect(http.StatusFound, "/books") // 或重定向到借阅历史页面
}

// (可选) 显示用户借阅历史页面
func ShowBorrowingHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	records, err := models.GetBorrowingRecordsByUser(userID.(uint))
	if err != nil {
		handleError(c, http.StatusInternalServerError, "获取借阅历史失败")
		return
	}

	c.HTML(http.StatusOK, "history.html", gin.H{
		"nickname": c.GetString("nickname"),
		"gender":   c.GetString("gender"),
		"records":  records,
	})
}
