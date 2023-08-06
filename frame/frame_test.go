package frame_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"scheduleme/frame"
	"testing"
)

type myRouteInfo struct {
	Value string
}

func (ri myRouteInfo) ContextKey() string {
	return "myRouteInfo"
}

func TestNewContextWith(t *testing.T) {
	ctx := context.Background()
	oldRi := &myRouteInfo{Value: "test"}

	ctx = frame.NewContextWith(ctx, oldRi)

	ori := ctx.Value("myRouteInfo").(*myRouteInfo)
	if got, want := ori.Value, "test"; got != want {
		t.Errorf("Expected Value to be '%s', got '%s'", want, got)
	}
}

func TestFromAnyContext(t *testing.T) {
	// ctx := context.WithValue(context.Background(), "myRouteInfo", &myRouteInfo{Value: "test"})
	ctx := context.WithValue(context.Background(), myRouteInfo{}, &myRouteInfo{Value: "test"})
	myri := frame.FromContextAnyKey[myRouteInfo](ctx)
	if got, want := myri.Value, "test"; got != want {
		t.Errorf("Expected Value to be '%s', got '%s'", want, got)
	}
}

func TestFromContext(t *testing.T) {
	// ctx := context.WithValue(context.Background(), "myRouteInfo", &myRouteInfo{Value: "test"})
	ctx := context.WithValue(context.Background(), myRouteInfo{}.ContextKey(), &myRouteInfo{Value: "test"})
	myri := frame.FromContext[myRouteInfo](ctx)
	if got, want := myri.Value, "test"; got != want {
		t.Errorf("Expected Value to be '%s', got '%s'", want, got)
	}
}

func TestModifyContextWith(t *testing.T) {
	ctx := frame.ModifyContextWith(context.Background(), func(myri *myRouteInfo) {
		myri.Value = "test"
	})
	//Retrieve the RouteInfo from the context
	myri := frame.FromContext[myRouteInfo](ctx)

	if got, want := myri.Value, "test"; got != want {
		t.Errorf("Expected Value to be '%s', got '%s'", want, got)
	}

}

func TestServeWithNewContextInfo(t *testing.T) {
	// Create a context without RouteInfo
	ctx := context.Background()

	// Create a new info
	myri := myRouteInfo{
		Value: "test",
	}

	// Create a new request
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(ctx)

	// Create a new response writer
	w := httptest.NewRecorder()

	// Create a new handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ri := frame.FromContext[myRouteInfo](r.Context())
		if got, want := ri.Value, "test"; got != want {
			t.Errorf("Expected value '%s', got '%s'", want, got)
		}
	})

	// Call ServeWithNewContextInfo
	frame.ServeWithNewContextInfo(w, r, next, &myri)
}

// test state comparison
func TestCompareStates(t *testing.T) {
	//compare two states
	state1 := frame.CtxState("state1")
	blankState := frame.CtxState("")
	if got, want := state1.CompareStates("state2"), false; got != want {
		t.Fatalf("CompareStates=%v, want %v", got, want)
	}
	if got, want := blankState.CompareStates(""), false; got != want {
		t.Fatalf("CompareStates=%v, want %v", got, want)
	}
	if got, want := blankState.CompareStates("state2"), false; got != want {
		t.Fatalf("CompareStates=%v, want %v", got, want)
	}
	if got, want := state1.CompareStates(""), false; got != want {
		t.Fatalf("CompareStates=%v, want %v", got, want)
	}
	if got, want := state1.CompareStates("state1"), true; got != want {
		t.Fatalf("CompareStates=%v, want %v", got, want)
	}
}
