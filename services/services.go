package services

import (
	"net/http"
	"scheduleme/config"
	"scheduleme/middlewares"
	"scheduleme/models"
	"scheduleme/oauthhandlers"
	"scheduleme/reqhandlers"
	"scheduleme/resources"
	"scheduleme/secure_cookie"
	"scheduleme/sqlite"

	"golang.org/x/oauth2"
)

type Services struct {
	SecureCookie  secure_cookie.SecureCookieInterface
	Middlewares   *middlewares.Middlewares
	Handlers      *reqhandlers.Handlers //TODO rename more specific, since it only uses secure cookie
	Resources     *resources.Resources
	OAuth         *oauthhandlers.OAuth
	Repo          *models.Repo
	ServiceConfig *config.ConfigStruct
}

func BuildServices(
	sc secure_cookie.SecureCookieInterface,
	repo *models.Repo, //Should this be more explicit passing of required services?
	//TODO:  combine googleAuth and authConfig; which is more generic and works with diff provider
	googleAuth *oauth2.Config,
	cfg *config.ConfigStruct,
) *Services {
	midd := middlewares.NewMiddleware(sc)
	hand := reqhandlers.NewHandlers(sc, googleAuth, repo)
	resrc := resources.NewResources(repo)
	auf := &models.AuthUserFlow{AuthService: repo.AuthService, UserService: repo.UserService}
	oa := oauthhandlers.NewOAuth(sc, googleAuth, auf, &http.Client{})
	return &Services{sc, midd, hand, resrc, oa, repo, cfg}
}

// func TopService(cfg *config.ConfigStruct) *Services {
// 	db, err := sqlite.NewOpenDB(cfg.Dsn)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return serviceFactory(cfg, db)
// }

// func TopServiceForTesting(cfg *config.ConfigStruct, db *sqlite.Db) *Services {
// 	return serviceFactory(cfg, db)
// }

func TopServices(cfg *config.ConfigStruct, db *sqlite.Db) *Services {
	googleAuth := oauthhandlers.NewGoogleOAuth(
		cfg.GoogleRedirectURL,
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
	)
	repo := models.NewRepo(db, googleAuth)
	seccook := secure_cookie.NewSecureCookie(cfg.Secret, "session")

	return BuildServices(seccook, repo, googleAuth, cfg)
}
