package aws_s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// UploadImageToS3 uploads an image to the specified S3 bucket and returns the URL.
func UploadImageToS3(ctx context.Context, imageBytes []byte, bucketName, key string) (string, error) {
	// Load the shared AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"TEB2ISTYBQMY3XKVU4OJ",
			"pBgyz0eXTpd3Aac67Cg5ClU7Fw1wvjGV6Hn3zDwR",
			"",
		)),
	)
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Create an S3 client with custom endpoint
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointResolver = s3.EndpointResolverFromURL("https://object.pscloud.io")
		o.UsePathStyle = true // Use path-style addressing
	})

	// Upload the file to S3 using PutObject
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(imageBytes),
		//ContentType: aws.String("image/jpeg"),        // Adjust as needed
		ACL: types.ObjectCannedACLPublicRead, // Optional: Set ACL to make the object publicly readable
	}

	_, err = client.PutObject(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	// Construct the URL of the uploaded object
	url := fmt.Sprintf("https://%s.object.pscloud.io/%s", bucketName, key)

	return url, nil
}
