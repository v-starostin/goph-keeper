package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	chunkSize = 5 * 1024 * 1024 // 5MB per chunk (minimum size for multipart)
)

func main() {
	bucket := "your-s3-bucket"
	key := "your/large-file-key"
	filePath := "/path/to/your/large-file"

	err := uploadFileInChunks(bucket, key, filePath)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
}

// uploadFileInChunks uploads a large file to S3 in chunks using multipart upload
func uploadFileInChunks(bucket, key, filePath string) error {
	// Load AWS credentials and S3 config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Initiate the multipart upload
	createResp, err := s3Client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to initiate multipart upload: %w", err)
	}
	uploadID := createResp.UploadId
	fmt.Printf("Multipart upload initiated. Upload ID: %s\n", *uploadID)

	var completedParts []types.CompletedPart
	partNumber := 1

	// Read and upload each chunk
	for {
		partBuffer := make([]byte, chunkSize)
		bytesRead, err := file.Read(partBuffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file: %w", err)
		}
		if bytesRead == 0 {
			break // End of file
		}

		partBuffer = partBuffer[:bytesRead]

		// Upload the part
		uploadResp, err := s3Client.UploadPart(context.TODO(), &s3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(key),
			UploadId:   uploadID,
			PartNumber: int32(partNumber),
			Body:       bytes.NewReader(partBuffer),
		})
		if err != nil {
			// If upload part fails, abort the multipart upload
			_, abortErr := s3Client.AbortMultipartUpload(context.TODO(), &s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucket),
				Key:      aws.String(key),
				UploadId: uploadID,
			})
			if abortErr != nil {
				log.Printf("Failed to abort multipart upload: %v", abortErr)
			}
			return fmt.Errorf("failed to upload part: %w", err)
		}

		fmt.Printf("Uploaded part %d, ETag: %s\n", partNumber, *uploadResp.ETag)
		completedParts = append(completedParts, types.CompletedPart{
			ETag:       uploadResp.ETag,
			PartNumber: int32(partNumber),
		})

		partNumber++
	}

	// Complete the multipart upload
	_, err = s3Client.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: uploadID,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to complete multipart upload: %w", err)
	}

	fmt.Println("File uploaded successfully!")
	return nil
}
