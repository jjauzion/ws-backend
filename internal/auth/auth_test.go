package auth

import (
	"github.com/golang/mock/gomock"
	conf2 "github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/db/mock_db"
	"github.com/jjauzion/ws-backend/internal/logger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	signKey = "abcdef0123456789"
	userID  = "JBDmH5vuR48nA4py"
)

func Test_Scene1(t *testing.T) {
	var testAuth Auth
	var log logger.Logger
	var conf conf2.Configuration
	var err error

	ctrl := gomock.NewController(t)
	dbal := mock_db.NewMockDbal(ctrl)

	t.Run("config...", func(t *testing.T) {
		conf, err = conf2.GetConfig()
		assert.NoError(t, err)

		conf.WS_ES_HOST = "http://localhost"
		conf.WS_ES_PORT = "9200"
		log = logger.MockLogger()
	})

	t.Run("NewAuth", func(t *testing.T) {
		testAuth, err = NewAuth(dbal, log, signKey, 10)
		assert.NoError(t, err)
		assert.IsType(t, &auth{}, testAuth)
	})

	token := ""
	t.Run("GenerateToken", func(t *testing.T) {
		token, err = testAuth.GenerateToken(userID)
		assert.NoError(t, err)
	})

	t.Run("Middleware", func(t *testing.T) {
		dbal.EXPECT().GetUserByID(gomock.Any(), userID).Return(db.User{ID: userID}, nil)
		mw := testAuth.Middleware()

		req := httptest.NewRequest("POST", "http://testing", nil)
		req.Header.Add("auth", token)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Run("UserFromContext", func(t *testing.T) {
				res, err := testAuth.UserFromContext(r.Context())
				assert.NoError(t, err)
				assert.Equal(t, userID, res.ID)
			})
		})

		next := mw(nextHandler)
		next.ServeHTTP(httptest.NewRecorder(), req)
	})
}
