package db

import (
	"fmt"
	"github.com/jjauzion/ws-backend/internal"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/jjauzion/ws-backend/graph/model"
)

func (es *esHandler) GetUserByEmail(email string) (user model.User, err error) {
	logger := internal.GetLogger()
	search := `{
	  "query": {
		"match_all": {}
	  }
    }`
	res, err := es.search([]string{userIndex}, strings.NewReader(search))
	if err != nil {
		logger.Error("", zap.Error(err))
	}
	fmt.Println(res.String())
	user = model.User{
		ID: "1",
		Email: "f",
		Login: "l",
		CreatedAt: time.Now(),
	}
	return
}
