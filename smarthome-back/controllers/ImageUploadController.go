package controllers

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
)

var (
	awsRegion          = "eu-central-1"
	awsAccessKeyID     = "AKIAXTEDOKGSGESVDNWJ"
	awsSecretAccessKey = "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"
	s3Bucket           = "images-bucket-sh"
)

type ImageUploadController struct {
}

func NewImageUploadController() ImageUploadController {
	return ImageUploadController{}
}

func (iup ImageUploadController) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fileBytes, err := iup.readFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	err = uploadToS3(fileBytes, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}

func (iup ImageUploadController) readFile(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(src)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func uploadToS3(fileBytes []byte, fileName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	uploader := s3.New(sess)

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s3Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(http.DetectContentType(fileBytes)),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %v", err)
	}

	return nil
}
