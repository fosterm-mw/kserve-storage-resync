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
				log.Print(obj.Name)
				modelFiles = append(modelFiles, obj.Name)
			}
		}

		diffFiles := compareDirectories(&localFiles, modelFiles, modelPath)
		if len(diffFiles) > 0 {
			statusChan <- fmt.Sprintf("Found new models, pulling...")
			if err = pullModels(ctx, modelBucket, diffFiles, destination); err != nil {
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
	diffList := make(map[string]bool)
	for _, file := range *localDir {
		diffList[file] = false
	}
	for _, model := range bucketDir {
		diffList[model] = true
	}

	var localDirCopy []string
	for _, file := range *localDir {
		if diffList[file] == false {
			delete(diffList, file)
		}
		if diffList[file] == true {
			localDirCopy = append(localDirCopy, file)
			delete(diffList, file)
		}
	}
	*localDir = localDirCopy

	var pullList []string
	for file := range diffList {
		if file != modelPath {
			pullList = append(pullList, file)
			log.Print(file)
		}
	}
	// if (len(pullList) > 0) {
	// 	return pullList[1:]
	// } else {
	// 	return pullList
	// }
	return pullList
}

func pullModels(ctx context.Context, modelBucket *storage.BucketHandle, pullFiles []string, destination string) error {
	statusChan <- "Pulling Models"
	for _, model := range pullFiles {
		f, err := os.Create(destination + "/" + model)
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}
		rc, err := modelBucket.Object(model).NewReader(ctx)
		if err != nil {
			return fmt.Errorf("Object(%q).NewReader: %w", model, err)
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

