package timeindex

import (
    "sort"
    "time"
)

type TimeSlice []time.Time

func (ts TimeSlice) Len() int {
    return len(ts)
}

func (ts TimeSlice) Less(i, j int) bool {
    return ts[i].Before(ts[j])
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

func (ts TimeSlice) Add(t time.Time) (newTs TimeSlice) {
    i := ts.Search(t)
    if i < len(ts) && ts[i] == t {
        return
    }

    right := append([]time.Time {t}, ts[i:]...)
    newTs = append(ts[:i], right...)

    return newTs
}

func SearchTimes(ts TimeSlice, t time.Time) int {
    p := func(i int) bool {
        return ts[i].After(t) || ts[i] == t
//        return t.Before(ts[i])
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
