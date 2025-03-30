package main

import (
	"strings"
	"errors"
	"os"
	"time"
	"fmt"
)

const CLS_CMD_CHAR = "\033[H\033[2J"

const RES_160_X_48 = "160:48"
const RES_80_X_24 = "80:24"
const RES_720_X_480 = "720:480" 

const EXTRACTED_IMAGES_OUTPUT_DIR = "extracted_images"
const EXTRACTED_AUDIO_OUTPUT_DIR = "extracted_audio"
const AUDIO_FILE_NAME = "audio.mp3"
const IMAGE_FILE_NAME = "frame_%06d.png"
const FFMPEG_PATH = "bin/ffmpeg"

func isValidMp4File(fileName string) bool {
	// checking the fileName for .mp4 extension
	if !strings.HasSuffix(fileName, ".mp4") {
		return false
	}

	_, err := os.Stat(fileName)
	return !errors.Is(err, os.ErrNotExist)
}

func PrintAndWait(msg string, sleepTime int) {
	fmt.Println(msg)
	time.Sleep(time.Duration(sleepTime) * time.Second)
}