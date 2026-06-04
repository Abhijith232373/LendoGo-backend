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

// UploadFileToS3 securely sends a file from Fiber to your AWS bucket using AWS SDK V2
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

// GeneratePresignedURL creates a 15-minute secure link using AWS SDK V2
func GeneratePresignedURL(rawURL string) string {
    if rawURL == "" {
        return ""
    }

    // 1. Extract the exact file path (Object Key) from your full URL
    parts := strings.SplitN(rawURL, ".amazonaws.com/", 2)
    if len(parts) != 2 {
        return rawURL // If it's not an AWS URL, return it normally
    }
    objectKey := parts[1]

    bucket := os.Getenv("AWS_BUCKET_NAME")
    region := os.Getenv("AWS_REGION")

    // 2. Load the AWS V2 Configuration
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        return rawURL
    }

    // 3. Create the S3 client and the V2 Presign Client
    client := s3.NewFromConfig(cfg)
    presignClient := s3.NewPresignClient(client)

    // 4. Generate the 15-minute key
    req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(objectKey),
    }, func(opts *s3.PresignOptions) {
        opts.Expires = 15 * time.Minute
    })

    if err != nil {
        return rawURL
    }

    return req.URL
}