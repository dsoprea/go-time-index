An implementation of a `time` slice (`[]time.Time`) that contains semantics for searching, sorting a slice, and incrementally sorting when starting from an empty slice (insertion sort).

Features:

- `timeindex`.`TimeSlice` (type alias for `[]time.Time`) providing `Search`, `Sort`, and `Add` (in addition to fulfilling the `sort.Interface` interface).
- `timeindex`.`TimeSlice` also provides `SearchNearest` method to invoke a callback for all of the times near a given time and range of tolerance (expressed as a `time.Duration`).
- `timeindex`.`AbsoluteDistance`: Returns the absolute difference between two times.

See the unit-tests for examples.
