package ctxkeys

type CtxKey int8

const (
	CtxKeyUser CtxKey = iota
	CtxKeyRequsetID
)
