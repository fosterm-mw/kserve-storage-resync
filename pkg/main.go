package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"context"
	"strconv"

	"cloud.google.com/go/storage"
)

var (
	addr = flag.String("addr", ":8080", "http service address")

	gcsBucketURI = os.Getenv("GCS_BUCKET")
	envResyncInterval = getEnv("INTERVAL", "300")
	destination = getEnv("DEST", "tmp/mnt")
)

var statusChan chan string

var (
	InfoLogger *log.Logger
	DebugLogger *log.Logger
	WarningLogger *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		ErrorLogger.Fatal("Cannot authenticate new client for GCP.")
	}

	bucketName, modelPath := parseBucketURI(gcsBucketURI)

	modelBucket := client.Bucket(bucketName)

	resyncInterval, err := strconv.Atoi(envResyncInterval)
	if err != nil {
		ErrorLogger.Fatal("Cannot convert INTERVAL to int")
	}

	statusChan = make(chan string)
	go syncBucket(ctx, modelBucket, modelPath, resyncInterval)

	flag.Parse()
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		InfoLogger.Println("/status", *addr, http.StatusOK)
		status := <- statusChan
		fmt.Fprintf(w, "Status: %s", status)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		InfoLogger.Println("/healthz", *addr, http.StatusOK)
		fmt.Fprintf(w, "Healthy!")
	})
	
	ErrorLogger.Fatal(http.ListenAndServe(*addr, nil))
}
