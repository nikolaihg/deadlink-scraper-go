package stats

import (
	"log"
	"strconv"
	"sync"
)

type LinkStats struct {
	mu           sync.Mutex
	Total        int
	Internal     int
	External     int
	Alive        int
	Dead         int
	Skipped      int
	ByStatusCode map[string]int
}

func New() *LinkStats {
	return &LinkStats{
		ByStatusCode: make(map[string]int),
	}
}

func (s *LinkStats) UpdatePageLink() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Total++
	s.Skipped++
}

func (s *LinkStats) UpdateUnknown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Total++
	s.Skipped++
}

func (s *LinkStats) UpdateEmptyURL() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Total++
	s.Dead++
}

func (s *LinkStats) UpdateInternal() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Total++
	s.Internal++
}

func (s *LinkStats) UpdateExternal() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Total++
	s.External++
}

func (s *LinkStats) UpdateResult(code int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if there was an error, count as dead
	if err != nil {
		s.Dead++
		return
	}

	codeStr := strconv.Itoa(code)
	s.ByStatusCode[codeStr]++
	if code >= 400 {
		s.Dead++
	} else {
		s.Alive++
	}
}

func (s *LinkStats) Print() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Scan complete:")
	log.Printf("Total:    %d\n", s.Total)
	log.Printf("Internal: %d\n", s.Internal)
	log.Printf("External: %d\n", s.External)
	log.Printf("Alive:    %d\n", s.Alive)
	log.Printf("Dead:     %d\n", s.Dead)
	log.Printf("Skipped:  %d\n", s.Skipped)
	log.Println("Status codes distribution:")
	for code, count := range s.ByStatusCode {
		log.Printf("  %s: %d\n", code, count)
	}
}
