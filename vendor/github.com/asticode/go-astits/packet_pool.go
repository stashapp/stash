package astits

import (
	"sort"
	"sync"
)

// packetAccumulator keeps track of packets for a single PID and decides when to flush them
type packetAccumulator struct {
	parser     PacketsParser
	pid        uint16
	programMap *programMap
	q          []*Packet
}

// newPacketAccumulator creates a new packet queue for a single PID
func newPacketAccumulator(pid uint16, parser PacketsParser, programMap *programMap) *packetAccumulator {
	return &packetAccumulator{
		parser:     parser,
		pid:        pid,
		programMap: programMap,
	}
}

// add adds a new packet for this PID to the queue
func (b *packetAccumulator) add(p *Packet) (ps []*Packet) {
	mps := b.q

	// Empty buffer if we detect a discontinuity
	if hasDiscontinuity(mps, p) {
		mps = []*Packet{}
	}

	// Throw away packet if it's the same as the previous one
	if isSameAsPrevious(mps, p) {
		return
	}

	// Flush buffer if new payload starts here
	if p.Header.PayloadUnitStartIndicator {
		ps = mps
		mps = []*Packet{p}
	} else {
		mps = append(mps, p)
	}

	// Check if PSI payload is complete
	if b.programMap != nil &&
		(b.pid == PIDPAT || b.programMap.exists(b.pid)) {
		// TODO Use partial data parsing instead
		if _, err := parseData(mps, b.parser, b.programMap); err == nil {
			ps = mps
			mps = nil
		}
	}

	b.q = mps
	return
}

// packetPool represents a queue of packets for each PID in the stream
type packetPool struct {
	b map[uint16]*packetAccumulator // Indexed by PID
	m *sync.Mutex

	parser     PacketsParser
	programMap *programMap
}

// newPacketPool creates a new packet pool with an optional parser and programMap
func newPacketPool(parser PacketsParser, programMap *programMap) *packetPool {
	return &packetPool{
		b: make(map[uint16]*packetAccumulator),
		m: &sync.Mutex{},

		parser:     parser,
		programMap: programMap,
	}
}

// add adds a new packet to the pool
func (b *packetPool) add(p *Packet) (ps []*Packet) {
	// Throw away packet if error indicator
	if p.Header.TransportErrorIndicator {
		return
	}

	// Throw away packets that don't have a payload until we figure out what we're going to do with them
	// TODO figure out what we're going to do with them :D
	if !p.Header.HasPayload {
		return
	}

	// Lock
	b.m.Lock()
	defer b.m.Unlock()

	// Make sure accumulator exists
	if _, ok := b.b[p.Header.PID]; !ok {
		b.b[p.Header.PID] = newPacketAccumulator(p.Header.PID, b.parser, b.programMap)
	}

	// Add to the accumulator
	return b.b[p.Header.PID].add(p)
}

// dump dumps the packet pool by looking for the first item with packets inside
func (b *packetPool) dump() (ps []*Packet) {
	b.m.Lock()
	defer b.m.Unlock()
	var keys []int
	for k := range b.b {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
		ps = b.b[uint16(k)].q
		delete(b.b, uint16(k))
		if len(ps) > 0 {
			return
		}
	}
	return
}

// hasDiscontinuity checks whether a packet is discontinuous with a set of packets
func hasDiscontinuity(ps []*Packet, p *Packet) bool {
	return (p.Header.HasAdaptationField && p.AdaptationField.DiscontinuityIndicator) ||
		(len(ps) > 0 && p.Header.HasPayload && p.Header.ContinuityCounter != (ps[len(ps)-1].Header.ContinuityCounter+1)%16) ||
		(len(ps) > 0 && !p.Header.HasPayload && p.Header.ContinuityCounter != ps[len(ps)-1].Header.ContinuityCounter)
}

// isSameAsPrevious checks whether a packet is the same as the last packet of a set of packets
func isSameAsPrevious(ps []*Packet, p *Packet) bool {
	return len(ps) > 0 && p.Header.HasPayload && p.Header.ContinuityCounter == ps[len(ps)-1].Header.ContinuityCounter
}
