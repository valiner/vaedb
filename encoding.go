/**
 * User: coder.sdp@gmail.com
 * Date: 2023/8/23
 * Time: 15:26
 */

package vaedb

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

const (
	crcBytes             = 4
	timestampSizeInBytes = 8
	keySizeInBytes       = 2
	valuesSizeInBytes    = 8
	headersSizeInBytes   = crcBytes + timestampSizeInBytes + keySizeInBytes + valuesSizeInBytes // Number of bytes used for all headers

	defaultBufferSize = 1024
)

type entry struct {
	crc         []byte
	timeStamp   uint64
	keyLength   uint16
	valueLength uint64
	key         []byte
	value       []byte
	length      int64
}

func (e *entry) Clone() entry {
	newEntry := *e
	newEntry.crc = CopySlice(e.crc)
	newEntry.key = CopySlice(e.key)
	newEntry.value = CopySlice(e.value)
	return newEntry
}

func CopySlice(originalSlice []byte) []byte {
	newSlice := make([]byte, len(originalSlice), cap(originalSlice))
	copy(newSlice, originalSlice)
	return newSlice
}

func readEntryByPos(file *os.File, valueSize int, valuePos int64) (*entry, error) {
	e := new(entry)
	blob := make([]byte, valueSize)
	_, err := file.ReadAt(blob, valuePos)
	if err != nil {
		return e, err
	}

	e.crc = blob[:crcBytes]
	e.timeStamp = binary.LittleEndian.Uint64(blob[crcBytes:])
	e.keyLength = binary.LittleEndian.Uint16(blob[crcBytes+timestampSizeInBytes:])
	e.valueLength = binary.LittleEndian.Uint64(blob[crcBytes+timestampSizeInBytes+keySizeInBytes:])
	e.key = blob[headersSizeInBytes : headersSizeInBytes+e.keyLength]
	e.value = blob[headersSizeInBytes+e.keyLength:]
	e.length = int64(valueSize)
	return e, nil
}

// 顺序读取数据并操作
func readEntryFormFile(file *os.File, handler entryHandle) {
	bf := bufio.NewReader(file)
	header := make([]byte, headersSizeInBytes)
	body := make([]byte, defaultBufferSize)
	e := new(entry)
	for {
		_, err := io.ReadFull(bf, header)
		//n, err := bf.Read(header)
		if err != nil {
			return
		}

		e.crc = header[:crcBytes]
		e.timeStamp = binary.LittleEndian.Uint64(header[crcBytes:])
		e.keyLength = binary.LittleEndian.Uint16(header[crcBytes+timestampSizeInBytes:])
		e.valueLength = binary.LittleEndian.Uint64(header[crcBytes+timestampSizeInBytes+keySizeInBytes:])
		blobLength := headersSizeInBytes + int(e.keyLength) + int(e.valueLength)
		//fmt.Printf("%+v ,Key:%s,kl:%d,vl:%d \n", header, e.Key, e.keyLength, e.valueLength)
		bodyLength := blobLength - headersSizeInBytes
		if bodyLength > cap(body) {
			body = make([]byte, bodyLength)
		} else {
			body = body[:bodyLength]
		}
		_, err = io.ReadFull(bf, body)
		//_, err = bf.Read(body)
		if err != nil {
			return
		}
		e.key = body[:e.keyLength]
		e.value = body[e.keyLength:]
		e.length = int64(blobLength)
		handler(e)
	}
}

func wrapEntry(ts int64, key string, value []byte, buffer *[]byte, hash Hasher) []byte {
	keyLength := len(key)
	valueLength := len(value)
	blobLength := headersSizeInBytes + keyLength + valueLength
	if blobLength > len(*buffer) {
		*buffer = make([]byte, blobLength)
	}
	blob := *buffer

	binary.LittleEndian.PutUint64(blob[crcBytes:], uint64(ts))
	binary.LittleEndian.PutUint16(blob[crcBytes+timestampSizeInBytes:], uint16(keyLength))
	binary.LittleEndian.PutUint64(blob[crcBytes+timestampSizeInBytes+keySizeInBytes:], uint64(valueLength))
	copy(blob[headersSizeInBytes:], key)
	copy(blob[headersSizeInBytes+keyLength:], value)
	copy(blob, hash.Hash(blob[crcBytes:]))

	return blob[:blobLength]
}
