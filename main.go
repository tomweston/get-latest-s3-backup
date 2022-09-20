// AWS SDK for Go v2 - S3 List Objects V2 Paginator
// This example demonstrates how to use the S3 ListObjectsV2Paginator to list
// objects in a bucket.
//
// Usage:
// go run s3_list_objects_v2_paginator.go
package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	bucket = "bucket-name-goes-here"
)

// Usage:
// go run s3_list_objects_v2_paginator.go
func main() {

	// get the current time in UTC
	now := time.Now().UTC()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get next page, %v", err)
		}

		for _, obj := range page.Contents {

			if obj.LastModified != nil {

				// for each object, get the last modified time
				lastModified := obj.LastModified.UTC()

				// calculate the difference between the current time and the last modified time
				diff := now.Sub(lastModified)

				// if the difference is less than 24 hours, print the object name
				if diff.Hours() < 24 {
					// download the object
					downloadObject(client, aws.ToString(obj.Key))

				}
			}

		}

	}
}

// downloadObject downloads an object from an Amazon S3 bucket.
// Inputs:
//
//	svc is an Amazon S3 service client
//	bucket is the name of the bucket
//	item is the name of the item
//
// Output:
//
//	If success, nil
//	Otherwise, an error from the call to GetObject
func downloadObject(svc *s3.Client, item string) error {
	// Get the object from the item and bucket.
	result, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})
	if err != nil {
		return err
	}

	// Write the contents of the object to a file.
	err = writeToFile(item, result.Body)
	if err != nil {
		return err
	}

	return nil
}

// writeToFile writes the contents of a stream to a file.
// Inputs:
//
//	filename is the name of the file
//	stream is the stream to write to the file
//
// Output:
//
//	If success, nil
//	Otherwise, an error from the call to ioutil.ReadAll
func writeToFile(filename string, stream io.ReadCloser) error {
	// Read the body into a byte slice.
	buf, err := ioutil.ReadAll(stream)
	if err != nil {
		return err
	}

	// Write the bytes to a file.
	err = ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		return err
	}

	return nil
}
