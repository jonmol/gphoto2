package main

import (
	"fmt"

	"github.com/jonmol/gphoto2"
)

func printInfo(c *gphoto2.Camera) {
	var cameraInfo = []string{
		"manufacturer",
		"cameramodel",
		"deviceversion",
		"lensname",
		"serialnumber",    // non-canon (at least nikon)
		"eosserialnumber", // canon being canon
	}

	fmt.Println("Fetching camera info...")
	cameraLock.Lock()
	for _, name := range cameraInfo {
		if w, err := c.GetSetting(name); err != nil || w == nil {
			if err != nil {

			} else {
				fmt.Printf("%s is nil\n", name)
			}
		} else {
			if val, err := w.Get(); err == nil {
				if str, ok := val.(string); ok {
					fmt.Println(name, " = ", str)
				}
			}
		}
	}
	cameraLock.Unlock()
}
