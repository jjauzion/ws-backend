package graph

import (
	"github.com/jjauzion/ws-backend/conf"
	"github.com/jjauzion/ws-backend/db"
	"github.com/jjauzion/ws-backend/internal/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log     *logger.Logger
	DB      db.DatabaseHandler
	Config	conf.Configuration
	ApiPort string
}
