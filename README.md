[![Build Status](https://travis-ci.org/dsoprea/go-time-index.svg?branch=master)](https://travis-ci.org/dsoprea/go-time-index)
[![Coverage Status](https://coveralls.io/repos/github/dsoprea/go-time-index/badge.svg?branch=master)](https://coveralls.io/github/dsoprea/go-time-index?branch=master)
[![GoDoc](https://godoc.org/github.com/dsoprea/go-time-index?status.svg)](https://godoc.org/github.com/dsoprea/go-time-index)


# Overview

An implementation of a `time` slice (`[]time.Time`) that contains semantics for searching, sorting a slice, and incrementally sorting when starting from an empty slice (insertion sort).


## Features

- `timeindex`.`TimeSlice` (type alias for `[]time.Time`) providing `Search`, `Sort`, and `Add` (in addition to fulfilling the `sort.Interface` interface).
- `timeindex`.`TimeSlice` also provides `SearchNearest` method to invoke a callback for all of the times near a given time and range of tolerance (expressed as a `time.Duration`).
- `timeindex`.`AbsoluteDistance`: Returns the absolute difference between two times.

See the unit-tests for examples.
