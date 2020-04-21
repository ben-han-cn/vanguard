package iterator

import (
	"container/list"
	"math"
	"net"
	"time"
)

const MaxTimeoutCount = 3
const TimeoutServerRetryInterval = 60 * time.Second
const HostInitRtt = 0 * time.Second

type Host net.IP

type HostSelector interface {
	SetRtt(host Host, rtt time.Duration)
	SetTimeout(host Host, timeout time.Duration)
	SelectHost(hosts []Host) (Host, bool)
}

type HostState struct {
	host         Host
	rtt          time.Duration
	timeoutCount int
	wakeupTime   *time.Time
}

func newState(host Host, rtt time.Duration) *HostState {
	return &HostState{
		host:         host,
		rtt:          rtt,
		timeoutCount: 0,
		wakeupTime:   nil,
	}
}

func newStateWithTimeout(host Host, timeout time.Duration) *HostState {
	return &HostState{
		host:         host,
		rtt:          timeout,
		timeoutCount: 1,
		wakeupTime:   nil,
	}
}

func (s *HostState) getHost() Host {
	return s.host
}

func (s *HostState) SetRtt(rtt time.Duration) {
	if s.timeoutCount > 0 {
		s.timeoutCount = 0
		s.wakeupTime = nil
	}

	s.rtt = calcuateRtt(s.rtt, rtt)
}

func (s *HostState) SetTimeout(timeout time.Duration) {
	if s.timeoutCount < MaxTimeoutCount {
		s.timeoutCount += 1
		s.rtt = calcuateRtt(s.rtt, timeout)
	}

	if s.timeoutCount == MaxTimeoutCount {
		wakeupTime := time.Now().Add(TimeoutServerRetryInterval)
		s.wakeupTime = &wakeupTime
	}
}

func (s *HostState) isUsable() bool {
	if s.wakeupTime != nil {
		return time.Now().After(*s.wakeupTime)
	} else {
		return true
	}
}

func (s *HostState) GetRtt() time.Duration {
	if s.isUsable() {
		return s.rtt
	} else {
		return math.MaxInt64 * time.Nanosecond
	}
}

func calcuateRtt(last, now time.Duration) time.Duration {
	ln := last.Nanoseconds()
	nn := now.Nanoseconds()
	return time.Duration((ln*7+nn*3)/10) * time.Nanosecond
}

type RttBasedHostSelector struct {
	cap   int
	hosts map[string]*list.Element
	ll    *list.List
}

var _ HostSelector = &RttBasedHostSelector{}

func newRttBasedHostSelector(cap int) *RttBasedHostSelector {
	return &RttBasedHostSelector{
		cap:   cap,
		hosts: make(map[string]*list.Element),
		ll:    list.New(),
	}
}

func (s *RttBasedHostSelector) SetRtt(host Host, rtt time.Duration) {
	key := string(host)
	if elem, ok := s.hosts[key]; ok {
		state := elem.Value.(*HostState)
		state.SetRtt(rtt)
		s.ll.MoveToFront(elem)
	} else if s.ll.Len() < s.cap {
		elem := s.ll.PushFront(newState(host, rtt))
		s.hosts[key] = elem
	} else {
		elem := s.ll.Back()
		os := elem.Value.(*HostState)
		delete(s.hosts, string(os.host))
		*os = *newState(host, rtt)
		s.hosts[key] = elem
		s.ll.MoveToFront(elem)
	}
}

func (s *RttBasedHostSelector) SetTimeout(host Host, timeout time.Duration) {
	key := string(host)
	if elem, ok := s.hosts[key]; ok {
		state := elem.Value.(*HostState)
		state.SetTimeout(timeout)
		s.ll.MoveToFront(elem)
	} else if s.ll.Len() < s.cap {
		elem := s.ll.PushFront(newStateWithTimeout(host, timeout))
		s.hosts[key] = elem
	} else {
		elem := s.ll.Back()
		os := elem.Value.(*HostState)
		delete(s.hosts, string(os.getHost()))
		*os = *newStateWithTimeout(host, timeout)
		s.hosts[key] = elem
		s.ll.MoveToFront(elem)
	}
}

func (s *RttBasedHostSelector) getRtt(host Host) time.Duration {
	if elem, ok := s.hosts[string(host)]; ok {
		return elem.Value.(*HostState).GetRtt()
	} else {
		return HostInitRtt
	}
}

func (s *RttBasedHostSelector) SelectHost(hosts []Host) (Host, bool) {
	min := math.MaxInt64 * time.Nanosecond
	minIndex := -1
	for i, host := range hosts {
		rtt := s.getRtt(host)
		if rtt < min {
			min = rtt
			minIndex = i
		}
	}

	if minIndex == -1 {
		return Host(net.IP{}), false
	} else {
		return hosts[minIndex], true
	}
}
