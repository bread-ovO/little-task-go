package main

import (
	"little-task-go/db"       // 数据库初始化模块
	"little-task-go/handlers" // 处理器模块
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
	http.HandleFunc("/add_user", handlers.AddUser)  // 增加用户功能

	// 启动服务器
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 配置其他路由和服务
	http.ListenAndServe(":8080", nil)
}
