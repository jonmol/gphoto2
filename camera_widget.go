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

// #include <gphoto2/gphoto2.h>
// #include <stdio.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func (camera *Camera) getRootWidget() (*C.CameraWidget, error) {
	var rootWidget *C.CameraWidget

	if res := C.gp_camera_get_config(camera.gpCamera, (**C.CameraWidget)(unsafe.Pointer(&rootWidget)), camera.Ctx.gpContext); res != GPOK {
		return nil, newError("cannot initialize camera settings tree", int(res))
	}
	return rootWidget, nil
}

func (camera *Camera) getChildWidget(name *string) (*C.CameraWidget, error) {
	var rootWidget, childWidget *C.CameraWidget
	var err error
	if rootWidget, err = camera.getRootWidget(); err != nil {
		return nil, err
	}

	gpChildWidgetName := C.CString(*name)
	defer C.free(unsafe.Pointer(gpChildWidgetName))

	if res := C.gp_widget_get_child_by_name(rootWidget, gpChildWidgetName, (**C.CameraWidget)(unsafe.Pointer(&childWidget))); res != GPOK {
		return nil, newError(fmt.Sprintf("Could not retrieve child widget with name %s", *name), int(res))
	}
	return childWidget, nil
}

func (camera *Camera) freeChildWidget(input *C.CameraWidget) {
	var rootWidget *C.CameraWidget
	C.gp_widget_get_root(input, (**C.CameraWidget)(unsafe.Pointer(&rootWidget)))
	C.free(unsafe.Pointer(rootWidget))
}

func (camera *Camera) getWidgetInfo(input *C.CameraWidget, parent *CameraWidget) (*CameraWidget, error) {
	var gpInfo *C.char
	var gpLabel *C.char
	var gpName *C.char
	var gpWidgetType C.CameraWidgetType
	var child *C.CameraWidget
	var readonly C.int

	if res := C.gp_widget_get_info(input, (**C.char)(unsafe.Pointer(&gpInfo))); res != GPOK {
		return nil, newError("Failed to get info", int(res))
	}
	if res := C.gp_widget_get_label(input, (**C.char)(unsafe.Pointer(&gpLabel))); res != GPOK {
		return nil, newError("Failed to get label", int(res))
	}
	if res := C.gp_widget_get_name(input, (**C.char)(unsafe.Pointer(&gpName))); res != GPOK {
		return nil, newError("Failed to get name", int(res))
	}
	if res := C.gp_widget_get_type(input, (*C.CameraWidgetType)(unsafe.Pointer(&gpWidgetType))); res != GPOK {
		return nil, newError("Failed to get type", int(res))
	}
	if res := C.gp_widget_get_readonly(input, &readonly); res != GPOK {
		return nil, newError("Failed to get readonly", int(res))
	}

	widget := CameraWidget{
		widgetType:   widgetType(gpWidgetType),
		gpWidgetType: uint(gpWidgetType),
		label:        C.GoString(gpLabel),
		info:         C.GoString(gpInfo),
		name:         C.GoString(gpName),
		readonly:     (readonly == 1),
		camera:       camera,
		parent:       parent,
		children:     make([]*CameraWidget, 0, 0),
	}

	childrenCount := int(C.gp_widget_count_children(input))
	for n := 0; n < childrenCount; n++ {
		if res := C.gp_widget_get_child(input, C.int(n), (**C.CameraWidget)(unsafe.Pointer(&child))); res != GPOK {
			return nil, newError("Failed to get Child", int(res))
		}
		if childW, err := camera.getWidgetInfo(child, &widget); err != nil || childW == nil {
			return nil, err
		} else {
			widget.children = append(widget.children, childW)
		}
	}
	return &widget, nil
}
