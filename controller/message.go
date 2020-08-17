package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/taoshihan1991/imaptool/models"
	"time"
)
// @Summary 发送消息接口
// @Produce  json
// @Accept multipart/form-data
// @Param from_id formData   string true "来源uid"
// @Param toID formData   string true "目标uid"
// @Param content formData   string true "内容"
// @Param type formData   string true "类型|kefu,visitor"
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /message [post]
func SendMessage(c *gin.Context) {
	fromId := c.PostForm("from_id")
	toId := c.PostForm("toID")
	content := c.PostForm("content")
	cType := c.PostForm("type")
	if content==""{
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "内容不能为空",
		})
		return
	}

	var kefuInfo models.User
	var visitorInfo models.Visitor
	if cType=="kefu"{
		kefuInfo=models.FindUser(fromId)
		visitorInfo=models.FindVisitorByVisitorId(toId)
	}else if cType=="visitor"{
		visitorInfo=models.FindVisitorByVisitorId(fromId)
		kefuInfo=models.FindUser(toId)
	}

	if kefuInfo.ID==0 ||visitorInfo.ID==0{
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户不存在",
		})
		return
	}
	models.CreateMessage(kefuInfo.Name,visitorInfo.VisitorId,content,cType)

	if cType=="kefu"{
		guest,ok:=clientList[visitorInfo.VisitorId]
		if guest==nil||!ok{
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "ok",
			})
			return
		}
		conn := guest.conn

		msg := TypeMessage{
			Type: "message",
			Data: ClientMessage{
				Name:  kefuInfo.Nickname,
				Avatar:   kefuInfo.Avatar,
				Id:    kefuInfo.Name,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
				ToId: visitorInfo.VisitorId,
				Content:  content,
			},
		}
		str, _ := json.Marshal(msg)
		conn.WriteMessage(websocket.TextMessage,str)
	}
	if cType=="visitor"{
		kefuConns,ok := kefuList[kefuInfo.Name]
		if kefuConns==nil||!ok{
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "ok",
			})
			return
		}
		msg := TypeMessage{
			Type: "message",
			Data: ClientMessage{
				Avatar: visitorInfo.Avatar,
				Id:     visitorInfo.VisitorId,
				Name:   visitorInfo.Name,
				ToId:       kefuInfo.Name,
				Content:    content,
				Time:        time.Now().Format("2006-01-02 15:04:05"),
			},
		}
		str, _ := json.Marshal(msg)
		for _,kefuConn:=range kefuConns{
			kefuConn.WriteMessage(websocket.TextMessage,str)
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
