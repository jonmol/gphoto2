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
// #include "callbacks.h"
import "C"
import (
	"fmt"
)

// ContextLogCallback defineds a function used to log info associated to lobgphoto2 context
type ContextLogCallback func(string)

// LogCallback defines a generic libgphoto2 logging function
type LogCallback func(int, string, string)

// ContextIdleCallback should take Context, interface{} as its arguments
type ContextIdleCallback func(Context)

// ContextInfoCallback is the function logging info logs from  libgphoto2 context.
// By default it logs everything to standard outout. You can assign your own method to this var
var ContextInfoCallback ContextLogCallback

// ContextErrorCallback is the function logging error logs from  libgphoto2 context.
// By default it logs everything to standard outout. You can assign your own method to this var
var ContextErrorCallback ContextLogCallback

// LoggerCallback is the libgphoto2 logging function. Currently there is no possibility to add multiple log function like it is possible in
// native C library implementation. Default implementation log everything to standard output with log level set to DEBUG
var LoggerCallback LogCallback

// IdleCallback is the libgphoto2 context idle function. It's being called when the library is waiting for the camera
// so small client side updates can be done when it's called
var IdleCallback ContextIdleCallback

func defaultLoggerCallback(debugLevel int, domain, data string) {
	if debugLevel < 2 {
		fmt.Println(fmt.Sprintf("LOGGING : level [%d] domain [%s] data [%s]", debugLevel, domain, data))
	}
}
func defaultInfoCallback(data string) {
	fmt.Println("INFO: " + data)
}

func defaultErrorCallback(data string) {
	fmt.Println("ERROR : " + data)
}

func defaultIdleCallback(c Context) {
	fmt.Println("Idle called")
}

//export wrapperInfoCallback
func wrapperInfoCallback(input *C.char) {
	if ContextInfoCallback != nil {
		ContextInfoCallback(C.GoString(input))
	}
}

//export wrapperErrorCallback
func wrapperErrorCallback(input *C.char) {
	if ContextErrorCallback != nil {
		ContextErrorCallback(C.GoString(input))
	}
}

//export wrapperLoggingCallback
func wrapperLoggingCallback(logLevel int, domain, data *C.char) {
	if LoggerCallback != nil {
		LoggerCallback(logLevel, C.GoString(domain), C.GoString(data))
	}
}

//export wrapperIdleCallback
func wrapperIdleCallback(context *C.GPContext) {
	if IdleCallback != nil {
		IdleCallback(Context{context})
	}
}

func SetIdleFunc(ctx *Context, fun ContextIdleCallback) {
	IdleCallback = fun
	C.gp_context_set_idle_func(ctx.gpContext, (*[0]byte)(C.ctx_idle_func), nil)
}

func SetLoggerFunc(fun *LogCallback) {
	if fun == nil {
		LoggerCallback = defaultLoggerCallback
	}
	C.gp_log_add_func(LogError, (*[0]byte)(C.logger_func), nil)
}

// make the general logging better. Make min log level settable somehow
func init() {
	ContextInfoCallback = defaultInfoCallback
	ContextErrorCallback = defaultErrorCallback
	LoggerCallback = defaultLoggerCallback
	IdleCallback = defaultIdleCallback

}
