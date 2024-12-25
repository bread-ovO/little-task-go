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
	pageSize := 20

	books, total, err := models.GetBooksByUser(userID.(uint), keyword, page, pageSize)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "获取书籍失败")
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "books.html", gin.H{
		"nickname":    c.GetString("nickname"),
		"gender":      c.GetString("gender"),
		"books":       books,
		"keyword":     keyword,
		"currentPage": page,
		"totalPages":  totalPages,
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
		if !strings.HasSuffix(file.Filename, ".jpg") && !strings.HasSuffix(file.Filename, ".png") {
			handleError(c, http.StatusBadRequest, "仅支持上传jpg或png格式的图片")
			return
		}
		coverImagePath := fmt.Sprintf("%s%s_%d_%s", UploadDir, uuid.New().String(), time.Now().Unix(), file.Filename)
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

	if err := models.CreateBook(&book); err != nil {
		handleError(c, http.StatusInternalServerError, "添加书籍失败")
		return
	}

	c.Redirect(http.StatusFound, "/books")
}

// 更新书籍
func UpdateBookHandler(c *gin.Context) {
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
	if err := database.DB.Where("id = ? AND user_id = ?", bookID, userID.(uint)).First(&book).Error; err != nil {
		handleError(c, http.StatusNotFound, "书籍不存在")
		return
	}

	// 更新字段
	book.BookNumber = c.PostForm("book_number")
	book.Title = c.PostForm("title")
	book.Author = c.PostForm("author")
	book.Publisher = c.PostForm("publisher")
	book.Description = c.PostForm("description")

	// 处理封面图片上传
	file, err := c.FormFile("cover_image")
	if err == nil { // 如果用户上传了文件
		if !strings.HasSuffix(file.Filename, ".jpg") && !strings.HasSuffix(file.Filename, ".png") {
			handleError(c, http.StatusBadRequest, "仅支持上传jpg或png格式的图片")
			return
		}
		coverImagePath := fmt.Sprintf("%s%s_%d_%s", UploadDir, uuid.New().String(), time.Now().Unix(), file.Filename)
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
	if err := models.UpdateBook(&book); err != nil {
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
	if err := database.DB.Where("id = ? AND user_id = ?", bookID, userID.(uint)).First(&book).Error; err != nil {
		handleError(c, http.StatusNotFound, "书籍不存在")
		return
	}

	// 删除书籍记录
	if err := models.DeleteBook(uint(bookID), userID.(uint)); err != nil {
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
	userID, exists := c.Get("user_id")
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
	if err := database.DB.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
		c.String(http.StatusNotFound, "书籍不存在")
		return
	}

	c.HTML(http.StatusOK, "book_edit.html", gin.H{
		"Book": book,
	})
}
