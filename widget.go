package gphoto2

/** \file
 *
 * \author Copyright 2020 Jon Molin
 *
 * \note
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2 of the License, or (at your option) any later version.
 *
 * \note
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * \note
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the
 * Free Software Foundation, Inc., 51 Franklin Street, Fifth Floor,
 * Boston, MA  02110-1301  USA
 */

// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include <gphoto2/gphoto2-setting.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

type WidgetType string

type CameraWidget struct {
	id           int
	label        string
	name         string
	info         string
	widgetType   WidgetType
	gpWidgetType uint
	readonly     bool

	parent   *CameraWidget
	children []*CameraWidget
	camera   *Camera
}

func (w *CameraWidget) String() string {
	res := fmt.Sprintf("%s %s %s ", w.name, w.label, w.widgetType)

	if opts, err := w.Options(); err == nil {
		res += " [ " + strings.Join(opts, " | ") + " ] "
	}
	if val, err := w.Get(); err == nil {
		str, ok := val.(string)
		if ok {
			res += str + " "
		}
	}
	return res
}

func (w *CameraWidget) Find(name string) *CameraWidget {
	if name == "" {
		return nil
	}
	if w.name == name {
		return w
	}
	if w.children != nil && len(w.children) > 0 {
		for _, c := range w.children {
			if m := c.Find(name); m != nil {
				return m
			}
		}
	}
	return nil
}

func (w *CameraWidget) FindByPath(path string) *CameraWidget {
	if path == "" {
		return w
	}
	path = strings.TrimPrefix(path, "/")
	pathComponents := strings.Split(path, "/")
	if w.name == pathComponents[0] {
		return w.FindByPath(strings.Join(pathComponents[1:], "/"))
	}

	for _, v := range w.children {
		if v.name == pathComponents[0] {
			return v.FindByPath(strings.Join(pathComponents[1:], "/"))
		}
	}
	return nil
}

func (w CameraWidget) ID() int {
	return w.id
}
func (w CameraWidget) Type() WidgetType {
	return w.widgetType
}

func (w CameraWidget) Label() string {
	return w.label
}
func (w CameraWidget) Name() string {
	return w.name
}
func (w CameraWidget) Info() string {
	return w.Info()
}
func (w CameraWidget) ReadOnly() bool {
	return w.readonly
}
func (w CameraWidget) Parent() *CameraWidget {
	return w.parent
}
func (w CameraWidget) Children() []*CameraWidget {
	return w.children
}

func (w CameraWidget) saveSetting(gpCW *C.CameraWidget) error {
	gpName := C.CString(w.name)
	if res := C.gp_camera_set_single_config(w.camera.gpCamera, gpName, gpCW, w.camera.Ctx.gpContext); res != GPOK {
		return newError("Could not set widget value with gp_camera_set_single_config", int(res))
	}
	return nil
}

func (w CameraWidget) Options() ([]string, error) {
	switch w.widgetType {
	case WidgetToggle:
		return []string{"on", "off"}, nil
	case WidgetRadio, WidgetMenu:
		return w.radioOptions()
	}
	return nil, newError("Widget has no options", ErrorWidgetHasNoOptions)
}

func (w CameraWidget) radioOptions() ([]string, error) {
	var gpWidget *C.CameraWidget
	var err error
	if gpWidget, err = w.camera.getChildWidget(&w.name); err != nil {
		return nil, err
	}
	defer w.camera.freeChildWidget(gpWidget)

	choicesList := []string{}
	numChoices := C.gp_widget_count_choices(gpWidget)
	for i := 0; i < int(numChoices); i++ {
		var gpChoice *C.char
		C.gp_widget_get_choice(gpWidget, C.int(i), (**C.char)(unsafe.Pointer(&gpChoice)))
		choicesList = append(choicesList, C.GoString(gpChoice))
	}
	return choicesList, nil
}

func (w CameraWidget) Set(input interface{}) error {
	if w.readonly {
		return newError("Widget is readonly", ErrorReadOnly)
	}

	switch w.widgetType {
	case WidgetRadio, WidgetMenu:
		return w.setRadio(input.(string))
	case WidgetText:
		return w.setText(input.(string))
	case WidgetToggle:
		return w.setToggle(input.(bool))
	}
	return nil
}

// TODO: use interface and scrap the three set functions?
func (w CameraWidget) setToggle(input bool) error {
	to := C.int(1)
	if !input {
		to = C.int(0)
	}
	var gpWidget *C.CameraWidget
	var err error
	if gpWidget, err = w.camera.getChildWidget(&w.name); err != nil {
		return err
	}
	defer w.camera.freeChildWidget(gpWidget)
	if res := C.gp_widget_set_value(gpWidget, unsafe.Pointer(&to)); res != GPOK {
		return newError("Could not set widget value", int(res))
	}
	return w.saveSetting(gpWidget)
}

func (w CameraWidget) setText(input string) error {
	var err error
	gpText := C.CString(input)
	defer C.free(unsafe.Pointer(gpText))

	var gpWidget *C.CameraWidget
	if gpWidget, err = w.camera.getChildWidget(&w.name); err != nil {
		return err
	}
	defer w.camera.freeChildWidget(gpWidget)
	if res := C.gp_widget_set_value(gpWidget, unsafe.Pointer(gpText)); res != GPOK {
		return newError("Could not set widget value", int(res))
	}
	return w.saveSetting(gpWidget)
}

func (w CameraWidget) SetInt(input int) error {
	var err error
	gpInt := C.int(input)

	var gpWidget *C.CameraWidget
	if gpWidget, err = w.camera.getChildWidget(&w.name); err != nil {
		return err
	}
	defer w.camera.freeChildWidget(gpWidget)
	if res := C.gp_widget_set_value(gpWidget, unsafe.Pointer(&gpInt)); res != GPOK {
		return newError("Could not set widget value", int(res))
	}
	return w.saveSetting(gpWidget)
}

func (w CameraWidget) setRadio(input string) error {
	if choices, err := w.Options(); err == nil {
		for _, item := range choices {
			if item == input {
				return w.setText(input)
			}
		}
		return newError("Could not find provided value in alloved values list", ErrorWidgetIllegalOption)
	} else {
		return err
	}
}

func (w CameraWidget) Get() (interface{}, error) {
	var gpWidget *C.CameraWidget
	var err error
	if gpWidget, err = w.camera.getChildWidget(&w.name); err != nil {
		return nil, err
	}
	defer w.camera.freeChildWidget(gpWidget)

	switch w.widgetType {
	case WidgetRadio, WidgetMenu, WidgetText:
		var data *C.char
		if res := C.gp_widget_get_value(gpWidget, (unsafe.Pointer(&data))); res != GPOK {
			return nil, newError("Cannot read widget property value", int(res))
		}
		val := C.GoString(data)
		return val, nil
	case WidgetToggle:
		var data int
		if res := C.gp_widget_get_value(gpWidget, (unsafe.Pointer(&data))); res != GPOK {
			return nil, newError("Cannot read widget property value", int(res))
		}
		val := (data == 1)
		return val, nil
	}
	return nil, nil
}

func widgetType(gpWidgetType C.CameraWidgetType) WidgetType {
	switch int(gpWidgetType) {
	case gpWidgetButton:
		return WidgetButton
	case gpWidgetDate:
		return WidgetDate
	case gpWidgetMenu:
		return WidgetMenu
	case gpWidgetRadio:
		return WidgetRadio
	case gpWidgetRange:
		return WidgetRange
	case gpWidgetSection:
		return WidgetSection
	case gpWidgetText:
		return WidgetText
	case gpWidgetToggle:
		return WidgetToggle
	case gpWidgetWindow:
		return WidgetWindow
	}
	panic("should not be here")
}
