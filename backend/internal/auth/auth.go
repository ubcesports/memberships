package auth

import (
	"context"
	"database/sql"
	"os"

	"github.com/google/uuid"
	"github.com/thecodearcher/limen"
	sqladapter "github.com/thecodearcher/limen/adapters/sql"
	"github.com/thecodearcher/limen/plugins/oauth"
	oauthgoogle "github.com/thecodearcher/limen/plugins/oauth-google"
)

type uuidGenerator struct{}

func (g *uuidGenerator) GetColumnType() limen.ColumnType { return limen.ColumnTypeUUID }
func (g *uuidGenerator) Generate(_ context.Context) (any, error) {
	return uuid.New().String(), nil
}

func NewLimen(db *sql.DB) (*limen.Limen, error) {
	schema := limen.NewDefaultSchemaConfig(
		limen.WithSchemaIDGenerator(&uuidGenerator{}),
	)

	return limen.New(&limen.Config{
		Database: sqladapter.NewPostgreSQL(db),
		Schema:   schema,
		Plugins: []limen.Plugin{
			oauth.New(
				oauth.WithProviders(oauthgoogle.New(
					oauthgoogle.WithRedirectURL(os.Getenv("GOOGLE_CALLBACK_URI")),
				)),
				oauth.WithMapProfileToUser(func(info *limen.OAuthAccountProfile) map[string]any {
					return map[string]any{
						"full_name": info.Name,
					}
				}),
			),
		},
	})
}
