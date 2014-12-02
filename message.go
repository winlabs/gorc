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
	"github.com/winlabs/gowin32/wrappers"

	"sort"
	"syscall"
	"unsafe"
)

type idSlice []uint32

func (ids idSlice) Len() int {
	return len(ids)
}

func (ids idSlice) Less(i, j int) bool {
	return ids[i] < ids[j]
}

func (ids idSlice) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}

func EncodeMessageTable(messages map[uint32]string) []byte {
	ids := make(idSlice, 0, len(messages))
	for id := range messages {
		ids = append(ids, id)
	}
	sort.Sort(ids)

	// partition the ID list into blocks of consecutive numbers
	blockLengths := []int{}
	var curBlockLength int
	var idLast uint32
	for _, id := range ids {
		if curBlockLength == 0 || id != idLast+1 {
			if curBlockLength > 0 {
				blockLengths = append(blockLengths, curBlockLength)
			}
			curBlockLength = 0
		}
		idLast = id
		curBlockLength++
	}
	if curBlockLength > 0 {
		blockLengths = append(blockLengths, curBlockLength)
	}

	// build an entry data structure for each message
	entries := make([][]byte, 0, len(ids))
	for _, id := range ids {
		messageText := messages[id] + "\u000D\u000A"
		textByteLength := 2 * len(messageText)
		var nullByteCount uint16
		if textByteLength%4 == 0 {
			nullByteCount = 4
		} else {
			nullByteCount = 2
		}
		var entryHeader wrappers.MESSAGE_RESOURCE_ENTRY
		entryHeader.Length = uint16(unsafe.Sizeof(entryHeader)) + uint16(textByteLength) + nullByteCount
		entryHeader.Flags = 1
		entry := make([]byte, entryHeader.Length)
		wrappers.RtlMoveMemory(
			&entry[0],
			(*byte)(unsafe.Pointer(&entryHeader)),
			unsafe.Sizeof(entryHeader))
		wrappers.RtlMoveMemory(
			&entry[unsafe.Sizeof(entryHeader)],
			(*byte)(unsafe.Pointer(syscall.StringToUTF16Ptr(messageText))),
			uintptr(textByteLength))
		entries = append(entries, entry)
	}

	// build the file header
	header := wrappers.MESSAGE_RESOURCE_DATA{
		NumberOfBlocks: uint32(len(blockLengths)),
	}

	// build the block headers
	blocks := make([]wrappers.MESSAGE_RESOURCE_BLOCK, len(blockLengths))
	headerLength := uint32(unsafe.Sizeof(header) + unsafe.Sizeof(blocks[0])*uintptr(len(blocks)))
	offset := headerLength
	entryIndex := 0
	for i := range blocks {
		blockLength := blockLengths[i]
		blocks[i].LowId = ids[entryIndex]
		blocks[i].HighId = ids[entryIndex+blockLength-1]
		blocks[i].OffsetToEntries = offset
		for j := 0; j < blockLength; j++ {
			offset += uint32(len(entries[entryIndex+j]))
		}
		entryIndex += blockLength
	}

	// concatenate everything
	data := make([]byte, headerLength)
	offset = 0
	wrappers.RtlMoveMemory(&data[offset], (*byte)(unsafe.Pointer(&header)), unsafe.Sizeof(header))
	offset += uint32(unsafe.Sizeof(header))
	for _, block := range blocks {
		wrappers.RtlMoveMemory(&data[offset], (*byte)(unsafe.Pointer(&block)), unsafe.Sizeof(block))
		offset += uint32(unsafe.Sizeof(block))
	}
	for _, entry := range entries {
		data = append(data, entry...)
	}
	return data
}
