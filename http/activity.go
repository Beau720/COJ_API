package http

import (
	api "COJ_API/service/activity"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Grouping all routers links into one through the http.go
func activity_group(rg *gin.RouterGroup, a *api.Service) {
	activity_group := rg.Group("/")

	//http request to find all users in the database
	activity_group.GET("/list", func(c *gin.Context) {
		activities, activity_error := a.List()
		if activity_error != api.ActivityErrorNil {
			log.Println("List activities failed:", activity_error, "data[", activities, "]")
			c.JSON(http.StatusInternalServerError, gin.H{"error": string(activity_error)})
			return
		}

		c.JSON(http.StatusOK, activities)
	})

	//http request to find user by ID
	activity_group.GET("/activityid/:id", func(c *gin.Context) {
		id_param := c.Param("id")
		id, err := strconv.Atoi(id_param)
		if err != nil {
			log.Println("Select activity, invalid id")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
			return
		}

		//switch case error handling if activity is not registered
		activity, activity_error := a.SelectByActivityId(id)
		if activity_error != api.ActivityErrorNil {
			log.Println("Select activity id failed:", activity_error, "data[", activity, "]")
			switch activity_error {
			case api.ActivityErrorActivityNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": string(api.ActivityErrorActivityNotFound)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(activity_error)})
				return
			}
		}

		c.JSON(http.StatusOK, activity)
	})

	//http request to find user by email
	activity_group.GET("/activitytype/:activitytype ", func(c *gin.Context) {
		activitytype_param := c.Param("activitytype")

		//switch case error handling if user is not registered
		activity, activity_error := a.SelectByactivitytype(activitytype_param)
		if activity_error != api.ActivityErrorNil {
			log.Println("Select activity type failed:", activity_error, "data[", activity, "]")
			switch activity_error {
			case api.ActivityErrorActivityNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": string(api.ActivityErrorActivityNotFound)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(activity_error)})
				return
			}
		}

		c.JSON(http.StatusOK, activity)
	})

}
