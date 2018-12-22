package timeindex

import (
	"fmt"
	"sort"
	"time"

	"github.com/dsoprea/go-logging"
)

type TimeInterval struct {
	From  time.Time
	To    time.Time
	Items []interface{}
}

// TimeIntervalSlice Stores a list of two-tuples, representing "from" and "to"
// times, sorted by the first (almost identical to TimeSlice but more
// convenient).
type TimeIntervalSlice []TimeInterval

func (tis TimeIntervalSlice) Len() int {
	return len(tis)
}

func (tis TimeIntervalSlice) Less(i, j int) bool {
	return tis[i].From.Before(tis[j].From)
}

func (tis TimeIntervalSlice) Swap(i, j int) {
	tis[i], tis[j] = tis[j], tis[i]
}

// Sort is a convenience method.
func (tis TimeIntervalSlice) Sort() {
	sort.Sort(tis)
}

func (tis TimeIntervalSlice) search(from time.Time, to time.Time) int {
	return SearchTimeIntervals(tis, from, to)
}

func (tis TimeIntervalSlice) SearchAndReturn(t time.Time) (matches []TimeInterval) {
	matches = make([]TimeInterval, 0)

	cb := func(ti TimeInterval) (err error) {
		matches = append([]TimeInterval{ti}, matches...)
		return nil
	}

	if err := tis.Search(t, cb); err != nil {
		log.Panic(err)
	}

	return matches
}

// Search Call the callback with all intervals that contain the given time.
func (tis TimeIntervalSlice) Search(t time.Time, cb func(ti TimeInterval) error) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	len_ := len(tis)
	if len_ == 0 {
		return nil
	}

	// If there is more than one match, this will point to the first.
	i := SearchStartTimes(tis, t)

	if i >= len_ {
		// We're using the search to find a starting point in the list. It doesn't
		// have to be perfect. Therefore, if we're past the end, move to the left
		// by one.

		i--
	} else {
		// The interval we're pointing at might have a start-time equal to what
		// we're searching for. If there is more than one then step to the end
		// (we need to return all intervals that match).

		for ; i < (len_ - 1); i++ {
			if tis[i].From.After(t) {
				break
			}
		}
	}

	// This won't run if no matches (i == len_).
	for ; i >= 0; i-- {
		if tis[i].From.After(t) {
			// We're too far right in the list (the whole interval is greater
			// than our query).

			continue
		} else if tis[i].To.Before(t) {
			// We've moved too far left in the list (the whole intervals is
			// smaller than our query).

			break
		}

		// We're working our way to the front of the sorted list, so prepend
		// the results (so that we maintain order).
		if err := cb(tis[i]); err != nil {
			log.Panic(err)
		}
	}

	return nil
}

// getInsertLocation returns two integers: If the exact interval already exists,
// the first integer is the position and the second is (-1). Else, the second
// integer is the position to insert at and the first integer is (-1).
func (tis TimeIntervalSlice) getInsertLocation(from time.Time, to time.Time) (foundAt int, insertAt int) {
	len_ := len(tis)
	if len_ == 0 {
		return -1, 0
	}

	i := tis.search(from, to)

	// We were told to append.
	if i >= len_ {
		last := tis[len_-1]

		// Either the new entry's start-time is greater than the start-time of
		// the last element or they're equal and the new-entry's stop-time is
		// greater.
		if last.From.Before(from) || last.From == from && last.To.Before(to) {
			return -1, len_
		}

		// Else, step to the left (to a valid entry rather than being past the
		// existing list) so that we can do the comparisons below.
		i--
	}

	// We were told to insert in the middle of the list.
	for ; i > 0; i-- {
		if tis[i].From == from {
			// The current entry's start-time matches the start-time of the new
			// entry.

			if tis[i].To == to {
				// The current entry's stop-time matches the stop-time of the
				// new-entry. Entry already exists.

				return i, -1
			} else if tis[i].To.Before(to) {
				// The current entry's stop-time is less thanthe new entry's
				// stop-time. Insert after.

				return -1, i + 1
			}
		} else if tis[i].From.Before(from) {
			// The current entry's start-time is less-than the start-time of
			// the new entry (the search will only return an `i` of an entry
			// that is greater-than-or-equal-to the entry to be added).

			return -1, i + 1
		}
	}

	first := tis[0]

	if first.From == from && first.To == to {
		// The first entry is identical.

		return 0, -1
	} else if first.From.Before(from) || first.From == from && first.To.Before(to) {
		// Either the first entry in the list has a start-time that's less
		// than the new-entry's or the start-times are equal and the new-
		// entry's stop-time is greater.

		return -1, 1
	}

	// If we get here, all elements in [0, `i`] had the same start-time and
	// greater stop-times. Insert at the beginning.
	return -1, 0
}

func (tis TimeIntervalSlice) Add(from time.Time, to time.Time, data interface{}) (newTis TimeIntervalSlice) {
	if from.Before(to) == false {
		log.Panic(fmt.Errorf("interval is invalid"))
	}

	foundAt, insertAt := tis.getInsertLocation(from, to)

	// Already exists.
	if insertAt == -1 {
		if data != nil {
			tis[foundAt].Items = append(tis[foundAt].Items, data)
		}

		return tis
	}

	ti := TimeInterval{
		From:  from,
		To:    to,
		Items: []interface{}{},
	}

	if data != nil {
		ti.Items = []interface{}{data}
	}

	right := append(TimeIntervalSlice{ti}, tis[insertAt:]...)
	newTis = append(tis[:insertAt], right...)

	return newTis
}

func SearchTimeIntervals(tis TimeIntervalSlice, from time.Time, to time.Time) int {
	p := func(i int) bool {
		return tis[i].From.After(from) || from == tis[i].From && to == tis[i].To
	}

	return search(len(tis), p)
}

func SearchStartTimes(tis TimeIntervalSlice, t time.Time) int {
	p := func(i int) bool {
		return tis[i].From.After(t) || t == tis[i].From
	}

	return search(len(tis), p)
}

func search(n int, f func(int) bool) int {
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
