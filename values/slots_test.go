package values_test

import (
	"fmt"
	"scheduleme/values"
	"testing"
	"time"
)

func TestSlots(t *testing.T) {
	start, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 00:13:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2022-01-01 01:15:00")
	interval := time.Minute * 30
	nearest := time.Minute * 15

	dateSlots := values.Slots(start, end, interval, &nearest)

	if dateSlots[0].Start != values.RoundUp(start, time.Minute*15) {
		t.Errorf("Expected first slot to start at %v, got %v", values.RoundUp(start, time.Minute*15), dateSlots[0].Start)
	}

	for _, ds := range dateSlots {
		fmt.Printf("%+v\n", ds)
	}
	// Assuming DateSlots is a slice of DateSlot
	// And DateSlot is a struct with Start and End time.Time fields
	if len(dateSlots) != 2 { // there are 24 hours from start till end
		t.Errorf("Expected 24 slots, got %v", len(dateSlots))
	}

	for i, slot := range dateSlots { // check if each slot is 1 hour long and rounded to nearest 15 minutes
		expectedStart := start.Add(time.Duration(i) * interval).Round(nearest)
		expectedEnd := expectedStart.Add(interval)

		if !slot.Start.Equal(expectedStart) {
			t.Errorf("Expected slot %v to start at %v, got %v", i, expectedStart, slot.Start)
		}
		if !slot.End.Equal(expectedEnd) {
			t.Errorf("Expected slot %v to end at %v, got %v", i, expectedEnd, slot.End)
		}
	}
}
