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

var baseURL = "https://d13e-125-161-205-3.ngrok-free.app"

const (
	QR_INTERVAL                = 30 * time.Second
	QR_LIFETIME                = 2 * time.Minute
	REDIS_LOGIN_LIFETIME       = 1 * time.Minute
	EVENT_LOGIN_SUCCESS_PREFIX = "login_success:"
	EVENT_LOGIN_QRCODE_PREFIX  = "qr_event:"
)

func generateQRCodeBase64(data string) (string, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64Img := base64.StdEncoding.EncodeToString(png)
	return "data:image/png;base64," + base64Img, nil
}

func PushQRCodeToChannel(clientToken string) {
	token := uuid.New().String()
	redisKey := redisutil.REDIS_QR_LOGIN_PREFIX + token

	err := config.Cfg.Redis.Set(context.Background(), redisKey, clientToken, REDIS_LOGIN_LIFETIME).Err()
	if err != nil {
		panic(err)
	}
	content := baseURL + "/auth/verify?token=" + token

	qrBase64, err := generateQRCodeBase64(content)
	log.Println(content)
	if err != nil {
		log.Println(err)
		return
	}

	err = config.Cfg.PusherClient.Trigger(
		pusherutil.QR_CHANNEL,
		EVENT_LOGIN_QRCODE_PREFIX+clientToken,
		map[string]string{
			"qr_code": qrBase64,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func SendQRLogin(c echo.Context) (err error) {
	clientToken := c.QueryParam("uuid")
	if len(clientToken) == 0 {
		return c.String(http.StatusBadRequest, "client token can't be empty")
	}

	go func() {
		start := time.Now()
		for {
			PushQRCodeToChannel(clientToken)
			if time.Since(start) > QR_LIFETIME {
				return
			}
			time.Sleep(QR_INTERVAL)
		}
	}()

	log.Println("qr sent")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "qr sent",
	})
}

func VerifyQRLogin(c echo.Context) (err error) {
	token := c.QueryParam("token")
	redisKey := redisutil.REDIS_QR_LOGIN_PREFIX + token

	// fetch the client token
	clientToken, err := config.Cfg.Redis.Get(context.Background(), redisKey).Result()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	
	}
	fmt.Println("yey berhasil", clientToken)
	newToken := "newToken"
	err = config.Cfg.PusherClient.Trigger(
		pusherutil.QR_CHANNEL,
		EVENT_LOGIN_SUCCESS_PREFIX+clientToken,
		map[string]string{
			"status": "true",
			"token":  newToken,
		},
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	deleted, err := config.Cfg.Redis.Del(context.Background(), redisKey).Result()
	if err != nil {
		log.Println(err)
	}

	log.Println("deleted keys", deleted)
	// disconnect client
	return
}
