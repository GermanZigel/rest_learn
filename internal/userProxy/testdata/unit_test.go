package testdata

import (
	"rest/internal/config"
	"rest/internal/userProxy"
	"testing"
)

/*
type PgTest struct{}

func (p *PgTest) Query(ctx context.Context, args *pgx.QueryArgs) (*pgx.Rows, error) {
	return
}*/

func TestSetUser(t *testing.T) {
	cfg := config.GetConfig()
	got := userProxy.Setter()
	if got.Id < cfg.User.MinId {
		t.Errorf("MinId test failed got: %d; want id more than %d", got.Id, cfg.User.MinId)
	}
}
