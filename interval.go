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
