package main

import (
	"fmt"
	"os"
)

func main() {

	// checking for user provided arguments
	if (len(os.Args)) < 2 {
		fmt.Println("Please provide a file name!")
		os.Exit(1)
	}

	videoFileName := os.Args[1]
	fmt.Println("Processing video file : ", videoFileName)

	// validating the file
	if !isValidMp4File(videoFileName) {
		fmt.Println("Error!", videoFileName, "is not a valid existing mp4 file!")
		os.Exit(1)
	}

	fmt.Println("Converting video file to images...")



	err := extractSeparateMediaFromVideo(videoFileName, EXTRACTED_IMAGES_OUTPUT_DIR, EXTRACTED_AUDIO_OUTPUT_DIR)
	if err != nil {
		fmt.Println("Error extracting image-frames and sound from video!", err)
		cleanUp()
		os.Exit(1)
	}

	fps := getOriginalFPS(videoFileName)
	PrintAndWait(fmt.Sprint("🎉 Media extracted successfully! Playing animation with fps= %d", fps), 1)
	
	err = PlayCLIAnimation(fps)
	if err != nil {
		fmt.Println("Error playing animation!", err)
		cleanUp()
		os.Exit(1)
	}

	cleanUp()
}
