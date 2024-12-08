package services

import (
	"ProductManagement/db"
	
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/streadway/amqp"
)

type ImageProcessingTask struct {
	ProductID int      `json:"product_id"`
	ImageURLs []string `json:"image_urls"`
}

func StartImageProcessor(channel *amqp.Channel) error {
	msgs, err := channel.Consume(
		"image_processing_queue",
		"",    
		true,  
		false, 
		false,
		false, 
		nil,   
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	go func() {
		for msg := range msgs {
			var task ImageProcessingTask
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				log.Printf("Failed to parse task: %v", err)
				continue
			}
			processImages(task)
		}
	}()

	return nil
}

func processImages(task ImageProcessingTask) {
	var compressedImageLinks []string

	for _, url := range task.ImageURLs {
		compressedURL, err := compressAndUploadImage(url, task.ProductID)
		if err != nil {
			log.Printf("Failed to process image %s: %v", url, err)
			continue
		}
		compressedImageLinks = append(compressedImageLinks, compressedURL)
	}

	if len(compressedImageLinks) > 0 {
		compressedLinks := strings.Join(compressedImageLinks, ",")
		_, err := db.DB.Exec(
			"UPDATE products SET compressed_product_images = ? WHERE id = ?",
			compressedLinks, task.ProductID,
		)
		if err != nil {
			log.Printf("Failed to update product %d with compressed images: %v", task.ProductID, err)
		} else {
			log.Printf("Successfully updated product %d with compressed images", task.ProductID)
		}
	}
}

func compressAndUploadImage(imageURL string, productID int) (string, error) {
	
	resp, err := http.Get(imageURL)
	if err != nil {

		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	compressedImage := new(bytes.Buffer)
	if _, err := io.Copy(compressedImage, resp.Body); err != nil {
		return "", fmt.Errorf("failed to compress image: %w", err)
	}

	imageReader := bytes.NewReader(compressedImage.Bytes())

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"), 
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := s3.New(sess)

	filename := fmt.Sprintf("compressed_images/product_%d_%s", productID, filepath.Base(imageURL))
	bucketName := "product-management-bucket-2024"

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		Body:        imageReader,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, filename)
	return s3URL, nil
}

