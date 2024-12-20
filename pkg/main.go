package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"context"

	"cloud.google.com/go/storage"
)

var (
	addr = flag.String("addr", ":8080", "http service address")

	gcsBucketURI = os.Getenv("GCS_BUCKET")
	resyncInterval = os.Getenv("RESYNC_INTERVAL")
)

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Cannot authenticate new client for GCP.")
	}

	bucketName, modelPath := parseBucketURI(gcsBucketURI)
	
	modelBucket := client.Bucket(bucketName)

	go syncBucket(ctx, modelBucket, modelPath)

	flag.Parse()
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		// respond with status of the pull, if there were any errors while syncing 
		// if there was a successful sync
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/healthz", *addr, http.StatusOK)
		fmt.Fprintf(w, "Healthy!")
	})
	
	log.Fatal(http.ListenAndServe(*addr, nil))
}

