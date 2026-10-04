package middleware

import (
	"fmt"
	"testing"

	"github.com/tung-dnt/pagoda/pkg/context"
	pgdb "github.com/tung-dnt/pagoda/pkg/postgres/db"
	"github.com/tung-dnt/pagoda/pkg/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUser(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	ctx.SetParamNames("user")
	ctx.SetParamValues(fmt.Sprintf("%d", usr.ID))
	_ = tests.ExecuteMiddleware(ctx, LoadUser(c.Queries))
	ctxUsr, ok := ctx.Get(context.UserKey).(*pgdb.User)
	require.True(t, ok)
	assert.Equal(t, usr.ID, ctxUsr.ID)
}
