package middleware

import "context"

type Middleware struct {
	ctx       context.Context
	PublicKey string
}

func InitMiddleware(ctx context.Context, publicKey string) *Middleware {
	m := &Middleware{
		ctx:       ctx,
		PublicKey: publicKey,
	}
	return m
}
