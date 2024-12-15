package wanruncmd

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wanrun-develop/wanrun/configs"
	"github.com/wanrun-develop/wanrun/internal"
	authRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/repository"
	authScopeRepository "github.com/wanrun-develop/wanrun/internal/auth/adapters/scoperepository"
	authController "github.com/wanrun-develop/wanrun/internal/auth/controller"
	authHandler "github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	authMW "github.com/wanrun-develop/wanrun/internal/auth/middleware"
	cmsAWS "github.com/wanrun-develop/wanrun/internal/cms/adapters/aws"
	cmsRepository "github.com/wanrun-develop/wanrun/internal/cms/adapters/repository"
	cmsController "github.com/wanrun-develop/wanrun/internal/cms/controller"
	cmsHandler "github.com/wanrun-develop/wanrun/internal/cms/core/handler"
	"github.com/wanrun-develop/wanrun/internal/db"
	dogRepository "github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	dogController "github.com/wanrun-develop/wanrun/internal/dog/controller"
	dogHandler "github.com/wanrun-develop/wanrun/internal/dog/core/handler"
	dogOwnerRepository "github.com/wanrun-develop/wanrun/internal/dogowner/adapters/repository"
	dogOwnerScopeRepository "github.com/wanrun-develop/wanrun/internal/dogowner/adapters/scoperepository"
	dogOwnerController "github.com/wanrun-develop/wanrun/internal/dogowner/controller"
	dogOwnerHandler "github.com/wanrun-develop/wanrun/internal/dogowner/core/handler"
	"github.com/wanrun-develop/wanrun/internal/dogrun/adapters/googleplace"
	dogrunR "github.com/wanrun-develop/wanrun/internal/dogrun/adapters/repository"
	dogrunC "github.com/wanrun-develop/wanrun/internal/dogrun/controller"
	dogrunH "github.com/wanrun-develop/wanrun/internal/dogrun/core/handler"
	dogrunFacade "github.com/wanrun-develop/wanrun/internal/dogrun/facade"
	interactionR "github.com/wanrun-develop/wanrun/internal/interaction/adapters/repository"
	interactionC "github.com/wanrun-develop/wanrun/internal/interaction/controller"
	interactionH "github.com/wanrun-develop/wanrun/internal/interaction/core/handler"
	interactionFacade "github.com/wanrun-develop/wanrun/internal/interaction/facade"
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

	// 最大リクエストボディサイズの指定
	e.Use(middleware.BodyLimit("10M")) // 最大10MB

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

	// dogOwner関連
	dogOwnerController := newDogOwner(dbConn)
	dogOwner := e.Group("dogowner")
	dogOwner.POST("/signUp", dogOwnerController.DogOwnerSignUp)

	// auth関連
	authController := newAuth(dbConn)
	auth := e.Group("auth")
	auth.POST("/token", authController.LogIn)
	auth.POST("/revoke", authController.Revoke)
	// auth.GET("/google/oauth", authController.GoogleOAuth)

	//interaction関連
	interactionController := newInteraction(dbConn)
	bookmark := e.Group("bookmark")
	bookmark.POST("/dogrun", interactionController.AddBookmark)
	bookmark.DELETE("/dogrun", interactionController.DeleteBookmarks)

	// cms関連
	cmsController := newCms(dbConn)
	cms := e.Group("cms")
	cms.POST("/upload/file", cmsController.UploadFile)

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
	//facadeの準備
	interactionRepository := interactionR.NewBookmarkRepository(dbConn)
	dogrunFacade := interactionFacade.NewBookmarkFacade(interactionRepository)

	dogrunRest := googleplace.NewRest()
	dogrunRepository := dogrunR.NewDogrunRepository(dbConn)
	dogrunHandler := dogrunH.NewDogrunHandler(dogrunRest, dogrunRepository, dogrunFacade)
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

func newInteraction(dbConn *gorm.DB) interactionC.IBookmarkController {
	//facadeの準備
	dogrunRepository := dogrunR.NewDogrunRepository(dbConn)
	dogrunFacade := dogrunFacade.NewDogrunFacade(dogrunRepository)

	interactionRepository := interactionR.NewBookmarkRepository(dbConn)
	interactionHandler := interactionH.NewBookmarkHandler(interactionRepository, dogrunFacade)

	return interactionC.NewBookmarkController(interactionHandler)
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

func newCms(dbConn *gorm.DB) cmsController.ICmsController {
	cmsRepository := cmsRepository.NewCmsRepository(dbConn)
	// aws設定
	sdkCfg, err := loadAWSConfig()

	if err != nil {
		log.Fatalf("AWSのクレデンシャル取得に失敗: %v", err)
	}
	cmsAWS := cmsAWS.NewS3Provider(sdkCfg)
	cmsHandler := cmsHandler.NewCmsHandler(cmsAWS, cmsRepository)
	cmsController := cmsController.NewCmsController(cmsHandler)
	return cmsController
}

func loadAWSConfig() (aws.Config, error) {
	// local
	if configs.FetchConfigStr("ENV") == "local" {
		return config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cmsAWS.DEFAULT_REGION),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					configs.FetchConfigStr("aws.access.key"),
					configs.FetchConfigStr("aws.secret.access.key"),
					"",
				),
			),
		)
	}

	// クラウド
	return config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cmsAWS.DEFAULT_REGION),
	)
}
