package main

import (
	"reflect"
	"testing"
)

func TestCompareExtraFiles(t *testing.T) {
	localDir := []string{
		"llama.yaml",
		"llama.gguf",
		"nemo.yaml",
		"nemo.gguf",
	}
	bucketDir := []string{
		"llama.yaml",
		"llama.yaml",
		"llama-70b.yaml",
		"llama-70b.gguf",
		"nemo.yaml",
		"nemo.gguf",
	}

	wantPullList := []string{
		"llama-70b.yaml",
		"llama-70b.gguf",
	}
	wantLocalDir := []string{
		"llama.yaml",
		"llama.yaml",
		"nemo.yaml",
		"nemo.gguf",
	}

	gotPullList := compareDirectories(&localDir, bucketDir)

	if !reflect.DeepEqual(localDir, wantLocalDir) {
		t.Fatalf(`Got %v, wanted %v`, localDir, wantLocalDir)
	}
	if !reflect.DeepEqual(gotPullList, wantPullList) {
		t.Fatalf(`Got %v, wanted %v`, gotPullList, wantPullList)
	}
}

func TestCompareSameFiles(t *testing.T) {
	localDir := []string{
		"llama.yaml",
		"llama.gguf",
		"nemo.yaml",
		"nemo.gguf",
	}
	bucketDir := []string{
		"llama.yaml",
		"llama.yaml",
		"nemo.yaml",
		"nemo.gguf",
	}

	wantPullList := []string{
	}
	wantLocalDir := []string{
		"llama.yaml",
		"llama.yaml",
		"nemo.yaml",
		"nemo.gguf",
	}

	gotPullList := compareDirectories(&localDir, bucketDir)

	if !reflect.DeepEqual(localDir, wantLocalDir) {
		t.Fatalf(`Got %v, wanted %v`, localDir, wantLocalDir)
	}
	if !reflect.DeepEqual(gotPullList, wantPullList) {
		t.Fatalf(`Got %v, wanted %v`, gotPullList, wantPullList)
	}

}

func TestCompareDifferentFiles(t *testing.T) {
	localDir := []string{
		"llama.yaml",
		"llama.gguf",
		"nemo.yaml",
		"nemo.gguf",
	}
	bucketDir := []string{
		"llama-70b.yaml",
		"llama-70b.gguf",
		"nemo.yaml",
		"nemo.gguf",
	}

	wantPullList := []string{
		"llama-70b.yaml",
		"llama-70b.gguf",
	}
	wantLocalDir := []string{
		"nemo.yaml",
		"nemo.gguf",
	}

	gotPullList := compareDirectories(&localDir, bucketDir)

	if !reflect.DeepEqual(localDir, wantLocalDir) {
		t.Fatalf(`Got %v, wanted %v`, localDir, wantLocalDir)
	}
	if !reflect.DeepEqual(gotPullList, wantPullList) {
		t.Fatalf(`Got %v, wanted %v`, gotPullList, wantPullList)
	}
}

