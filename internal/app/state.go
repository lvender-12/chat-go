package app

import (
	"chat-go/internal/config"
	"database/sql"
)

type State struct {
	DB     *sql.DB
	Config *config.Config
}
