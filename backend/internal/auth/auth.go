package auth

import (
	"context"
	"database/sql"
	"os"

	"github.com/google/uuid"
	"github.com/thecodearcher/limen"
	sqladapter "github.com/thecodearcher/limen/adapters/sql"
	"github.com/thecodearcher/limen/plugins/oauth"

	oauthgeneric "github.com/thecodearcher/limen/plugins/oauth-generic"
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
				oauth.WithProviders(
					oauthgeneric.New(
						oauthgeneric.WithName("jasperlabs"),
						oauthgeneric.WithClientID(os.Getenv("OAUTH_CLIENT_ID")),
						oauthgeneric.WithClientSecret(os.Getenv("OAUTH_CLIENT_SECRET")),
						oauthgeneric.WithRedirectURL(os.Getenv("OAUTH_CALLBACK_URL")),
						oauthgeneric.WithScopes("openid", "profile", "email"),

						oauthgeneric.WithDiscoveryURL("https://auth.jasperlabs.net/.well-known/openid-configuration"),

						oauthgeneric.WithMapUserInfo(func(raw map[string]any) (*oauth.ProviderUserInfo, error) {
							id, _ := raw["sub"].(string)
							name, _ := raw["name"].(string)
							email, _ := raw["email"].(string)

							return &oauth.ProviderUserInfo{
								ID:    id,
								Name:  name,
								Email: email,
							}, nil
						}),
					),
				),
				oauth.WithMapProfileToUser(func(info *limen.OAuthAccountProfile) map[string]any {
					return map[string]any{
						"full_name": info.Name,
					}
				}),
			),
		},
	})
}
