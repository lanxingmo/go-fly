package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/taoshihan1991/imaptool/models"
	"github.com/taoshihan1991/imaptool/tmpl"
	"github.com/taoshihan1991/imaptool/tools"
	"net/http"
)

func ActionMain(w http.ResponseWriter, r *http.Request) {
	sessionId := tools.GetCookie(r, "session_id")
	info := AuthCheck(sessionId)
	if len(info) == 0 {
		http.Redirect(w, r, "/login", 302)
		return
	}
	render := tmpl.NewRender(w)
	render.Display("main", render)
}
func MainCheckAuth(c *gin.Context) {
	id,_:=c.Get("kefu_id")
	userinfo:=models.FindUserRole("user.avatar,user.name,user.id, role.name role_name",id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "验证成功",
		"result":gin.H{
			"avatar":userinfo.Avatar,
			"name":userinfo.Name,
			"role_name":userinfo.RoleName,
		},
	})
}
