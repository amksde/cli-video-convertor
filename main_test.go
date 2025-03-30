package main

import (
	"os"
	"testing"
)

// test for the isValidMp4File function
func TestValidFileExistsMethod(t *testing.T) {
	cases := []struct {
		fileName string
		output   bool
	}{
		{"resources/sample.txt", false},
		{"resources/sample1.mp4", true}}

	for _, test_case := range cases {
		isValid := isValidMp4File(test_case.fileName)
		if isValid != test_case.output {
			t.Errorf("File " + test_case.fileName + " is not a valid MP4 file!")
		}

	}
}

func TestMain(m *testing.M) {
	// setup code here
	// run the tests
	code := m.Run()
	// teardown code here
	os.Exit(code)
}