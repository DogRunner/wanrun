package google

import (
	// "context"

	// "github.com/labstack/echo/v4"
	_ "github.com/wanrun-develop/wanrun/configs"

	// _wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	_ "github.com/wanrun-develop/wanrun/pkg/log"
	// "golang.org/x/oauth2"

	// _googleOAuth "golang.org/x/oauth2/google"
	// v2 "google.golang.org/api/oauth2/v2"
	_ "google.golang.org/api/option"
)

const (
	GOOGLE_TOKEN_BASE_ENDPOINT = "https://oauth2.googleapis.com/token"
	GOOGLE_REQ_BODY_GRANT_TYPE = "authorization_code"
	GOOGLE_OAUTH_SCOPES_PREFIX = "https://www.googleapis.com/auth/"
)

// type IOAuthGoogle interface {
// 	GetGoogleUserInfo(c echo.Context, token *oauth2.Token, ctx context.Context) (*v2.Tokeninfo, error)
// 	GetAccessToken(c echo.Context, authorizationCode string, ctx context.Context) (*oauth2.Token, error)
// }

// type oauthGoogle struct {
// 	config *oauth2.Config
// }

// func NewOAuthGoogle() IOAuthGoogle {
// 	return &oauthGoogle{
// 		config: &oauth2.Config{
// 			ClientID:     configs.FetchConfigStr("os.secret.key"),
// 			ClientSecret: configs.FetchConfigStr("gcp.client.secret"),
// 			RedirectURL:  configs.FetchConfigStr("gcp.redirect.uri"),
// 			Scopes: []string{
// 				GOOGLE_OAUTH_SCOPES_PREFIX + "userinfo.email",
// 			},
// 			Endpoint: googleOAuth.Endpoint,
// 		},
// 	}
// }

// func (og *oauthGoogle) GetLoginURL(state string) (clientID string) {
// 	return og.config.AuthCodeURL(state)
// }

/*
AccessTokenの取得
*/
// func (og *oauthGoogle) GetAccessToken(c echo.Context, authorizationCode string, ctx context.Context) (*oauth2.Token, error) {
// 	logger := log.GetLogger(c).Sugar()

// 	// アクセストークンなどを取得
// 	token, err := og.config.Exchange(ctx, authorizationCode)
// 	if err != nil {
// 		wrErr := wrErrors.NewWRError(
// 			err,
// 			"failed to exchange authorization code for token",
// 			wrErrors.NewAuthClientErrorEType()) // Question 認証エラーかな？

// 		logger.Errorf("OAuth Exchange error: %v", wrErr)

// 		return nil, wrErr
// 	}

// 	logger.Info("token info: %v", token)

// 	return token, nil
// }

// /*
// Googleのユーザー情報の取得
// */
// func (og *oauthGoogle) GetGoogleUserInfo(c echo.Context, token *oauth2.Token, ctx context.Context) (*v2.Tokeninfo, error) {
// 	logger := log.GetLogger(c).Sugar()

// 	httpClient := og.config.Client(ctx, token)

// 	// Google API サービスの初期化
// 	service, err := v2.NewService(ctx, option.WithHTTPClient(httpClient))
// 	if err != nil {
// 		wrErr := wrErrors.NewWRError(
// 			err,
// 			"failed to initialize Google OAuth2 service",
// 			wrErrors.NewAuthServerErrorEType(),
// 		) // Question

// 		logger.Error("Google OAuth2 service initialization error: %v", wrErr)

// 		return nil, wrErr
// 	}

// 	// Googleユーザーの取得
// 	userInfo, err := service.Tokeninfo().AccessToken(token.AccessToken).Context(ctx).Do()
// 	if err != nil {
// 		wrErr := wrErrors.NewWRError(
// 			err,
// 			"failed to retrieve user info from Google OAuth2 service",
// 			wrErrors.NewAuthServerErrorEType(),
// 		)
// 		logger.Errorf("Google user info retrieval error: %v", wrErr)
// 		return nil, wrErr
// 	}

// 	logger.Infof("google user info: %v", userInfo)
// 	return userInfo, nil
// }
