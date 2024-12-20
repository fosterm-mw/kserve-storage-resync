package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func syncBucket(ctx context.Context, modelBucket *storage.BucketHandle, modelPath string) {
	for true {
		// TODO change this I don't like it and by this I mean dir
		// I want to be able to set this to a env var
		dir := "tmp/mnt/models"
		localFiles, err := getFileNames(dir)
		if err != nil {
			log.Println("Unable to read dir", dir, "error: ", err)
		}
		var modelFiles []string
		iter := modelBucket.Objects(ctx, nil)
		for {
			obj, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Unable to read from bucket, error: %s", err)
			}
			if strings.Contains(obj.Name, modelPath)  {
				modelFiles = append(modelFiles, obj.Name)
			}
		}

		diffFiles := compareDirectories(&localFiles, modelFiles)
		if len(diffFiles) > 0 {
			if err = pullModels(ctx, modelBucket, diffFiles, modelPath, dir); err != nil {
				log.Fatalf("Error pulling model: %s", err)
			}
		}
		time.Sleep(180)
	}
}

func compareDirectories(localDir *[]string, bucketDir []string) []string {
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
	}
	return pullList
}

func pullModels(ctx context.Context, modelBucket *storage.BucketHandle, pullFiles []string, modelPath string, destination string) error {
	for _, model := range pullFiles {
		f, err := os.Create(destination + model)
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}
		rc, err := modelBucket.Object(modelPath + model).NewReader(ctx)
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

