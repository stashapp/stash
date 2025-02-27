package dlna

import (
	"slices"
	"sync"
	"time"
)

// only keep the 10 most recent IP addresses
const recentListLength = 10

const wildcard = "*"

type tempIPWhitelist struct {
	pattern string
	until   *time.Time
}

type ipWhitelistManager struct {
	recentIPAddresses []string
	config            Config
	tempWhitelist     []tempIPWhitelist
	mutex             sync.Mutex
}

// addRecent adds the provided address to the recent IP addresses list if it
// was not already present. Returns true if it was already present.
func (m *ipWhitelistManager) addRecent(addr string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	i := slices.Index(m.recentIPAddresses, addr)
	if i != -1 {
		if i == 0 {
			// don't do anything if it's already at the start
			return true
		}

		// remove from the list
		m.recentIPAddresses = append(m.recentIPAddresses[:i], m.recentIPAddresses[i+1:]...)
	}

	// add to the top of the list
	m.recentIPAddresses = append([]string{addr}, m.recentIPAddresses...)

	if len(m.recentIPAddresses) > recentListLength {
		m.recentIPAddresses = m.recentIPAddresses[:recentListLength]
	}

	return i != -1
}

func (m *ipWhitelistManager) getRecent() []string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.recentIPAddresses
}

func (m *ipWhitelistManager) getTempAllowed() []*Dlnaip {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var ret []*Dlnaip

	now := time.Now()
	removeExpired := false
	for _, a := range m.tempWhitelist {
		if a.until != nil && now.After(*a.until) {
			removeExpired = true
			continue
		}

		ret = append(ret, &Dlnaip{
			IPAddress: a.pattern,
			Until:     a.until,
		})
	}

	if removeExpired {
		m.removeExpiredWhitelists()
	}

	return ret
}

func (m *ipWhitelistManager) ipAllowed(addr string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, a := range m.config.GetDLNADefaultIPWhitelist() {
		if a == wildcard {
			return true
		}

		if addr == a {
			return true
		}
	}

	now := time.Now()
	removeExpired := false
	for _, a := range m.tempWhitelist {
		if a.until != nil && now.After(*a.until) {
			removeExpired = true
			continue
		}

		if a.pattern == wildcard {
			return true
		}

		if addr == a.pattern {
			return true
		}
	}

	if removeExpired {
		m.removeExpiredWhitelists()
	}

	return false
}

func (m *ipWhitelistManager) removeExpiredWhitelists() {
	// assumes mutex is already held
	var newList []tempIPWhitelist
	now := time.Now()

	for _, a := range m.tempWhitelist {
		if a.until != nil && now.After(*a.until) {
			continue
		}

		newList = append(newList, a)
	}

	m.tempWhitelist = newList
}

func (m *ipWhitelistManager) allowPattern(pattern string, duration *time.Duration) {
	if pattern == "" {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// overwrite existing
	var newList []tempIPWhitelist
	found := false

	var until *time.Time

	if duration != nil {
		u := time.Now().Add(*duration)
		until = &u
	}

	for _, a := range m.tempWhitelist {
		if a.pattern == pattern {
			a.until = until
			found = true
		}

		newList = append(newList, a)
	}

	if !found {
		newList = append(newList, tempIPWhitelist{
			pattern: pattern,
			until:   until,
		})
	}

	m.tempWhitelist = newList
}

func (m *ipWhitelistManager) removePattern(pattern string) bool {
	if pattern == "" {
		return false
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	var newList []tempIPWhitelist
	found := false

	for _, a := range m.tempWhitelist {
		if a.pattern == pattern {
			found = true
			continue
		}

		newList = append(newList, a)
	}

	m.tempWhitelist = newList
	return found
}
