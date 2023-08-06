package test

import (
	"net/http"
	"scheduleme/config"
	"scheduleme/frame"
	"scheduleme/values"
	"time"

	"flag"
	"testing"
)

// // Convenience function for integration tests for now, may require more details in more descriminate tests
func TestConfig() *config.ConfigStruct {
	return &config.ConfigStruct{
		GoogleClientSecret: "test-google-secret",
		GoogleClientID:     "test-google-id",
		GoogleRedirectURL:  "http://example/google_redirect_url",
		GoogleRedirectPath: "/google_redirect_path",
		Port:               "8080",
		ENV:                "test",
		Dsn:                "mock-dsn",
		Secret:             "mock-secret",
	}
}

func init() {
	testing.Init()
	flag.Parse()
	// config.SetConfigForTestingOnly(TestConfig())
}

func IsTestRun() bool {
	testFlag := flag.Lookup("test.v")
	if testFlag == nil {
		return false
	}
	return testFlag.Value.(flag.Getter).Get().(bool)
}

func TimeNow() time.Time {
	mockTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	return mockTime
}

func MakeMask() *values.AvailMask {
	durs := values.DurationSlots{
		values.DurationSlot{
			Start: time.Duration(0),
			End:   time.Duration(1),
		},
		values.DurationSlot{
			Start: time.Duration(2),
			End:   time.Duration(3),
		},
	}
	dats := values.DateSlots{
		values.DateSlot{
			Start: TimeNow(),
			End:   TimeNow().Add(time.Hour),
		},
		values.DateSlot{
			Start: TimeNow().Add(time.Hour),
			End:   TimeNow().Add(2 * time.Hour),
		},
	}
	days := values.DaySlots{
		values.DaySlot{
			Day: time.Monday,
			DurationSlot: &values.DurationSlot{
				Start: time.Duration(0),
				End:   time.Duration(1),
			},
		},
		values.DaySlot{
			Day: time.Tuesday,
			DurationSlot: &values.DurationSlot{
				Start: time.Duration(0),
				End:   time.Duration(1),
			},
		},
	}
	return values.NewIncMask(&durs, &dats, &days)
}

func MockContextMiddleware(mockData *frame.Contextable) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			frame.ServeWithNewContextInfo(w, r, next, mockData)
		})
	}
}
