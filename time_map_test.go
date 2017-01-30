package time_map

import (
	"testing"
	"time"
)

func makeSecs(s int64) *time.Time {
	t := time.Unix(s, 0)
	return &t
}

var (
	a *time.Time = makeSecs(1)
	b *time.Time = makeSecs(2)
	c *time.Time = makeSecs(3)
	d *time.Time = makeSecs(4)
	e *time.Time = makeSecs(5)
	f *time.Time = makeSecs(6)
	g *time.Time = makeSecs(7)
)

func TestTimeMap_AddInterval(t *testing.T) {
	tm := NewTimeMap()

	if _, err := tm.AddInterval(nil, nil, 0); err != nil {
		t.Error(err.Error())
	}

	if _, err := tm.AddInterval(a, c, 1); err != nil {
		t.Error(err.Error())
	}

	if _, err := tm.AddInterval(a, c, 2); err != nil {
		t.Error(err.Error())
	}

	if _, err := tm.AddInterval(c, f, 3); err != nil {
		t.Error(err.Error())
	}

	if _, err := tm.AddInterval(e, f, 4); err != nil {
		t.Error(err.Error())
	}


	if val, ok := tm.GetOk(*b); ok {
		if val != 2 {
			t.Errorf("%d expected to be %d", val, 2)
		}
	} else {
		t.Error("Should have found something")
	}

	if val, ok := tm.GetOk(*d); ok {
		if val != 3 {
			t.Errorf("%d expected to be %d", val, 3)
		}
	} else {
		t.Error("Should have found something")
	}

	if val, ok := tm.GetOk(*g); ok {
		if val != 0 {
			t.Errorf("%d expected to be %d", val, 0)
		}
	} else {
		t.Error("Should have found something")
	}
}
