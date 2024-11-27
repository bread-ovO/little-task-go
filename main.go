package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	ID            int
	Title         string
	Author        string
	PublishedYear int
	Status        string
}

var db *sql.DB
var tmpl *template.Template

func main() {
	// 初始化数据库连接
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/LibraryManagement") // 替换为实际用户名和密码
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 加载模板
	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	// 路由设置
	http.HandleFunc("/", homePage)
	http.HandleFunc("/add", addBook)
	http.HandleFunc("/delete", deleteBook)
	http.HandleFunc("/update", updateBook)

	// 启动服务器
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 首页：展示书籍列表
func homePage(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM Books")
	if err != nil {
		http.Error(w, "Failed to query books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear, &book.Status)
		if err != nil {
			http.Error(w, "Failed to parse book data", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// 渲染模板并传递数据
	tmpl.ExecuteTemplate(w, "index.html", books)
}

// 添加书籍
func addBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		author := r.FormValue("author")
		publishedYear := r.FormValue("publishedYear")

		_, err := db.Exec("INSERT INTO Books (Title, Author, PublishedYear) VALUES (?, ?, ?)", title, author, publishedYear)
		if err != nil {
			http.Error(w, "Failed to add book", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// 删除书籍
func deleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")

		_, err := db.Exec("DELETE FROM Books WHERE BookID = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete book", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// 更新书籍状态
func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		status := r.FormValue("status")

		_, err := db.Exec("UPDATE Books SET Status = ? WHERE BookID = ?", status, id)
		if err != nil {
			http.Error(w, "Failed to update book status", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
