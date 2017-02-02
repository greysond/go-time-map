package timeMap

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrInvalidIntervalEndBefore is the error returned when
	// AddInterval is called with an end time that's before the start time
	ErrInvalidIntervalEndBefore error = errors.New("Invalid Interval: end time cannot be before start time")

	// ErrIntervalIDNotFound is the error returned when
	// RemoveIntervalID is called with a non-existent interval id
	ErrIntervalIDNotFound error = errors.New("Interval ID Not Found")

	// ErrInvalidIntervalID is the error returned when
	// RemoveIntervalID is called with a zero-value interval id
	ErrInvalidIntervalID error = errors.New("Invalid Interval ID: cannot be zero")
)

// New creates a new, empty TimeMap
func New() *TimeMap {
	return &TimeMap{}
}

// TimeMap is able to store interfaces within intervals of time,
// which then may be retrieved by supplying a time which is
// contained in a previously supplied interval
type TimeMap struct {
	counter   uint64
	intervals intervals
	mu        sync.RWMutex
}

// AddInterval adds an interface which may be looked up with
// a time that lies within the supplied interval (inclusive).
// Returns an ID for the interval which may be used to remove it later.
func (tm *TimeMap) AddInterval(start, end *time.Time, obj interface{}) (uint64, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if start != nil && end != nil && end.Before(*start) {
		return 0, ErrInvalidIntervalEndBefore
	}

	var startT, endT *time.Time

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

// Get looks up an interface for the supplied time.
// If multiple time intervals match, the interval that
// was defined latest will match.
// Returns nil if not found
func (tm *TimeMap) Get(a time.Time) interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	item, _ := tm.GetOk(a)

	return item
}

// GetOk works the same as Get(time.Time), but also returns true or false
// if there was a matching interval.
func (tm *TimeMap) GetOk(a time.Time) (interface{}, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var (
		item  interface{}
		found bool
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

// RemoveIntervalID removes the interval referenced by
// the id that was returned by AddInterval
func (tm *TimeMap) RemoveIntervalID(id uint64) error {
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

// Clear removes all intervals from the TimeMap.
// Does not reset the counter, so interval IDs remain
// unique within the map
func (tm *TimeMap) Clear() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.intervals = nil
}
