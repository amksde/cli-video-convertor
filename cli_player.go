package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"time"
	"syscall"
	"github.com/gdamore/tcell/v2"
)

func PlayCLIAnimation(fps int) error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("failed to initialize terminal screen: %v", err)
	}
	defer screen.Fini() // Ensure cleanup on exit

	if err := screen.Init(); err != nil {
		return fmt.Errorf("failed to start terminal screen: %v", err)
	}
	
	// get the frames
	files, err := os.ReadDir(EXTRACTED_IMAGES_OUTPUT_DIR)
	if err != nil {
		return fmt.Errorf("failed to read frames from path %s: %v", EXTRACTED_IMAGES_OUTPUT_DIR, err)
	}

	// sort files to ensure correct playback order
	var framePaths []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".png" {
			framePaths = append(framePaths, filepath.Join(EXTRACTED_IMAGES_OUTPUT_DIR, file.Name()))
		}
	}
	sort.Strings(framePaths)
	frameDelay := time.Duration(1000/fps) * time.Millisecond

	// Start playing audio in background
	audioCmd := exec.Command("ffplay", "-nodisp", "-autoexit", filepath.Join(EXTRACTED_AUDIO_OUTPUT_DIR, AUDIO_FILE_NAME))
	err = audioCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to play audio: %v", err)
	}

	// for user interruption
	quit := make(chan struct{})
	go func() {
		for {
			ev := screen.PollEvent()
			switch ev.(type) {
			case *tcell.EventKey:
				close(quit)
				return
			}
		}
	}()
	
	// Handle Ctrl+C (SIGINT) to stop audio
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		close(quit)
	}()

	const batchSize = 20
	// load images in batches
	for i:=0; i < len(framePaths); i+=batchSize {
		end := i+ batchSize
		if end > len(framePaths) {
			end = len(framePaths)
		}
		// end is non-inclusive
		batchedPaths := framePaths[i:end]

		images, err := loadImages(batchedPaths)
		if err != nil {
			return fmt.Errorf("error loading images in batch %d: %w", i, err)
		}

		for _, img := range images {
			select {
				case <- quit: 
					fmt.Println("Stopping video and audio!")
					audioCmd.Process.Kill()
					return nil
				default:
			}

			screen.Clear()
			drawFrame(screen, img)
			screen.Show()

			time.Sleep(frameDelay)
		}
	}

	audioCmd.Wait() // Wait for audio to finish
	return nil
}

func drawFrame(screen tcell.Screen, img image.Image) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// convert image to ascii
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// r, g, b = r>>8, g>>8, b>>8 // Convert to 8-bit color

			style := tcell.StyleDefault.Background(tcell.NewRGBColor(int32(r), int32(g), int32(b)))
			screen.SetContent(x, y, ' ', nil, style)
		}
	}
}

func loadImages(imagePathsToLoad []string) ([]image.Image, error) {
	var images []image.Image
	for _, imagePath := range imagePathsToLoad {
		file, err := os.Open(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open image with path %s: %w", imagePath, err)
		}

		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image with path %s: %w", imagePath, err)
		}

		images = append(images, img)
	}

	return images, nil
}