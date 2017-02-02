package timeMap

import (
	"fmt"
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
	tm := New()

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

func ExampleNew() {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	tm := New()

	// No begin time
	tm.AddInterval(nil, &tomorrow, "up to tomorrow")

	str := tm.Get(today)
	fmt.Println(str) // up to tomorrow

	// Starting yesterday
	tm.AddInterval(&yesterday, nil, "since yesterday")

	str, ok := tm.GetOk(today)
	if ok {
		fmt.Println(str) // since yesterday
	}

	// Create an explicit interval
	id, _ := tm.AddInterval(&yesterday, &tomorrow, "these days")

	str, ok = tm.GetOk(today)
	if ok {
		fmt.Println(str) // these days
	}

	// Delete an interval, referenced by id
	err := tm.RemoveIntervalID(id)

	if err == nil {
		fmt.Println("id removed")
	} else {
		fmt.Println("id not found") // Not reached in this case
	}

	// Verify that the old value is in effect
	str, ok = tm.GetOk(today)
	if ok {
		fmt.Println(str) // since yesterday
	}

	// Output:
	// up to tomorrow
	// since yesterday
	// these days
	// id removed
	// since yesterday
}

func ExampleTimeMap_AddInterval() {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	tm := New()

	// [−∞, +∞]
	_, err := tm.AddInterval(nil, nil, "all of time")

	if err != nil {
		fmt.Println("invalid interval") // Prints nothing
	}

	// [yesterday, tomorrow]
	_, err = tm.AddInterval(&yesterday, &tomorrow, "finite range")

	if err != nil {
		fmt.Println("invalid interval") // Prints nothing
	}

	// [tomorrow, +∞]
	_, err = tm.AddInterval(&tomorrow, nil, "tomorrow and forever")

	if err != nil {
		fmt.Println("invalid interval") // Prints nothing
	}

	// [tomorrow, yesterday] (invalid)
	_, err = tm.AddInterval(&tomorrow, &yesterday, "tomorrow and forever")

	if err != nil {
		fmt.Println("invalid interval") // Will print error
	}

	// Output:
	// invalid interval
}

func ExampleTimeMap_Clear() {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	tm := New()

	// No begin time
	tm.AddInterval(&yesterday, &tomorrow, "these days")

	if str, ok := tm.GetOk(today); ok {
		fmt.Println(str) // these days
	} else {
		fmt.Println("not found")
	}

	tm.Clear()

	if str, ok := tm.GetOk(today); ok {
		fmt.Println(str) // these days
	} else {
		fmt.Println("not found")
	}

	// Output:
	// these days
	// not found
}
