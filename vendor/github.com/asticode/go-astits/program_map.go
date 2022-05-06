package astits

import "sync"

// programMap represents a program ids map
type programMap struct {
	m *sync.Mutex
	p map[uint16]uint16 // map[ProgramMapID]ProgramNumber
}

// newProgramMap creates a new program ids map
func newProgramMap() programMap {
	return programMap{
		m: &sync.Mutex{},
		p: make(map[uint16]uint16),
	}
}

// exists checks whether the program with this pid exists
func (m programMap) exists(pid uint16) (ok bool) {
	m.m.Lock()
	defer m.m.Unlock()
	_, ok = m.p[pid]
	return
}

// set sets a new program id
func (m programMap) set(pid, number uint16) {
	m.m.Lock()
	defer m.m.Unlock()
	m.p[pid] = number
}

func (m programMap) unset(pid uint16) {
	m.m.Lock()
	defer m.m.Unlock()
	delete(m.p, pid)
}

func (m programMap) toPATData() *PATData {
	m.m.Lock()
	defer m.m.Unlock()

	d := &PATData{
		Programs:          []*PATProgram{},
		TransportStreamID: uint16(PSITableIDPAT),
	}

	for pid, pnr := range m.p {
		d.Programs = append(d.Programs, &PATProgram{
			ProgramMapID:  pid,
			ProgramNumber: pnr,
		})
	}

	return d
}
