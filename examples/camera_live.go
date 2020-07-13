package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"github.com/jonmol/gphoto2"
)

// Capture images in preview mode
func previewMode(c *gphoto2.Camera, finished chan bool) {
	if _, err := os.Stat(previewDir); os.IsNotExist(err) {
		os.Mkdir(previewDir, 0775)
	}
	if _, err := os.Stat(previewDir); err != nil {
		fmt.Println("Failed to create directory at", previewDir, "giving up!", err)
		finished <- true
		return
	}

	buf := new(bytes.Buffer)
	for i := 0; i < previewAmount; i++ {
		time.Sleep(200 * time.Millisecond)

		fmt.Println("Taking shot", i+1, "of", previewAmount)
		cameraLock.Lock()
		if err := c.CapturePreview(buf); err != nil {
			fmt.Println("Failed to take preview, make sure your camera is in Manual mode", err)
			cameraLock.Unlock()
			continue
		}
		cameraLock.Unlock()

		fmt.Println("Converting preview to jpeg")
		img, err := jpeg.Decode(buf)
		if err != nil {
			fmt.Println("Failed to make jpeg out of it. Buffer is ", buf.Len(), "bytes. The go library can't handle when the camera makes multiple copies, raw+jpeg for instance. You're welcome to patch, or set mode to only jpeg")
			buf = new(bytes.Buffer)
			continue
		}
		buf = new(bytes.Buffer)

		fn := filepath.Join(previewDir, fmt.Sprintf("%05d.jpeg", i+1))
		if f, err := os.Create(fn); err != nil {
			fmt.Println("Failed to create file", fn, "with error:", err)
			continue
		} else {
			defer f.Close()
			jpeg.Encode(f, img, nil)
			fmt.Println("Wrote to file ", fn)
		}
	}
	finished <- true
	fmt.Println("Preview done")
}
