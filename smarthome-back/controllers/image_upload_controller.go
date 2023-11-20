package controllers

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
)

type ImageUploadController struct {
}

func NewImageUploadController() ImageUploadController {
	return ImageUploadController{}
}

func (iup ImageUploadController) GetImageURL(c *gin.Context) {
	// TODO : change this later
	username := "examplegmail.com"
	filename := c.Param("file-name")

	fullName, err := findFullFileName(filename)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while searching for file"})
		return
	}

	// Replace 'your-s3-bucket-name' with your actual S3 bucket name
	s3URL := fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", awsRegion, username, fullName)
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
	// TODO : replace this with actual email
	userBucketExists, err := doesBucketExist(s3Bucket)
	if err != nil {
		fmt.Println("Error: ", "Failed to check user's bucket")
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user's bucket"})
		return
	}

	if !userBucketExists {
		// TODO : replace with actual email
		// create new bucket
		err := createBucket(s3Bucket)
		if err != nil {
			fmt.Println("Error: ", "Failed to create user's bucket")
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user's bucket"})
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
	// TODO : change this later
	username := "examplegmail.com"
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
		if strings.Split(*item.Key, ".")[0] == fileName {
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
		Bucket: aws.String(s3Bucket),
		// TODO : change key for real estate name
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(http.DetectContentType(fileBytes)),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %v", err)
	}

	return nil
}

func doesBucketExist(bucketName string) (bool, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return false, err
	}

	svc := s3.New(sess)

	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		fmt.Println("Error checking bucket:", err)
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func createBucket(bucketName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	return nil
}
