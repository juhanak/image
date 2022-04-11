package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/juhanak/image/libs/imageProcessor"
	"mime/multipart"
	"net/http"
)

type ImageParams struct {
	MaxWidth  int    `form:"maxWidth" binding:"gte=0,lte=2560"`
	MaxHeight int    `form:"maxHeight" binding:"gte=0,lte=1600"`
	File      string `form:"file" binding:"required,endswith=.jpg"`
}

func abortWithErrorCode(c *gin.Context, code int, description string) {
	c.AbortWithStatusJSON(code,
		gin.H{
			"error": description,
		})
}

func GetImage(c *gin.Context) {

	var imageParams ImageParams

	if c.ShouldBindQuery(&imageParams) != nil {
		abortWithErrorCode(c, http.StatusBadRequest, "unexpected parameters")
		return
	}

	if dstImage, err := imageProcessor.GetDefault().Resize(imageParams.File, imageParams.MaxWidth, imageParams.MaxHeight); err != nil {
		abortWithErrorCode(c, http.StatusBadRequest, "failed to resize image")
	} else {
		c.File(dstImage)
	}

}

func Post(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		abortWithErrorCode(c, http.StatusBadRequest, "failed to save file")
		return
	}

	var fileHandle multipart.File
	if fileHandle, err = file.Open(); err != nil {
		abortWithErrorCode(c, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer fileHandle.Close()

	if err = imageProcessor.GetDefault().Validate(fileHandle); err != nil {
		abortWithErrorCode(c, http.StatusBadRequest, "failed to save file")
		return
	}

	nameWithPath, name := imageProcessor.GetDefault().NewImageName()

	if err = c.SaveUploadedFile(file, nameWithPath); err != nil {
		abortWithErrorCode(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	c.JSON(http.StatusOK, gin.H{"original": name})
}
