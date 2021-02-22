package graph

import (
	"github.com/jjauzion/ws-backend/db"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log *zap.Logger
	DB  db.DatabaseHandler
}
