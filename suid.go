package suid

import (
	"fmt"
	"sync"
	"time"
)

var maxSeq = 1<<10 - 1 //  63 bit(total) - 41 bit(ms) - 12 bit(appId) = 10

type suid struct {
	appId     int64
	seq       int
	currentMs int64
	sync.Mutex
}

func NewSUID(appId int) *suid {
	if appId >= 2048 {
		panic("App Id cannot be more than 4096")
	}

	return &suid{
		appId: int64(appId) << 12,
		seq:   0,
	}
}

func (s *suid) Generate() (int64, error) {
	var id, ms int64
	ms = time.Now().UnixNano() / 1e6
	// ms goes to head
	id = ms << 22 // 63 bit - 41 bit(ms)
	// set appId into middle
	id |= s.appId
	// generate sequence
	seq, err := s.nextSeq(ms)
	if err != nil {
		return int64(0), err
	}

	// generated sequence goes to the end
	id |= seq

	return id, nil
}

func (s *suid) nextSeq(ms int64) (int64, error) {
	s.Lock()
	if s.currentMs > ms {
		return int64(0), fmt.Errorf("Time goes backward in this machine")
	}

	if s.currentMs < ms {
		s.currentMs = ms
		s.seq = -1
	}

	s.seq++
	if s.seq > maxSeq {
		return int64(0), fmt.Errorf("You created more than %d ids in one milisecond", maxSeq)
	}

	s.Unlock()
	return int64(s.seq), nil
}
