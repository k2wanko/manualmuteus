package game

import "sync"

type Session interface {
	GetCrewmateUsers() []string
	AddCrewmateUser(id string)
	DeleteCrewmateUser(id string) bool
	GetImposterUsers() []string
	AddImposterUser(id string)
	DeleteImposterUser(id string) bool
	GetDeadUsers() []string
	AddDeadUser(id string)
	DeleteDeadUser(id string) bool
	Reset()
}

type session struct {
	mutex         sync.RWMutex
	crewmateUsers []string
	imposterUsers []string
	deadUsers     []string
}

func NewSession() Session {
	return &session{}
}

func (s *session) GetCrewmateUsers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.crewmateUsers
}

func (s *session) AddCrewmateUser(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.crewmateUsers = append(s.crewmateUsers, id)
}

func (s *session) DeleteCrewmateUser(id string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range s.crewmateUsers {
		if s.crewmateUsers[i] == id {
			s.crewmateUsers[i] = s.crewmateUsers[len(s.crewmateUsers)-1]
			s.crewmateUsers = s.crewmateUsers[:len(s.crewmateUsers)-1]
			return true
		}
	}

	return false
}

func (s *session) GetImposterUsers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.imposterUsers
}

func (s *session) AddImposterUser(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.imposterUsers = append(s.imposterUsers, id)
}

func (s *session) DeleteImposterUser(id string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range s.imposterUsers {
		if s.imposterUsers[i] == id {
			s.imposterUsers[i] = s.imposterUsers[len(s.imposterUsers)-1]
			s.imposterUsers = s.imposterUsers[:len(s.imposterUsers)-1]
			return true
		}
	}

	return false
}

func (s *session) GetDeadUsers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.deadUsers
}

func (s *session) AddDeadUser(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.deadUsers = append(s.deadUsers, id)
}

func (s *session) DeleteDeadUser(id string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range s.deadUsers {
		if s.deadUsers[i] == id {
			s.deadUsers[i] = s.deadUsers[len(s.deadUsers)-1]
			s.deadUsers = s.deadUsers[:len(s.deadUsers)-1]
			return true
		}
	}

	return false
}

func (s *session) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.crewmateUsers = append(s.crewmateUsers, s.imposterUsers...)
	s.crewmateUsers = append(s.crewmateUsers, s.deadUsers...)
	s.imposterUsers = nil
	s.deadUsers = nil
}
