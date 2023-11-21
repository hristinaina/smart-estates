package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smarthome-back/services"
)

type ImageController struct {
	service services.ImageService
}

func NewImageController() ImageController {
	return ImageController{service: services.NewImageService()}
}

func (ic ImageController) Get(c *gin.Context) {
	fileName := c.Param("file-name")
	imageURL, err := ic.service.GetImageURL(fileName)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while searching for file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"imageUrl": imageURL})
}

func (ic ImageController) Post(c *gin.Context) {
	name := c.Param("real-estate-name")
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = ic.service.UploadImage(name, file)
	if err != nil {
		c.JSON(400, "Error while uploading image")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}
