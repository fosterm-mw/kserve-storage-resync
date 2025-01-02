package main

import (
	"context"
	"fmt"
	"os"
	"io"
	"strings"
	"time"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func syncBucket(ctx context.Context, modelBucket *storage.BucketHandle, modelPath string, resyncInterval int) {
	statusChan <- "Initializing"
	for true {
		localFiles, err := getFileNames(destination)
		if err != nil {
			statusChan <- fmt.Sprintf("Unable to read dir %s, error: %e", destination, err)
		}
		var modelFiles []string
		iter := modelBucket.Objects(ctx, nil)
		for {
			obj, err := iter.Next()
			if err == iterator.Done {
				log.Print("Finished reading from bucket!")
				break
			}
			if err != nil {
				statusChan <- fmt.Sprintf("Unable to read from bucket, error: %e", err)
				log.Fatalf("Unable to read from bucket, error: %s", err)
			} 
			if strings.Contains(obj.Name, modelPath) {
				if obj.Name != modelPath{
					modelFile := strings.Replace(obj.Name, modelPath, "", -1)
					log.Print(modelFile)
					modelFiles = append(modelFiles, modelFile)
				}
			}
		}

		diffFiles := compareDirectories(&localFiles, modelFiles, modelPath)
		if len(diffFiles) > 0 {
			statusChan <- fmt.Sprintf("Found new models, pulling...")
			if err = pullModels(ctx, modelBucket, diffFiles, modelPath, destination); err != nil {
				log.Fatalf("Error pulling model: %s", err)
			}
		}
		log.Print("Sleep")
		statusChan <- "Successful pull, sleeping"
		time.Sleep(time.Duration(resyncInterval) * time.Second)
	}
}

func compareDirectories(localDir *[]string, bucketDir []string, modelPath string) []string {
	log.Print("Comparing Directories...")
	var localDirCopy []string
	pullFiles := bucketDir
	offset := 0

	for _, localFile := range *localDir{
		for idx, bucketFile := range bucketDir{
			if bucketFile == modelPath {
				break
			}
			if localFile == bucketFile{
				if idx == 0 {
					pullFiles = pullFiles[1:]
					offset += 1
				} else {
					pullFiles = append(pullFiles[:idx-offset], pullFiles[idx+1-offset:]...)
				}
				localDirCopy = append(localDirCopy, localFile)
				break
			}
		}
	}
	*localDir = localDirCopy

	return pullFiles
}

func pullModels(ctx context.Context, modelBucket *storage.BucketHandle, pullFiles []string, modelPath, destination string) error {
	statusChan <- "Pulling Models"
	for _, model := range pullFiles {
		f, err := os.Create(destination + "/" + model)
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}
		rc, err := modelBucket.Object(modelPath + model).NewReader(ctx)
		if err != nil {
			return fmt.Errorf("Object(%q).NewReader: %w", modelPath +  model, err)
		}
		defer rc.Close()
		if _, err := io.Copy(f, rc); err != nil {
			return fmt.Errorf("io.Copy: %w", err)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("f.Close: %w", err)
		}
		log.Printf("Successfully downloaded: %s", model)
	}
	return nil
}

