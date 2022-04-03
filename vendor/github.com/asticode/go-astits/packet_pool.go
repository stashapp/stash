package astits

import (
	"sort"
	"sync"
)

// packetPool represents a pool of packets
type packetPool struct {
	b map[uint16][]*Packet // Indexed by PID
	m *sync.Mutex
}

// newPacketPool creates a new packet pool
func newPacketPool() *packetPool {
	return &packetPool{
		b: make(map[uint16][]*Packet),
		m: &sync.Mutex{},
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

	// Init buffer
	var mps []*Packet
	var ok bool
	if mps, ok = b.b[p.Header.PID]; !ok {
		mps = []*Packet{}
	}

	// Empty buffer if we detect a discontinuity
	if hasDiscontinuity(mps, p) {
		mps = []*Packet{}
	}

	// Throw away packet if it's the same as the previous one
	if isSameAsPrevious(mps, p) {
		return
	}

	// Add packet
	if len(mps) > 0 || (len(mps) == 0 && p.Header.PayloadUnitStartIndicator) {
		mps = append(mps, p)
	}

	// Check payload unit start indicator
	if p.Header.PayloadUnitStartIndicator && len(mps) > 1 {
		ps = mps[:len(mps)-1]
		mps = []*Packet{p}
	}

	// Assign
	b.b[p.Header.PID] = mps
	return
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
		ps = b.b[uint16(k)]
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
