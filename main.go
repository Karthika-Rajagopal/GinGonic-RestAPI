package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// create a new Gin router
var router *gin.Engine

func main() {
	// set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// create a new Gin router
	router = gin.Default()

	// serve static files from the "static" directory
	router.Static("/static", "./static")

	// setup routes
	setupRoutes()

	// start the server
	router.Run(":8080")
}

func setupRoutes() {
	// render the main page
	router.GET("/", func(c *gin.Context) {
		render(c, "index.html", gin.H{})
	})

	// handle file uploads
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}

		// save the uploaded file to disk
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "./static/images/"+filename); err != nil {
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		c.String(http.StatusOK, "File uploaded successfully")
	})
}

// render renders a template with the given name and data
func render(c *gin.Context, name string, data gin.H) {
	// get the absolute path to the template file
	tmplPath := filepath.Join("templates", name)

	// parse the template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	// execute the template with the given data
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}
}
