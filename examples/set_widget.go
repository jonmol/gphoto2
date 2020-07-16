package main

import (
	"fmt"

	"github.com/jonmol/gphoto2"
)

func setShutterSpeed(name string, c *gphoto2.Camera) {
	cameraLock.Lock()
	if w, err := c.GetSetting(name); err != nil || w == nil {
		fmt.Println("Failed to get", name, " ", err)
	} else {
		fmt.Println(name, "before trying to change (written as 'internal name in camera' 'Display name' 'type of widget' '[ list of valid values]' 'current value'")
		fmt.Println(w)
		fmt.Println("Trying to set", name, "to", shutterSpeed)
		fmt.Println("**************************************************")

		w.Set(shutterSpeed)
		fmt.Println("Value after setting:")
		fmt.Println(w)
		fmt.Print("\n\n\n\n\n")
	}
	cameraLock.Unlock()
}
