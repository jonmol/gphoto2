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

const (
	libGphotoErrStart = -1
	libGphotoErrEnd   = -200
)

// libgphoto2 error codes. See http://www.gphoto.org/doc/api/gphoto2-result_8h.html for errors < -100 and http://www.gphoto.org/doc/api/gphoto2-port-result_8h.html for errors <= 0 >= -100
const (
	//GPOK means no error
	GPOK = 0
	//Error is a Generic Error
	Error = -1
	//ErrorBadParameters : Bad parameters passed
	ErrorBadParameters = -2
	//ErrorNoMemory : Out of memory
	ErrorNoMemory = -3
	//ErrorLibrary : Error in the camera driver
	ErrorLibrary = -4
	//ErrorUnknownPort : Unknown libgphoto2 port passed
	ErrorUnknownPort = -5
	//ErrorNotSupported : Functionality not supported
	ErrorNotSupported = -6
	//ErrorIO : Generic I/O error
	ErrorIO = -7
	//ErrorFixedLimitExceeded : Buffer overflow of internal structure
	ErrorFixedLimitExceeded = -8
	//ErrorTimeout : Operation timed out
	ErrorTimeout = -10
	//ErrorIOSupportedSerial : Serial ports not supported
	ErrorIOSupportedSerial = -20
	//ErrorIOSupportedUsb : USB ports not supported
	ErrorIOSupportedUsb = -21
	//ErrorIOInit : Error initialising I/O
	ErrorIOInit = -31
	//ErrorIORead : I/O during read
	ErrorIORead = -34
	//ErrorIOWrite : I/O during write
	ErrorIOWrite = -35
	//ErrorIOUpdate : I/O during update of settings
	ErrorIOUpdate = -37
	//ErrorIOSerialSpeed : Specified serial speed not possible.
	ErrorIOSerialSpeed = -41
	//ErrorIOUSBClearHalt : Error during USB Clear HALT
	ErrorIOUSBClearHalt = -51
	//ErrorIOUSBFind : Error when trying to find USB device
	ErrorIOUSBFind = -52
	//ErrorIOUSBClaim : Error when trying to claim the USB device
	ErrorIOUSBClaim = -53
	//ErrorIOLock : Error when trying to lock the device
	ErrorIOLock = -60
	//ErrorHal : Unspecified error when talking to HAL
	ErrorHal = -70

	//ErrorCorruptedData : This error is reported by camera drivers if corrupted data has been received that
	//can not be automatically handled. Normally, drivers will do everything possible to automatically recover from this error
	ErrorCorruptedData = -102
	//ErrorFileExists : An operation failed because a file existed. This error is reported for example when the user tries to create a file that already exists.
	ErrorFileExists = -103
	//ErrorModelNotFound : The specified model could not be found. This error is reported when the user specified a model that does not seem to be supported by any driver.
	ErrorModelNotFound = -105
	//ErrorDirectoryNotFound : The specified directory could not be found. This error is reported when the user specified a directory that is non-existent.
	ErrorDirectoryNotFound = -107
	//ErrorFileNotFound : The specified file could not be found. This error is reported when the user wants to access a file that is non-existent.
	ErrorFileNotFound = -108
	//ErrorDirectoryExists : The specified directory already exists. This error is reported for example when the user wants to create a directory that already exists.
	ErrorDirectoryExists = -109
	//ErrorCameraBusy : Camera I/O or a command is in progress.
	ErrorCameraBusy = -110
	//ErrorPathNotAbsolute : The specified path is not absolute. This error is reported when the user specifies paths that are not absolute, i.e. paths like "path/to/directory". As a rule of thumb, in gphoto2, there is nothing like relative paths.
	ErrorPathNotAbsolute = -111
	//ErrorCancel : A cancellation requestion by the frontend via progress callback and GP_CONTEXT_FEEDBACK_CANCEL was successful and the transfer has been aborted.
	ErrorCancel = -112
	//ErrorCameraError : The camera reported some kind of error. This can be either a photographic error, such as failure to autofocus, underexposure, or violating storage permission, anything else that stops the camera from performing the operation.
	ErrorCameraError = -113
	//ErrorosFailure : There was some sort of OS error in communicating with the camera,
	// e.g. lack of permission for an operation.
	ErrorosFailure = -114
	//ErrorNoSpace : There was not enough free space when uploading a file.
	ErrorNoSpace = -115

	ErrorReadOnly             = -201
	ErrorWidgetHasNoOptions   = -202
	ErrorWidgetIllegalOption  = -203
	ErrorWidgetNotImplemented = -204
)

//Log level
const (
	//LogError : Log message is an error infomation
	LogError = iota
	//LogVerbose : Log message is an verbose debug infomation
	LogVerbose
	//LogDebug : Log message is an debug infomation
	LogDebug
	//LogData : Log message is a data hex dump
	LogData
)

//File types
const (
	//FileTypePreview is a preview of an image
	FileTypePreview = iota
	//FileTypeNormal is regular normal data of a file
	FileTypeNormal
	//FileTypeRaw usually the same as FileTypeNormal for modern cameras ( left for compatibility purposes)
	FileTypeRaw
	//FileTypeAudio is a audio view of a file. Perhaps an embedded comment or similar
	FileTypeAudio
	//FileTypeExif is the  embedded EXIF data of an image
	FileTypeExif
	//FileTypeMetadata is the metadata of a file, like Metadata of files on MTP devices
	FileTypeMetadata
)

//widget types
const (
	gpWidgetWindow = iota //(0)
	gpWidgetSection
	gpWidgetText
	gpWidgetRange
	gpWidgetToggle
	gpWidgetRadio
	gpWidgetMenu
	gpWidgetButton
	gpWidgetDate
)

//widget types
const (
	//WidgetWindow is the toplevel configuration widget. It should likely contain multiple #WidgetSection entries.
	WidgetWindow WidgetType = "window"
	//WidgetSection : Section widget (think Tab)
	WidgetSection WidgetType = "section"
	//WidgetText : Text widget (string)
	WidgetText WidgetType = "text"
	//WidgetRange : Slider widget (float)
	WidgetRange WidgetType = "range"
	//WidgetToggle : Toggle widget (think check box) (int)
	WidgetToggle WidgetType = "toggle"
	//WidgetRadio : Radio button widget (string)
	WidgetRadio WidgetType = "radio"
	//WidgetMenu : Menu widget (same as RADIO) (string)
	WidgetMenu WidgetType = "menu"
	//WidgetButton : Button press widget ( CameraWidgetCallback )
	WidgetButton WidgetType = "button"
	//WidgetDate : Date entering widget (int)
	WidgetDate WidgetType = "date"
)

const (
	SettingProgram      = "expprogram"
	SettingFocusMode    = "drivemode"
	SettingAperture     = "aperture"
	SettingFocalLength  = "focallength"
	Settingshutterspeed = "shutterspeed"
	SettingISO          = "iso"
	SettingWB           = "whitebalance"
	SettingQuality      = "imagequality"
	SettingSize         = "size"

	SettingLiveViewZoomRatio = "d1a3"
)

type settingMapping struct {
	primaryName   string
	secondaryName string
}

// CommonSettings contains common settings, some have different names since camera makers don't believe in standards
var CommonSettings = map[string]settingMapping{
	SettingProgram:      {primaryName: SettingProgram},
	SettingFocusMode:    {primaryName: "focusmode", secondaryName: SettingFocusMode},
	SettingAperture:     {primaryName: "f-number", secondaryName: SettingAperture},
	SettingFocalLength:  {primaryName: SettingFocalLength},
	Settingshutterspeed: {primaryName: "shutterspeed2", secondaryName: Settingshutterspeed},
	SettingISO:          {primaryName: SettingISO},
	SettingWB:           {primaryName: SettingWB},
	SettingQuality:      {primaryName: SettingQuality},
	SettingSize:         {primaryName: SettingSize, secondaryName: "eoszoom"},
}
