package timeindex

import (
	"errors"
	"sort"
	"time"

	"github.com/dsoprea/go-logging"
)

var (
	ErrNotFound = errors.New("not found")
)

type TimeEntry struct {
	Time  time.Time
	Items []interface{}
}

func (te TimeEntry) IsZero() bool {
	return te.Time.IsZero()
}

type TimeSlice []TimeEntry

func (ts TimeSlice) Len() int {
	return len(ts)
}

func (ts TimeSlice) Less(i, j int) bool {
	return ts[i].Time.Before(ts[j].Time)
}

func (ts TimeSlice) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

// Sort is a convenience method.
func (ts TimeSlice) Sort() {
	sort.Sort(ts)
}

func (ts TimeSlice) Search(t time.Time) int {
	return SearchTimes(ts, t)
}

func (ts TimeSlice) Add(t time.Time, data interface{}) (newTs TimeSlice) {
	i := ts.Search(t)
	if i < len(ts) && ts[i].Time == t {
		if data != nil {
			ts[i].Items = append(ts[i].Items, data)
		}

		return
	}

	items := []interface{}{}
	if data != nil {
		items = []interface{}{data}
	}

	newTimeEntry := TimeEntry{
		Time:  t,
		Items: items,
	}

	right := append(TimeSlice{newTimeEntry}, ts[i:]...)
	newTs = append(ts[:i], right...)

	return newTs
}

func SearchTimes(ts TimeSlice, t time.Time) int {
	p := func(i int) bool {
		return ts[i].Time.After(t) || t == ts[i].Time
	}

	return SearchTime(len(ts), p)
}

func SearchTime(n int, f func(int) bool) int {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, n
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if !f(h) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}

	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

func AbsoluteDistance(a, b time.Time) time.Duration {
	if b.Before(a) {
		return a.Sub(b)
	} else {
		return b.Sub(a)
	}
}

func (ts TimeSlice) SearchNearest(t time.Time, tolerance time.Duration, cb func(t time.Time) error) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	if len(ts) == 0 {
		return ErrNotFound
	}

	i := ts.Search(t)

	if i >= len(ts) {
		i--
	}

	// `i` is currently equal-to-or-just-greater to the time we searched for.
	// Step the left to the earliest time that falls within tolerance of the
	// search time. If there is none, return.
	if i > 0 {
		didMove := false
		for ; i >= 0; i-- {
			if AbsoluteDistance(t, ts[i].Time) > tolerance {
				break
			}

			didMove = true
		}

		// We're out of tolerance.
		if AbsoluteDistance(t, ts[i].Time) > tolerance {
			if didMove {
				// We found at least one match but then moved out of tolerance
				// to the left. Step back to the right.
				i++
			} else {
				// No matches.
				return nil
			}
		}
	}

	// We now have the earliest timestamp in the list that is within tolerance.
	// Step forward until we find the other end.
	for ; i < len(ts); i++ {
		if AbsoluteDistance(ts[i].Time, t) > tolerance {
			break
		}

		if err := cb(ts[i].Time); err != nil {
			panic(err)
		}
	}

	return nil
}
