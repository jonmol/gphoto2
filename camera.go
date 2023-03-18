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
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

//Camera struct represents a camera connected to the computer
type Camera struct {
	gpCamera *C.Camera
	Ctx      *Context
	Settings *CameraWidget
}

//Exit Closes a connection to the camera and therefore gives other application
//the possibility to access the camera, too. It is recommended that you call
//this function when you currently don't need the camera. The camera will get
//reinitialized by gp_camera_init() automatically if you try to access the camera again.
func (c Camera) Exit() error {
	if c.gpCamera != nil {
		if res := C.gp_camera_exit(c.gpCamera, c.Ctx.gpContext); res != GPOK {
			return newError("", int(res))
		}
	}
	return nil
}

func (c Camera) Free() error {
	if err := c.Exit(); err != nil {
		return err
	}
	if res := C.gp_camera_unref(c.gpCamera); res != GPOK {
		return newError("", int(res))
	}
	c.Ctx.free()
	return nil
}

// ResetCamera resets the camera port, can be needed at times
// https://github.com/gphoto/gphoto2/blob/7a48ea37832bcd19e17b80afef2f7f2d426419f3/gphoto2/main.c#L1675
func (c *Camera) Reset() error {
	if err := c.Exit(); err != nil {
		return err
	}

	var port *C.GPPort
	var info C.GPPortInfo
	if res := C.gp_port_new(&port); res != GPOK {
		return newError("", int(res))
	}
	if res := C.gp_camera_get_port_info(c.gpCamera, &info); res != GPOK {
		return newError("", int(res))
	}
	if res := C.gp_port_set_info(port, info); res != GPOK {
		return newError("", int(res))
	}
	if res := C.gp_port_open(port); res != GPOK {
		return newError("", int(res))
	}
	if res := C.gp_port_reset(port); res != GPOK {
		return newError("", int(res))
	}
	if res := C.gp_port_close(port); res != GPOK {
		return newError("", int(res))
	}
	if res := C.gp_port_free(port); res != GPOK {
		return newError("", int(res))
	}
	return nil
}

// NewCamera tries to connect to a camera with name name, if name is empty it tries with the first connected camera. It returns a new Camera struct.
func NewCamera(name string) (*Camera, error) {
	ctx, err := NewContext()
	if err != nil {
		return nil, err
	}
	var gpCamera *C.Camera

	if res := C.gp_camera_new((**C.Camera)(unsafe.Pointer(&gpCamera))); res != GPOK {
		return nil, newError("Cannot initialize camera pointer", int(res))
	} else if gpCamera == nil {
		return nil, newError("Cannot initialize camera pointer", Error)
	}

	if res := C.gp_camera_init(gpCamera, ctx.gpContext); res != GPOK {
		C.gp_camera_exit(gpCamera, ctx.gpContext)
		C.gp_camera_unref(gpCamera)
		ctx.free()
		return nil, newError("", int(res))
	}

	return &Camera{gpCamera: gpCamera, Ctx: ctx}, nil
}

func ListCameras() ([]string, error) {
	//ctx, err := NewContext()
	names := make([]string, 0)
	var cameraList *C.CameraList
	C.gp_list_new(&cameraList)
	defer C.free(unsafe.Pointer(cameraList))

	size := int(C.gp_list_count(cameraList))

	if size < 0 {
		return nil, newError("Cannot get camera list", size)
	}

	for i := 0; i < size; i++ {
		var cKey *C.char

		C.gp_list_get_name(cameraList, C.int(i), &cKey)
		defer C.free(unsafe.Pointer(cKey))

		names = append(names, C.GoString(cKey))

	}
	return names, nil
}
