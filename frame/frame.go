package frame

import (
	"context"
	"net/http"

	"scheduleme/values"
)

type CtxState = values.CtxState

// Return a unique string key for context
type Contextable interface {
	ContextKey() string
}

func NewContextWith[T Contextable](ctx context.Context, m *T) context.Context {
	//lint:ignore SA1029 Collision is not probabable, contextables keys are named after their struct
	return context.WithValue(ctx, (*m).ContextKey(), m)
	//what if : return context.WithValue(ctx, &m{}, m)

}

func FromContextKey[T Contextable](ctxK string, ctx context.Context) *T {
	var n T
	m, ok := ctx.Value(ctxK).(*T)
	if !ok {
		return &n
	}
	return m
}

func FromContext[T Contextable](ctx context.Context) *T {
	var n T
	m, ok := ctx.Value(n.ContextKey()).(*T)
	if !ok {
		return &n
	}
	return m
}

func FromContextAnyKey[T any](ctx context.Context) *T {
	var n T
	m, ok := ctx.Value(n).(*T)
	if !ok {
		return &n
	}
	return m
}

// Used to modify and then serve a requests Contextable member of its context
// A shortcut for:
// next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), infoStruct.ContextKey(), infoStruct)))
func ServeWithNewContextInfo[T Contextable](
	w http.ResponseWriter,
	r *http.Request,
	next http.Handler,
	info *T,
) {
	next.ServeHTTP(w, r.WithContext(NewContextWith(r.Context(), info)))
}

//	func ModifyContextAndServeWith[T Contextable](w http.ResponseWriter, r *http.Request, next http.Handler, fn func(*T)) {
//		next.ServeHTTP(w, r.WithContext(ModifyContextWith(r.Context(), fn)))
//	}
func ModifyContextWith[T Contextable](ctx context.Context, fn func(*T)) context.Context {
	newInfo := FromContext[T](ctx)
	fn(newInfo)
	return NewContextWith(ctx, newInfo)
}
