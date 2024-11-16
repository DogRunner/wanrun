package googleplace

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/configs"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"go.uber.org/zap"
)

var apiKey string

const (
	H_G_FIELD_MASK string = "X-Goog-FieldMask"
)

func init() {
	apiKey = configs.FetchConfigStr("google.place.api.key")
}

type IRest interface {
	GETPlaceInfo(echo.Context, string, IFieldMask) ([]byte, error)
	POSTSearchNearby(echo.Context, SearchNearbyPayLoad, IFieldMask) ([]byte, error)
	POSTSearchText(echo.Context, SearchTextPayLoad, IFieldMask) ([]byte, error)
	GETPhotoByName(echo.Context, string, string, string) ([]byte, error)
}
type rest struct{}

func NewRest() IRest {
	return &rest{}
}

/*
GET
google place details apiの実行
*/
func (r *rest) GETPlaceInfo(c echo.Context, placeId string, field IFieldMask) ([]byte, error) {
	logger := log.GetLogger(c).Sugar()

	url := urlPlacesWPlaceId(placeId)

	req, err := fetchGETReq(c, url)
	if err != nil {
		return nil, err
	}
	req.Header.Set(H_G_FIELD_MASK, field.getValue())
	logger.Info("field mask:", field.getValue())

	resp, err := exec(c, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// レスポンスの処理
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
GET
google place photo media apiの実行
*/
func (r *rest) GETPhotoByName(c echo.Context, name, widthPx, heightPx string) ([]byte, error) {
	logger := log.GetLogger(c).Sugar()

	url := urlPlacesPhotoWName(name)

	req, err := fetchGETReq(c, url)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("maxWidthPx", widthPx)
	q.Add("maxHeightPx", heightPx)
	q.Add("skipHttpRedirect", "true")
	req.URL.RawQuery = q.Encode()

	logger.Info("request url : ", req.URL.String())

	resp, err := exec(c, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// レスポンスの処理
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
POST
google place search nearby apiの実行
*/
func (r *rest) POSTSearchNearby(c echo.Context, payload SearchNearbyPayLoad, field IFieldMask) ([]byte, error) {
	logger := log.GetLogger(c).Sugar()

	url := urlPlacesWSearchNearBy()

	req, err := fetchPOSTReq(c, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set(H_G_FIELD_MASK, field.getValueWPlaces())
	logger.Info("field mask:", field.getValueWPlaces())

	resp, err := exec(c, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// レスポンスの処理
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
POST
google place search text apiの実行
*/
func (r *rest) POSTSearchText(c echo.Context, payload SearchTextPayLoad, field IFieldMask) ([]byte, error) {
	logger := log.GetLogger(c).Sugar()

	url := urlPlacesWSearchText()

	req, err := fetchPOSTReq(c, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set(H_G_FIELD_MASK, field.getValueWPlacesAndNextPageToken())
	logger.Info("field mask:", field.getValueWPlacesAndNextPageToken())

	resp, err := exec(c, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// レスポンスの処理
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
標準のGET用http.Requestの生成
*/
func fetchGETReq(c echo.Context, url string) (*http.Request, error) {
	logger := log.GetLogger(c)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.NewWRError(err, "リクエストの生成に失敗しました", errors.NewDogrunServerErrorEType())
		logger.Sugar().Error(err)
		return nil, err
	}
	attachBaseHeader(req)
	logger.Debug("http GET requestの生成:",
		zap.String("method", req.Method), zap.String("url", req.URL.String()), zap.Any("header", req.Header))
	return req, nil
}

/*
標準のPOST用http.Requestの生成
*/
func fetchPOSTReq(c echo.Context, url string, payload any) (*http.Request, error) {
	logger := log.GetLogger(c)

	// payloadをJSONに変換
	jsonData, err := json.Marshal(payload)
	if err != nil {
		err = errors.NewWRError(err, "payloadのJSON変換に失敗しました", errors.NewDogrunServerErrorEType())
		logger.Sugar().Error(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		err = errors.NewWRError(err, "リクエストの生成に失敗しました", errors.NewDogrunServerErrorEType())
		logger.Sugar().Error(err)
		return nil, err
	}
	attachBaseHeader(req)
	logger.Debug("http POST requestの生成:",
		zap.String("method", req.Method), zap.String("url", req.URL.String()), zap.Any("header", req.Header), zap.Any("payload", req.Body))
	return req, nil
}

/*
共通のHeaderをhttp.Requestに付与
*/
func attachBaseHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", apiKey)
	req.Header.Set("Accept-Language", "ja")
}

/*
HTTPリクエストの実行
*/
func exec(c echo.Context, req *http.Request) (*http.Response, error) {
	logger := log.GetLogger(c)

	client := &http.Client{
		Timeout: 60 * time.Second, // タイムアウトの設定
	}

	resp, err := client.Do(req)

	if err != nil {
		logger.Sugar().Error(err)
		return nil, errors.NewWRError(err, "リクエスト失敗", errors.NewDogrunServerErrorEType())
	}

	// レスポンスのステータスコードに応じてエラーハンドリング
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // レスポンスボディを読み取る
		bodyString := string(bodyBytes)

		//ハンドリング
		switch resp.StatusCode {
		case http.StatusBadRequest:
			logger.Error("google place api RESPONSE is Bad Request", zap.Int("api_response_status_code", resp.StatusCode), zap.String("api_response_body", bodyString))
		case http.StatusForbidden:
			logger.Error("google place api RESPONSE is Request forbidden", zap.Int("api_response_status_code", resp.StatusCode), zap.String("api_response_body", bodyString))
		case http.StatusNotFound:
			logger.Error("google place api RESPONSE is Not Found", zap.Int("api_response_status_code", resp.StatusCode), zap.String("api_response_body", bodyString))
		case http.StatusInternalServerError:
			logger.Error("google place api RESPONSE is Server error", zap.Int("api_response_status_code", resp.StatusCode), zap.String("rapi_esponse_body", bodyString))
		default:
			logger.Error("google place api RESPONSE is Unexpected error", zap.Int("api_response_status_code", resp.StatusCode), zap.String("api_response_body", bodyString))
		}
		return nil, errors.NewWRError(nil, "Google API リクエスト失敗", errors.NewDogrunServerErrorEType())
	} else {
		// ステータスコードが200 OKの場合
		logger.Info("google place api Request is successful")
		return resp, nil
	}
}
