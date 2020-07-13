# gphoto2

Partial Go bindings for http://www.gphoto.org/proj/libgphoto2/

## History

Much of the code is copied from and/or insipired of https://github.com/szank/gphoto which seems abandonend and had a couple of bugs. I waanted to write a simple stop motion program and this is an artefact from it. So it's claiming to or trying to cover the whole library but you can do a lot, and I'm very open to suggestions and/or merge requests. I only have access to one Nikon camera so it's not tested with any other brand, and it seems like all camera manufacturers have their own quirks.

## Installlation

To build the library you need to have libgphoto2-6 and libgphoto2-port12 or later installed.

## Usage

Under examples are there a couple of basic usage examples for setting shutter speed, taking a photo, entering live view and copying images from there. A very simple example of taking a shot would be:

```Go
func main() {
	camera, err := gphoto2.NewCamera("")
	if err != nil {
		panic(fmt.Sprintf("%s: %s", "Failed to connect to camera, make sure it's around!", err))
	}
    snapFile := "/tmp/testshot.jpeg"
	if f, err := os.Create(snapFile); err != nil {
		fmt.Println("Failed to create temp file", snapFile, "giving up!", err)
	} else {
		fmt.Println("Taking shot, then copy to", snapFile)
		if err := c.CaptureDownload(f, false); err != nil {
			fmt.Println("Failed to capture!", err)
		}
	}
	camera.Exit()
	camera.Free()
}
```


## Notes

For Nikon (I don't know about other cameras as I haven't been able to test with them) you want to put the DSLR in manual mode, and turn auto focus off typically. Manual mode is necessary to be able to change many of the settings, and to be able to enter live view mode for instance. Manual focus makes taking shots a lot more reliable and fast as the camera won't have to focus for each shot.

