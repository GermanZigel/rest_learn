package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"rest/internal/user"
	"testing"
)

/*
type PgTest struct{}

func (p *PgTest) Query(ctx context.Context, args *pgx.QueryArgs) (*pgx.Rows, error) {
	return
}*/

func DeleteUserTest(t *testing.T) {
	(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	resultCode := user.Handler{}.DeleteUser()
	if resultCode != 204 {
		t.Errorf("DeleteUserTest failed, expected 204, got %d", resultCode)
	}
}
