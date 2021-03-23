package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	userCtxKey    = "user"
	userIDClaim   = "user_id"
	notAuthorized = "not authorized"
	authHeader    = "auth"
)

type Auth interface {
	Middleware() func(http.Handler) http.Handler
	GenerateToken(userID string) (string, error)
	UserFromContext(ctx context.Context) (db.User, error)
}

type auth struct {
	signinKey []byte
	log       logger.Logger
	dbal      db.Dbal
}

func NewAuth(dbal db.Dbal, log logger.Logger, signinKey string) (Auth, error) {
	return &auth{
		signinKey: []byte(signinKey),
		log:       log,
		dbal:      dbal,
	}, nil
}

// Middleware decodes the share session and packs the session into context
func (m *auth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.log.Debug("start middleware...")
			forbidden := func(fields ...zap.Field) {
				m.log.Warn("access denied: ", fields...)
				http.Error(w, notAuthorized, http.StatusForbidden)
				return
			}

			tokenHeader := r.Header.Get(authHeader)
			if tokenHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("wrong signin method")
				}

				return m.signinKey, nil
			})
			if err != nil {
				forbidden(
					zap.String("reason", "cannot parse jwt token"),
					zap.String("jwt", tokenHeader),
					zap.Error(err),
				)
				return
			}

			if !token.Valid {
				forbidden(zap.String("reason", "jwt token invalid"))
				return
			}

			userID, ok := claims[userIDClaim].(string)
			if !ok {
				forbidden(
					zap.String("reason", "token user_id claim is not a string"),
					zap.String("type", fmt.Sprintf("%T", token.Header[userIDClaim])),
				)
				return
			}

			user, err := m.dbal.GetUserByID(r.Context(), userID)
			if err != nil {
				forbidden(
					zap.String("reason", "cannot find user on database"),
					zap.String("user_id", userID),
					zap.Error(err),
				)
				return
			}

			m.log.Debug("...end middleware", zap.String("user_email", user.Email))
			ctx := context.WithValue(r.Context(), userCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *auth) UserFromContext(ctx context.Context) (db.User, error) {
	user, ok := ctx.Value(userCtxKey).(db.User)
	if !ok {
		return db.User{}, fmt.Errorf("wrong value on context")
	}

	return user, nil
}

func (m *auth) GenerateToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims[userIDClaim] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	ret, err := token.SignedString(m.signinKey)
	if err != nil {
		m.log.Error("cannot generate token", zap.Error(err))
		return "", err
	}

	return ret, nil
}
