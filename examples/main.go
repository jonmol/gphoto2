package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/jonmol/gphoto2"
)

var (
	snap          bool   = false
	preview       bool   = false
	previewSnap   bool   = false
	printSettings bool   = false
	snapFile      string = "/tmp/snap_test.jpeg"
	previewDir    string = "/tmp/gphoto2_preview"
	shutterSpeed  string
	previewAmount int = 10
)

var cameraLock sync.Mutex // we need to lock the camera so we're not trying two operations at the same time

func init() {
	flag.IntVar(&previewAmount, "preview-shots", previewAmount, "Amount of shots to take in preview mode")

	flag.BoolVar(&snap, "snap", snap, "Take a picture and write to -snap-file")
	flag.BoolVar(&printSettings, "print-settings", printSettings, "Print camera settings")
	flag.BoolVar(&preview, "preview", preview, "Enter preview mode")
	flag.BoolVar(&previewSnap, "preview-snap", previewSnap, "Take -preview-shots amount of pictures in preview mode and write them to --preview-dir")

	flag.StringVar(&shutterSpeed, "shutter-speed", shutterSpeed, "Set shutter speed")
	flag.StringVar(&snapFile, "snap-file", snapFile, "File to write to when taking a picture")
	flag.StringVar(&previewDir, "preview-dir", previewDir, "Directory to write to when taking a series of preview pictures")
}

func main() {
	flag.Parse()
	camera, err := gphoto2.NewCamera("")
	if err != nil {
		panic(fmt.Sprintf("%s: %s", "Failed to connect to camera, make sure it's around!", err))
	}
	printInfo(camera)
	pFinished := make(chan bool)
	if preview {
		go previewMode(camera, pFinished)
	}

	if shutterSpeed != "" {
		setShutterSpeed("shutterspeed", camera)
	}

	if snap {
		snapPhoto(camera)
	}
	if preview {
		<-pFinished
	}
	camera.Exit()
	camera.Free()
}
