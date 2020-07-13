package main

import (
	"fmt"
	"os"

	"github.com/jonmol/gphoto2"
)

// take a photo
func snapPhoto(c *gphoto2.Camera) {
	cameraLock.Lock()
	if f, err := os.Create(snapFile); err != nil {
		fmt.Println("Failed to create temp file", snapFile, "giving up!", err)
	} else {
		fmt.Println("Taking shot, then copy to", snapFile)
		if err := c.CaptureDownload(f, false); err != nil {
			fmt.Println("Failed to capture!", err)
		}
	}
	cameraLock.Unlock()
}
