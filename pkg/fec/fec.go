// Package fec implements Reed Solomon 9/3 forward error correction,
//  intended to be sent as 9 pieces where 3 uncorrupted parts allows assembly of the message
package fec

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	
	"github.com/vivint/infectious"
)

var (
	rsTotal    = 9
	rsRequired = 3
	rsFEC      = func() *infectious.FEC {
		var e error
		var c *infectious.FEC
		if c, e = infectious.NewFEC(rsRequired, rsTotal); !E.Chk(e) {
			return c
		}
		panic(e)
	}()
)

// padData appends a 2 byte length prefix, and pads to a multiple of rsTotal.
// Max message size is limited to 1<<32 but in our use will never get near
// this size through higher level protocols breaking packets into sessions
func padData(data []byte) (out []byte) {
	dataLen := len(data)
	prefixBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(prefixBytes, uint32(dataLen))
	data = append(prefixBytes, data...)
	dataLen = len(data)
	chunkLen := (dataLen) / rsTotal
	chunkMod := (dataLen) % rsTotal
	if chunkMod != 0 {
		chunkLen++
	}
	padLen := rsTotal*chunkLen - dataLen
	out = append(data, make([]byte, padLen)...)
	return
}

// Encode turns a byte slice into a set of shards with first byte containing the shard number.
func Encode(data []byte) (chunks [][]byte, e error) {
	// First we must pad the data
	data = padData(data)
	shares := make([]infectious.Share, rsTotal)
	output := func(s infectious.Share) {
		shares[s.Number] = s.DeepCopy()
	}
	e = rsFEC.Encode(data, output)
	if e != nil {
		E.Ln(e)
		return
	}
	for i := range shares {
		// Append the chunk number to the front of the chunk
		chunk := append([]byte{byte(shares[i].Number)}, shares[i].Data...)
		// Checksum includes chunk number byte so we know if its checksum is incorrect so could the chunk number be
		checksum := crc32.Checksum(chunk, crc32.MakeTable(crc32.Castagnoli))
		checkBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(checkBytes, checksum)
		chunk = append(chunk, checkBytes...)
		chunks = append(chunks, chunk)
	}
	return
}

// Decode takes a set of shards and if there is sufficient to reassemble,
// returns the corrected data
func Decode(chunks [][]byte) (data []byte, e error) {
	var shares []infectious.Share
	for i := range chunks {
		bodyLen := len(chunks[i])
		// log.SPEW(chunks[i])
		body := chunks[i][:bodyLen-4]
		checksum := crc32.Checksum(body, crc32.MakeTable(crc32.Castagnoli))
		var share infectious.Share
		if binary.LittleEndian.Uint32(chunks[i][bodyLen-4:]) == checksum {
			share = infectious.Share{
				Number: int(body[0]),
				Data:   body[1:],
			}
		} else {
			share = infectious.Share{
				Number: i,
				Data:   nil,
			}
		}
		shares = append(shares, share)
	}
	data, e = rsFEC.Decode(nil, shares)
	if len(data) > 4 {
		prefix := data[:4]
		data = data[4:]
		dataLen := int(binary.LittleEndian.Uint32(prefix))
		if len(data) != dataLen {
			e = fmt.Errorf("somehow data is corrupted though everything else about it is correct")
		} else {
			data = data[:dataLen]
		}
	}
	return
}
