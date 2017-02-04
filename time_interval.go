package timeindex

import (
    "time"
    "sort"
    "fmt"

    "github.com/dsoprea/go-logging"
)

type TimeInterval [2]time.Time

// TimeIntervalSlice Stores a list of two-tuples, representing "from" and "to" 
// times, sorted by the first (almost identical to TimeSlice but more 
// convenient).
type TimeIntervalSlice []TimeInterval

func (tis TimeIntervalSlice) Len() int {
    return len(tis)
}

func (tis TimeIntervalSlice) Less(i, j int) bool {
    return tis[i][0].Before(tis[j][0])
}

func (tis TimeIntervalSlice) Swap(i, j int) {
    tis[i], tis[j] = tis[j], tis[i]
}

// Sort is a convenience method.
func (tis TimeIntervalSlice) Sort() {
    sort.Sort(tis)
}

func (tis TimeIntervalSlice) search(ti TimeInterval) int {
    return SearchTimeIntervals(tis, ti)
}

func (tis TimeIntervalSlice) SearchAndReturn(t time.Time) (matches []TimeInterval) {
    matches = make([]TimeInterval, 0)

    cb := func(ti TimeInterval) (err error) {
        matches = append([]TimeInterval { ti }, matches...)
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
            if tis[i][0].After(t) {
                break
            }
        }
    }

    // This won't run if no matches (i == len_).
    for ; i >= 0; i-- {
        if tis[i][0].After(t) {
            // We're too far right in the list (the whole interval is greater 
            // than our query).

            continue
        } else if tis[i][1].Before(t) {
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

func (tis TimeIntervalSlice) getInsertLocation(ti TimeInterval) (i int) {
    len_ := len(tis)
    if len_ == 0 {
        return 0
    }

    i = tis.search(ti)

    // We were told to append.
    if i >= len_ {
        last := tis[len_- 1]

        // Either the new entry's start-time is greater than the start-time of 
        // the last element or they're equal and the new-entry's stop-time is 
        // greater.
        if last[0].Before(ti[0]) || last[0] == ti[0] && last[1].Before(ti[1]) {
            return len_
        }

        // Else, step to the left (to a valid entry rather than being past the 
        // existing list) so that we can do the comparisons below.
        i--
    }

    // We were told to insert in the middle of the list.
    for ; i > 0; i-- {
        if tis[i][0] == ti[0] {
            // The current entry's start-time matches the start-time of the new
            // entry.

            if tis[i][1] == ti[1] {
                // The current entry's stop-time matches the stop-time of the 
                // new-entry. Entry already exists.

                return -1
            } else if tis[i][1].Before(ti[1]) {
                // The current entry's stop-time is less thanthe new entry's 
                // stop-time. Insert after.

                return i + 1
            }
        } else if tis[i][0].Before(ti[0]) {
            // The current entry's start-time is less-than the start-time of 
            // the new entry (the search will only return an `i` of an entry 
            // that is greater-than-or-equal-to the entry to be added).

            return i + 1
        }
    }

    first := tis[0]

    if first == ti {
        // The first entry is identical.

        return -1
    } else if first[0].Before(ti[0]) || first[0] == ti[0] && first[1].Before(ti[1]) {
        // Either the first entry in the list has a start-time that's less 
        // than the new-entry's or the start-times are equal and the new-
        // entry's stop-time is greater.

        return 1
    }

    // If we get here, all elements in [0, `i`] had the same start-time and 
    // greater stop-times. Insert at the beginning.
    return 0
}

func (tis TimeIntervalSlice) Add(ti TimeInterval) (newTis TimeIntervalSlice) {
    if ti[0].Before(ti[1]) == false {
        log.Panic(fmt.Errorf("interval is invalid"))
    }

    i := tis.getInsertLocation(ti)

    // Already exists.
    if i == -1 {
        return tis
    }

    right := append(TimeIntervalSlice { ti }, tis[i:]...)
    newTis = append(tis[:i], right...)

    return newTis
}

func SearchTimeIntervals(tis TimeIntervalSlice, ti TimeInterval) int {
    p := func(i int) bool {
        return tis[i][0].After(ti[0]) || ti == tis[i]
    }
    
    return search(len(tis), p)
}

func SearchStartTimes(tis TimeIntervalSlice, t time.Time) int {
    p := func(i int) bool {
        return tis[i][0].After(t) || t == tis[i][0]
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
