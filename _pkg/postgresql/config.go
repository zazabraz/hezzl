package postgresql

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ClientConfig struct {
	MaxConnections                               int
	MaxConnectionAttempts                        int
	WaitingDuration                              time.Duration
	Username, Password, Host, Port, DatabaseName string
}

func NewConfigFromClientConfig(cc *ClientConfig) (*pgxpool.Config, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s target_session_attrs=read-write ",
		cc.Host, cc.Port, cc.DatabaseName, cc.Username, cc.Password,
	)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	return config, nil
}
