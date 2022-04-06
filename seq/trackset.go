package seq

import "github.com/chinenual/synergize/logger"

// https://www.davidkaya.com/sets-in-golang/

var exists = struct{}{}

type trackset struct {
	m map[int]struct{}
}

func (s *trackset) Init() {
	s.m = make(map[int]struct{})
}

func (s *trackset) Add(value int) {
	logger.Debugf("Add %d to %v\n", value, s)
	s.m[value] = exists
}

func (s *trackset) Remove(value int) {
	delete(s.m, value)
}

func (s *trackset) Contains(value int) bool {
	_, c := s.m[value]
	return c
}

func (s *trackset) Clear() {
	for k, _ := range s.m {
		delete(s.m, k)
	}
}

func (s *trackset) Contents() (result []int) {
	for k, _ := range s.m {
		result = append(result, k)
	}
	return
}
