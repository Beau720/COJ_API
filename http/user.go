package http

import (
	api "COJ_API/service/user"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Grouping all routers links into one through the http.go
func user_group(rg *gin.RouterGroup, u *api.Service) {
	user_group := rg.Group("/")

	//http request to find all users in the database
	user_group.GET("/list", func(c *gin.Context) {
		users, user_error := u.List()
		if user_error != api.UserErrorNil {
			log.Println("List users failed:", user_error, "data[", users, "]")
			c.JSON(http.StatusInternalServerError, gin.H{"error": string(user_error)})
			return
		}

		c.JSON(http.StatusOK, users)
	})

	//http request to find user by ID
	user_group.GET("/userid/:id", func(c *gin.Context) {
		id_param := c.Param("id")
		id, err := strconv.Atoi(id_param)
		if err != nil {
			log.Println("Select user, invalid id")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
			return
		}

		//switch case error handling if user is not registered
		user, user_error := u.SelectById(id)
		if user_error != api.UserErrorNil {
			log.Println("Select user id failed:", user_error, "data[", user, "]")
			switch user_error {
			case api.UserErrorUserNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": string(api.UserErrorUserNotFound)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(user_error)})
				return
			}
		}

		c.JSON(http.StatusOK, user)
	})

	//http request to find user by email
	user_group.GET("/email/:email", func(c *gin.Context) {
		email_param := c.Param("email")

		//switch case error handling if user is not registered
		user, user_error := u.SelectByEmail(email_param)
		if user_error != api.UserErrorNil {
			log.Println("Select user email failed:", user_error, "data[", user, "]")
			switch user_error {
			case api.UserErrorUserNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": string(api.UserErrorUserNotFound)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(user_error)})
				return
			}
		}

		c.JSON(http.StatusOK, user)
	})

	//http request to find user by RSA-ID
	user_group.GET("/rsaid/:said", func(c *gin.Context) {
		rsaid_param := c.Param("said")

		//switch case error handling if user is not registered
		user, user_error := u.SelectByRSAID(rsaid_param)
		if user_error != api.UserErrorNil {
			log.Println("Select user rsais failed:", user_error, "data[", user, "]")
			switch user_error {
			case api.UserErrorUserNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": string(api.UserErrorUserNotFound)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(user_error)})
				return
			}
		}

		c.JSON(http.StatusOK, user)
	})

	//http request to create user
	user_group.POST("/create", func(c *gin.Context) {
		user := api.User{}
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Println("Create user failed to bind to JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data sent"})
			return
		}

		//switch case error handling for user who already registered
		user_resp, user_error := u.Create(&user)
		if user_error != api.UserErrorNil {
			log.Println("Create user failed:", user_error, "data[", user, "]")
			switch user_error {
			case api.UserErrorUserEmailExists:
				c.JSON(http.StatusConflict, gin.H{"error": string(api.UserErrorUserEmailExists)})
				return
			case api.UserErrorUserRSAIDExists:
				c.JSON(http.StatusConflict, gin.H{"error": string(api.UserErrorUserRSAIDExists)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(user_error)})
				return
			}
		}

		c.JSON(http.StatusCreated, user_resp)
	})

	//http request to update user by id
	user_group.PUT("/update/:id", func(c *gin.Context) {
		id_param := c.Param("id")
		id, err := strconv.Atoi(id_param)

		//error handling if user doesn't exist
		if err != nil {
			log.Println("Select user, invalid id")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
			return
		}
		//error handling if user failed to be updated
		user := api.User{}
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Println("Update user failed to bind to JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data sent"})
			return
		}

		user.ID = id

		user_resp, user_error := u.Update(&user)
		if user_error != api.UserErrorNil {
			c.JSON(http.StatusConflict, gin.H{"error": string(user_error)})
			return
		}

		c.JSON(http.StatusOK, user_resp)
	})
}
