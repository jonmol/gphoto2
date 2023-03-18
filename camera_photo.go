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
// #include "callbacks.h"
// #include <stdlib.h>
// #include <stdio.h>
import "C"
import (
	"io"
	"unsafe"
)

// TODO: Can't handle when the camera is making multiple copies of the image. IE when
// raw+jpeg. gphoto2 can handle that so it's clearly doable but I don't have that use case
// if anyone needs it at some point it should be implemented

// CaptureDownload is for the lazy, capture and download in one call
func (c *Camera) CaptureDownload(buffer io.Writer, leaveOnCamera bool) error {
	file, err := c.CaptureImage()
	if err != nil {
		return err
	}
	return file.DownloadImage(buffer, leaveOnCamera)
}

//CaptureImage captures image with current setings into camera's internal storage
//call CameraFilePath.DownloadImage to
func (c *Camera) CaptureImage() (*CameraFilePath, error) {
	cameraPath := cameraFilePathInternal{}
	if res := C.gp_camera_capture(c.gpCamera, C.GP_CAPTURE_IMAGE, (*C.CameraFilePath)(unsafe.Pointer(&cameraPath)), c.Ctx.gpContext); res != GPOK {
		return nil, newError("Cannot capture photo", int(res))
	}
	return newCameraFilePath(&cameraPath, c), nil
}

//CapturePreview  captures image preview and saves it in provided buffer
func (c *Camera) CapturePreview(buffer io.Writer) error {
	gpFile, err := newGpFile()
	if err != nil {
		return err
	}
	if res := C.gp_camera_capture_preview(c.gpCamera, gpFile, c.Ctx.gpContext); res != GPOK {
		return newError("Cannot capture preview", int(res))

	}
	err = copyFile(gpFile, buffer)
	C.gp_file_free(gpFile)
	return err
}
