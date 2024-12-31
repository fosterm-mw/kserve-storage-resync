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
	for true {
		localFiles, err := getFileNames(destination)
		if err != nil {
			// WarningLogger.Println("Unable to read dir", destination, "error: ", err)
			log.Println("Unable to read dir", destination, "error: ", err)
		}
		var modelFiles []string
		iter := modelBucket.Objects(ctx, nil)
		for {
			obj, err := iter.Next()
			if err == iterator.Done {
				// InfoLogger.Print("Finished reading from bucket!")
				log.Print("Finished reading from bucket!")
				break
			}
			if err != nil {
				// ErrorLogger.Fatalf("Unable to read from bucket, error: %s", err)
				log.Fatalf("Unable to read from bucket, error: %s", err)
			} 
			if strings.Contains(obj.Name, modelPath) {
				// DebugLogger.Print(obj.Name)
				log.Print(obj.Name)
				modelFiles = append(modelFiles, obj.Name)
			}
		}

		diffFiles := compareDirectories(&localFiles, modelFiles)
		if len(diffFiles) > 0 {
			// InfoLogger.Printf("Found new Models, pulling...")
			log.Printf("Found new Models, pulling...")
			if err = pullModels(ctx, modelBucket, diffFiles, modelPath, destination); err != nil {
				// ErrorLogger.Fatalf("Error pulling model: %s", err)
				log.Fatalf("Error pulling model: %s", err)
			}
		}
		// InfoLogger.Print("Sleep")
		log.Print("Sleep")
		time.Sleep(time.Duration(resyncInterval) * time.Second)
	}
}

func compareDirectories(localDir *[]string, bucketDir []string) []string {
	// InfoLogger.Print("Comparing Directories...")
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
		pullList = append(pullList, file)
		// DebugLogger.Print(file)
		log.Print(file)
	}
	// if (len(pullList) > 0) {
	// 	return pullList[1:]
	// } else {
	// 	return pullList
	// }
	return pullList[1:]
}

func pullModels(ctx context.Context, modelBucket *storage.BucketHandle, pullFiles []string, modelPath string, destination string) error {
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
		// InfoLogger.Printf("Successfully downloaded: %s", model)
		log.Printf("Successfully downloaded: %s", model)
	}
	return nil
}

