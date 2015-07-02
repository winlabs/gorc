/*
 * Copyright (c) 2014 MongoDB, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the license is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/winlabs/gowin32"
	"github.com/winlabs/gowin32/wrappers"

	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Resource struct {
	Type gowin32.ResourceType
	Id	 uint
	Data []byte
}

type stringFileInfoField struct {
	JsonName string
	WinName  string
}

var stringFileInfoFields = []stringFileInfoField{
	{JsonName: "comments",         WinName: "Comments"},
	{JsonName: "companyName",      WinName: "CompanyName"},
	{JsonName: "fileDescription",  WinName: "FileDescription"},
	{JsonName: "fileVersion",      WinName: "FileVersion"},
	{JsonName: "internalName",     WinName: "InternalName"},
	{JsonName: "legalCopyright",   WinName: "LegalCopyright"},
	{JsonName: "legalTrademarks",  WinName: "LegalTrademarks"},
	{JsonName: "originalFilename", WinName: "OriginalFilename"},
	{JsonName: "privateBuild",     WinName: "PrivateBuild"},
	{JsonName: "productName",      WinName: "ProductName"},
	{JsonName: "productVersion",   WinName: "ProductVersion"},
	{JsonName: "specialBuild",     WinName: "SpecialBuild"},
}

func parseFileFlags(fileFlagsObj interface{}) (uint32, error) {
	fileFlagsArray, ok := fileFlagsObj.([]interface{})
	if !ok {
		return 0, errors.New("field fileFlags must specify a list of strings")
	}
	var fileFlags uint32
	for _, flagNameObj := range fileFlagsArray {
		flagName, ok := flagNameObj.(string)
		if !ok {
			return 0, errors.New("field fileFlags must specify a list of strings")
		}
		switch flagName {
		case "VS_FF_DEBUG":
			fileFlags |= wrappers.VS_FF_DEBUG
		case "VS_FF_PRERELEASE":
			fileFlags |= wrappers.VS_FF_PRERELEASE
		case "VS_FF_PATCHED":
			fileFlags |= wrappers.VS_FF_PATCHED
		case "VS_FF_PRIVATEBUILD":
			fileFlags |= wrappers.VS_FF_PRIVATEBUILD
		case "VS_FF_INFOINFERRED":
			fileFlags |= wrappers.VS_FF_INFOINFERRED
		case "VS_FF_SPECIALBUILD":
			fileFlags |= wrappers.VS_FF_SPECIALBUILD
		default:
			return 0, errors.New(fmt.Sprintf("invalid file flag: %s", flagName))
		}
	}
	return fileFlags, nil
}

func parseFileOS(fileOSObj interface{}) (uint32, error) {
	fileOSName, ok := fileOSObj.(string)
	if !ok {
		return 0, errors.New("field fileOS must specify a string")
	}
	switch fileOSName {
	case "VOS_UNKNOWN":
		return wrappers.VOS_UNKNOWN, nil
	case "VOS_DOS":
		return wrappers.VOS_DOS, nil
	case "VOS_OS216":
		return wrappers.VOS_OS216, nil
	case "VOS_OS232":
		return wrappers.VOS_OS232, nil
	case "VOS_NT":
		return wrappers.VOS_NT, nil
	case "VOS__WINDOWS16":
		return wrappers.VOS__WINDOWS16, nil
	case "VOS__PM16":
		return wrappers.VOS__PM16, nil
	case "VOS__PM32":
		return wrappers.VOS__PM32, nil
	case "VOS__WINDOWS32":
		return wrappers.VOS__WINDOWS32, nil
	case "VOS_DOS_WINDOWS16":
		return wrappers.VOS_DOS_WINDOWS16, nil
	case "VOS_DOS_WINDOWS32":
		return wrappers.VOS_DOS_WINDOWS32, nil
	case "VOS_OS216_PM16":
		return wrappers.VOS_OS216_PM16, nil
	case "VOS_OS232_PM32":
		return wrappers.VOS_OS232_PM32, nil
	case "VOS_NT_WINDOWS32":
		return wrappers.VOS_NT_WINDOWS32, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid file OS: %s", fileOSName))
	}
}

func parseFileType(fileTypeObj interface{}) (uint32, error) {
	fileTypeName, ok := fileTypeObj.(string)
	if !ok {
		return 0, errors.New("field fileType must specify a string")
	}
	switch fileTypeName {
	case "VFT_UNKNOWN":
		return wrappers.VFT_UNKNOWN, nil
	case "VFT_APP":
		return wrappers.VFT_APP, nil
	case "VFT_DLL":
		return wrappers.VFT_DLL, nil
	case "VFT_DRV":
		return wrappers.VFT_DRV, nil
	case "VFT_FONT":
		return wrappers.VFT_FONT, nil
	case "VFT_VXD":
		return wrappers.VFT_VXD, nil
	case "VFT_STATIC_LIB":
		return wrappers.VFT_STATIC_LIB, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid file type: %s", fileTypeName))
	}
}

func parseFileSubtype(fileSubtypeObj interface{}) (uint32, error) {
	fileSubtypeName, ok := fileSubtypeObj.(string)
	if !ok {
		return 0, errors.New("field fileSubtype must specify a string")
	}
	switch fileSubtypeName {
	case "VFT2_UNKNOWN":
		return wrappers.VFT2_UNKNOWN, nil
	case "VFT2_DRV_PRINTER":
		return wrappers.VFT2_DRV_PRINTER, nil
	case "VFT2_DRV_KEYBOARD":
		return wrappers.VFT2_DRV_KEYBOARD, nil
	case "VFT2_DRV_LANGUAGE":
		return wrappers.VFT2_DRV_LANGUAGE, nil
	case "VFT2_DRV_DISPLAY":
		return wrappers.VFT2_DRV_DISPLAY, nil
	case "VFT2_DRV_MOUSE":
		return wrappers.VFT2_DRV_MOUSE, nil
	case "VFT2_DRV_NETWORK":
		return wrappers.VFT2_DRV_NETWORK, nil
	case "VFT2_DRV_SYSTEM":
		return wrappers.VFT2_DRV_SYSTEM, nil
	case "VFT2_DRV_INSTALLABLE":
		return wrappers.VFT2_DRV_INSTALLABLE, nil
	case "VFT2_DRV_SOUND":
		return wrappers.VFT2_DRV_SOUND, nil
	case "VFT2_DRV_COMM":
		return wrappers.VFT2_DRV_COMM, nil
	case "VFT2_DRV_VERSIONED_PRINTER":
		return wrappers.VFT2_DRV_VERSIONED_PRINTER, nil
	case "VFT2_FONT_RASTER":
		return wrappers.VFT2_FONT_RASTER, nil
	case "VFT2_FONT_VECTOR":
		return wrappers.VFT2_FONT_VECTOR, nil
	case "VFT2_FONT_TRUETYPE":
		return wrappers.VFT2_FONT_TRUETYPE, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid file subtype: %s", fileSubtypeName))
	}
}

func parseVersionResource(versionJson map[string]interface{}, language gowin32.Language) (*Resource, error) {
	fixedFileInfo := wrappers.VS_FIXEDFILEINFO{
		Signature:     0xFEEF04BD,
		FileFlagsMask: 0x0000003F,
	}
	if fileVersionObj, ok := versionJson["fileVersion"]; ok {
		if fileVersionStr, ok := fileVersionObj.(string); ok {
			if fileVersionNumber, err := gowin32.StringToFileVersionNumber(fileVersionStr); err != nil {
				return nil, errors.New(fmt.Sprintf("invalid version number: %s", fileVersionStr))
			} else {
				fixedFileInfo.FileVersionMS = wrappers.MAKELONG(
					uint16(fileVersionNumber.Minor),
					uint16(fileVersionNumber.Major))
				fixedFileInfo.FileVersionLS = wrappers.MAKELONG(
					uint16(fileVersionNumber.Revision),
					uint16(fileVersionNumber.Build))
			}
		} else {
			return nil, errors.New("field fileVersion must specify a string")
		}
	}
	if productVersionObj, ok := versionJson["productVersion"]; ok {
		if productVersionStr, ok := productVersionObj.(string); ok {
			if productVersionNumber, err := gowin32.StringToFileVersionNumber(productVersionStr); err != nil {
				return nil, errors.New(fmt.Sprintf("invalid version number: %s", productVersionStr))
			} else {
				fixedFileInfo.ProductVersionMS = wrappers.MAKELONG(
					uint16(productVersionNumber.Minor),
					uint16(productVersionNumber.Major))
				fixedFileInfo.ProductVersionLS = wrappers.MAKELONG(
					uint16(productVersionNumber.Revision),
					uint16(productVersionNumber.Build))
			}
		}
	}
	if fileFlagsObj, ok := versionJson["fileFlags"]; ok {
		if fileFlags, err := parseFileFlags(fileFlagsObj); err != nil {
			return nil, err
		} else {
			fixedFileInfo.FileFlags = fileFlags
		}
	}
	if fileOSObj, ok := versionJson["fileOS"]; ok {
		if fileOS, err := parseFileOS(fileOSObj); err != nil {
			return nil, err
		} else {
			fixedFileInfo.FileOS = fileOS
		}
	}
	if fileTypeObj, ok := versionJson["fileType"]; ok {
		if fileType, err := parseFileType(fileTypeObj); err != nil {
			return nil, err
		} else {
			fixedFileInfo.FileType = fileType
		}
	}
	if fileSubtypeObj, ok := versionJson["fileSubtype"]; ok {
		if fileSubtype, err := parseFileSubtype(fileSubtypeObj); err != nil {
			return nil, err
		} else {
			fixedFileInfo.FileSubtype = fileSubtype
		}
	}
	stringFileInfo := make([]VersionString, 0, 10)
	if stringFileInfoObj, ok := versionJson["stringFileInfo"]; ok {
		if stringFileInfoJson, ok := stringFileInfoObj.(map[string]interface{}); ok {
			for _, field := range stringFileInfoFields {
				if fieldValObj, ok := stringFileInfoJson[field.JsonName]; ok {
					if fieldVal, ok := fieldValObj.(string); ok {
						stringFileInfo = append(stringFileInfo, VersionString{
							Key:   field.WinName,
							Value: fieldVal,
						})
					} else {
						return nil, errors.New(fmt.Sprintf("field %s must specify a string"))
					}
				}
			}
		} else {
			return nil, errors.New("field stringFileInfo must specify an object")
		}
	}
	return &Resource{
		Type: gowin32.ResourceTypeVersion,
		Id:   1,
		Data: EncodeVersionInfo(&fixedFileInfo, language, 1200, stringFileInfo),
	}, nil
}

func parseMessageSeverity(severityObj interface{}) (uint32, error) {
	severityName, ok := severityObj.(string)
	if !ok {
		return 0, errors.New("field severity must specify a string")
	}
	switch severityName {
	case "Success":
		return 0x00000000, nil
	case "Informational":
		return 0x40000000, nil
	case "Warning":
		return 0x80000000, nil
	case "Error":
		return 0xC0000000, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid severity: %s", severityName))
	}
}

func parseMessageTableResource(messageTableJson []interface{}) (*Resource, error) {
	messages := make(map[uint32]string)
	for _, messageObj := range messageTableJson {
		messageJson, ok := messageObj.(map[string]interface{})
		if !ok {
			return nil, errors.New("field messageTable must specify a list of objects")
		}
		idObj, ok := messageJson["id"]
		if !ok {
			return nil, errors.New("field id is required")
		}
		idFloat, ok := idObj.(float64)
		if !ok {
			return nil, errors.New("field id must specify an integer")
		}
		id := uint32(idFloat)
		severityObj, ok := messageJson["severity"]
		if !ok {
			return nil, errors.New("field severity is required")
		}
		severity, err := parseMessageSeverity(severityObj)
		if err != nil {
			return nil, err
		}
		id |= severity
		if _, ok := messages[id]; ok {
			return nil, errors.New(fmt.Sprintf("field duplicate message with ID %x", id))
		}
		messageTextObj, ok := messageJson["messageText"]
		if !ok {
			return nil, errors.New("field messageText is required")
		}
		messageText, ok := messageTextObj.(string)
		if !ok {
			return nil, errors.New("field messageText must specify a string")
		}
		messages[id] = messageText
	}
	return &Resource{
		Type: gowin32.ResourceTypeMessageTable,
		Id:   1,
		Data: EncodeMessageTable(messages),
	}, nil
}

func loadManifestResource(manifestFileName string) (*Resource, error) {
	f, err := os.Open(manifestFileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not open file '%s'", manifestFileName))
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not read file '%s'", manifestFileName))
	}
	return &Resource{
		Type: gowin32.ResourceTypeManifest,
		Id:   uint(wrappers.CREATEPROCESS_MANIFEST_RESOURCE_ID),
		Data: data,
	}, nil
}

func ParseResources(jsonData map[string]interface{}, sourceDir string) (gowin32.Language, []*Resource, error) {
	locale := gowin32.LocaleNeutral
	if languageObj, ok := jsonData["language"]; ok {
		if languageName, ok := languageObj.(string); ok {
			if localeId, err := gowin32.LocaleFromLocaleName(languageName, 0); err != nil {
				return 0, nil, errors.New(fmt.Sprintf("invalid language %s", languageName))
			} else {
				locale = localeId
			}
		} else {
			return 0, nil, errors.New("field language must specify a string")
		}
	}
	resources := make([]*Resource, 0)
	for key, value := range jsonData {
		switch key {
		case "version":
			if versionJson, ok := value.(map[string]interface{}); ok {
				if versionRes, err := parseVersionResource(versionJson, locale.Language()); err != nil {
					return 0, nil, err
				} else {
					resources = append(resources, versionRes)
				}
			} else {
				return 0, nil, errors.New("field version must specify an object")
			}
		case "messageTable":
			if messageJson, ok := value.([]interface{}); ok {
				if messageRes, err := parseMessageTableResource(messageJson); err != nil {
					return 0, nil, err
				} else {
					resources = append(resources, messageRes)
				}
			} else {
				return 0, nil, errors.New("field messageTable must specify a list of objects")
			}
		case "manifest":
			if manifestFileName, ok := value.(string); ok {
				if manifestRes, err := loadManifestResource(filepath.Join(sourceDir, manifestFileName)); err != nil {
					return 0, nil, err
				} else {
					resources = append(resources, manifestRes)
				}
			} else {
				return 0, nil, errors.New("field manifest must specify a file name")
			}
		case "language":
			// handled above
		default:
			return 0, nil, errors.New(fmt.Sprintf("invalid resource type %s", key))
		}
	}
	return locale.Language(), resources, nil
}
