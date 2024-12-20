package main

import (
	"context"
	"log"
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
			if strings.Contains(obj.Name, "localai/models")  {
				modelFiles = append(modelFiles, obj.Name)
			}
		}

		diffFiles := compareDirectories(&localFiles, modelFiles)
		if len(diffFiles) > 0 {
			pullModels(diffFiles)
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

func pullModels(pullFiles []string) {
	// copy files from bucket directory to local directory
	// final check

}
