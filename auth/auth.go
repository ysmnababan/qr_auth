package auth

import (
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

type AuthHandler struct {
	cache      redisutil.ICache
	pusher     pusherutil.IPusher
	qrInterval time.Duration
	qrLifetime time.Duration
}

func NewAuthHandler(c *config.Config, interval, lifetime time.Duration) *AuthHandler {
	return &AuthHandler{
		cache:      c.Redis,
		pusher:     c.PusherClient,
		qrInterval: interval,
		qrLifetime: lifetime,
	}
}

func generateQRCodeBase64(data string) (string, error) {
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64Img := base64.StdEncoding.EncodeToString(png)
	return "data:image/png;base64," + base64Img, nil
}

func (h *AuthHandler) pushQRCodeToChannel(clientToken string) error {
	token := uuid.New().String()
	redisKey := redisutil.REDIS_QR_LOGIN_PREFIX + token

	err := h.cache.Set(redisKey, clientToken, REDIS_LOGIN_LIFETIME)
	if err != nil {
		log.Println(err)
		return err
	}
	content := baseURL + "/auth/verify?token=" + token

	qrBase64, err := generateQRCodeBase64(content)
	log.Println(content)
	if err != nil {
		log.Println(err)
		return err
	}

	err = h.pusher.Trigger(
		pusherutil.QR_CHANNEL,
		EVENT_LOGIN_QRCODE_PREFIX+clientToken,
		map[string]string{
			"qr_code": qrBase64,
		},
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (h *AuthHandler) SendQRLogin(c echo.Context) (err error) {
	clientToken := c.QueryParam("uuid")
	if len(clientToken) == 0 {
		return c.String(http.StatusBadRequest, "client token can't be empty")
	}

	go func() {
		start := time.Now()
		for {
			err := h.pushQRCodeToChannel(clientToken)
			if err != nil {
				log.Println(err)
				return
			}
			if time.Since(start) > h.qrLifetime {
				return
			}
			time.Sleep(h.qrInterval)
		}
	}()

	log.Println("qr sent")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "qr sent",
	})
}

func (h *AuthHandler) VerifyQRLogin(c echo.Context) (err error) {
	token := c.QueryParam("token")
	redisKey := redisutil.REDIS_QR_LOGIN_PREFIX + token

	// fetch the client token
	clientToken, err := h.cache.Get(redisKey)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("yey berhasil", clientToken)
	newToken := "newToken"
	err = h.pusher.Trigger(
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

	deleted, err := h.cache.Delete(redisKey)
	if err != nil {
		log.Println(err)
	}

	log.Println("deleted keys", deleted)
	// disconnect client
	return c.JSON(http.StatusOK, map[string]string{
		"message": "login success",
	})
}
