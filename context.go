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
import "C"

type Context struct {
	gpContext *C.GPContext
}

func (ctx *Context) free() {
	if ctx.gpContext != nil {
		C.gp_context_unref(ctx.gpContext)
		ctx = nil
	}
}

func NewContext() (*Context, error) {
	var gpContext *C.GPContext
	gpContext = C.gp_context_new()

	if gpContext == nil {
		return nil, newError("Could not initialize libgphoto2 context", Error)
	}
	return &Context{gpContext}, nil
}
