package logging

import "context"

type contextKey string

const RequestIDKey contextKey = "request_id"

func WithRequestID(
	ctx context.Context,
	requestID string,
) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func RequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}
