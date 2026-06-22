package auth

import (
	"context"
	"database/sql"
	"fmt"
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
	jasperLabsProvider := oauthgeneric.New(
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
			emailVerified, _ := raw["email_verified"].(bool)
			avatarURL, _ := raw["picture"].(string)
			if avatarURL == "" {
				avatarURL, _ = raw["avatar_url"].(string)
			}

			return &oauth.ProviderUserInfo{
				ID:            id,
				Name:          name,
				Email:         email,
				EmailVerified: emailVerified,
				AvatarURL:     avatarURL,
			}, nil
		}),
	)

	return limen.New(&limen.Config{
		Database: sqladapter.NewPostgreSQL(db),
		HTTP: limen.NewDefaultHTTPConfig(
			limen.WithHTTPTrustedOrigins(TrustedFrontendOrigins()),
		),
		Schema: schema,
		Plugins: []limen.Plugin{
			oauth.New(
				oauth.WithProviders(jasperLabsProvider),
				oauth.WithGetUserInfo(func(ctx context.Context, provider string, token *oauth.TokenResponse) (*oauth.ProviderUserInfo, error) {
					info, err := jasperLabsProvider.GetUserInfo(ctx, token)
					if err != nil {
						return nil, err
					}
					if provider == "jasperlabs" && info != nil {
						if err := syncOAuthUserProfile(ctx, db, provider, info); err != nil {
							return nil, err
						}
					}
					return info, nil
				}),
				oauth.WithMapProfileToUser(func(info *limen.OAuthAccountProfile) map[string]any {
					return map[string]any{
						"full_name":  info.Name,
						"avatar_url": info.AvatarURL,
					}
				}),
			),
		},
	})
}

func syncOAuthUserProfile(ctx context.Context, db *sql.DB, provider string, info *oauth.ProviderUserInfo) error {
	if info.ID == "" || info.Email == "" {
		return nil
	}

	const query = `
WITH target_user AS (
    SELECT a.user_id AS id
    FROM accounts a
    WHERE a.provider = $1
      AND a.provider_account_id = $2
    UNION
    SELECT u.id
    FROM users u
    WHERE u.email = $3
      AND NOT EXISTS (
          SELECT 1
          FROM accounts a
          WHERE a.provider = $1
            AND a.provider_account_id = $2
      )
    LIMIT 1
)
UPDATE users
SET
    email = $3,
    full_name = COALESCE(NULLIF($4, ''), full_name),
    avatar_url = NULLIF($5, ''),
    email_verified_at = CASE
        WHEN $6 THEN
            CASE
                WHEN email IS DISTINCT FROM $3 OR email_verified_at IS NULL THEN NOW()
                ELSE email_verified_at
            END
        ELSE NULL
    END,
    updated_at = NOW()
WHERE id = (SELECT id FROM target_user)`

	if _, err := db.ExecContext(ctx, query, provider, info.ID, info.Email, info.Name, info.AvatarURL, info.EmailVerified); err != nil {
		return fmt.Errorf("sync oauth user profile: %w", err)
	}

	return nil
}
