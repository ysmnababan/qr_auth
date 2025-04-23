package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"qr_auth/config"
	"qr_auth/pusherutil"
	"qr_auth/redisutil"
	"time"

	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/skip2/go-qrcode"
)

var baseURL = "https://663a-203-128-66-62.ngrok-free.app"

func generateQRCodeBase64(data string) (string, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	base64Img := base64.StdEncoding.EncodeToString(png)
	return "data:image/png;base64," + base64Img, nil
}

func SendQRLogin(c echo.Context) (err error) {
	token := uuid.New().String()
	redisKey := redisutil.REDIS_QR_LOGIN_PREFIX + token

	err = config.Cfg.Redis.Set(context.Background(), redisKey, "true", time.Minute*5).Err()
	if err != nil {
		panic(err)
	}
	content := baseURL + "/auth/verify?token=" + token

	qrBase64, err := generateQRCodeBase64(content)
	log.Println(content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not generate QR code")
	}
	log.Println("qr sent")
	return c.JSON(http.StatusOK, map[string]string{
		"qr_code": qrBase64,
	})
}

func VerifyQRLogin(c echo.Context) (err error) {
	token := c.QueryParam("token")
	redisKey := redisutil.REDIS_QR_LOGIN_PREFIX + token

	// Check if key exists
	exists, err := config.Cfg.Redis.Exists(context.Background(), redisKey).Result()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("Key exists?", exists > 0)
	if exists == 0 {
		return c.String(http.StatusNotFound, "qr token is not found")
	}
	fmt.Println("yey berhasil")
	err = config.Cfg.PusherClient.Trigger(
		pusherutil.QR_CHANNEL,
		"login_success",
		map[string]string{
			"status": "true",
		},
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// delete the item
	// create token
	// send redirect link to channel
	// disconnect client
	return
}
