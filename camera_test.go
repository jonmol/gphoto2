package gphoto2

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func estNewCamera(t *testing.T) {
	t.Log("Start")
	t.Log("Connect first")
	c, err := NewCamera("")
	if c == nil || err != nil {
		fmt.Println("Failed first")
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	fmt.Println("Moving along")
	t.Log("First connected, about to exit", err)
	c.Exit()
	t.Log("First exit, connecting next", err)
	c2, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	c2.Exit()
	c.Free()
	c2.Free()
}

func estNewCameraPhotos(t *testing.T) {
	t.Log("Start")
	t.Log("Connect first")
	c, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	t.Log("First connected, about to exit", err)
	c.Exit()
	t.Log("First exit, connecting next", err)
	c2, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	c2.Exit()
	//	c.CaptureImage()
	c.Exit()
	//	c2.CaptureImage()
	c.Free()
	c2.Free()
}

/*
func estSettings(t *testing.T) {
	c, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}

	if s, err := c.GetSetting(SettingProgram); err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to get settings", err.Error()))
		t.Fail()
	} else {
		t.Log(s.Name(), s.Label(), s.Info(), s.Type())
		switch widgetAccessor := s.(type) {
		case CameraWidgetRadio:
			if choices, err := widgetAccessor.GetChoices(); err == nil {
				t.Log(fmt.Sprintf("      Choices : %+v", choices))
			} else {
				t.Log("Not knowing")
			}
		}
	}
	c.Free()

}
*/

func estSettings3(t *testing.T) {
	c, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}

	for k := range CommonSettings {
		t.Log("Doing ", k)
		if k == "size" {
			t.Log("Skipping because segfault") //TODO
			continue
		}
		if s, err := c.GetSetting(k); err != nil || s == nil {
			t.Log("Failed to get", k, err)
		} else {
			g, _ := s.Get()
			t.Log(s.Label(), s.Name(), s.Type(), s.ReadOnly(), g)
			o, e := s.Options()
			t.Log("Got a setting", strings.Join(o, "-"), e)
		}
	}
	c.Free()
}
func estSettings2(t *testing.T) {
	c, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}

	if s, err := c.GetSetting(SettingAperture); err != nil {
		t.Log("Failed to get 	1")
		t.Fail()
		return
	} else {
		fooType := reflect.TypeOf(s)
		fmt.Println("It's a ", s.Type())
		fmt.Println("Name=", fooType.Name())
		for i := 0; i < fooType.NumMethod(); i++ {
			method := fooType.Method(i)
			fmt.Println(method.Name)
		}
		o, e := s.Options()
		t.Log("Got a setting", o, e)
	}
	t.Log("Now going exit and checking iso")
	if s, err := c.GetSetting("iso"); err != nil {
		t.Log("Failed to get 	1")
		t.Fail()
		return
	} else {
		fmt.Println("It's a ", s.Type())

		fooType := reflect.TypeOf(s)
		fmt.Println("Name=", fooType.Name())
		for i := 0; i < fooType.NumMethod(); i++ {
			method := fooType.Method(i)
			fmt.Println(method.Name)
		}
		o, e := s.Options()
		t.Log("Got a setting", o, e)
	}
	t.Log("Now going exit and checking iso")
	c.Exit()

	if _, err := c.GetSetting(SettingISO); err != nil {
		t.Log("Failed to get 	1")
		t.Fail()
		return
	}
	c.Free()
}

func estSettings5(t *testing.T) {
	c, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	res := make([]*CameraWidget, 0, len(CommonSettings))
	bla := *c
	for name := range CommonSettings {
		fmt.Println("************************************************** STARTING TO LOOK FOR ", name)
		if w, err := bla.GetSetting(name); err != nil || w == nil {
			fmt.Println("getCameraWidgets skipping", name, err)
		} else {
			fmt.Println(w)
			fmt.Println("Tried finding", name)
			fmt.Println("getCameraWidgets adding ", name, w.name)
			w.Name()
			res = append(res, w)
		}
	}
	c.Free()
}

func estListSettings(t *testing.T) {
	names := []string{"serialnumber", "manufacturer", "cameramodel", "deviceversion", "vendorextension", "lensname", "batterylevel", "orientation", "orientation2", "acpower", "flashopen"}

	c, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}

	for _, name := range names {
		if w, err := c.GetSetting(name); err != nil || w == nil {
			fmt.Println("getCameraWidgets skipping", name, err)
		} else {
			o, _ := w.Options()
			v, _ := w.Get()
			fmt.Println(fmt.Sprintf("[%s][%s][%+v][%s]", w.label, w.name, o, v))
		}

	}
	c.Free()

}

func estSettingd1a3(t *testing.T) {
	c, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	res := make([]*CameraWidget, 0, len(CommonSettings))
	if w, err := c.GetSetting(SettingLiveViewZoomRatio); err != nil || w == nil {
		fmt.Println("getCameraWidgets skipping", SettingLiveViewZoomRatio, err)
	} else {
		fmt.Println(w)
		fmt.Println("Tried finding", SettingLiveViewZoomRatio)
		fmt.Println("getCameraWidgets adding ", SettingLiveViewZoomRatio, w.name)
		buf := new(bytes.Buffer)
		c.CapturePreview(buf)

		opts, _ := w.Options()
		fmt.Println("Options are", opts)
		val, _ := w.Get()
		fmt.Println("Value is", val)
		w.Set("3")
		w2, _ := c.GetSetting(SettingLiveViewZoomRatio)
		opts, _ = w2.Options()
		fmt.Println("Options are", opts)
		val, _ = w2.Get()
		fmt.Println("Value is", val)

		res = append(res, w)
	}
	c.Free()
}

func estPrint(t *testing.T) {
	c, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	c.LoadWidgets()
	fmt.Println(c.Settings)
	names := []string{"serialnumber", "manufacturer", "cameramodel", "deviceversion", "vendorextension", "lensname", "batterylevel", "orientation", "orientation2", "acpower", "flashopen"}
	for _, name := range names {
		if w, err := c.GetSetting(name); err != nil || w == nil {
			fmt.Println("Failed", name, err)
		} else {
			fmt.Println(w)
		}
	}
	c.Free()
}

func _TestDisco(t *testing.T) {
	names := []string{"serialnumber", "manufacturer", "cameramodel", "deviceversion", "vendorextension", "lensname", "batterylevel", "orientation", "orientation2", "acpower", "flashopen"}
	t.Log("Start")
	t.Log("Connect first")
	c, err := NewCamera("")
	if c == nil || err != nil {
		fmt.Println("Failed first")
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	fmt.Println("Moving along")
	for _, name := range names {
		if w, err := c.GetSetting(name); err != nil || w == nil {
			fmt.Println("getCameraWidgets skipping", name, err)
		} else {
			o, _ := w.Options()
			v, _ := w.Get()
			fmt.Println(fmt.Sprintf("[%s][%s][%+v][%s]", w.label, w.name, o, v))
		}

	}

	t.Log("First connected, about to exit", err)
	c.Exit()
	t.Log("First exit, connecting next", err)
	c2, err := NewCamera("")
	if c == nil || err != nil {
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	for _, name := range names {
		if w, err := c2.GetSetting(name); err != nil || w == nil {
			fmt.Println("getCameraWidgets skipping", name, err)
		} else {
			o, _ := w.Options()
			v, _ := w.Get()
			fmt.Println(fmt.Sprintf("[%s][%s][%+v][%s]", w.label, w.name, o, v))
		}

	}

	c2.Exit()
	c.Free()
	c2.Free()

}

func TestDisco(t *testing.T) {
	c, err := NewCamera("")
	if c == nil || err != nil {
		fmt.Println("Failed connect")
		t.Log(fmt.Sprintf("%s: %s", "Failed to connect", err.Error()))
		t.Fail()
		return
	}
	SetIdleFunc(c.Ctx, defaultIdleCallback)

	previewDir := "/tmp/gphoto2_preview"
	if _, err := os.Stat(previewDir); os.IsNotExist(err) {
		os.Mkdir(previewDir, 0775)
	}
	if _, err := os.Stat(previewDir); err != nil {
		fmt.Println("Failed to create directory at", previewDir, "giving up!", err)
		return
	}

	buf := new(bytes.Buffer)

	for i := 0; i < 10; i++ {
		time.Sleep(150 * time.Millisecond)

		fmt.Println("Taking shot", i+1, "of", 10)
		if err := c.CapturePreview(buf); err != nil {
			fmt.Println("Failed to take preview, make sure your camera is in Manual mode", err)
			continue
		}

		if i == 3 {
			snapFile := "/tmp/snap_test.jpeg"
			if f, err := os.Create(snapFile); err != nil {
				fmt.Println("Failed to create temp file", snapFile, "giving up!", err)
			} else {
				fmt.Println("Taking shot, then copy to", snapFile)
				if err := c.CaptureDownload(f, true); err != nil {
					fmt.Println("Failed to capture!", err)
				}
			}
			time.Sleep(250 * time.Millisecond)
		}

		fmt.Println("Converting preview to jpeg")
		img, err := jpeg.Decode(buf)
		if err != nil {
			fmt.Println("Failed to make jpeg out of it. Buffer is ", buf.Len(), "bytes")
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
	c.Exit()
	c.Free()

}
