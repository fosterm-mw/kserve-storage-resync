package main

import (
	"testing"
)

func TestParseBucketUri(t *testing.T) {
	bucketUris := []string{
		"gs://test-bucket-gogo/models",
		"gs://test-bucket-gogo/localai/models",
	}
	wantModelDirs := []string{
		"models/",
		"localai/models/",
	}
	wantBucket := "test-bucket-gogo"

	for i, bucketUri := range bucketUris {
		gotBucket, gotModelDir := parseBucketURI(bucketUri)

		if (gotBucket != wantBucket) {
			t.Fatalf(`Got %s, wanted %s`, gotBucket, wantBucket)
		}
		if (gotModelDir != wantModelDirs[i]) {
			t.Fatalf(`Got %s, wanted %s`, gotModelDir, wantModelDirs[i])
		}
	}
}
