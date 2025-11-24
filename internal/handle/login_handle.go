package handle

import (
	"net/http"

	"github.com/l-jessie/test-im/internal/utils"

	"github.com/gin-gonic/gin"
)

func LoginHandle(c *gin.Context) {
	var body map[string]string
	err := c.ShouldBindBodyWithJSON(&body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}

	username := body["username"]
	if username == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "用户名不能为空",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "登录成功",
		"data": gin.H{
			"id":       utils.GenerateUUID(),
			"username": username,
		},
	})
}
