package http

import (
	api "COJ_API/service/form"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Grouping all routers links into one through the http.go
func form_group(rg *gin.RouterGroup, f *api.Service) {
	form_group := rg.Group("/")

	//http request to find all forms in the database
	form_group.GET("/list", func(c *gin.Context) {
		forms, form_error := f.List()
		if form_error != api.FormErrorNil {
			log.Println("List forms failed:", form_error, "data[", forms, "]")
			c.JSON(http.StatusInternalServerError, gin.H{"error": string(form_error)})
			return
		}

		c.JSON(http.StatusOK, forms)
	})

	//Get method for finding user by ID
	form_group.GET("/formrefNo/:refNo", func(c *gin.Context) {
		refNo_param := c.Param("refNo")
		refNo, err := strconv.Atoi(refNo_param)
		if err != nil {
			log.Println("Select form, invalid refNo")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid refNo"})
			return
		}

		//switch case error handling if form is not registered
		form, form_error := f.SelectByRefNo(refNo)
		if form_error != api.FormErrorNil {
			log.Println("Select form id failed:", form_error, "data[", form, "]")
			switch form_error {
			case api.FormErrorFormNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": string(api.FormErrorFormNotFound)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(form_error)})
				return
			}
		}

		c.JSON(http.StatusOK, form)
	})

	//Http get method for finding form by inspector officer
	form_group.GET("/inspector/:inspector_name", func(c *gin.Context) {
		inspector_name_param := c.Param("inspector_officer")

		//switch case error handling if form is not registered
		form, form_error := f.FindByInspectorName(inspector_name_param)
		if form_error != api.FormErrorNil {
			log.Println("Select form name failed:", form_error, "data[", form, "]")
			switch form_error {
			case api.FormErrorFormNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": string(api.FormErrorFormNotFound)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(form_error)})
				return
			}
		}

		c.JSON(http.StatusOK, form)
	})

	//Http get method for finding form by date
	form_group.GET("/date/:arrival", func(c *gin.Context) {
		arrival_param := c.Param("dateOfArrival")

		//switch case error handling if form is not registered
		form, form_error := f.FindbydateofArrival(arrival_param)
		if form_error != api.FormErrorNil {
			log.Println("Select form name failed:", form_error, "data[", form, "]")
			switch form_error {
			case api.FormErrorFormNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": string(api.FormErrorFormNotFound)})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": string(form_error)})
				return
			}
		}

		c.JSON(http.StatusOK, form)
	})

	//Http Post method for finding form by name
	form_group.POST("/create", func(c *gin.Context) {
		form := api.Form{}
		if err := c.ShouldBindJSON(&form); err != nil {
			log.Println("Create form failed to bind to JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data sent"})
			return
		}

		form_resp, form_error := f.Create(&form)
		if form_error != api.FormErrorNil {
			log.Println("Create user failed:", form_error, "data[", form, "]")
			c.JSON(http.StatusInternalServerError, gin.H{"error": string(form_error)})
			return
		}

		c.JSON(http.StatusCreated, form_resp)
	})

	//Http Post method for uodating form by id
	form_group.PUT("/update/:refNo", func(c *gin.Context) {
		refNo_param := c.Param("refNo")
		refNo, err := strconv.Atoi(refNo_param)
		if err != nil {
			log.Println("Select form, invalid refNo")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid refNo"})
			return
		}

		form := api.Form{}
		if err := c.ShouldBindJSON(&form); err != nil {
			log.Println("Update form failed to bind to JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data sent"})
			return
		}

		form.ReferenceNumber = refNo

		form_resp, form_error := f.Update(&form)
		if form_error != api.FormErrorNil {
			c.JSON(http.StatusConflict, gin.H{"error": string(form_error)})
			return
		}

		c.JSON(http.StatusOK, form_resp)
	})
}
