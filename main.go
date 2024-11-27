package main

import (
	"little-task-go/db"       // 数据库初始化模块
	"little-task-go/handlers" // 处理器模块
	"log"
	"net/http"
)

func main() {
	// 初始化数据库
	db.Init()
	defer db.DB.Close()

	// 路由配置
	http.HandleFunc("/login", handlers.Login)       // 登录页面
	http.HandleFunc("/books", handlers.Books)       // 书籍管理页面
	http.HandleFunc("/user", handlers.User)         // 用户信息管理页面
	http.HandleFunc("/add", handlers.AddBook)       // 添加书籍
	http.HandleFunc("/delete", handlers.DeleteBook) // 删除书籍
	http.HandleFunc("/update", handlers.UpdateBook) // 更新书籍信息

	// 静态文件处理
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 启动服务器
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
