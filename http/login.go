package http

import (
	logi "COJ_API/service/login"
	api "COJ_API/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

var cred *api.User

func login_group(rg *gin.RouterGroup, l *logi.Service) {

	login_group := rg.Group("/")

	/*var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}*/

	login_group.POST("/signIn", func(c *gin.Context) {
		username_param := c.Param("username")

		/*if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}*/

		input, login_error := l.GetUserByUsername(username_param)
		if login_error != logi.LoginErrorNil || !l.ValidatePassword(cred.Password, input.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": string(logi.LoginErrorCredintialsNotFound)})
			return
		}
		c.JSON(http.StatusOK, input)
	})

}
