package main

import (
	"fmt"
	"net/http"
	"task4/controllers"
	"task4/database"
	"task4/models"
	"text/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDB()
	defer database.DB.Close()

	// In main.go or database.InitDB()
	database.DB.AutoMigrate(&models.User{}, &models.Book{}, &models.BorrowingRecord{}, &models.ReservationRecord{}) // 添加 ReservationRecord

	// 创建 Gin 实例
	router := gin.Default()

	// 注册自定义模板函数
	router.SetFuncMap(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"add": func(a, b int) int { return a + b },
	})

	// 设置模板路径
	router.LoadHTMLGlob("templates/*")

	// 设置静态文件路径
	router.Static("/static", "./static")

	// 设置会话
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// 根页面重定向到登录页面
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// 未登录页面分组：无需验证用户身份
	publicGroup := router.Group("/")
	{
		publicGroup.GET("/login", controllers.ShowLoginPage)
		publicGroup.POST("/login", controllers.Login)
		publicGroup.GET("/register", func(c *gin.Context) {
			c.HTML(http.StatusOK, "register.html", nil)
		})
		publicGroup.POST("/register", controllers.Register)
	}

	// 登录状态中间件
	authMiddleware := func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 从数据库获取用户信息
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", user.ID)
		c.Set("nickname", user.Nickname)
		c.Set("gender", user.Gender)
		fmt.Printf("User authenticated: %d - %s\n", user.ID, user.Nickname)
		c.Next()
	}

	// 登录后可访问的页面分组：需要验证用户身份
	authGroup := router.Group("/")
	authGroup.Use(authMiddleware)
	{
		// 书籍管理
		authGroup.GET("/books", controllers.ShowBooksPage)
		authGroup.GET("/books/add", controllers.ShowAddBookPage)
		authGroup.POST("/books/add", controllers.AddBook)
		authGroup.GET("/books/update/:id", controllers.ShowUpdateBookPage)
		authGroup.POST("/books/update/:id", controllers.UpdateBookHandler)
		authGroup.GET("/books/delete/:id", controllers.DeleteBookHandler)

		// 借阅与归还 (新添加)
		authGroup.GET("/books/borrow/:id", controllers.BorrowBook)   // 使用 GET 简化示例，POST 更佳
		authGroup.GET("/books/return/:id", controllers.ReturnBook)   // 使用 GET 简化示例，POST 更佳
		authGroup.GET("/history", controllers.ShowBorrowingHistory)  // (可选)
		authGroup.GET("/books/reserve/:id", controllers.ReserveBook) // 新增
		authGroup.GET("/books/renew/:id", controllers.RenewBook)     // 新增
		authGroup.GET("/myborrows", controllers.ShowMyBorrowsPage)   // 新增

		// 用户管理
		authGroup.GET("/user/edit", controllers.ShowUserEditPage)
		authGroup.POST("/user/edit", controllers.UpdateUser)

		// 登出
		authGroup.GET("/logout", controllers.Logout)
	}

	// 启动服务器
	router.Run(":8080")
}
