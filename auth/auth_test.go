package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"qr_auth/pusherutil"
	mp "qr_auth/pusherutil/mocks"
	"qr_auth/redisutil"
	mr "qr_auth/redisutil/mocks"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPushQRCodeToChannel_Success(t *testing.T) {
	//setup
	mockRedis := mr.NewMockCache()
	mockPush := mp.NewMockPush()
	mockAuthHandler := &AuthHandler{
		cache:  mockRedis,
		pusher: mockPush,
	}
	mockPush.Wg = &sync.WaitGroup{}

	//assertion
	mockPush.Wg.Add(1)
	err := mockAuthHandler.pushQRCodeToChannel(mock.Anything)
	mockPush.Wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, len(mockPush.TriggerCalls), 1)
	assert.Equal(t, len(mockRedis.SetCalls), 1)
}

func TestVerifyQRLogin_Success(t *testing.T) {
	// setup
	e := echo.New()
	token := "some-token"
	clientToken := "clientToken"
	q := make(url.Values)
	q.Set("token", token)
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	mockRedis := mr.NewMockCache()
	mockPush := mp.NewMockPush()
	mockRedis.Set(redisutil.REDIS_QR_LOGIN_PREFIX+token, clientToken, 0)
	mockAuthHandler := &AuthHandler{
		cache:  mockRedis,
		pusher: mockPush,
	}

	//assert
	if assert.NoError(t, mockAuthHandler.VerifyQRLogin(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "login success")
		mp := mockPush.TriggerCalls[0]
		assert.Equal(t, len(mockPush.TriggerCalls), 1)
		assert.Equal(t, mp.Channel, pusherutil.QR_CHANNEL)
		assert.Equal(t, mp.EventName, EVENT_LOGIN_SUCCESS_PREFIX+clientToken)
		data := mp.Data.(map[string]string)
		assert.Equal(t, data["token"], "newToken")
		assert.Equal(t, len(mockRedis.GetCalls), 1)
		assert.Equal(t, len(mockRedis.DeleteCalls), 1)
		assert.Equal(t, mockRedis.GetCalls[0], redisutil.REDIS_QR_LOGIN_PREFIX+token)
		assert.Equal(t, mockRedis.DeleteCalls[0], redisutil.REDIS_QR_LOGIN_PREFIX+token)
	}
}

func TestSendQRLogin(t *testing.T) {
	// setup
	e := echo.New()
	q := make(url.Values)
	token := "some-token"
	q.Set("uuid", token)
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	mockRedis := mr.NewMockCache()
	mockPush := mp.NewMockPush()
	mockPush.Wg = &sync.WaitGroup{}
	qrInterval := 1
	qrLifetime := 4
	mockAuthHandler := &AuthHandler{
		cache:      mockRedis,
		pusher:     mockPush,
		qrInterval: time.Duration(qrInterval) * time.Second,
		qrLifetime: time.Duration(qrLifetime) * time.Second,
	}

	// assert
	mockPush.Wg.Add(qrLifetime/qrInterval + 1)
	err := mockAuthHandler.SendQRLogin(c)
	mockPush.Wg.Wait()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "qr sent")
	assert.Equal(t, len(mockPush.TriggerCalls), qrLifetime/qrInterval+1)
	assert.Equal(t, len(mockRedis.SetCalls), qrLifetime/qrInterval+1)
}
