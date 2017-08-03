package nsuschedule

import "sync"

type Schedule struct {
	Mux      sync.Mutex
	Schedule map[string]string
}

func GetAllSchedule() (Schedule, error) {
	return Schedule{}, nil
}

func NewSchedule() Schedule {
	return Schedule{
		Schedule: make(map[string]string),
	}
}

func (s *Schedule) Update() {
	s.Mux.Lock()
	defer s.Mux.Unlock()
}
