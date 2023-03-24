package utils

import (
	"context"
	"fmt"
)

const (
	CtxUserIDKey = "ctx_user_id"
)

func SetUserIDToCtx(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, CtxUserIDKey, userID)
}

func GetUserIDFromCtx(ctx context.Context) (string, error) {
	userIDAny := ctx.Value(CtxUserIDKey)
	if userIDAny == nil {
		return "", fmt.Errorf("user not found in ctx")
	}
	if userIDStr, ok := userIDAny.(string); !ok {
		return "", fmt.Errorf("userID is not str: %T, value: %v", userIDAny, userIDAny)
	} else {
		return userIDStr, nil
	}
}
