package wanruncmd

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wanrun-develop/wanrun/internal"
	authRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	authScopeRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/scopeRepository"
	authController "github.com/wanrun-develop/wanrun/internal/auth/controller"
	authHandler "github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	authMW "github.com/wanrun-develop/wanrun/internal/auth/middleware"
	"github.com/wanrun-develop/wanrun/internal/db"
	dogRepository "github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	dogController "github.com/wanrun-develop/wanrun/internal/dog/controller"
	dogHandler "github.com/wanrun-develop/wanrun/internal/dog/core/handler"
	dogOwnerRepository "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/repository"
	dogOwnerScopeRepository "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/scopeRepository"
	dogOwnerController "github.com/wanrun-develop/wanrun/internal/dogOwner/controller"
	dogOwnerHandler "github.com/wanrun-develop/wanrun/internal/dogOwner/core/handler"
	"github.com/wanrun-develop/wanrun/internal/dogrun/adapters/googleplace"
	dogrunR "github.com/wanrun-develop/wanrun/internal/dogrun/adapters/repository"
	dogrunC "github.com/wanrun-develop/wanrun/internal/dogrun/controller"
	dogrunH "github.com/wanrun-develop/wanrun/internal/dogrun/core/handler"
	"github.com/wanrun-develop/wanrun/internal/transaction"

	"github.com/wanrun-develop/wanrun/pkg/errors"
	logger "github.com/wanrun-develop/wanrun/pkg/log"
	"gorm.io/gorm"
)

func init() {

}

func Main() {
	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatalln(err)
	}

	defer db.CloseDB(dbConn)

	e := echo.New()

	// グローバルロガーの初期化
	zap := logger.NewWanRunLogger()
	logger.SetLogger(zap) // グローバルロガーを設定
	// アプリケーション終了時にロガーを同期
	defer zap.Sync()

	// CORSの設定
	e.Use(middleware.CORS())

	// ミドルウェアを登録
	e.Use(middleware.RequestID())
	e.HTTPErrorHandler = errors.HttpErrorHandler
	e.Use(logger.RequestLoggerMiddleware(zap))

	// JWTミドルウェアの設定
	authMiddleware := newAuthMiddleware(dbConn)
	e.Use(authMiddleware.NewJwtValidationMiddleware())

	// Router設定
	newRouter(e, dbConn)
	e.GET("/test", internal.Test)

	e.Logger.Fatal(e.Start(":8080"))
}

func newRouter(e *echo.Echo, dbConn *gorm.DB) {
	// dog関連
	dogController := newDog(dbConn)
	dog := e.Group("dog")
	dog.GET("/all", dogController.GetAllDogs)
	dog.GET("/detail/:dogID", dogController.GetDogByID)
	dog.GET("/owned/:dogOwnerId", dogController.GetDogByDogOwnerID)
	dog.GET("/mst/dogType", dogController.GetDogTypeMst)
	dog.POST("", dogController.CreateDog)
	dog.PUT("", dogController.UpdateDog)
	dog.DELETE("", dogController.DeleteDog)
	// dog.PUT("/:dogID", dogController.UpdateDog)

	// dogrun関連
	dogrunController := newDogrun(dbConn)
	dogrun := e.Group("dogrun")
	dogrun.GET("/detail/:placeId", dogrunController.GetDogrunDetail)
	dogrun.GET("/:id", dogrunController.GetDogrun)
	dogrun.GET("/photo/src", dogrunController.GetDogrunPhoto)
	dogrun.GET("/mst/tag", dogrunController.GetDogrunTagMst)
	dogrun.POST("/search", dogrunController.SearchAroundDogruns)

	// auth関連
	authController := newAuth(dbConn)
	dogOwnerController := newDogOwner(dbConn)
	auth := e.Group("auth")
	// auth.GET("/google/oauth", authController.GoogleOAuth)
	// auth.POST("/signUp", authController.SignUp)
	auth.POST("/dogOwner/signUp", dogOwnerController.DogOwnerSignUp)
	auth.POST("/token", authController.LogIn)
	auth.POST("/revoke", authController.Revoke)

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}

// dogの初期化
func newDog(dbConn *gorm.DB) dogController.IDogController {
	dogRepository := dogRepository.NewDogRepository(dbConn)
	dogOwnerRepository := dogOwnerRepository.NewDogRepository(dbConn)
	dogHandler := dogHandler.NewDogHandler(dogRepository, dogOwnerRepository)
	dogController := dogController.NewDogController(dogHandler)
	return dogController
}

func newDogrun(dbConn *gorm.DB) dogrunC.IDogrunController {
	dogrunRest := googleplace.NewRest()
	dogrunRepository := dogrunR.NewDogrunRepository(dbConn)
	dogrunHandler := dogrunH.NewDogrunHandler(dogrunRest, dogrunRepository)
	return dogrunC.NewDogrunController(dogrunHandler)
}

func newAuth(dbConn *gorm.DB) authController.IAuthController {
	authRepository := authRepository.NewAuthRepository(dbConn)
	// googleOAuth := google.NewOAuthGoogle()
	// authHandler := authHandler.NewAuthHandler(authRepository, googleOAuth)
	authHandler := authHandler.NewAuthHandler(authRepository)
	authController := authController.NewAuthController(authHandler)
	return authController
}

func newAuthMiddleware(dbConn *gorm.DB) authMW.IAuthJwt {
	authRepository := authRepository.NewAuthRepository(dbConn)
	return authMW.NewAuthJwt(authRepository)
}

// dogOwnerの初期化
func newDogOwner(dbConn *gorm.DB) dogOwnerController.IDogOwnerController {
	// repository層
	dogOwnerRepository := dogOwnerRepository.NewDogRepository(dbConn)
	authRepository := authRepository.NewAuthRepository(dbConn)

	// transaction層
	transactionManager := transaction.NewTransactionManager(dbConn)

	// scopeRepository層
	dogOwnerScopeRepository := dogOwnerScopeRepository.NewDogOwnerScopeRepository()
	authScopeRepository := authScopeRepository.NewAuthScopeRepository()

	// handler層
	authHandler := authHandler.NewAuthHandler(authRepository)
	dogOwnerHandler := dogOwnerHandler.NewDogOwnerHandler(
		dogOwnerScopeRepository,
		transactionManager,
		authScopeRepository,
		dogOwnerRepository,
		authRepository,
	)

	// controller層
	return dogOwnerController.NewDogOwnerController(dogOwnerHandler, authHandler)
}
