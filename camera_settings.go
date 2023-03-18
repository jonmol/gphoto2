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
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include <gphoto2/gphoto2-setting.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

func (c *Camera) LoadWidgets() error {
	var rootWidget *C.CameraWidget
	var err error

	if rootWidget, err = c.getRootWidget(); err != nil {
		return err
	}
	defer C.free(unsafe.Pointer(rootWidget))

	if c.Settings, err = c.getWidgetInfo(rootWidget, nil); err != nil {
		return err
	}
	return nil
}

// GetSetting returns a CameraWidget if it can be found among the camera settings
//
func (c *Camera) GetSetting(name string) (*CameraWidget, error) {
	if c.Settings == nil {
		if err := c.LoadWidgets(); err != nil {
			return nil, err
		}
	}

	if com, ok := CommonSettings[name]; ok {
		if res := c.Settings.Find(com.primaryName); res != nil {
			if res == nil {
				return nil, nil
			}
			return res, nil
		} else if com.secondaryName != "" {
			if res := c.Settings.Find(com.secondaryName); res != nil {
				return res, nil
			}
		}
	} else {
		return c.Settings.Find(name), nil
	}
	return nil, nil
}
