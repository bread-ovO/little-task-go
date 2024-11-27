package handlers

import (
	"html/template"
	"little-task-go/db"
	"net/http"
	"strconv"
)

type Book struct {
	ID          int
	UserID      int
	ISBN        string
	Title       string
	Author      string
	Publisher   string
	Description string
	CoverImage  string
}

func getSessionUserID(r *http.Request) int {
	// 模拟获取会话中的用户ID，假设为1
	return 1
}

func Books(w http.ResponseWriter, r *http.Request) {
	// 获取分页参数
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageSize := 20
	offset, _ := strconv.Atoi(page)
	offset = (offset - 1) * pageSize

	// 获取搜索关键字
	keyword := r.URL.Query().Get("keyword")
	userID := getSessionUserID(r)

	// 查询书籍
	var books []Book
	query := "SELECT id, user_id, isbn, title, author, publisher, description, cover_image FROM books WHERE user_id = ?"
	args := []interface{}{userID}

	if keyword != "" {
		query += " AND (title LIKE ? OR description LIKE ?)"
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Failed to load books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.UserID, &book.ISBN, &book.Title, &book.Author, &book.Publisher, &book.Description, &book.CoverImage); err != nil {
			http.Error(w, "Failed to parse book data", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// 渲染模板
	tmpl := template.Must(template.ParseFiles("templates/books.html"))
	tmpl.Execute(w, books)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID := getSessionUserID(r)
		isbn := r.FormValue("isbn")
		title := r.FormValue("title")
		author := r.FormValue("author")
		publisher := r.FormValue("publisher")
		description := r.FormValue("description")
		coverImage := r.FormValue("cover_image")

		_, err := db.DB.Exec("INSERT INTO books (user_id, isbn, title, author, publisher, description, cover_image) VALUES (?, ?, ?, ?, ?, ?, ?)",
			userID, isbn, title, author, publisher, description, coverImage)
		if err != nil {
			http.Error(w, "Failed to add book", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/books", http.StatusSeeOther)
	}
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bookID := r.FormValue("id")
		title := r.FormValue("title")
		author := r.FormValue("author")
		publisher := r.FormValue("publisher")
		description := r.FormValue("description")
		coverImage := r.FormValue("cover_image")
		userID := getSessionUserID(r)

		_, err := db.DB.Exec("UPDATE books SET title = ?, author = ?, publisher = ?, description = ?, cover_image = ? WHERE id = ? AND user_id = ?",
			title, author, publisher, description, coverImage, bookID, userID)
		if err != nil {
			http.Error(w, "Failed to update book", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/books", http.StatusSeeOther)
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bookID := r.FormValue("id")
		userID := getSessionUserID(r)

		_, err := db.DB.Exec("DELETE FROM books WHERE id = ? AND user_id = ?", bookID, userID)
		if err != nil {
			http.Error(w, "Failed to delete book", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/books", http.StatusSeeOther)
	}
}
