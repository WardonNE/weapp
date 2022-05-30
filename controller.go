package weapp

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	*Component `inject:"component"`
}

func (c *Controller) Init() {
}

func (c *Controller) View(ctx *gin.Context, template string, data any) {
	ctx.HTML(http.StatusOK, template, data)
}

func (c *Controller) Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "ok",
		"data":    data,
	})
}

func (c *Controller) BadRequest(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"status":  "error",
		"message": "bad request",
		"data":    data,
	})
}

func (c *Controller) Unauthorized(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"status":  "error",
		"message": "unauthorized",
		"data":    data,
	})
}

func (c *Controller) Forbidden(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusForbidden, gin.H{})
}

func (c *Controller) NotFound(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"status":  "error",
		"message": "not found",
		"data":    data,
	})
}

func (c *Controller) InvalidParams(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"code":    http.StatusUnprocessableEntity,
		"status":  "error",
		"message": "invalid params",
		"data":    data,
	})
}

func (c *Controller) Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"status":  "fail",
		"message": "server error",
		"data": map[string]string{
			"error": err.Error(),
		},
	})
}

func (c *Controller) Attachment(ctx *gin.Context, file *os.File, filename string, contentType string) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		c.Error(ctx, err)
	} else {
		ctx.Data(http.StatusOK, contentType, content)
	}
}

func (c *Controller) Download(ctx *gin.Context, filepath string, filename string) {
	ctx.FileAttachment(filepath, filename)
}
