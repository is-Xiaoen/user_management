package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"user-management-system/errors"
	"user-management-system/models"
	"user-management-system/session"
)

// funcMap 定义模板函数
var funcMap = template.FuncMap{
	"upper": strings.ToUpper,
}

// RenderHomePage 渲染主页
func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	//获取当前用户
	currentUser, _ := session.GetCurrentUser(r)

	data := struct {
		CurrentUser *models.User
	}{
		CurrentUser: currentUser,
	}

	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles("views/layout", "views/index.html")

	if err != nil {
		log.Printf("模板解析错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("模板执行错误: %v", err)
		errors.HandleError(w, r, errors.NewInternalError(err))
	}

}
