# Kserve Resync Model Script

This repository contains a script to expose an endpoint and continuously sync a GCS bucket with the models downloaded to the pod locally.

## Env
`GCS_BUCKET` - Google Cloud Storage Bucket Uri
`INTERVAL` - Time to check bucket for any changes (seconds)
`DEST` - Destination of files to be synced from GCS bucket

## Run
Locally:
```
go run .
```
Navigate to http://localhost:8080 to see running web server.

## Testing
Unit Tests: `go test -v`
e2e:
* Create GCS Bucket with files
* `gcloud auth application-default login`
* Add files to GCS Bucket different from `tmp` dir
* Check if the files have been added

## Contributing

## Maintainers
* fostemi
