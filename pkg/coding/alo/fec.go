// Package alo implements a forward error correction scheme using Reed Solomon Erasure Coding. It segments the data into
// 1kb chunks in 16kb segments and contains a function for use in a network handler to process received packets
package alo

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/templexxx/reedsolomon"
)

type Segments [][]byte
type ShardedSegments []Segments

const (
	SegmentSize = 2 << 13
	ShardSize   = 2 << 9
)

func getEmptyShards(size, count int) (out Segments) {
	out = make(Segments, count)
	for i := range out {
		out[i] = make([]byte, size)
	}
	return
}

// GetShards returns a bundle of segments to be sent or stored in a 1kb segment size with redundancy shards added to
// each segment's shards that can reconstruct the original message by derivation via the available parity shards.
func GetShards(buf []byte, redundancy int, ) (out ShardedSegments) {
	prefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(prefix, uint32(len(buf)))
	// the following was eliminated to avoid a second copy of the buffer for 4 bytes
	// buf = append(prefix, buf...)
	segments := SegmentBytes(buf, SegmentSize)
	sl := len(segments)
	sharded := make(ShardedSegments, sl)
	for i := range segments {
		sharded[i] = SegmentBytes(segments[i], ShardSize)
	}
	// the foregoing operations should not have required any memory allocations since they were creating new slices so
	// this should be the only (necessary) allocation of the data the RS codec will work on in situ. Effectively using
	// Go's slice syntax to create a copy map.
	out = make(ShardedSegments, sl)
	for i := range sharded {
		// add 4 bytes for the shard identifier prefix (segment/segments/shard/shards) and 4 bytes for the total data
		// payload length which gives required also from a segment (last segments can differ in length), and 8 bytes
		// means no alignment cost for the copy
		out[i] = getEmptyShards(ShardSize, ShardsPerRedundancy(len(sharded[i]), redundancy))
	}
	for i := range sharded {
		for j := range sharded[i] {
			// copy the data out of the segments into place with the additional segments prepared for the RS codec
			copy(out[i][j], sharded[i][j])
		}
	}
	for i := range out {
		dataLen := len(sharded[i])
		parityLen := len(out[i]) - dataLen
		if rs, err := reedsolomon.New(dataLen, parityLen); !Check(err) {
			if err = rs.Encode(out[i]); Check(err) {
			}
		}
	}
	// put shard metadata in front of the shards
	for i := range out {
		for j := range out[i] {
			p := make([]byte, 8)
			// the segment number
			p[0] = byte(i)
			// number of segments
			p[1] = byte(len(sharded))
			// shard number
			p[2] = byte(j)
			// number of shards
			p[3] = byte(len(out[i]))
			// required shards can be computed based on the length of the payload thus shortening this header to a round
			// 8 bytes. Grouping is handled by using the nonce of GCM-AES encryption for puncture detection and tamper
			// resistance to associate packets in the decoder
			copy(p[4:8], prefix)
			out[i][j] = append(p, out[i][j]...)
		}
	}
	// fmt.Println()
	return
}

func ShardsPerRedundancy(nShards, redundancy int) int {
	return nShards * (redundancy + 100) / 100
}

// TODO: this might be a useful thing with a closure for debugging library
// st := ""
// for i := range out {
//	for j := range out[i] {
//		st += fmt.Sprintln(i, j, len(out[i][j]))
//	}
// }
// Debug(st)

// SegmentBytes breaks a chunk of data into requested sized limit chunks
func SegmentBytes(buf []byte, lim int) (out [][]byte) {
	p := Pieces(len(buf), lim)
	chunks := make([][]byte, p)
	for i := range chunks {
		if len(buf) < lim {
			chunks[i] = buf
		} else {
			chunks[i], buf = buf[:lim], buf[lim:]
		}
	}
	return chunks
}

// Pieces computes the number of pieces based on a given chunk size
func Pieces(dLen, size int) (s int) {
	sm := dLen % size
	if sm > 0 {
		return dLen/size + 1
	}
	return dLen / size
}

// GetParams reads the shard's prefix to provide the correct parameters for the RS codec the packet requires
// based on the prefix on a shard (presumably to create the codec when a new packet/group of shards arrives)
func GetParams(data []byte) (
	p ShardPrefix, err error,
) {
	if len(data) <= 3 {
		err = errors.New("provided data is not long enough to be a shard")
		return
	}
	p.segment, p.totalSegments, p.shard, p.totalShards = int(data[0]), int(data[1]), int(data[2]), int(data[3])
	p.length = int(binary.LittleEndian.Uint32(data[4:8]))
	p.requiredShards = Pieces(p.length, ShardSize)
	return
}

// PartialSegment is a max 16kb long segment with arbitrary redundancy parameters when all of the data segments are
// successfully received hasAll indicates the segment may be ready to reassemble
type PartialSegment struct {
	data, parity int
	segment      Segments
	hasAll       bool
}

// GetShardCount returns the number of shards already acquired
func (p PartialSegment) GetShardCount() (count int) {
	for i := range p.segment {
		if p.segment[i] == nil || len(p.segment[i]) > 0 {
			count++
		}
	}
	return
}

// Partials is a structure for storing a new inbound packet
type Partials struct {
	totalSegments, length int
	segments              []PartialSegment
	decoded               bool
}

func (p *Partials) IsDecoded() bool {
	return p.decoded
}

type ShardPrefix struct {
	segment, totalSegments, shard, totalShards, requiredShards, length int
}

// NewPacket creates a new structure to store a collection of incoming shards when the first of a new packet arrives
func NewPacket(firstShard []byte) (o *Partials, err error) {
	o = &Partials{}
	var p ShardPrefix
	if p, err = GetParams(firstShard); Check(err) {
	}
	o.totalSegments = p.totalSegments
	o.length = p.length
	o.segments = make([]PartialSegment, o.totalSegments)
	o.segments[p.segment] = PartialSegment{
		data:    p.requiredShards,
		parity:  p.totalShards - p.requiredShards,
		segment: make(Segments, p.totalShards),
	}
	if o.segments[p.segment].segment == nil {
		o.segments[p.segment].segment = make(Segments, p.totalShards)
	}
	o.segments[p.segment].segment[p.shard] = firstShard[8:]
	return
}

// AddShard adds a newly received shard to a Partials, ensuring that it has matching parameters (if the HMAC on the
// packet's wrapper passes it should be unless someone is playing silly buggers)
func (p *Partials) AddShard(newShard []byte) (err error) {
	var params ShardPrefix
	if params, err = GetParams(newShard); Check(err) {
	}
	if p.totalSegments != params.totalSegments {
		return errors.New("shard has incorrect segment count for bundle")
	}
	if p.length != params.length {
		return errors.New("shard specifies different length from the bundle")
	}
	p.segments[params.segment].data = params.requiredShards
	p.segments[params.segment].parity = params.totalShards - params.requiredShards
	if p.segments[params.segment].segment == nil {
		p.segments[params.segment].segment = make(Segments, params.totalShards)
	}
	p.segments[params.segment].segment[params.shard] = newShard[8:]
	// as the pieces are likely to arrive more or less in order, check when the data shards are done and mark the
	// segment as ready to decode
	if !p.segments[params.segment].hasAll {
		var count int
		// if we count all of the data shards are present mark the segment as ready to decode
		for i := range p.segments[params.segment].segment {
			if p.segments[params.segment].segment[i] != nil || len(p.segments[params.segment].segment[i]) != 0 {
				count++
			} else {
				// if we encounter empty shards, stop counting
				break
			}
		}
		if count >= p.segments[params.segment].data {
			p.segments[params.segment].hasAll = true
		}
	}
	return
}

// HasAllDataShards returns true if all data shards are present in a Partials
func (p *Partials) HasAllDataShards() bool {
	for i := range p.segments {
		if i > p.segments[i].data {
			break
		}
		s := p.segments[i]
		if s.segment == nil {
			return false
		}
	}
	return true
}

// HasMinimum returns true if there may be enough data to decode
func (p *Partials) HasMinimum() bool {
	// first check if all segments have all data shards already
	if p.HasAllDataShards() {
		return true
	}
	for i := range p.segments {
		// if the segment hasn't got all of the data shards, count total number of shards, otherwise move to the next
		if !p.segments[i].hasAll {
			// if the number of shards in the segment is above the required move to the next segment
			if p.segments[i].GetShardCount() >= p.segments[i].data {
				continue
			}
			// if we encounter a segment with less than required we can return the packet has not got minimum
			return false
		}
	}
	return true
}

// GetRatio is used after the receive delay period expires to determine how successful a packet was. If it was exactly
// enough with no surplus, return 0, if there is more than the minimum, return the proportion compared to the total
// redundancy of the packet, if there is less, return a negative proportion versus the amount, -1 means zero received,
// and a fraction of 1 indicates the proportion that was received compared to the minimum
func (p *Partials) GetRatio() (out float64) {
	var count, max, min int
	for i := range p.segments {
		max += p.segments[i].data + p.segments[i].parity
		min += p.segments[i].data
		for j := range p.segments[i].segment {
			if p.segments[i].segment[j] != nil || len(p.segments[i].segment[j]) != 0 {
				count++
			}
		}
	}
	excess := float64(count - min)
	beyond := float64(max - min)
	switch {
	case excess > 0:
		// if we had enough, the proportion towards all (all being equal to 1) is returned
		out = excess / beyond
	case excess == 0:
		// if we had exactly enough, we get 0
	case excess < 0:
		// if we had less than enough, return the proportion compared to the minimum, will be negative proportion of
		// the minimum. This value can be used to scale the response for requesting increasing redundancy in case of
		// failure for the retransmit
		out = excess / float64(min)
	}
	return
}

// Decode the received message if we have sufficient pieces
func (p *Partials) Decode() (final []byte, err error) {
	final = make([]byte, p.length)
	if p.HasAllDataShards() {
		var parts [][]byte
		// if all data shards were received we can just join them together and return the original data
		for i := range p.segments {
			for j := range p.segments[i].segment {
				if j < p.segments[i].data {
					parts = append(parts, p.segments[i].segment[j])
				} else {
					break
				}
			}
		}
		var cursor int
		for i := range parts {
			if cursor > p.length {
				break
			}
			l := cursor + p.length
			copy(final[cursor:], parts[i])
			cursor = l
		}
		p.decoded = true
		final = final[:p.length]
		return
	}
	if !p.HasMinimum() {
		return nil, fmt.Errorf(
			"not enough shards, have %f less than required",
			-p.GetRatio(),
		)
	}
	// if we don't have all data shards but have total shards equal to the original we can reconstruct the original
	var needReconst, dpHas []int
	var rs *reedsolomon.RS
	for i := range p.segments {
		if rs, err = reedsolomon.New(p.segments[i].data, p.segments[i].parity); !Check(err) {
			for j := range p.segments[i].segment {
				// if the segment is empty it wasn't received or deciphered but we only need reconstruction on the
				// data shards, append to the list for the reconstruction
				if j < p.segments[i].data && p.segments[i].segment[j] == nil {
					needReconst = append(needReconst, j)
				} else {
					dpHas = append(dpHas, j)
				}
			}
			Info(dpHas, needReconst)
			if err = rs.Reconst(p.segments[i].segment, dpHas, needReconst); Check(err) {
				return
			}
		}
	}
	var parts [][]byte
	for i := range p.segments {
		for j := range p.segments[i].segment {
			if j < p.segments[i].data {
				parts = append(parts, p.segments[i].segment[j])
			} else {
				break
			}
		}
	}
	var cursor int
	for i := range parts {
		lp := len(parts[i])
		copy(final[cursor:cursor+lp], parts[i])
		cursor += lp
	}
	p.decoded = true
	final = final[:p.length]
	return
}
