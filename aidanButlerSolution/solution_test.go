package main

import "testing"
import "os"

func TestRemoveDuplicates(t *testing.T) {
	original := []int{1, 2, 3, 2}
	removed := removeDuplicates(original)
	if len(original) == len(removed) {
		t.Errorf("Duplciate not removed")
	}
}

func TestGetJsonStrings(t *testing.T) {
	// TODO: mock test to get json strings from server
}

func TestBuildJsonStrings(t *testing.T) {
	// TODO: test to combine JSON strings according to schema
}

func TestIsValidJsonString(t *testing.T) {
	// TODO: test to check string is valid against schema
}

func TestPostJson(t *testing.T) {
	// TODO: mock test to post JSON strings to URL
}

func TestGetEnvironmentalVariableValue(t *testing.T) {
	// TODO: combine test to combine JSON strings according to schema
	url := "http://www.example.com/"
	os.Setenv("STORE_ADDR", url)
	val := getEnvironmentVariableValue()
	if val != url {
		t.Errorf("Error reading environmental variable.")
	}
}
