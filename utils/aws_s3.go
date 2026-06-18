package utils

import (
    "context"
    "fmt"
    "mime/multipart"
    "os"
    "path/filepath"
    "strings"
    "time"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/google/uuid"
)

func UploadFileToS3(fileHeader *multipart.FileHeader) (string, error) {
    file, err := fileHeader.Open()
    if err != nil {
        return "", fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()
    extension := filepath.Ext(fileHeader.Filename)
    uniqueFileName := fmt.Sprintf("kyc_docs/%s%s", uuid.New().String(), extension)
    bucket := os.Getenv("AWS_BUCKET_NAME")
    region := os.Getenv("AWS_REGION")
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        return "", fmt.Errorf("failed to load AWS config: %v", err)
    }
    client := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(client)
    result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(uniqueFileName),
        Body:   file,
    })
    if err != nil {
        return "", fmt.Errorf("failed to upload to S3: %v", err)
    }
    return result.Location, nil
}
func GeneratePresignedURL(rawURL string) string {
    if rawURL == "" {
        return ""
    }
    parts := strings.SplitN(rawURL, ".amazonaws.com/", 2)
    if len(parts) != 2 {
        return rawURL 
    }
    objectKey := parts[1]
    bucket := os.Getenv("AWS_BUCKET_NAME")
    region := os.Getenv("AWS_REGION")
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        return rawURL
    }
    client := s3.NewFromConfig(cfg)
    presignClient := s3.NewPresignClient(client)
    req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(objectKey),
    }, func(opts *s3.PresignOptions) {
        opts.Expires = 7 * 24 * time.Hour
    })
    if err != nil {
        return rawURL
    }
    return req.URL
}