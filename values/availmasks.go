package values

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"errors"
	"fmt"
	"scheduleme/hof"
	"time"
)

const AvailMaskEXC AvailMaskType = "EXC"
const AvailMaskINC AvailMaskType = "INC"

func BusyTimesToMask(ds *DateSlots) *AvailMask {
	return &AvailMask{Type: AvailMaskEXC, Dates: ds}
}

func NewIncMask(durations *DurationSlots, times *DateSlots, days *DaySlots) *AvailMask {
	return &AvailMask{Type: AvailMaskINC, Durations: durations, Dates: times, Days: days}
}

func NewMask(t AvailMaskType, durations *DurationSlots, times *DateSlots, days *DaySlots) *AvailMask {
	return &AvailMask{Type: t, Durations: durations, Dates: times, Days: days}
}

func (am *AvailMask) IsExc() bool {
	return am.Type == AvailMaskEXC
}

func (am *AvailMask) IsInc() bool {
	return am.Type == AvailMaskINC
}

type AvailMaskType string

type AvailMasks []*AvailMask

// Exact date range
type DateSlot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Duration on everyday of the week
// time, 2pm-4pm for example
type DurationSlot struct {
	Start time.Duration `json:"start"`
	End   time.Duration `json:"end"`
}

// Duration Slot on a specified day of week
type DaySlot struct {
	Day time.Weekday `json:"day"`
	*DurationSlot
}

type AvailMask struct {
	Type      AvailMaskType  `json:"type"`
	Durations *DurationSlots `json:"durations"`
	Dates     *DateSlots     `json:"dates"`
	Days      *DaySlots      `json:"days"`
	// FloatDuration time.Duration  `json:"float_duration"`
	// Floating      bool           `json:"floating"`
}

func (ds *DaySlot) Validate() (err error) {
	if (ds.Start > ds.End) || (ds.Start == ds.End) {
		return fmt.Errorf("DaySlot.Validate(), Start: %v > End: %v", ds.Start, ds.End)
	}
	return
}

type DurationSlots []DurationSlot
type DateSlots []DateSlot
type DaySlots []DaySlot

// Exact Date Range, Start: Oct 21, 2023 2pm, End: Oct 21, 2023 4pm
type DateRange = DateSlot

// Exact Date Ranges
type DateRanges = DateSlots

func (ds *DateSlot) overlaps(dr *DateRange) bool {
	return ds.Start.Before(dr.End) && dr.Start.Before(ds.End)
}

func (dss *DateSlots) Withholds(d *DateRange) bool {
	//TODO: ranges should be merged
	for _, ds := range *dss {
		if d.within(&ds) {
			return true
		}
	}
	return false
}

func (ds *DateSlots) overlaps(dr *DateRange) bool {
	for i := 0; i < len(*ds); i++ {
		if (*ds)[i].overlaps(dr) {
			return true
		}
	}
	return false
}

func (dss *DurationSlots) overlaps(dr *DateRange) bool {
	for i := 0; i < len(*dss); i++ {
		if (*dss)[i].overlaps(dr) {
			return true
		}
	}
	return false
}

func (ds *DurationSlot) Withholds(dr *DateRange) bool {
	// Check if the time ranges overlap
	dateRangeStart := dr.Start.Hour()*60 + dr.Start.Minute()
	dateRangeEnd := dr.End.Hour()*60 + dr.End.Minute()
	slotStart := int(ds.Start.Minutes())
	slotEnd := int(ds.End.Minutes())
	return dateRangeStart >= slotStart && dateRangeEnd <= slotEnd
}

func (ds *DurationSlots) Withholds(dr *DateRange) bool {
	for i := 0; i < len(*ds); i++ {
		if (*ds)[i].Withholds(dr) {
			return true
		}
	}
	return false
}

func (ds *DurationSlot) overlaps(dr *DateRange) bool {
	// Check if the time ranges overlap
	dateRangeStart := dr.Start.Hour()*60 + dr.Start.Minute()
	dateRangeEnd := dr.End.Hour()*60 + dr.End.Minute()
	slotStart := int(ds.Start.Minutes())
	slotEnd := int(ds.End.Minutes())
	return dateRangeStart < slotEnd && slotStart < dateRangeEnd
}

func (dss *DaySlots) withholds(dr *DateRange) bool {
	return hof.Any(*dss, func(ds DaySlot) bool {
		return ds.Withholds(dr)
	})
}

func (dss *DaySlots) overlaps(dr *DateRange) bool {
	return hof.Any(*dss, func(ds DaySlot) bool {
		return ds.overlaps(dr)
	})
}

func (ds *DateSlots) ToMask(t AvailMaskType) *AvailMask {
	return &AvailMask{Dates: ds, Type: t}
}

func (am *AvailMasks) Overlaps(dr *DateRange) bool {
	return hof.Any(*am, func(a *AvailMask) bool {
		return (*a).Overlaps(dr)
	})
}

func (am *AvailMasks) FilterType(amt AvailMaskType) *AvailMasks {
	hof.Filter(*am, func(a *AvailMask) bool {
		return a.Type == amt
	})
	//filter am by type
	var masks AvailMasks
	for i := 0; i < len(*am); i++ {
		if (*am)[i].Type == amt {
			masks = append(masks, (*am)[i])
		}
	}
	if len(masks) == 0 {
		return &AvailMasks{}
	}
	return &masks
}

func (am *AvailMask) Overlaps(dr *DateRange) bool {
	if am.Dates != nil && am.Dates.overlaps(dr) {
		return true
	}

	if am.Days != nil && am.Days.overlaps(dr) {
		return true
	}

	if am.Durations != nil && am.Durations.overlaps(dr) {
		return true
	}
	return false
}

func (dr *DateRange) within(d *DateRange) bool {
	return (dr.Start == d.Start || dr.Start.After(d.Start)) &&
		(dr.End.Before(d.End) || dr.End == d.End)
}

func (ams *AvailMasks) Withholds(dr *DateRange) bool {
	return hof.Any(*ams, func(am *AvailMask) bool {
		return am.Withholds(dr)
	})
}
func (am *AvailMask) Withholds(dr *DateRange) bool {
	if am.Dates != nil && am.Dates.Withholds(dr) {
		return true
	}

	if am.Days != nil && am.Days.withholds(dr) {
		return true
	}

	if am.Durations != nil && am.Durations.Withholds(dr) {
		return true
	}
	return false
}

func (am *AvailMasks) AppendDateSlot(t AvailMaskType, ds *DateSlot) {
	*am = append(*am, NewMask(t, nil, &DateSlots{*ds}, nil))
}

/*
Given the following:

Slots(12:26, 13:35, 15min, 15min)
  - Given the start of 12:26 and end of 13:55
    nearest 15min is 12:30
  - 12:30
  - 12:45
  - 13:00
  - 13:15
  - 13:30

Slots(12:46, 13:33, 30min, 30min)
Slots(13:55, 14:55, 45min, 45min)
*/

func RoundUp(t time.Time, dur time.Duration) time.Time {
	// Set a default "rounder" as zero duration. It represents the extra duration we need to add
	// to the time 't' in order to "round it up" to the nearest multiple of 'dur'.
	rounder := time.Duration(0)

	// Calculate the remainder of the time 't' divided by the duration 'dur'
	remainder := t.UnixNano() % dur.Nanoseconds()

	// If the remainder is not zero, it means that 't' is not exactly on a multiple of 'dur'
	// So, we need to calculate the additional duration to add to 't' to "round it up"
	// That would be 'dur - remainder'
	// For example, if 't' is 15 minutes past an hour, 'dur' is 1 hour, and we want to round up 't' to the next hour mark,
	// remainder would be 15min (900 seconds), and 'rounder' would be 'dur - remainder' = 60min - 15min = 45min
	// So, we need to add 45mins more to 't' to make it a whole hour.
	if remainder != 0 {
		rounder = dur - time.Duration(remainder)
	}

	// Add the "rounder" on to 't' to round it up to the next multiple of 'dur'
	// If there was no remainder (meaning 't' was already at a multiple of 'dur'), 'rounder' is 0 and 't' is already correct
	// If there was a remainder, 'rounder' is the extra duration we need to add to make 't' reach the next 'dur'
	return t.Add(rounder)
}

func Slots(start time.Time, end time.Time, interval time.Duration, nearest *time.Duration) DateSlots {
	var dateSlots DateSlots
	t := start
	for i := 0; t.Before(end); t = t.Add(interval) {
		if nearest != nil {
			t = RoundUp(t, *nearest)
		}
		nextEnd := t.Add(interval)
		if nextEnd.After(end) {
			break
		}
		dateSlots = append(dateSlots, DateSlot{Start: t, End: nextEnd})
		i++
	}
	return dateSlots
}

func (am *AvailMasks) GetDateSlots(d time.Duration, start time.Time, end time.Time) *DateSlots {

	incMasks := am.FilterType(AvailMaskINC)
	excMasks := am.FilterType(AvailMaskEXC)

	dateSlots := Slots(start, end, d, &d)

	// If INC masks, filter out all slots that don't explicitly overlap. Otherwise, use all slots.
	var incExcSlots DateSlots
	if len(*incMasks) != 0 {
		for _, ds := range dateSlots {
			// if incMasks.Overlaps(&ds) && incMasks.Withholds(&ds) {
			if incMasks.Withholds(&ds) {
				incExcSlots = append(incExcSlots, ds)
			}
		}
	} else {
		incExcSlots = dateSlots
	}

	// If EXC masks, filter out all slots that overlap.
	if len(*excMasks) != 0 {
		var filteredSlots DateSlots
		for _, ds := range incExcSlots {
			// if !excMasks.Overlaps(&ds) || excMasks.Withholds(&ds) {
			if !excMasks.Overlaps(&ds) {
				filteredSlots = append(filteredSlots, ds)
			}
		}
		incExcSlots = filteredSlots
	}

	return &incExcSlots // No need to return pointer as slices are already reference types
}

func (am *AvailMasks) Value() (driver.Value, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(am)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (am *AvailMasks) Scan(source interface{}) error {
	src, ok := source.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}
	var a AvailMasks
	gob.NewDecoder(bytes.NewReader(src)).Decode(&a)
	*am = a
	return nil
}
