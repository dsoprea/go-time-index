package timeindex

import (
	"testing"
	"time"

	"github.com/dsoprea/go-logging"
)

func TestAdd(t *testing.T) {
	ts := make(TimeSlice, 0)

	time1, err := time.Parse(time.RFC3339, "2016-12-02T08:05:44Z")
	log.PanicIf(err)

	time2, err := time.Parse(time.RFC3339, "2016-12-02T09:05:44Z")
	log.PanicIf(err)

	time3, err := time.Parse(time.RFC3339, "2016-12-02T10:05:44Z")
	log.PanicIf(err)

	ts = ts.Add(time2, nil)
	ts = ts.Add(time3, nil)
	ts = ts.Add(time1, nil)

	if ts[0].Time != time1 || ts[1].Time != time2 || ts[2].Time != time3 {
		t.Fatalf("Times not sorted correctly.")
	}
}

func TestSearch(t *testing.T) {
	ts := make(TimeSlice, 0)

	time1, err := time.Parse(time.RFC3339, "2016-12-02T08:05:44Z")
	log.PanicIf(err)

	time2, err := time.Parse(time.RFC3339, "2016-12-02T09:05:44Z")
	log.PanicIf(err)

	time3, err := time.Parse(time.RFC3339, "2016-12-02T10:05:44Z")
	log.PanicIf(err)

	ts = ts.Add(time2, nil)
	ts = ts.Add(time3, nil)
	ts = ts.Add(time1, nil)

	i := SearchTimes(ts, time1)
	j := SearchTimes(ts, time2)
	k := SearchTimes(ts, time3)

	if i != 0 || j != 1 || k != 2 {
		t.Fatalf("Times didn't search correctly")
	}
}

func TestSort(t *testing.T) {
	ts := make(TimeSlice, 0)

	time1, err := time.Parse(time.RFC3339, "2016-12-02T08:05:44Z")
	log.PanicIf(err)

	time2, err := time.Parse(time.RFC3339, "2016-12-02T09:05:44Z")
	log.PanicIf(err)

	time3, err := time.Parse(time.RFC3339, "2016-12-02T10:05:44Z")
	log.PanicIf(err)

	ts = ts.Add(time2, nil)
	ts = ts.Add(time3, nil)
	ts = ts.Add(time1, nil)

	ts = TimeSlice{
		TimeEntry{Time: time3},
		TimeEntry{Time: time2},
		TimeEntry{Time: time1},
	}

	ts.Sort()

	if ts[0].Time != time1 || ts[1].Time != time2 || ts[2].Time != time3 {
		t.Fatalf("Times didn't sort correctly.")
	}
}

func TestAbsoluteDistance(t *testing.T) {
	time1, err := time.Parse(time.RFC3339, "2016-12-02T08:05:44Z")
	log.PanicIf(err)

	time2, err := time.Parse(time.RFC3339, "2016-12-02T09:05:44Z")
	log.PanicIf(err)

	if AbsoluteDistance(time1, time2) != AbsoluteDistance(time2, time1) {
		t.Fatalf("Absolute-distance calculation failed.")
	}
}

func TestSearchNearest(t *testing.T) {
	time1, err := time.Parse(time.RFC3339, "2016-12-02T08:05:44Z")
	log.PanicIf(err)

	time2, err := time.Parse(time.RFC3339, "2016-12-02T08:06:45Z")
	log.PanicIf(err)

	time3, err := time.Parse(time.RFC3339, "2016-12-02T08:07:46Z")
	log.PanicIf(err)

	time4, err := time.Parse(time.RFC3339, "2016-12-02T08:12:47Z")
	log.PanicIf(err)

	time5, err := time.Parse(time.RFC3339, "2016-12-02T08:13:48Z")
	log.PanicIf(err)

	time6, err := time.Parse(time.RFC3339, "2016-12-02T08:14:49Z")
	log.PanicIf(err)

	time7, err := time.Parse(time.RFC3339, "2016-12-02T08:26:50Z")
	log.PanicIf(err)

	time8, err := time.Parse(time.RFC3339, "2016-12-02T08:27:51Z")
	log.PanicIf(err)

	time9, err := time.Parse(time.RFC3339, "2016-12-02T08:28:52Z")
	log.PanicIf(err)

	ts := make(TimeSlice, 0)

	ts = ts.Add(time1, nil)
	ts = ts.Add(time2, nil)
	ts = ts.Add(time3, nil)
	ts = ts.Add(time4, nil)
	ts = ts.Add(time5, nil)
	ts = ts.Add(time6, nil)
	ts = ts.Add(time7, nil)
	ts = ts.Add(time8, nil)
	ts = ts.Add(time9, nil)

	least, _, n := getNearest(ts, time1, time.Minute*1)
	if n != 1 || least != time1 {
		t.Fatalf("Search failed.")
	}

	least, most, n := getNearest(ts, time1, time.Minute*2)
	if n != 2 || least != time1 || most != time2 {
		t.Fatalf("Wider search failed.")
	}

	// We should receive the whole series of adjacent times.
	least, most, n = getNearest(ts, time1, time.Minute*5)
	if n != 3 || least != time1 || most != time3 {
		t.Fatalf("Block search failed.")
	}

	// We should now receive the first time of the next series.
	least, most, n = getNearest(ts, time1, time.Minute*8)
	if n != 4 || least != time1 || most != time4 {
		t.Fatalf("Expanded block search for adjacent time failed.")
	}

	// Search in-between two valid times, where both should be within
	// tolerance.
	q, err := time.Parse(time.RFC3339, "2016-12-02T08:10:00Z")
	log.PanicIf(err)
	least, most, n = getNearest(ts, q, time.Minute*3)
	if n != 2 || least != time3 || most != time4 {
		t.Fatalf("In-between search failed.")
	}

	// Search in-between two valid times, where only the right should be within
	// tolerance.
	q, err = time.Parse(time.RFC3339, "2016-12-02T08:11:00Z")
	log.PanicIf(err)
	least, most, n = getNearest(ts, q, time.Minute*3)
	if n != 2 || least != time4 || most != time5 {
		t.Fatalf("In-between search (with bias to right) failed.")
	}
}

func getNearest(ts TimeSlice, t time.Time, tolerance time.Duration) (least, most time.Time, n int) {
	cb := func(t time.Time) error {
		n++

		if least.After(t) || least.IsZero() {
			least = t
		}

		if most.Before(t) || most.IsZero() {
			most = t
		}

		return nil
	}

	if err := ts.SearchNearest(t, tolerance, cb); err != nil {
		log.Panic(err)
	}

	return least, most, n
}
