package timeindex

import (
    "testing"
    "time"
)

func addNowToIndex(ts TimeSlice) (newTs TimeSlice, now time.Time) {
    now = time.Now()
    ts = ts.Add(now)
    time.Sleep(1 * time.Millisecond)

    return ts, now
}

func TestAdd(t *testing.T) {
    ts := make(TimeSlice, 0)
    
    ts, now := addNowToIndex(ts)
    if ts[0] != now {
        t.Errorf("Time not stored correctly (1).")
    }

    ts, now2 := addNowToIndex(ts)
    if ts[1] != now2 {
        t.Errorf("Time not sorted correctly (2).")
    }

    ts, now3 := addNowToIndex(ts)
    if ts[2] != now3 {
        t.Errorf("Time not sorted correctly (3).")
    }
}

func TestSearch(t *testing.T) {
    ts := make(TimeSlice, 0)

    ts, now1 := addNowToIndex(ts)
    ts, now2 := addNowToIndex(ts)
    ts, now3 := addNowToIndex(ts)
    
    i := SearchTimes(ts, now1)
    j := SearchTimes(ts, now2)
    k := SearchTimes(ts, now3)

    if i != 0 || j != 1 || k != 2 {
        t.Errorf("Times didn't insert correctly")
    }
}

func TestSort(t *testing.T) {
    now1 := time.Now()
    time.Sleep(1 * time.Millisecond)

    now2 := time.Now()
    time.Sleep(1 * time.Millisecond)

    now3 := time.Now()

    ts := TimeSlice { now3, now2, now1 }
    ts.Sort()

    if ts[0] != now1 || ts[1] != now2 || ts[2] != now3 {
        t.Errorf("Times didn't sort correctly.")
    }
}
