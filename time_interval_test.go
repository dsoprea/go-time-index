package timeindex

import (
    "testing"
    "time"
    "fmt"

    "github.com/dsoprea/go-logging"
)

func TestTimeIntervalAddDuplicate(t *testing.T) {
    left1, err := time.Parse(time.RFC3339, "2016-12-03T07:23:50Z")
    log.PanicIf(err)

    right1, err := time.Parse(time.RFC3339, "2016-12-04T07:23:50Z")
    log.PanicIf(err)

    ti1 := TimeInterval { left1, right1 }

    tis := make(TimeIntervalSlice, 0)
    
    tis = tis.Add(ti1)
    tis = tis.Add(ti1)

    if len(tis) != 1 {
        t.Fatalf("Duplicate added.")
    } else if tis[0] != ti1 {
        t.Fatalf("TIS is out of order (with duplicate?).")
    }
}

func checkTestAddOrder(description string, t *testing.T, unsortedIntervals []TimeInterval, sortedIntervals []TimeInterval) {
    if len(unsortedIntervals) != len(sortedIntervals) {
        log.Panic(fmt.Errorf("unsorted and sorted intervals do not have the same length."))
    }

    tis := make(TimeIntervalSlice, 0)
    for _, ti := range unsortedIntervals {
        tis = tis.Add(ti)
    }

    len1 := len(tis)
    len2 := len(unsortedIntervals)
    if len1 != len2 {
        log.Panic(fmt.Errorf("%s: Size of new slice does not equal size of dataset: (%d) != (%d)", description, len1, len2))
    }

    for i, ti := range sortedIntervals {
        if tis[i] != ti {
            t.Errorf("%s: Slice out of order at (%d).\n", description, i)
            t.Errorf("  UNSORTED: %v", unsortedIntervals)
            t.Errorf("    SORTED: %v", sortedIntervals)
            t.Errorf("    ACTUAL: %v", tis)

            log.Panic(fmt.Errorf("Slice not sorted after add."))
        }
    }
}

func TestTimeIntervalAddDifferentStartTimes(t *testing.T) {
    left1, err := time.Parse(time.RFC3339, "2016-12-03T07:23:50Z")
    log.PanicIf(err)

    left2, err := time.Parse(time.RFC3339, "2016-12-03T07:24:50Z")
    log.PanicIf(err)

    left3, err := time.Parse(time.RFC3339, "2016-12-03T07:25:50Z")
    log.PanicIf(err)

    right1, err := time.Parse(time.RFC3339, "2016-12-04T07:23:50Z")
    log.PanicIf(err)

    right2, err := time.Parse(time.RFC3339, "2016-12-08T07:23:50Z")
    log.PanicIf(err)

    right3, err := time.Parse(time.RFC3339, "2016-12-06T07:23:50Z")
    log.PanicIf(err)

    ti1 := TimeInterval { left1, right1 }
    ti2 := TimeInterval { left2, right2 }
    ti3 := TimeInterval { left3, right3 }

    sortedIntervals := []TimeInterval { ti1, ti2, ti3 }

    checkTestAddOrder("1", t, []TimeInterval { ti1, ti2, ti3 }, sortedIntervals)
    checkTestAddOrder("2", t, []TimeInterval { ti1, ti3, ti2 }, sortedIntervals)
    checkTestAddOrder("3", t, []TimeInterval { ti2, ti1, ti3 }, sortedIntervals)
    checkTestAddOrder("3", t, []TimeInterval { ti2, ti3, ti1 }, sortedIntervals)
    checkTestAddOrder("3", t, []TimeInterval { ti3, ti2, ti1 }, sortedIntervals)
    checkTestAddOrder("3", t, []TimeInterval { ti3, ti1, ti2 }, sortedIntervals)
}

func TestTimeIntervalAddSameStartTime(t *testing.T) {
    left1, err := time.Parse(time.RFC3339, "2016-12-03T07:23:50Z")
    log.PanicIf(err)
/*
    left2, err := time.Parse(time.RFC3339, "2016-12-03T07:24:50Z")
    log.PanicIf(err)

    left3, err := time.Parse(time.RFC3339, "2016-12-03T07:25:50Z")
    log.PanicIf(err)
*/
    right1, err := time.Parse(time.RFC3339, "2016-12-04T07:23:50Z")
    log.PanicIf(err)

    right2, err := time.Parse(time.RFC3339, "2016-12-06T07:23:50Z")
    log.PanicIf(err)

    right3, err := time.Parse(time.RFC3339, "2016-12-08T07:23:50Z")
    log.PanicIf(err)

    ti1 := TimeInterval { left1, right1 }
    ti2 := TimeInterval { left1, right2 }
    ti3 := TimeInterval { left1, right3 }

    sortedIntervals := []TimeInterval { ti1, ti2, ti3 }

    checkTestAddOrder("1", t, []TimeInterval { ti1, ti2, ti3 }, sortedIntervals)
    checkTestAddOrder("2", t, []TimeInterval { ti1, ti3, ti2 }, sortedIntervals)
    checkTestAddOrder("3", t, []TimeInterval { ti3, ti2, ti1 }, sortedIntervals)
    checkTestAddOrder("4", t, []TimeInterval { ti2, ti3, ti1 }, sortedIntervals)
}

func TestTimeIntervalAddSameAndDifferentStartTimes(t *testing.T) {
    // Minutes are different.

    left1, err := time.Parse(time.RFC3339, "2016-12-03T07:23:50Z")
    log.PanicIf(err)

    left2, err := time.Parse(time.RFC3339, "2016-12-03T07:24:50Z")
    log.PanicIf(err)

    left3, err := time.Parse(time.RFC3339, "2016-12-03T07:25:50Z")
    log.PanicIf(err)

    left4, err := time.Parse(time.RFC3339, "2016-12-03T07:26:50Z")
    log.PanicIf(err)

    left5, err := time.Parse(time.RFC3339, "2016-12-03T07:27:50Z")
    log.PanicIf(err)

    // Days are different.

    right1, err := time.Parse(time.RFC3339, "2016-12-04T07:23:50Z")
    log.PanicIf(err)

    right2, err := time.Parse(time.RFC3339, "2016-12-06T07:23:50Z")
    log.PanicIf(err)

    right3, err := time.Parse(time.RFC3339, "2016-12-08T07:23:50Z")
    log.PanicIf(err)

    ti1 := TimeInterval { left1, right1 }

    ti2 := TimeInterval { left2, right1 }
    ti3 := TimeInterval { left3, right2 }
    ti4 := TimeInterval { left4, right3 }

    ti5 := TimeInterval { left5, right1 }

    sortedIntervals := []TimeInterval { ti1, ti2, ti3, ti4, ti5 }

    // Keep the entries with the same start-times together and in order.
    checkTestAddOrder("1", t, []TimeInterval { ti1, ti2, ti3, ti4, ti5 }, sortedIntervals)
    checkTestAddOrder("2", t, []TimeInterval { ti2, ti3, ti4, ti5, ti1 }, sortedIntervals)
    checkTestAddOrder("3", t, []TimeInterval { ti5, ti1, ti2, ti3, ti4 }, sortedIntervals)

    // Move the entries around the same as before, but shuffle the ones with 
    // the same start-time.

    checkTestAddOrder("5", t, []TimeInterval { ti1, ti2, ti4, ti3, ti5 }, sortedIntervals)
    checkTestAddOrder("6", t, []TimeInterval { ti1, ti3, ti2, ti4, ti5 }, sortedIntervals)
    checkTestAddOrder("7", t, []TimeInterval { ti1, ti3, ti4, ti2, ti5 }, sortedIntervals)
    checkTestAddOrder("8", t, []TimeInterval { ti1, ti4, ti2, ti3, ti5 }, sortedIntervals)
    checkTestAddOrder("9", t, []TimeInterval { ti1, ti4, ti3, ti2, ti5 }, sortedIntervals)

    checkTestAddOrder("11", t, []TimeInterval { ti2, ti4, ti3, ti5, ti1 }, sortedIntervals)
    checkTestAddOrder("12", t, []TimeInterval { ti3, ti2, ti4, ti5, ti1 }, sortedIntervals)
    checkTestAddOrder("13", t, []TimeInterval { ti3, ti4, ti2, ti5, ti1 }, sortedIntervals)
    checkTestAddOrder("14", t, []TimeInterval { ti4, ti2, ti3, ti5, ti1 }, sortedIntervals)
    checkTestAddOrder("15", t, []TimeInterval { ti4, ti3, ti2, ti5, ti1 }, sortedIntervals)

    checkTestAddOrder("17", t, []TimeInterval { ti5, ti1, ti2, ti4, ti3 }, sortedIntervals)
    checkTestAddOrder("18", t, []TimeInterval { ti5, ti1, ti3, ti2, ti4 }, sortedIntervals)
    checkTestAddOrder("19", t, []TimeInterval { ti5, ti1, ti3, ti4, ti2 }, sortedIntervals)
    checkTestAddOrder("20", t, []TimeInterval { ti5, ti1, ti4, ti2, ti3 }, sortedIntervals)
    checkTestAddOrder("21", t, []TimeInterval { ti5, ti1, ti4, ti3, ti2 }, sortedIntervals)
}

func TestTimeIntervalSearchSimple(t *testing.T) {
    left1, err := time.Parse(time.RFC3339, "2016-12-03T07:23:50Z")
    log.PanicIf(err)

    left2, err := time.Parse(time.RFC3339, "2016-12-05T07:24:50Z")
    log.PanicIf(err)

    left3, err := time.Parse(time.RFC3339, "2016-12-07T07:25:50Z")
    log.PanicIf(err)

    right1, err := time.Parse(time.RFC3339, "2016-12-04T07:23:50Z")
    log.PanicIf(err)

    right2, err := time.Parse(time.RFC3339, "2016-12-06T07:24:50Z")
    log.PanicIf(err)

    right3, err := time.Parse(time.RFC3339, "2016-12-08T07:25:50Z")
    log.PanicIf(err)

    ti1 := TimeInterval { left1, right1 }
    ti2 := TimeInterval { left2, right2 }
    ti3 := TimeInterval { left3, right3 }

    tis := make(TimeIntervalSlice, 0)
    
    tis = tis.Add(ti1)
    tis = tis.Add(ti2)
    tis = tis.Add(ti3)

    i := tis.search(ti1)
    if i != 0 {
        t.Fatalf("Search for item in entry 0 failed.")
    }

    i = tis.search(ti2)
    if i != 1 {
        t.Fatalf("Search for item in entry 1 failed.")
    }

    i = tis.search(ti3)
    if i != 2 {
        t.Fatalf("Search for item in entry 2 failed.")
    }
}

func searchTestIntervals(description string, t *testing.T, intervals []TimeInterval, q time.Time, expectedMatches []TimeInterval) {
    tis := make(TimeIntervalSlice, 0)
    for _, ti := range intervals {
        tis = tis.Add(ti)
    }

    matches := tis.SearchAndReturn(q)

    len1 := len(matches)
    len2 := len(expectedMatches)
    if len1 != len2 {
        log.Panic(fmt.Errorf("Match count is not expected: (%d) != (%d)", len1, len2))
    }

    for i, ti := range matches {
        if ti != expectedMatches[i] {
            t.Errorf("%s: Match at (%d) unexpected: [%s] != [%s]", description, i, ti, expectedMatches[i])
            t.Errorf("  CURRENT: %v", matches)
            t.Errorf("  EXPECTED: %v", expectedMatches)

            t.Fatalf("Matches not unexpected.")
        }
    }
}

func TestTimeIntervalSearchComplex(t *testing.T) {
    left1, err := time.Parse(time.RFC3339, "2016-01-01T02:00:00Z")
    log.PanicIf(err)

    left2, err := time.Parse(time.RFC3339, "2016-01-01T02:15:00Z")
    log.PanicIf(err)

    left3, err := time.Parse(time.RFC3339, "2016-01-01T02:45:00Z")
    log.PanicIf(err)

    left4, err := time.Parse(time.RFC3339, "2016-01-02T02:00:00Z")
    log.PanicIf(err)

    left5, err := time.Parse(time.RFC3339, "2016-01-02T02:00:00Z")
    log.PanicIf(err)

    left6, err := time.Parse(time.RFC3339, "2016-01-02T02:00:00Z")
    log.PanicIf(err)

    left7, err := time.Parse(time.RFC3339, "2016-01-03T02:00:00Z")
    log.PanicIf(err)

    left8, err := time.Parse(time.RFC3339, "2016-01-03T02:00:00Z")
    log.PanicIf(err)

    left9, err := time.Parse(time.RFC3339, "2016-01-03T02:00:00Z")
    log.PanicIf(err)

    right1, err := time.Parse(time.RFC3339, "2016-01-01T04:00:00Z")
    log.PanicIf(err)

    right2, err := time.Parse(time.RFC3339, "2016-01-01T06:00:00Z")
    log.PanicIf(err)

    right3, err := time.Parse(time.RFC3339, "2016-01-01T08:00:00Z")
    log.PanicIf(err)

    right4, err := time.Parse(time.RFC3339, "2016-01-02T04:00:00Z")
    log.PanicIf(err)

    right5, err := time.Parse(time.RFC3339, "2016-01-02T06:00:00Z")
    log.PanicIf(err)

    right6, err := time.Parse(time.RFC3339, "2016-01-02T08:00:00Z")
    log.PanicIf(err)

    right7, err := time.Parse(time.RFC3339, "2016-01-03T04:00:00Z")
    log.PanicIf(err)

    right8, err := time.Parse(time.RFC3339, "2016-01-03T06:00:00Z")
    log.PanicIf(err)

    right9, err := time.Parse(time.RFC3339, "2016-01-03T08:00:00Z")
    log.PanicIf(err)

    ti1 := TimeInterval { left1, right1 }
    ti2 := TimeInterval { left2, right2 }
    ti3 := TimeInterval { left3, right3 }
    ti4 := TimeInterval { left4, right4 }
    ti5 := TimeInterval { left5, right5 }
    ti6 := TimeInterval { left6, right6 }
    ti7 := TimeInterval { left7, right7 }
    ti8 := TimeInterval { left8, right8 }
    ti9 := TimeInterval { left9, right9 }

    intervals := []TimeInterval { ti1, ti2, ti3, ti4, ti5, ti6, ti7, ti8, ti9 }

    //// Generic tests.

    // Search before all intervals.

    q, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("1", t, intervals, q, []TimeInterval {})

    // Search between two intervals (hitting on neither).

    q, err = time.Parse(time.RFC3339, "2016-01-01T09:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("2", t, intervals, q, []TimeInterval {})

    // Search after all intervals.

    q, err = time.Parse(time.RFC3339, "2016-01-03T09:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("3", t, intervals, q, []TimeInterval {})

    //// Tests within the first three intervals.

    // Search within the first three intervals.

    q, err = time.Parse(time.RFC3339, "2016-01-01T02:10:00Z")
    log.PanicIf(err)

    searchTestIntervals("4", t, intervals, q, []TimeInterval { ti1 })

    // Search for the left side of the first three intervals.

    q, err = time.Parse(time.RFC3339, "2016-01-01T02:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("5", t, intervals, q, []TimeInterval { ti1 })

    // Search for the right side of the first three interval.

    q, err = time.Parse(time.RFC3339, "2016-01-01T08:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("6", t, intervals, q, []TimeInterval { ti3 })

    // Search such that first three intervals (with staggered start-times) will 
    // be returned.

    q, err = time.Parse(time.RFC3339, "2016-01-01T03:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("7", t, intervals, q, []TimeInterval { ti1, ti2, ti3 })

    // Search such that first three intervals (with identical start-times) will
    // be returned.

    q, err = time.Parse(time.RFC3339, "2016-01-02T02:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("8", t, intervals, q, []TimeInterval { ti4, ti5, ti6 })

    //// Tests within the second set of three intervals.

    // Search in the middle.

    q, err = time.Parse(time.RFC3339, "2016-01-02T03:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("9", t, intervals, q, []TimeInterval { ti4, ti5, ti6 })

    // Search for the left side.

    q, err = time.Parse(time.RFC3339, "2016-01-02T02:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("10", t, intervals, q, []TimeInterval { ti4, ti5, ti6 })

    // Search for the right side.

    q, err = time.Parse(time.RFC3339, "2016-01-02T08:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("11", t, intervals, q, []TimeInterval { ti6 })

    //// Tests within the third set of three intervals.

    // Search in the middle.

    q, err = time.Parse(time.RFC3339, "2016-01-03T03:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("12", t, intervals, q, []TimeInterval { ti7, ti8, ti9 })

    // Search for the left side.

    q, err = time.Parse(time.RFC3339, "2016-01-03T02:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("13", t, intervals, q, []TimeInterval { ti7, ti8, ti9 })

    // Search for the right side.

    q, err = time.Parse(time.RFC3339, "2016-01-03T08:00:00Z")
    log.PanicIf(err)

    searchTestIntervals("14", t, intervals, q, []TimeInterval { ti9 })
}
