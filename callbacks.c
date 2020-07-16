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

#include <gphoto2/gphoto2.h>
#include  <stdio.h>
extern void wrapperInfoCallback(char* p0);
extern void wrapperErrorCallback(char* p0);
extern void wrapperLoggingCallback(int  logLevel, char* domain, char* data);
extern void wrapperIdleCallback(GPContext *context);


void
ctx_error_func (GPContext *context, const char *str, void *data)
{
  wrapperErrorCallback((char*)str);
}

void
ctx_status_func (GPContext *context, const char *str, void *data)
{
  wrapperInfoCallback((char*)str);
}


void 
logger_func(GPLogLevel level, const char *domain, const char *str, void *data) 
{
  wrapperLoggingCallback((int)level, (char*) domain, (char*) str);
}

void
ctx_idle_func(GPContext *context, void *data) {
  wrapperIdleCallback(context);
}
