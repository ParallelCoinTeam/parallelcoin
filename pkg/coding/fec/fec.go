// Package fec implements Reed Solomon 9/3 forward error correction,
//  intended to be sent as 9 pieces where 3 uncorrupted parts allows assembly of the message
package fec

import (
	"encoding/binary"
	"errors"
	
	"github.com/vivint/infectious"
)

var (
	rsTotal    = 9
	rsRequired = 3
	rsFEC      = func() *infectious.FEC {
		fec, e := infectious.NewFEC(rsRequired, rsTotal)
		if e != nil {
			err.Ln(e)
		}
		return fec
	}()
)

// padData appends a 2 byte length prefix, and pads to a multiple of rsTotal. Max message size is limited to 1<<32 but
// in our use will never get near this size through higher level protocols breaking packets into sessions
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

// Encode turns a byte slice into a set of shards with first byte containing the shard number. Previously this code
// included a CRC32 but this is unnecessary since the shards will be sent wrapped in HMAC protected encryption
func Encode(data []byte) (chunks [][]byte, e error) {
	// First we must pad the data
	data = padData(data)
	shares := make([]infectious.Share, rsTotal)
	output := func(s infectious.Share) {
		shares[s.Number] = s.DeepCopy()
	}
	e = rsFEC.Encode(data, output)
	if e != nil {
		err.Ln(e)
		return
	}
	for i := range shares {
		// Append the chunk number to the front of the chunk
		chunk := append([]byte{byte(shares[i].Number)}, shares[i].Data...)
		// Checksum includes chunk number byte so we know if its checksum is incorrect so could the chunk number be
		//
		// checksum := crc32.Checksum(chunk, crc32.MakeTable(crc32.Castagnoli))
		// checkBytes := make([]byte, 4)
		// binary.LittleEndian.PutUint32(checkBytes, checksum)
		// chunk = append(chunk, checkBytes...)
		chunks = append(chunks, chunk)
	}
	// L.Spew(chunks)
	return
}

// Decode takes a set of shards and if there is sufficient to reassemble,
// returns the corrected data
func Decode(chunks [][]byte) (data []byte, e error) {
	var shares []infectious.Share
	if len(chunks) < 1 {
		dbg.Ln("nil chunks")
		return nil, errors.New("asked to decode nothing")
	}
	totalLen := 0
	for i := range chunks {
		// bodyLen := len(chunks[i])
		// log.SPEW(chunks[i])
		body := chunks[i] // [:bodyLen]
		share := infectious.Share{
			Number: int(body[0]),
			Data:   body[1:],
		}
		shares = append(shares, share)
		totalLen += len(share.Data)
	}
	if shares == nil {
		panic("shares are nil")
	}
	data = make([]byte, totalLen)
	dataLen := len(shares[0].Data)
	if e := rsFEC.Rebuild(
		shares, func(s infectious.Share) {
			copy(data[s.Number*dataLen:], s.Data)
		},
	); dbg.Chk(e) {
	}
	// data, e = rsFEC.Decode(nil, shares)
	if len(data) > 4 {
		prefix := data[:4]
		data = data[4:]
		dataLen := int(binary.LittleEndian.Uint32(prefix))
		data = data[:dataLen]
	}
	return
}
