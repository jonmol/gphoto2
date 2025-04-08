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
	"fmt"
	"sync"
	"unsafe"
)

var (
	abilitiesLoaded     sync.Once
	cameraAbilitiesList *C.CameraAbilitiesList // this is a static list of what libgphoto supports
)

// PortInfo loads the list of all currently connected cameras and returns the C.GPPortInfo for the given port
func PortInfo(port string) (C.GPPortInfo, error) {
	var portInfoList *C.GPPortInfoList
	defer C.gp_port_info_list_free(portInfoList)

	if res := C.gp_port_info_list_new(&portInfoList); res != GPOK {
		return nil, newError("failed to call gp_port_info_list_new", int(res))
	}
	if res := C.gp_port_info_list_load(portInfoList); res != GPOK {
		return nil, newError("failed to call gp_port_info_list_new", int(res))
	}

	cPort := C.CString(port)
	defer C.free(unsafe.Pointer(cPort))

	var portInfo C.GPPortInfo
	port_index := C.gp_port_info_list_lookup_path(portInfoList, cPort)
	if port_index < 0 {
		return nil, newError("port index not found", int(port_index))
	}

	if res := C.gp_port_info_list_get_info(portInfoList, port_index, &portInfo); res != GPOK {
		return nil, newError("failed to call gp_port_info_list_new", int(res))
	}

	return portInfo, nil
}

// loadCameraAbilitiesList is loading the static list, so only needed once
func loadCameraAbilitiesList() {
	abilitiesLoaded.Do(func() {
		ctx, _ := NewContext()
		if res := C.gp_abilities_list_new(&cameraAbilitiesList); res == GPOK {
			if res := C.gp_abilities_list_load(cameraAbilitiesList, ctx.gpContext); res != GPOK {
				panic(newError("Failed call gp_abilities_list_load", int(res)))
			}
		} else {
			panic(newError("Failed call gp_abilities_list_new", int(res)))
		}
	})
}

// CameraAbilities returns the abilites for a camera
func CameraAbilities(name string) (C.CameraAbilities, error) {
	loadCameraAbilitiesList()
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))

	var abilities C.CameraAbilities
	abil_index := C.gp_abilities_list_lookup_model(cameraAbilitiesList, cstr)
	if abil_index < 0 {
		return abilities, newError(fmt.Sprintf("couldn't find model '%s'", name), -1)
	}
	C.gp_abilities_list_get_abilities(cameraAbilitiesList, abil_index, &abilities)

	return abilities, nil
}
