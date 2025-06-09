package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"task4/database"
	"task4/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const UploadDir = "static/uploads"

// 统一错误处理函数
func handleError(c *gin.Context, statusCode int, message string) {
	log.Printf("Error: %s", message)
	c.JSON(statusCode, gin.H{"error": message})
}

// 展示书籍页面
func ShowBooksPage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	keyword := c.Query("keyword")
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize := 10 // 调整为 10 或其他合适值

	// 使用 GetBooks 获取所有书籍
	books, total, err := models.GetBooks(keyword, page, pageSize) //
	if err != nil {
		handleError(c, http.StatusInternalServerError, "获取书籍失败")
		return
	}

	// 获取当前用户正在借阅的书籍 ID 列表 (用于前端判断)
	var borrowedIDs []uint
	database.DB.Model(&models.BorrowingRecord{}).Where("user_id = ? AND status = ?", userID.(uint), models.StatusBorrowed).Pluck("book_id", &borrowedIDs) //

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "books.html", gin.H{
		"nickname":    c.GetString("nickname"),
		"gender":      c.GetString("gender"),
		"books":       books,
		"keyword":     keyword,
		"currentPage": page,
		"totalPages":  totalPages,
		"userID":      userID.(uint), // 传递用户ID
		"borrowedIDs": borrowedIDs,   // 传递已借书籍ID
	})
}

// 添加书籍
func AddBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	bookNumber := c.PostForm("book_number")
	title := c.PostForm("title")
	author := c.PostForm("author")
	publisher := c.PostForm("publisher")
	description := c.PostForm("description")

	coverImage := "default_cover.png" // 默认封面
	file, err := c.FormFile("cover_image")
	if err == nil { // 如果用户上传了文件
		if !strings.HasSuffix(file.Filename, ".jpg") && !strings.HasSuffix(file.Filename, ".png") { //
			handleError(c, http.StatusBadRequest, "仅支持上传jpg或png格式的图片")
			return
		}
		coverImagePath := fmt.Sprintf("%s%s_%d_%s", UploadDir, uuid.New().String(), time.Now().Unix(), file.Filename) //
		if err := c.SaveUploadedFile(file, coverImagePath); err != nil {
			handleError(c, http.StatusInternalServerError, "封面图片上传失败")
			return
		}
		coverImage = coverImagePath
	}

	book := models.Book{
		UserID:      userID.(uint),
		BookNumber:  bookNumber,
		Title:       title,
		Author:      author,
		Publisher:   publisher,
		Description: description,
		CoverImage:  coverImage,
	}

	if err := models.CreateBook(&book); err != nil { //
		handleError(c, http.StatusInternalServerError, "添加书籍失败")
		return
	}

	c.Redirect(http.StatusFound, "/books")
}

// 更新书籍
func UpdateBookHandler(c *gin.Context) {
	_, exists := c.Get("user_id")
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
	// =========== 修改开始 ===========
	// 原有代码检查了 user_id，导致只有创建者才能修改。
	// if err := database.DB.Where("id = ? AND user_id = ?", bookID, userID.(uint)).First(&book).Error; err != nil {
	// 现修改为只通过书籍ID查找，允许任何登录用户修改。
	if err := database.DB.First(&book, bookID).Error; err != nil {
		handleError(c, http.StatusNotFound, "书籍不存在")
		return
	}
	// =========== 修改结束 ===========

	// 更新字段
	book.BookNumber = c.PostForm("book_number")
	book.Title = c.PostForm("title")
	book.Author = c.PostForm("author")
	book.Publisher = c.PostForm("publisher")
	book.Description = c.PostForm("description")

	// 处理封面图片上传
	file, err := c.FormFile("cover_image")
	if err == nil { // 如果用户上传了文件
		if !strings.HasSuffix(file.Filename, ".jpg") && !strings.HasSuffix(file.Filename, ".png") { //
			handleError(c, http.StatusBadRequest, "仅支持上传jpg或png格式的图片")
			return
		}
		coverImagePath := fmt.Sprintf("%s%s_%d_%s", UploadDir, uuid.New().String(), time.Now().Unix(), file.Filename) //
		if err := c.SaveUploadedFile(file, coverImagePath); err != nil {
			handleError(c, http.StatusInternalServerError, "封面图片上传失败")
			return
		}
		// 删除旧封面文件（如果不是默认封面）
		if book.CoverImage != "default_cover.png" {
			os.Remove(book.CoverImage)
		}
		book.CoverImage = coverImagePath
	}

	// 更新书籍信息到数据库
	if err := models.UpdateBook(&book); err != nil { //
		handleError(c, http.StatusInternalServerError, "更新书籍失败")
		return
	}

	c.Redirect(http.StatusFound, "/books")
}

// 删除书籍
func DeleteBookHandler(c *gin.Context) {
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
	if err := database.DB.Where("id = ? AND user_id = ?", bookID, userID.(uint)).First(&book).Error; err != nil { //
		handleError(c, http.StatusNotFound, "书籍不存在")
		return
	}

	// 删除书籍记录
	if err := models.DeleteBook(uint(bookID), userID.(uint)); err != nil { //
		handleError(c, http.StatusInternalServerError, "删除书籍失败")
		return
	}

	// 删除旧封面文件（如果不是默认封面）
	if book.CoverImage != "default_cover.png" {
		os.Remove(book.CoverImage)
	}

	c.Redirect(http.StatusFound, "/books")
}

// 渲染添加书籍页面
func ShowAddBookPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_book.html", nil)
}

func createInitialBooksForUser(userID uint) error {
	var defaultBooks []models.Book
	for i := 0; i < 100; i++ {
		defaultBooks = append(defaultBooks, models.Book{
			UserID:      userID,
			BookNumber:  fmt.Sprintf("DEF%03d", i+1),
			Title:       fmt.Sprintf("默认书籍 %d", i+1),
			Author:      "系统默认",
			Publisher:   "系统默认",
			Description: "默认书籍内容",
			CoverImage:  "static/default_cover.png",
		})
	}
	return database.DB.Create(&defaultBooks).Error
}

func ShowUpdateBookPage(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	idStr := c.Param("id")
	bookID, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的书籍ID")
		return
	}

	var book models.Book
	// =========== 修改开始 ===========
	// 原有代码检查了 user_id，导致只有创建者才能修改。
	// if err := database.DB.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
	// 现修改为只通过书籍ID查找，允许任何登录用户修改。
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.String(http.StatusNotFound, "书籍不存在")
		return
	}
	// =========== 修改结束 ===========

	c.HTML(http.StatusOK, "book_edit.html", gin.H{
		"Book": book,
	})
}

func ReturnBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	bookIDStr := c.Param("id")
	bookID, _ := strconv.Atoi(bookIDStr)

	var record models.BorrowingRecord
	err := database.DB.Where("book_id = ? AND user_id = ? AND status = ?", bookID, userID.(uint), models.StatusBorrowed).First(&record).Error //
	if err != nil {
		handleError(c, http.StatusNotFound, "未找到借阅记录")
		return
	}

	tx := database.DB.Begin()

	// 查找是否有预定
	nextReservation, resErr := models.GetNextReservation(uint(bookID)) //

	newStatus := models.BookStatusAvailable //
	if resErr == nil && nextReservation.ID > 0 {
		// 如果有预定，将书状态改为 Reserved
		newStatus = models.BookStatusReserved //
		// 更新预定记录为 Fulfilled (或添加通知逻辑)
		nextReservation.Status = models.ReservationFulfilled //
		if err := tx.Save(&nextReservation).Error; err != nil {
			tx.Rollback()
			handleError(c, http.StatusInternalServerError, "更新预定记录失败")
			return
		}
		log.Printf("Book %d returned, now reserved for user %d", bookID, nextReservation.UserID)
	}

	// 更新书籍状态
	if err := tx.Model(&models.Book{}).Where("id = ?", bookID).Update("status", newStatus).Error; err != nil { //
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "更新书籍状态失败")
		return
	}

	// 更新借阅记录
	now := time.Now()
	record.ReturnDate = &now
	record.Status = models.StatusReturned //
	if err := tx.Save(&record).Error; err != nil {
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "更新借阅记录失败")
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, c.Request.Referer()) // 返回上一页
}

func ReserveBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	bookIDStr := c.Param("id")
	bookID, _ := strconv.Atoi(bookIDStr)

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		handleError(c, http.StatusNotFound, "书籍不存在")
		return
	}

	// 只能预定已借出的书，且不能是自己借的，且不能重复预定
	if book.Status != models.BookStatusBorrowed { //
		handleError(c, http.StatusBadRequest, "该书当前无需预定")
		return
	}
	if models.HasUserReserved(userID.(uint), uint(bookID)) { //
		handleError(c, http.StatusBadRequest, "您已预定此书")
		return
	}
	activeBorrow, _ := models.GetActiveBorrowingRecordByBook(uint(bookID)) //
	if activeBorrow != nil && activeBorrow.UserID == userID.(uint) {
		handleError(c, http.StatusBadRequest, "您正在借阅此书，无需预定")
		return
	}

	reservation := models.ReservationRecord{
		UserID:          userID.(uint),
		BookID:          uint(bookID),
		ReservationDate: time.Now(),
		Status:          models.ReservationActive, //
	}

	if err := models.CreateReservation(&reservation); err != nil { //
		handleError(c, http.StatusInternalServerError, "预定失败")
		return
	}

	// 可以加一个成功提示，然后重定向
	c.Redirect(http.StatusFound, "/books")
}

// 续借书籍
func RenewBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	bookIDStr := c.Param("id")
	bookID, _ := strconv.Atoi(bookIDStr)
	const maxRenewCount = 1 // 假设最多续借1次

	var record models.BorrowingRecord
	err := database.DB.Where("book_id = ? AND user_id = ? AND status = ?", bookID, userID.(uint), models.StatusBorrowed).First(&record).Error //
	if err != nil {
		handleError(c, http.StatusNotFound, "未找到借阅记录")
		return
	}

	// 检查续借条件
	if record.RenewCount >= maxRenewCount { //
		handleError(c, http.StatusBadRequest, "已达到最大续借次数")
		return
	}
	if time.Now().After(record.DueDate) {
		handleError(c, http.StatusBadRequest, "已超期，无法续借")
		return
	}
	_, resErr := models.GetNextReservation(uint(bookID)) //
	if resErr == nil {
		handleError(c, http.StatusBadRequest, "该书已被预定，无法续借")
		return
	}

	// 更新记录
	record.DueDate = record.DueDate.AddDate(0, 1, 0) // 再续借1个月
	record.RenewCount++

	if err := models.UpdateBorrowingRecord(&record); err != nil { //
		handleError(c, http.StatusInternalServerError, "续借失败")
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer()) // 返回上一页
}

// 显示我的借阅页面
func ShowMyBorrowsPage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var records []models.BorrowingRecord
	err := database.DB.Preload("Book").Where("user_id = ? AND status = ?", userID.(uint), models.StatusBorrowed).Order("due_date asc").Find(&records).Error //
	if err != nil {
		handleError(c, http.StatusInternalServerError, "获取我的借阅失败")
		return
	}

	// 检查是否有预定 (用于判断是否能续借)
	canRenew := make(map[uint]bool)
	for _, rec := range records {
		_, resErr := models.GetNextReservation(rec.BookID) //
		canRenew[rec.BookID] = (resErr != nil)             // 如果没有预定 (resErr != nil), 则可以续借
	}

	c.HTML(http.StatusOK, "my_borrows.html", gin.H{
		"nickname": c.GetString("nickname"),
		"gender":   c.GetString("gender"),
		"records":  records,
		"canRenew": canRenew,
	})
}
