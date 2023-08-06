package values_test

import (
	"reflect"
	"scheduleme/values"
	"testing"
	"time"
)

func ezTime(s string) time.Time {
	t, err := time.Parse("2006-Jan-02 15:04", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestGetTimeSlots(t *testing.T) {
	// basic func for creating values.DateSlot
	newDateSlot := func(start, end string) values.DateSlot {
		startTime, _ := time.Parse("2006-Jan-02 15:04", start)
		endTime, _ := time.Parse("2006-Jan-02 15:04", end)
		return values.DateSlot{
			Start: startTime,
			End:   endTime,
		}
	}

	// start, _ := time.Parse("2006-01-02 15:04:05", "2000-01-02 00:00:00")
	// end, _ := time.Parse("2006-01-02 15:04:05", "2022-01-03 00:00:00")

	tests := []struct {
		name     string
		mask     values.AvailMasks
		duration time.Duration
		start    time.Time
		end      time.Time
		want     values.DateSlots
	}{
		{
			name:     "Normal scenario",
			mask:     values.AvailMasks{(&values.DateSlots{newDateSlot("2000-Jan-02 17:45", "2000-Jan-02 21:00")}).ToMask(values.AvailMaskINC)},
			duration: 2 * time.Hour,
			start:    ezTime("2000-Jan-02 00:00"),
			end:      ezTime("2000-Jan-03 00:00"),
			want:     values.DateSlots{newDateSlot("2000-Jan-02 18:00", "2000-Jan-02 20:00")},
		},
		{
			name:     "Rejects because of precomputed slots fall on even times with 2 hr duration intervals from 00:00",
			mask:     values.AvailMasks{(&values.DateSlots{newDateSlot("2000-Jan-02 16:45", "2000-Jan-02 19:00")}).ToMask(values.AvailMaskINC)},
			duration: 2 * time.Hour,
			start:    ezTime("2000-Jan-02 00:00"),
			end:      ezTime("2000-Jan-03 00:00"),
			want:     values.DateSlots{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mask.GetDateSlots(tt.duration, tt.start, tt.end); !reflect.DeepEqual(*got, tt.want) {
				if len(*got) != 0 && len(tt.want) != 0 {
					t.Errorf("GetTimeSlots() = %v, want %v", *got, tt.want)
				}
			}
		})
	}
}
