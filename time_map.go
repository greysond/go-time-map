package time_map

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidIntervalEndBefore error = errors.New("Invalid Interval: end time cannot be before start time")
	ErrIntervalIDNotFound       error = errors.New("Interval ID Not Found")
	ErrInvalidIntervalID        error = errors.New("Invalid Interval ID: cannot be zero")
)

type TimeMap interface {
	AddInterval(start, end *time.Time, obj interface{}) (id int64, err error)
	Get(a time.Time) interface{}
	GetOk(a time.Time) (interface{}, bool)
	RemoveIntervalID(id int64) error
}

func NewTimeMap() TimeMap {
	return &timeMap{
		counter:   0,
		intervals: []*interval{},
		mu:        &sync.RWMutex{},
	}
}

type interval struct {
	id    int64
	start *time.Time
	end   *time.Time
	item  interface{}
}

func (i interval) contains(a time.Time) bool {
	result := true

	if i.start != nil {
		result = result && !a.Before(*i.start)
	}

	if i.end != nil {
		result = result && !a.After(*i.end)
	}

	return result
}

type intervals []*interval

type timeMap struct {
	counter   int64
	intervals intervals
	mu        *sync.RWMutex
}

func (tm *timeMap) Get(a time.Time) interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	item, _ := tm.GetOk(a)

	return item
}

func (tm *timeMap) GetOk(a time.Time) (interface{}, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var (
		item  interface{} = nil
		found bool        = false
	)

	// Find the highest priority interval that the time fits into
	for i, interval := range tm.intervals {
		if interval.contains(a) {
			found = true
			item = tm.intervals[i].item
		}
	}

	return item, found
}

func (tm *timeMap) RemoveIntervalID(id int64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if id == 0 {
		return ErrInvalidIntervalID
	}

	intervalIndex := -1

	for i, interval := range tm.intervals {
		if interval.id == id {
			intervalIndex = i
			break
		}
	}

	if intervalIndex == -1 {
		return ErrIntervalIDNotFound
	}

	tm.intervals = append(tm.intervals[:intervalIndex], tm.intervals[intervalIndex+1:]...)

	return nil
}

func (tm *timeMap) AddInterval(start, end *time.Time, obj interface{}) (int64, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if start != nil && end != nil && end.Before(*start) {
		return 0, ErrInvalidIntervalEndBefore
	}

	var (
		startT *time.Time = nil
		endT   *time.Time = nil
	)

	if start != nil {
		t := start.Add(0)
		startT = &t
	}

	if end != nil {
		t := end.Add(0)
		endT = &t
	}

	tm.counter++

	id := tm.counter

	tm.intervals = append(tm.intervals, &interval{
		id:    id,
		start: startT,
		end:   endT,
		item:  obj,
	})

	// TODO: Fix sorting to make lookups faster
	//sort.Sort(tm.intervals)

	return id, nil
}
