package handler

import (
	"github.com/labstack/echo/v4"
	authDTO "github.com/wanrun-develop/wanrun/internal/auth/core/dto"
	authFacade "github.com/wanrun-develop/wanrun/internal/auth/core/facade"
	authHandler "github.com/wanrun-develop/wanrun/internal/auth/core/handler"
	dogrunmgFacade "github.com/wanrun-develop/wanrun/internal/dogrunmg/core/facade"
	model "github.com/wanrun-develop/wanrun/internal/models"
	orgRepository "github.com/wanrun-develop/wanrun/internal/org/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/org/core/dto"
	"github.com/wanrun-develop/wanrun/internal/transaction"
	wrErrors "github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	wrUtil "github.com/wanrun-develop/wanrun/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IOrgHandler interface {
	OrgSignUp(c echo.Context, orgReq dto.OrgReq) (string, error)
}

type orgHandler struct {
	osr orgRepository.IOrgScopeRepository
	tm  transaction.ITransactionManager
	dmf dogrunmgFacade.IDogrunmgFacade
	af  authFacade.IAuthFacade
}

func NewOrgHandler(
	osr orgRepository.IOrgScopeRepository,
	tm transaction.ITransactionManager,
	dmf dogrunmgFacade.IDogrunmgFacade,
	af authFacade.IAuthFacade,
) IOrgHandler {
	return &orgHandler{
		osr: osr,
		tm:  tm,
		dmf: dmf,
		af:  af,
	}
}

func (oh *orgHandler) OrgSignUp(
	c echo.Context,
	orgReq dto.OrgReq,
) (string, error) {
	logger := log.GetLogger(c).Sugar()

	// パスワードのハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(orgReq.Password), bcrypt.DefaultCost) // 一旦costをデフォルト値

	if err != nil {
		wrErr := wrErrors.NewWRError(
			err,
			"パスワードに不正な文字列が入っています。",
			wrErrors.NewOrgClientErrorEType(),
		)
		logger.Error(wrErr)
		return "", wrErr
	}

	// JWT IDの生成
	jwtID, wrErr := authHandler.GenerateJwtID(c)

	if wrErr != nil {
		return "", wrErr
	}

	// requestからorgの構造体に詰め替え
	orgInfo := model.DogrunmgCredential{
		Email:    wrUtil.NewSqlNullString(orgReq.ContactEmail),
		Password: wrUtil.NewSqlNullString(string(hash)),
		IsAdmin:  wrUtil.NewSqlNullBool(true), // 初期adminユーザーのため
		AuthDogrunmg: model.AuthDogrunmg{
			JwtID: wrUtil.NewSqlNullString(jwtID),
			Dogrunmg: model.Dogrunmg{
				Name: wrUtil.NewSqlNullString("admin"),
				Organization: model.Organization{
					Name:         wrUtil.NewSqlNullString(orgReq.OrgName),
					ContactEmail: wrUtil.NewSqlNullString(orgReq.ContactEmail),
					PhoneNumber:  wrUtil.NewSqlNullString(orgReq.PhoneNumber),
					Address:      wrUtil.NewSqlNullString(orgReq.Address),
					Description:  wrUtil.NewSqlNullString(orgReq.Description),
				},
			},
		},
	}

	logger.Debugf("OrgCredential %v, Type: %T", orgInfo, orgInfo)

	ctx := c.Request().Context()

	// organizationの作成トランザクション
	if err := oh.tm.DoInTransaction(c, ctx, func(tx *gorm.DB) error {

		// organizationの作成
		if wrErr := oh.osr.CreateOrg(tx, c, &orgInfo.AuthDogrunmg.Dogrunmg); wrErr != nil {
			return wrErr
		}

		// dogrunmgの作成
		if wrErr := oh.dmf.CreateOrg(tx, c, &orgInfo.AuthDogrunmg); wrErr != nil {
			return wrErr
		}

		// dogrunmgのauthやcredentialsの作成
		if wrErr := oh.af.CreateOrg(tx, c, &orgInfo); wrErr != nil {
			return wrErr
		}

		// 正常に完了
		return nil

	}); err != nil {
		logger.Error("Transaction failed:", err)
		return "", err
	}

	// 作成したdogrunmgの情報をdto詰め替え
	dogrunmgrDetail := authDTO.UserAuthInfoDTO{
		UserID:   orgInfo.AuthDogrunmg.DogrunmgID.Int64,
		JwtID:    orgInfo.AuthDogrunmg.JwtID.String,
		RoleName: authHandler.DOGRUNMG_ROLE_NAME,
	}

	logger.Infof("dogrunmgDetail: %v", dogrunmgrDetail)

	// 署名済みのjwt token取得
	token, wrErr := authHandler.GetSignedJwt(c, dogrunmgrDetail)

	if wrErr != nil {
		return "", wrErr
	}

	return token, nil
}
