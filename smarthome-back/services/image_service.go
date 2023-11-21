package services

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"net/http"
	"strings"
)

var service = NewConfigService()

var (
	awsRegion             = "eu-central-1"
	awsAccessKeyID, _     = service.GetAccessKey("config/config.json")
	awsSecretAccessKey, _ = service.GetSecretAccessKey("config/config.json")
	s3Bucket              = "examplegmail.com"
	// TODO : replace this after A&A implementation
	username = "examplegmail.com"
)

type ImageService interface {
	GetImageURL(fileName string) (string, error)
	UploadImage(estateName string, file *multipart.FileHeader) error
	readFile(file *multipart.FileHeader) ([]byte, error)
	findFullFileName(fileName string) (string, error)
	uploadToS3(fileBytes []byte, fileName string) error
	doesFolderExist(folderName string) (bool, error)
	createFolder(folderName string) error
}

type ImageServiceImpl struct {
}

func NewImageService() ImageService {
	return &ImageServiceImpl{}
}

func (is *ImageServiceImpl) GetImageURL(fileName string) (string, error) {
	fullName, err := is.findFullFileName(fileName)
	if err != nil {
		return "", err
	}

	s3URL := fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", awsRegion, s3Bucket, fullName)
	return s3URL, nil
}

func (is *ImageServiceImpl) UploadImage(estateName string, file *multipart.FileHeader) error {
	fileBytes, err := is.readFile(file)
	if err != nil {
		fmt.Println("Error: ", "Failed to read file")
		fmt.Println(err)
		return err
	}

	// check if folder for logged user exists
	folderExists, err := is.doesFolderExist(username)
	if err != nil {
		return err
	}

	if !folderExists {
		// create new folder
		err := is.createFolder(username)
		if err != nil {
			return err
		}
	}

	err = is.uploadToS3(fileBytes, estateName)
	if err != nil {
		fmt.Println("Error: ", "Failed to upload file")
		fmt.Println(err)
		return err
	}

	return nil
}

func (is *ImageServiceImpl) readFile(file *multipart.FileHeader) ([]byte, error) {
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

func (is *ImageServiceImpl) findFullFileName(fileName string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		return "", err
	}
	svc := s3.New(sess)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s3Bucket),
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

func (is *ImageServiceImpl) uploadToS3(fileBytes []byte, fileName string) error {
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

func (is *ImageServiceImpl) doesFolderExist(folderName string) (bool, error) {
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

func (is *ImageServiceImpl) createFolder(folderName string) error {
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
