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

/*
 #include <stdlib.h>
 #include <gphoto2/gphoto2.h>
*/
import "C"

import (
	"fmt"
	//	"unsafe"
)

type GphotoError struct {
	msg   string
	extra string
	Code  int
}

func (g *GphotoError) Error() string {
	if g.extra != "" {
		return fmt.Sprintf("%s (%s): %d", g.msg, g.extra, g.Code)
	}
	return fmt.Sprintf("%s: %d", g.msg, g.Code)
}

func newError(e string, code int) error {
	msg := e
	extra := ""

	if code >= libGphotoErrEnd && code <= libGphotoErrStart {
		m := C.gp_result_as_string(C.int(code))
		// defer C.free(unsafe.Pointer(m)) // should this one be freed or not?
		msg = C.GoString(m)
		if msg == "" {
			msg = e
		} else {
			extra = e
		}
	}
	return &GphotoError{msg: msg, Code: code, extra: extra}
}
