package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"qr_auth/pusherutil"
	mp "qr_auth/pusherutil/mocks"
	"qr_auth/redisutil"
	mr "qr_auth/redisutil/mocks"
	"testing"

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

	//assertion
	err := mockAuthHandler.pushQRCodeToChannel(mock.Anything)
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
