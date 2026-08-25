package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tung-dnt/meme-app/config"
	"github.com/tung-dnt/meme-app/pkg/context"
	"github.com/tung-dnt/meme-app/pkg/tests"
)

func TestConfig(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	cfg := &config.Config{}
	err := tests.ExecuteMiddleware(ctx, Config(cfg))
	require.NoError(t, err)

	got, ok := ctx.Get(context.ConfigKey).(*config.Config)
	require.True(t, ok)
	assert.Same(t, got, cfg)
}
