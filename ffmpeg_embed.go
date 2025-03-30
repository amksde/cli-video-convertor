package main

import (
	"fmt"
	"os"
	"os/exec"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"
)

const ffmpegURL = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"

func ensureFFMPeg() error { 

	// Check if ffmpeg already exists
	if _, err := os.Stat(FFMPEG_PATH); err == nil {
		fmt.Println("✔ ffmpeg already exists")
		return nil
	}

	// Download ffmpeg
	fmt.Println("⬇ Downloading ffmpeg...")
	resp, err := http.Get(ffmpegURL)
	if err != nil {
		return fmt.Errorf("failed to download ffmpeg: %v", err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(FFMPEG_PATH)
	if err != nil {
		return fmt.Errorf("failed to create ffmpeg file: %v", err)
	}
	defer out.Close()

	// Copy contents to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save ffmpeg binary: %v", err)
	}

	// Make it executable
	err = os.Chmod(FFMPEG_PATH, 0755)
	if err != nil {
		return fmt.Errorf("failed to set ffmpeg as executable: %v", err)
	}

	fmt.Println("✔ ffmpeg downloaded successfully")
	return nil
}

func extractSeparateMediaFromVideo(videoFilePath string, imagesOutputPath string, audioOutputPath string) error {

	err := ensureFFMPeg()

	if err != nil {
		fmt.Println("Error downloading ffmpeg!", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(imagesOutputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	if err := os.MkdirAll(audioOutputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	err = extractImagesFromVideo(FFMPEG_PATH, videoFilePath, imagesOutputPath)
	if err != nil {
		return fmt.Errorf("failed to extract images from video: %v", err)
	}

	err = extractAudioFromVideo(FFMPEG_PATH, videoFilePath, audioOutputPath)
	if err != nil {
		return fmt.Errorf("failed to extract audio from video: %v", err)
	}

	return nil
}

func extractImagesFromVideo(ffMpegTempPath string, videoFilePath string, imagesOutputPath string) error {
	// run it from the temp location
	// refer the command here https://ffmpeg.org/ffmpeg.html#toc-Video-and-Audio-file-format-conversion
	cmd := exec.Command(ffMpegTempPath, "-i", videoFilePath, "-vf", "scale="+RES_160_X_48, filepath.Join(imagesOutputPath, IMAGE_FILE_NAME))
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run()
}

func extractAudioFromVideo(ffMpegTempPath string, videoFilePath string, audioOutputPath string) error {
	// run it from the temp location
	// refer the command here https://ffmpeg.org/ffmpeg.html#toc-Video-and-Audio-file-format-conversion
	cmd := exec.Command(ffMpegTempPath, "-i", videoFilePath, "-q:a", "0", "-map", "a", filepath.Join(audioOutputPath, AUDIO_FILE_NAME))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func getOriginalFPS(videoFilePath string) int {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", "-of", "default=noprint_wrappers=1:nokey=1", videoFilePath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting FPS:", err)
		return 30 // default FPS
	}

	fpsStr := string(output)
	fpsParts := strings.Split(fpsStr, "/")
	if len(fpsParts) != 2 {
		fmt.Println("Error parsing FPS:", fpsStr)
		return 30 // default FPS
	}

	num, err := strconv.Atoi(strings.TrimSpace(fpsParts[0]))
	if err != nil {
		fmt.Println("Error converting FPS numerator:", err)
		return 30 // default FPS
	}

	denom, err := strconv.Atoi(strings.TrimSpace(fpsParts[1]))
	if err != nil {
		fmt.Println("Error converting FPS denominator:", err)
		return 30 // default FPS
	}

	return num / denom
}