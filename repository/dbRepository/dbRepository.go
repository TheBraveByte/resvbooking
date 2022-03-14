package dbRepository

import (
	"database/sql"
	"github.com/dev-ayaa/resvbooking/pkg/config"
	"github.com/dev-ayaa/resvbooking/repository"
)

type PostgresDBRepository struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepository  return a struct which is of type interfaces
//which contains the database connection pool and the application configuration
func NewPostgresRepository(a *config.AppConfig, conn *sql.DB) repository.DatabaseRepository {
	return &PostgresDBRepository{App: a, DB: conn}
}
