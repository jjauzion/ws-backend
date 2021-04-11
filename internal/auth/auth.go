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
	userCtxKey  = "user"
	userIDClaim = "user_id"
	authHeader  = "auth"
)

type option int

const (
	OptAllowAll option = iota
	OptAuthenticatedUser
	OptOnlyAdmin
)

type Auth interface {
	Middleware() func(http.Handler) http.Handler
	GenerateToken(userID string) (string, error)
	UserFromContext(ctx context.Context, authOpt option) (db.User, error)
}

type auth struct {
	log  logger.Logger
	dbal db.Dbal

	signKey  []byte
	duration int
}

func NewAuth(dbal db.Dbal, log logger.Logger, signKey string, tokenDuration int) (Auth, error) {
	return &auth{
		signKey:  []byte(signKey),
		log:      log,
		dbal:     dbal,
		duration: tokenDuration,
	}, nil
}

// Middleware decodes the share session and packs the session into context
func (m *auth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.log.Debug("start middleware...")

			forbidden := func(fields ...zap.Field) {
				m.log.Debug("access denied: ", fields...)
				next.ServeHTTP(w, r)
				return
			}

			tokenHeader := r.Header.Get(authHeader)
			if tokenHeader == "" {
				m.log.Debug("unauthenticated user")
				next.ServeHTTP(w, r)
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("wrong signin method")
				}

				return m.signKey, nil
			})
			if err != nil {
				switch err.(type) {
				case *jwt.ValidationError:
					forbidden(
						zap.String("reason", "invalid token"),
						zap.Error(err),
					)
				default:
					forbidden(
						zap.String("reason", "cannot parse jwt token"),
						zap.String("jwt", tokenHeader),
						zap.Error(err),
					)
				}
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

			ctx := context.WithValue(r.Context(), userCtxKey, user)
			m.log.Debug("authenticated user", zap.String("user_email", user.Email))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *auth) UserFromContext(ctx context.Context, authOpt option) (db.User, error) {
	user, ok := ctx.Value(userCtxKey).(db.User)
	if !ok && authOpt != OptAllowAll {
		m.log.Error("auth failed: wrong user value on context")
		return db.User{}, fmt.Errorf("%d", http.StatusUnauthorized)
	}
	if authOpt == OptOnlyAdmin && !user.Admin {
		m.log.Error("auth failed: non admin user")
		return db.User{}, fmt.Errorf("%d", http.StatusForbidden)
	}
	return user, nil
}

func (m *auth) GenerateToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims[userIDClaim] = userID
	duration := time.Now().Add(time.Hour * time.Duration(m.duration)).Unix()
	claims["exp"] = duration

	ret, err := token.SignedString(m.signKey)
	if err != nil {
		return "", err
	}

	return ret, nil
}
