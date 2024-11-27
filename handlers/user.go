package handlers

import (
	"html/template"
	"little-task-go/db"
	"net/http"
)

// User 页面显示用户信息，并提供修改功能
func User(w http.ResponseWriter, r *http.Request) {
	// 确保用户已登录（伪代码，需根据具体会话管理实现）
	userID := getSessionUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		// 获取提交的表单数据
		newNickname := r.FormValue("nickname")
		newGender := r.FormValue("gender")
		newPassword := r.FormValue("password")

		// 更新用户信息
		_, err := db.DB.Exec("UPDATE users SET nickname = ?, gender = ?, password = ? WHERE id = ?", newNickname, newGender, newPassword, userID)
		if err != nil {
			http.Error(w, "Failed to update user information", http.StatusInternalServerError)
			return
		}

		// 更新成功后重定向回用户页面
		http.Redirect(w, r, "/user", http.StatusFound)
		return
	}

	// 查询当前用户信息
	var nickname, gender string
	err := db.DB.QueryRow("SELECT nickname, gender FROM users WHERE id = ?", userID).Scan(&nickname, &gender)
	if err != nil {
		http.Error(w, "Failed to load user information", http.StatusInternalServerError)
		return
	}

	// 渲染模板
	data := struct {
		Nickname string
		Gender   string
	}{
		Nickname: nickname,
		Gender:   gender,
	}

	tmpl := template.Must(template.ParseFiles("templates/user.html"))
	tmpl.Execute(w, data)
}
