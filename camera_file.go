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
// #include <string.h>
// #include <stdio.h>
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// DeleteFile tries to delete file from the camera, and returns error if it fails
func (camera *Camera) DeleteFile(path *CameraFilePath) error {
	fileDir := C.CString(path.Folder)
	defer C.free(unsafe.Pointer(fileDir))

	fileName := C.CString(path.Name)
	defer C.free(unsafe.Pointer(fileName))

	res := C.gp_camera_file_delete(camera.gpCamera, fileDir, fileName, camera.Ctx.gpContext)
	if res != GPOK {
		return newError("Cannot delete fine on camera", int(res))
	}
	return nil
}

// ListFiles returns a list of files and folders on the camera
func (camera *Camera) ListFiles() ([]CameraStorageInfo, error) {
	var gpCameraStorageInformation *C.CameraStorageInformation
	var storageCount C.int
	storageCount = 0
	returnedStorageInfo := []CameraStorageInfo{}

	res := C.gp_camera_get_storageinfo(camera.gpCamera, (**C.CameraStorageInformation)(unsafe.Pointer(&gpCameraStorageInformation)), &storageCount, camera.Ctx.gpContext)
	if res != GPOK {
		return nil, newError("Cannot get camera storage info", int(res))
	}
	defer C.free(unsafe.Pointer(gpCameraStorageInformation))

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(gpCameraStorageInformation)),
		Len:  int(storageCount),
		Cap:  int(storageCount),
	}
	nativeCameraFileSystemInfo := *(*[]C.CameraStorageInformation)(unsafe.Pointer(&hdr))
	for i := 0; i < int(storageCount); i++ {
		cameraStorage := CameraStorageInfo{
			Description: C.GoString((*C.char)(&nativeCameraFileSystemInfo[i].description[0])),
			Capacity:    uint64(nativeCameraFileSystemInfo[i].capacitykbytes),
			Free:        uint64(nativeCameraFileSystemInfo[i].freekbytes),
			FreeImages:  uint64(nativeCameraFileSystemInfo[i].freeimages),
			Children:    []CameraFilePath{},

			basedir: C.GoString((*C.char)(&nativeCameraFileSystemInfo[i].basedir[0])),
			camera:  camera,
		}

		if err := camera.recursiveListAllFiles(&cameraStorage.basedir, &cameraStorage.Children); err != nil {
			return nil, err
		}
		returnedStorageInfo = append(returnedStorageInfo, cameraStorage)
	}
	return returnedStorageInfo, nil
}

func (camera *Camera) recursiveListAllFiles(basedir *string, children *[]CameraFilePath) error {
	items, err := camera.findAllChildDirectories(basedir)
	if err != nil {
		return err
	}
	for _, dirName := range items {
		dirItem := CameraFilePath{
			Name:     dirName,
			Folder:   *basedir,
			Dir:      true,
			Children: []CameraFilePath{},
			camera:   camera,
		}
		childPath := *basedir + "/" + dirName
		if err := camera.recursiveListAllFiles(&childPath, &dirItem.Children); err != nil {
			return err
		}
		*children = append(*children, dirItem)

	}
	items, err = camera.findAllFilesInDir(basedir)
	if err != nil {
		return err
	}
	for _, fileName := range items {
		fileItem := CameraFilePath{
			Name:     fileName,
			Folder:   *basedir,
			Dir:      false,
			Children: nil,
			camera:   camera,
		}
		*children = append(*children, fileItem)
	}
	return nil
}

// Hmm, this could be reduced to one func, and a lambda passed as an arg
func (camera *Camera) findAllChildDirectories(basedirPath *string) ([]string, error) {
	var gpFileList *C.CameraList
	var err error
	returnedSlice := []string{}

	gpDirPath := C.CString(*basedirPath)
	defer C.free(unsafe.Pointer(gpDirPath))

	if gpFileList, err = newGphotoList(); err != nil {
		return nil, err
	}
	defer C.gp_list_free(gpFileList)

	if res := C.gp_camera_folder_list_folders(camera.gpCamera, gpDirPath, gpFileList, camera.Ctx.gpContext); res != GPOK {
		return nil, newError(fmt.Sprintf("Cannot get folder list from dir %s", *basedirPath), int(res))
	}

	listSize := int(C.gp_list_count(gpFileList))
	for i := 0; i < listSize; i++ {
		var gpListElementName *C.char
		C.gp_list_get_name(gpFileList, (C.int)(i), (**C.char)(&gpListElementName))
		returnedSlice = append(returnedSlice, C.GoString(gpListElementName))
	}
	return returnedSlice, nil
}

func (camera *Camera) findAllFilesInDir(basedirPath *string) ([]string, error) {
	var err error
	returnedSlice := []string{}

	gpDirPath := C.CString(*basedirPath)
	defer C.free(unsafe.Pointer(gpDirPath))

	gpFileList, err := newGphotoList()
	if err != nil {
		return nil, err
	}
	defer C.gp_list_free(gpFileList)

	if res := C.gp_camera_folder_list_files(camera.gpCamera, gpDirPath, gpFileList, camera.Ctx.gpContext); res != GPOK {
		return nil, newError(fmt.Sprintf("Cannot get file list from dir %s", *basedirPath), int(res))
	}

	listSize := int(C.gp_list_count(gpFileList))
	for i := 0; i < listSize; i++ {
		var gpListElementName *C.char
		C.gp_list_get_name(gpFileList, (C.int)(i), (**C.char)(&gpListElementName))
		returnedSlice = append(returnedSlice, C.GoString(gpListElementName))
	}
	return returnedSlice, nil
}
