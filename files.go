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

// #cgo LDFLAGS: -lgphoto2 -lgphoto2_port
// #include <gphoto2/gphoto2.h>
// #include <stdlib.h>
// #include <string.h>
// #include <stdio.h>
import "C"
import (
	"bytes"
	"io"
	"reflect"
	"unsafe"
)

//CamersStorageInfo is a struct describing one of the camera's storage spaces (SD or CF cards for example)
//Children is a directory tree present on the storage space
type CameraStorageInfo struct {
	Description string
	Capacity    uint64
	Free        uint64
	FreeImages  uint64
	Children    []CameraFilePath

	basedir string
	camera  *Camera
}

type cameraFilePathInternal struct {
	Name   [128]uint8
	Folder [1024]uint8
}

//CameraFilePath is a path to a file or dir on the camera file system
type CameraFilePath struct {
	Name     string
	Folder   string
	Dir      bool
	Children []CameraFilePath

	camera *Camera
}

//DownloadImage saves image pointed by path to the provided buffer. If leave on camera is set to false,the file will be deleted from the camera internal storage
func (file *CameraFilePath) DownloadImage(buffer io.Writer, leaveOnCamera bool) error {
	gpFile, err := newGpFile()
	if err != nil {
		return err
	}
	if file.camera.gpCamera == nil {
		return newError("Camera disconnected", Error)
	}
	defer C.gp_file_free(gpFile)

	fileDir := C.CString(file.Folder)
	defer C.free(unsafe.Pointer(fileDir))

	fileName := C.CString(file.Name)
	defer C.free(unsafe.Pointer(fileName))

	if res := C.gp_camera_file_get(file.camera.gpCamera, fileDir, fileName, FileTypeNormal, gpFile, file.camera.Ctx.gpContext); res != GPOK {
		return newError("Cannot download photo file", int(res))
	}

	err = copyFile(gpFile, buffer)
	if err == nil && !leaveOnCamera {
		C.gp_camera_file_delete(file.camera.gpCamera, fileDir, fileName, file.camera.Ctx.gpContext)
	}
	return err
}

func newCameraFilePath(input *cameraFilePathInternal, camera *Camera) *CameraFilePath {
	return &CameraFilePath{
		Name:     string(input.Name[:bytes.IndexByte(input.Name[:], 0)]),
		Folder:   string(input.Folder[:bytes.IndexByte(input.Folder[:], 0)]),
		Dir:      false,
		Children: nil,
		camera:   camera,
	}
}

func newGpFile() (*C.CameraFile, error) {
	var gpFile *C.CameraFile
	C.gp_file_new((**C.CameraFile)(unsafe.Pointer(&gpFile)))

	if gpFile == nil {
		return nil, newError("Cannot initialize camera file", Error)
	}
	return gpFile, nil
}

func copyFile(gpFileIn *C.CameraFile, bufferOut io.Writer) error {
	var fileData *C.char
	var fileLen C.ulong
	C.gp_file_get_data_and_size(gpFileIn, (**C.char)(unsafe.Pointer(&fileData)), &fileLen)

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(fileData)),
		Len:  int(fileLen),
		Cap:  int(fileLen),
	}
	goSlice := *(*[]byte)(unsafe.Pointer(&hdr))

	_, err := bufferOut.Write(goSlice)
	return err
}

func newGphotoList() (*C.CameraList, error) {
	var gpFileList *C.CameraList
	if res := C.gp_list_new((**C.CameraList)(unsafe.Pointer(&gpFileList))); res != GPOK {
		return nil, newError("Failed calling gp_list_new", int(res))
	}
	return gpFileList, nil
}
