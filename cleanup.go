package main

import (
	"fmt"
	"os"
)

func cleanUp() {
	fmt.Println("Cleaning up...")

	// clean up the images directory
	if err := os.RemoveAll(EXTRACTED_IMAGES_OUTPUT_DIR); err == nil {
		fmt.Println("✔ Images directory cleaned up")
	} else {
		fmt.Println("Error cleaning up images directory:", err)
	}

	// clean up the audio directory
	if err := os.RemoveAll(EXTRACTED_AUDIO_OUTPUT_DIR); err == nil {
		fmt.Println("✔ Audio directory cleaned up")
	} else {
		fmt.Println("Error cleaning up audio directory:", err)
	}

	fmt.Println("✅ Cleanup complete!")
}