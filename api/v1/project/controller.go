package project

import (
	"fmt"
	"zoho-center/core/response"
	"zoho-center/service"

	"github.com/gin-gonic/gin"
)

// @Summary 项目列表
// @Id 5
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param name query string false "项目编码"
// @Success 200 object response.ListRes{data=[]Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects [GET]
func GetProjectList(c *gin.Context) {
	var filter ProjectFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	projectService := NewProjectService()
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	count, list, err := projectService.GetProjectList(filter, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageId, filter.PageSize, count, list)
}

// @Summary 新建项目
// @Id 6
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param project_info body ProjectNew true "项目信息"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects [POST]
func NewProject(c *gin.Context) {
	var project ProjectNew
	if err := c.ShouldBindJSON(&project); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	project.User = claims.Username
	organizationID := claims.OrganizationID
	projectService := NewProjectService()
	new, err := projectService.NewProject(project, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取项目
// @Id 7
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目ID"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id [GET]
func GetProjectByID(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	organizationID := claims.OrganizationID
	projectService := NewProjectService()
	project, err := projectService.GetProjectByID(uri.ID, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, project)

}

// @Summary 根据ID更新项目
// @Id 8
// @Tags 项目管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "项目ID"
// @Param project_info body ProjectNew true "项目信息"
// @Success 200 object response.SuccessRes{data=Project} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /projects/:id [PUT]
func UpdateProject(c *gin.Context) {
	var uri ProjectID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var project ProjectNew
	if err := c.ShouldBindJSON(&project); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	fmt.Println(claims.Username)
	project.User = claims.Username
	organizationID := claims.OrganizationID
	projectService := NewProjectService()
	new, err := projectService.UpdateProject(uri.ID, project, organizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}
