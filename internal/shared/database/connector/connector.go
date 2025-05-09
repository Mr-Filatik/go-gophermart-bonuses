package connector

import (
	"errors"
	"time"

	migrations "github.com/Mr-Filatik/go-gophermart-bonuses"
	"github.com/Mr-Filatik/go-gophermart-bonuses/internal/shared/logger"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DatabaseConnectionCount = 10

type Connector struct {
	db  *gorm.DB
	log logger.Logger
}

func New(log logger.Logger) *Connector {
	return &Connector{
		log: log,
	}
}

func (c *Connector) Connect(dbURI string) error {
	db, openErr := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURI,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if openErr != nil {
		return errors.New(openErr.Error())
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(DatabaseConnectionCount)
	sqlDB.SetMaxOpenConns(DatabaseConnectionCount)
	sqlDB.SetConnMaxLifetime(time.Hour)

	pingErr := sqlDB.Ping()
	if pingErr != nil {
		return errors.New(pingErr.Error())
	}

	c.db = db
	c.log.Info("Database is connected")

	return nil
}

func (c *Connector) ApplyMigrations() error {
	sqlDB, _ := c.db.DB()

	goose.SetBaseFS(migrations.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return errors.New(err.Error())
	}

	if err := goose.Up(sqlDB, migrations.DirMigrations); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (c *Connector) GetDB() *gorm.DB {
	return c.db
}

func (c *Connector) Close() {
	sqlDB, _ := c.db.DB()
	closeErr := sqlDB.Close()
	if closeErr != nil {
		c.log.Error("Close connection error.", closeErr)
	}
}
