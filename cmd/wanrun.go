package wanruncmd

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wanrun-develop/wanrun/internal"
	authRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	authController "github.com/wanrun-develop/wanrun/internal/auth/controller"
	authHandler "github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	"github.com/wanrun-develop/wanrun/internal/db"
	dogRepository "github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	dogController "github.com/wanrun-develop/wanrun/internal/dog/controller"
	dogHandler "github.com/wanrun-develop/wanrun/internal/dog/core/handler"
	dogOwnerRepositoy "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/dogrun/adapters/googleplace"
	dogrunR "github.com/wanrun-develop/wanrun/internal/dogrun/adapters/repository"
	dogrunC "github.com/wanrun-develop/wanrun/internal/dogrun/controller"
	dogrunH "github.com/wanrun-develop/wanrun/internal/dogrun/core/handler"

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

	// ミドルウェアを登録
	e.Use(middleware.RequestID())
	e.HTTPErrorHandler = errors.HttpErrorHandler
	e.Use(logger.RequestLoggerMiddleware(zap))

	// CORSの設定
	e.Use(middleware.CORS())

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
	dog.GET("/ownered/:dogOwnerId", dogController.GetDogByDogOwnerID)
	dog.GET("/mst/dogType", dogController.GetDogTypeMst)
	dog.POST("", dogController.CreateDog)
	dog.PUT("", dogController.UpdateDog)
	dog.DELETE("", dogController.DeleteDog)
	// dog.PUT("/:dogID", dogController.UpdateDog)

	// dogrun関連
	dogrunConrtoller := newDogrun(dbConn)
	dogrun := e.Group("dogrun")
	dogrun.GET("/detail/:placeId", dogrunConrtoller.GetDogrunDetail)
	dogrun.GET("/:id", dogrunConrtoller.GetDogrun)
	dogrun.POST("/search", dogrunConrtoller.SearchAroundDogruns)
	dogrun.GET("/photo/src", dogrunConrtoller.GetDogrunPhoto)

	// auth関連
	authController := newAuth(dbConn)
	auth := e.Group("auth")
	// auth.GET("/google/oauth", authController.GoogleOAuth)
	auth.POST("/signUp", authController.SignUp)
	// auth.POST("/login", authController.LogIn)
}

// dogの初期化
func newDog(dbConn *gorm.DB) dogController.IDogController {
	dogRepository := dogRepository.NewDogRepository(dbConn)
	dogOwnerRepository := dogOwnerRepositoy.NewDogRepository(dbConn)
	dogHandler := dogHandler.NewDogHandler(dogRepository, dogOwnerRepository)
	dogController := dogController.NewDogController(dogHandler)
	return dogController
}

func newDogrun(dbConn *gorm.DB) dogrunC.IDogrunController {
	dogrunRest := googleplace.NewRest()
	dogrunRepository := dogrunR.NewDogrunRepository(dbConn)
	dogrunHanlder := dogrunH.NewDogrunHandler(dogrunRest, dogrunRepository)
	return dogrunC.NewDogrunController(dogrunHanlder)
}

func newAuth(dbConn *gorm.DB) authController.IAuthController {
	authRepository := authRepository.NewAuthRepository(dbConn)
	// googleOAuth := google.NewOAuthGoogle()
	// authHandler := authHandler.NewAuthHandler(authRepository, googleOAuth)
	authHandler := authHandler.NewAuthHandler(authRepository)
	authController := authController.NewAuthController(authHandler)
	return authController
}
