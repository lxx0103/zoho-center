package project

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.GET("/projects", GetProjectList)
	g.GET("/projects/:id", GetProjectByID)
	g.PUT("/projects/:id", UpdateProject)
	g.POST("/projects", NewProject)
}
