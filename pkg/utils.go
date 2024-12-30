package main

import (
	"strings"
	"os"
	"log"
)

func parseBucketURI(bucketUri string) (string, string) {
	parsedBucketUri := strings.FieldsFunc(bucketUri, func(r rune) bool {
		if r == ':' || r == '/' {
			return true
		}
		return false
	})

	var modelDir string
	for i, dir := range parsedBucketUri {
		if i > 1 {
			modelDir += dir + "/"
		}
	}

	return parsedBucketUri[1], modelDir
}

func getFileNames(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("No value set for %s, using %s", key, fallback)
	return fallback
}
