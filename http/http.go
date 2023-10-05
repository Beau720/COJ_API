package http

import (
	"COJ_API/service/activity"
	"COJ_API/service/form"
	"COJ_API/service/login"
	"COJ_API/service/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host            string
	Port            string
	UserService     *user.Service
	FormService     *form.Service
	LoginService    *login.Service
	ActivityService *activity.Service
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func Start(c *Config) {
	router := gin.Default()

	a := router.Group("activity")
	activity_group(a, c.ActivityService)

	l := router.Group("/login")
	login_group(l, c.LoginService)

	u := router.Group("/user")
	user_group(u, c.UserService)

	f := router.Group("/form")
	form_group(f, c.FormService)

	router.Run(c.Addr())
}
