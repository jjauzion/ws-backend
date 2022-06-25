package graph

import (
	"github.com/42-AI/ws-backend/conf"
	"github.com/42-AI/ws-backend/db"
	"github.com/42-AI/ws-backend/internal/auth"
	"github.com/42-AI/ws-backend/internal/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log     logger.Logger
	Dbal    db.Dbal
	Config  conf.Configuration
	ApiHost string
	ApiPort string
	Auth    auth.Auth
}
