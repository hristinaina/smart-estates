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
	"smarthome-back/services"
	"strings"
)

var service = services.NewConfigService()

var (
	awsRegion             = "eu-central-1"
	awsAccessKeyID, _     = service.GetAccessKey("config/config.json")
	awsSecretAccessKey, _ = service.GetSecretAccessKey("config/config.json")
	// TODO : replace this after A&A implementation
	s3Bucket = "examplegmail.com"
	username = "examplegmail.com"
)

type ImageUploadController struct {
}

func NewImageUploadController() ImageUploadController {
	return ImageUploadController{}
}

func (iup ImageUploadController) GetImageURL(c *gin.Context) {
	// TODO : change this later
	filename := c.Param("file-name")

	fullName, err := findFullFileName(filename)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while searching for file"})
		return
	}

	// Replace 'your-s3-bucket-name' with your actual S3 bucket name
	s3URL := fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", awsRegion, s3Bucket, fullName)
	fmt.Println("RETURNNN")
	fmt.Println(s3URL)
	c.JSON(http.StatusOK, gin.H{"imageUrl": s3URL})
}

func (iup ImageUploadController) UploadImage(c *gin.Context) {
	name := c.Param("real-estate-name")
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fileBytes, err := iup.readFile(file)
	if err != nil {
		fmt.Println("Error: ", "Failed to read file")
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// check if bucket exists
	folderExists, err := doesFolderExist(username)
	if err != nil {
		fmt.Println("Error: ", "Failed to check folder")
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user's folder"})
		return
	}

	if !folderExists {
		// create new bucket
		err := createFolder(username)
		if err != nil {
			fmt.Println("Error: ", "Failed to create user's folder")
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user's folder"})
			return
		}
	}

	err = uploadToS3(fileBytes, name)
	if err != nil {
		fmt.Println("Error: ", "Failed to upload file")
		fmt.Println(err)
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

func findFullFileName(fileName string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return "", err
	}
	svc := s3.New(sess)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(username),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		return "", err
	}

	// iterate through objects
	fmt.Println("Objects in the bucket:")
	for _, item := range result.Contents {
		lastDotIndex := strings.LastIndex(*item.Key, ".")
		comparation := *item.Key
		comparation = comparation[:lastDotIndex]
		if comparation == username+"/"+fileName {
			fmt.Println("FOUND!")
			return *item.Key, nil
		}
		fmt.Printf("Name: %s, Size: %d\n", *item.Key, *item.Size)
	}

	return "", nil
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
		Key:         aws.String(username + "/" + fileName),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(http.DetectContentType(fileBytes)),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %v", err)
	}

	return nil
}

func doesFolderExist(folderName string) (bool, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return false, err
	}

	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(s3Bucket),
		Prefix:    aws.String(folderName + "/"),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		return false, err
	}

	return len(resp.Contents) > 0, nil
}

func createFolder(folderName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(folderName + "/"), // Note the trailing "/"
		Body:   strings.NewReader(""),        // Empty content for a "folder"
	})

	if err != nil {
		return err
	}

	return nil
}
