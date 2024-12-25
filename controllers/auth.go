package controllers

import (
	"fmt"
	"log"
	"net/http"
	"task4/database"
	"task4/models"
	"task4/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	// 获取表单数据
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")
	gender := c.PostForm("gender")

	// 验证密码是否一致
	if password != confirmPassword {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "两次密码输入不一致",
		})
		return
	}

	// 验证性别字段是否有效
	validGenders := map[string]bool{"Male": true, "Female": true, "Other": true}
	if !validGenders[gender] {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "性别字段值无效",
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "用户名已被注册",
		})
		return
	}

	// 对密码进行加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "密码加密失败，请重试",
		})
		return
	}

	// 创建新用户
	newUser := models.User{
		Username: username,
		Password: hashedPassword,
		Gender:   gender,
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "注册失败，请重试",
		})
		return
	}

	// 注册成功，显示成功信息
	c.HTML(http.StatusOK, "login.html", gin.H{
		"success": "注册成功，请登录",
	})
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := models.GetUserByUsername(username)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 设置会话
	session := sessions.Default(c)
	session.Set("user_id", user.ID) // 将用户 ID 保存到会话中
	session.Save()                  // 保存会话到浏览器

	c.Redirect(http.StatusFound, "/books")
	fmt.Println("登录成功，重定向到 /books")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear() // 清除所有会话数据
	session.Save()  // 保存更改
	c.Redirect(http.StatusFound, "/login")
}

func ShowUserEditPage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user, err := models.GetUserByID(userID.(uint))
	if err != nil {
		c.String(http.StatusInternalServerError, "用户不存在")
		return
	}

	c.HTML(http.StatusOK, "user_edit.html", gin.H{
		"user": user,
	})
}

func UpdateUser(c *gin.Context) {
	// 检查用户是否已登录
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("Error fetching user: %v", err)
		c.String(http.StatusInternalServerError, "用户不存在")
		return
	}

	// 前端到后端的性别映射
	genderMap := map[string]string{
		"男":  "Male",
		"女":  "Female",
		"其他": "Other",
	}

	// 获取前端提交的表单数据
	submittedGender := c.PostForm("gender")
	gender, exists := genderMap[submittedGender]
	if !exists {
		log.Printf("Invalid gender value: %s", submittedGender)
		c.String(http.StatusBadRequest, "无效的性别值")
		return
	}

	nickname := c.PostForm("nickname")
	if len(nickname) == 0 {
		c.String(http.StatusBadRequest, "昵称不能为空")
		return
	}
	if len(nickname) < 3 || len(nickname) > 20 {
		c.String(http.StatusBadRequest, "昵称长度必须在3到20个字符之间")
		return
	}

	password := c.PostForm("password")
	if len(password) > 0 {
		if len(password) < 6 {
			c.String(http.StatusBadRequest, "密码长度必须至少为6位")
			return
		}
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			c.HTML(http.StatusOK, "register.html", gin.H{
				"error": "密码加密失败，请重试",
			})
			return
		}
		user.Password = hashedPassword
	}

	// 更新字段
	user.Gender = gender
	user.Nickname = nickname

	// 保存更新到数据库
	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("Error saving user: %v", err)
		c.String(http.StatusInternalServerError, "更新用户信息失败")
		return
	}

	// 更新成功，重定向到 /books 页面
	c.Redirect(http.StatusFound, "/books")
}
