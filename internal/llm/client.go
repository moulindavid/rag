package llm

import "context"

type Client interface {
	Complete(ctx context.Context, system string, user string) (string, error)
}
