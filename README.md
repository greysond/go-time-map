[![Build Status](https://travis-ci.org/greysond/go-time-map.svg?branch=master)](https://travis-ci.org/greysond/go-time-map)
[![GoDoc](https://godoc.org/github.com/greysond/go-time-map?status.svg)](http://godoc.org/github.com/greysond/go-time-map)

## go-time-map

A thread-safe data structure for mapping time intervals to interfaces.

### Installing

    $ go get github.com/greysond/go-time-map


### Example Usage

```go
package main

import (
	"time"
	"fmt"

	"github.com/greysond/go-time-map"
)

func main() {

	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	tm := time_map.NewTimeMap()

	// No begin time
	tm.AddInterval(nil, &tomorrow, "up to tomorrow")

	str := tm.Get(today)
	fmt.Println(str)
	// prints:
	// up to tomorrow

	// Starting yesterday
	tm.AddInterval(&yesterday, nil, "since yesterday")

	str, ok := tm.GetOk(today)
	if ok {
		fmt.Println(str)
	}
	// prints:
	// since yesterday

	// Create an explicit interval
	id, _ := tm.AddInterval(&yesterday, &tomorrow, "some value")

	str, ok = tm.GetOk(today)
	if ok {
		fmt.Println(str)
	}
	// prints:
	// some value

	// Delete an interval, referenced by id
	if err := tm.RemoveIntervalID(id); err != nil {
		// Shouldn't get here
	}

	str, ok = tm.GetOk(today)
	if ok {
		fmt.Println(str)
	}
	// prints:
	// since yesterday
}
```