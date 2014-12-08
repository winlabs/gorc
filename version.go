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

	"fmt"
	"syscall"
	"unsafe"
)

type VersionString struct {
	Key   string
	Value string
}

// The following structures define the file format for Win32 version resources.
// They are documented on MSDN but not included in any Win32 header file because of their variable size.

type vsVersionInfo struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [16]uint16
	Padding     uint16
	Value       wrappers.VS_FIXEDFILEINFO
}

type vsStringFileInfo struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [15]uint16
}

type vsStringTable struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [9]uint16
}

type vsString struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
}

type vsVarFileInfo struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [12]uint16
	Padding     uint16
}

type vsVar struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [12]uint16
	Padding     uint16
	Value       uint32
}

func stringToUTF16Bytes(text string) []byte {
	data := make([]byte, 2*len(text))
	wrappers.RtlMoveMemory(&data[0], (*byte)(unsafe.Pointer(syscall.StringToUTF16Ptr(text))), uintptr(len(data)))
	return data
}

func encodeString(key string, value string) []byte {
	var paddingBytes1, paddingBytes2 uint16
	if len(key) % 2 != 0 {
		paddingBytes1 = 2
	}
	if len(value) % 2 == 0 {
		paddingBytes2 = 2
	}
	var info vsString
	info.Length = uint16(unsafe.Sizeof(info)) + uint16(len(key)*2) + 2 + paddingBytes1 + uint16(len(value)*2) + 2
	info.ValueLength = uint16(len(value)) + 1
	info.Type = 1
	data := make([]byte, unsafe.Sizeof(info))
	wrappers.RtlMoveMemory(&data[0], (*byte)(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	data = append(data, stringToUTF16Bytes(key)...)
	data = append(data, make([]byte, 2 + paddingBytes1)...)
	data = append(data, stringToUTF16Bytes(value)...)
	return append(data, make([]byte, 2 + paddingBytes2)...)
}

func encodeStringTable(language gowin32.Language, codePage uint32, stringInfo []VersionString) []byte {
	extraData := make([]byte, 0, 1024)
	for _, pair := range stringInfo {
		extraData = append(extraData, encodeString(pair.Key, pair.Value)...)
	}
	var info vsStringTable
	info.Length = uint16(unsafe.Sizeof(info)) + uint16(len(extraData))
	info.ValueLength = 0
	info.Type = 1
	copy(info.Key[:], syscall.StringToUTF16(fmt.Sprintf("%04X%04X", uint16(language), uint16(codePage))))
	data := make([]byte, unsafe.Sizeof(info))
	wrappers.RtlMoveMemory(&data[0], (*byte)(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	return append(data, extraData...)
}

func encodeStringFileInfo(language gowin32.Language, codePage uint32, stringInfo []VersionString) []byte {
	extraData := encodeStringTable(language, codePage, stringInfo)
	var info vsStringFileInfo
	info.Length = uint16(unsafe.Sizeof(info)) + uint16(len(extraData))
	info.ValueLength = 0
	info.Type = 1
	copy(info.Key[:], syscall.StringToUTF16("StringFileInfo"))
	data := make([]byte, unsafe.Sizeof(info))
	wrappers.RtlMoveMemory(&data[0], (*byte)(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	return append(data, extraData...)
}

func encodeTranslation(language gowin32.Language, codePage uint32) []byte {
	var info vsVar
	info.Length = uint16(unsafe.Sizeof(info))
	info.ValueLength = uint16(unsafe.Sizeof(info.Value))
	info.Type = 0
	copy(info.Key[:], syscall.StringToUTF16("Translation"))
	info.Value = wrappers.MAKELONG(uint16(language), uint16(codePage))
	data := make([]byte, unsafe.Sizeof(info))
	wrappers.RtlMoveMemory(&data[0], (*byte)(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	return data
}

func encodeVarFileInfo(language gowin32.Language, codePage uint32) []byte {
	extraData := encodeTranslation(language, codePage)
	var info vsVarFileInfo
	info.Length = uint16(unsafe.Sizeof(info)) + uint16(len(extraData))
	info.ValueLength = 0
	info.Type = 1
	copy(info.Key[:], syscall.StringToUTF16("VarFileInfo"))
	data := make([]byte, unsafe.Sizeof(info))
	wrappers.RtlMoveMemory(&data[0], (*byte)(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	return append(data, extraData...)
}

func EncodeVersionInfo(fixedInfo *wrappers.VS_FIXEDFILEINFO, language gowin32.Language, codePage uint32, stringInfo []VersionString) []byte {
	stringData := encodeStringFileInfo(language, codePage, stringInfo)
	varData := encodeVarFileInfo(language, codePage)
	var info vsVersionInfo
	info.Length = uint16(unsafe.Sizeof(info)) + uint16(len(stringData)) + uint16(len(varData))
	info.ValueLength = uint16(unsafe.Sizeof(info.Value))
	info.Type = 0
	copy(info.Key[:], syscall.StringToUTF16("VS_VERSION_INFO"))
	info.Value = *fixedInfo
	data := make([]byte, unsafe.Sizeof(info))
	wrappers.RtlMoveMemory(&data[0], (*byte)(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	data = append(data, stringData...)
	return append(data, varData...)
}
