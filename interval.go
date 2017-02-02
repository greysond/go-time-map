package timeMap

import "time"

type (
	interval struct {
		id    uint64
		start *time.Time
		end   *time.Time
		item  interface{}
	}

	intervals []*interval
)

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

func (slice intervals) Len() int {
	return len(slice)
}

func (slice intervals) Less(i, j int) bool {
	a, b := slice[i], slice[j]

	is, ie := slice[i].start, slice[i].end
	js, je := slice[j].start, slice[j].end

	if is == js {
		// starts equal
		if ie == je {
			return slice[i].id < slice[j].id
		}
		// Ends not equal
		if ie == nil {

		}
	} else {
		// starts not equal
		if is == nil {
			// nil start is always less than non-nil start
			return true
		}
		if js == nil {
			return false
		}

		return is.Before(*js)
	}

	// a has no start, b does, a must be less
	if a.start == nil && b.start != nil {
		return true
	}

	// b has no start, a does, b must be less
	if b.start == nil && a.start != nil {
		return false
	}

	// Neither has a start, must compare ends
	if a.start == nil && b.start == nil {
		// a has no end, b does, b is less
		if a.end == nil && b.end != nil {
			return false
		}

		// b has no end, a does, a is less
		if b.end == nil && a.end != nil {
			return true
		}

		// Neither has a start nor an end, so resort to priority
		return a.id < b.id
	}

	// Both are non-nil

	if a.start.Before(*b.start) {
		return true
	}

	if b.start.Before(*a.start) {
		return false
	}

	// Starts are equal, must compare ends

	if a.end.Before(*b.end) {
		return true
	}

	if b.end.Before(*a.end) {
		return false
	}

	// Start and end are equal, resort to priority
	return a.id < b.id
}

func (slice intervals) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
