package handlers

import (
	"html/template"
	"little-task-go/db"
	"log"
	"net/http"
)

type Book struct {
	ID            int
	UserID        int
	Title         string
	Author        string
	Publisher     string
	Description   string
	CoverImage    string
	ISBN          string
	PublishedYear int
	Status        string
}

// 模拟从会话中获取当前用户ID
func getSessionUserID(r *http.Request) int {
	// 假设会话中的用户ID为1
	return 1
}

// 显示书籍列表
func Books(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	rows, err := db.DB.Query("SELECT ID, UserID, Title, Author, Publisher, Description, CoverImage, ISBN, PublishedYear, Status FROM books WHERE UserID = ?", userID)
	if err != nil {
		http.Error(w, "Failed to load books", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.UserID, &book.Title, &book.Author, &book.Publisher, &book.Description, &book.CoverImage, &book.ISBN, &book.PublishedYear, &book.Status); err != nil {
			http.Error(w, "Failed to parse book data", http.StatusInternalServerError)
			log.Println("Row scan error:", err)
			return
		}
		books = append(books, book)
	}

	tmpl := template.Must(template.ParseFiles("templates/books.html"))
	tmpl.Execute(w, books)
}

// 添加书籍
func AddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID := getSessionUserID(r)
		title := r.FormValue("title")
		author := r.FormValue("author")
		publisher := r.FormValue("publisher")
		description := r.FormValue("description")
		coverImage := r.FormValue("cover_image")
		isbn := r.FormValue("isbn")
		publishedYear := r.FormValue("published_year")

		_, err := db.DB.Exec("INSERT INTO books (UserID, Title, Author, Publisher, Description, CoverImage, ISBN, PublishedYear, Status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'Available')",
			userID, title, author, publisher, description, coverImage, isbn, publishedYear)
		if err != nil {
			http.Error(w, "Failed to add book", http.StatusInternalServerError)
			log.Println("Insert error:", err)
			return
		}

		http.Redirect(w, r, "/books", http.StatusSeeOther)
	}
}

// 更新书籍
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bookID := r.FormValue("id")
		userID := getSessionUserID(r)
		title := r.FormValue("title")
		author := r.FormValue("author")
		publisher := r.FormValue("publisher")
		description := r.FormValue("description")
		coverImage := r.FormValue("cover_image")
		publishedYear := r.FormValue("published_year")

		_, err := db.DB.Exec("UPDATE books SET Title = ?, Author = ?, Publisher = ?, Description = ?, CoverImage = ?, PublishedYear = ? WHERE ID = ? AND UserID = ?",
			title, author, publisher, description, coverImage, publishedYear, bookID, userID)
		if err != nil {
			http.Error(w, "Failed to update book", http.StatusInternalServerError)
			log.Println("Update error:", err)
			return
		}

		http.Redirect(w, r, "/books", http.StatusSeeOther)
	}
}

// 删除书籍
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bookID := r.FormValue("id")
		userID := getSessionUserID(r)

		_, err := db.DB.Exec("DELETE FROM books WHERE ID = ? AND UserID = ?", bookID, userID)
		if err != nil {
			http.Error(w, "Failed to delete book", http.StatusInternalServerError)
			log.Println("Delete error:", err)
			return
		}

		http.Redirect(w, r, "/books", http.StatusSeeOther)
	}
}
