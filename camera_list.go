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
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

type CameraList []CameraListInfo

type CameraListInfo struct {
	name string
	port string
}

func (c CameraListInfo) Name() string {
	return c.name
}

func (c CameraListInfo) Port() string {
	return c.port
}

func (c CameraListInfo) Abilities() (C.CameraAbilities, error) {
	return CameraAbilities(c.name)
}

func (c CameraListInfo) PortInfo() (C.GPPortInfo, error) {
	return PortInfo(c.port)
}

func (c CameraListInfo) Camera() (*Camera, error) {
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

	if a, err := c.Abilities(); err != nil {
		return nil, err
	} else {
		C.gp_camera_set_abilities(gpCamera, a)
	}

	if p, err := c.PortInfo(); err != nil {
		return nil, err
	} else {
		C.gp_camera_set_port_info(gpCamera, p)
	}

	return &Camera{gpCamera: gpCamera, Ctx: ctx}, nil

}

// ListCameras lists all connected cameras
func ListCameras() (CameraList, error) {
	ctx, err := NewContext()
	if err != nil {
		return nil, err
	}
	list := make(CameraList, 0)
	var cameraList *C.CameraList
	C.gp_list_new(&cameraList)
	defer C.free(unsafe.Pointer(cameraList))
	C.gp_camera_autodetect(cameraList, ctx.gpContext)
	size := int(C.gp_list_count(cameraList))

	if size < 0 {
		return nil, newError("Cannot get camera list", size)
	}
	for i := 0; i < size; i++ {
		var cName *C.char
		var cPort *C.char
		C.gp_list_get_name(cameraList, C.int(i), &cName)
		C.gp_list_get_value(cameraList, C.int(i), &cPort)
		defer C.free(unsafe.Pointer(cName))
		defer C.free(unsafe.Pointer(cPort))

		list = append(list, CameraListInfo{name: C.GoString(cName), port: C.GoString(cPort)})
	}
	return list, nil
}
