package postgresql

import (
	"fmt"
	"gym-badges-api/internal/repository/user"
	"gym-badges-api/internal/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	connectionFailedMsg = "postgres-gorm connection failed: %s"
)

var (
	Config GormConfiguration
)

type GormConfiguration struct {
	Host                      string `default:"127.0.0.1" envconfig:"CLOUDSQL_CONNECTION_NAME"`
	Port                      int    `default:"5432" envconfig:"CLOUDSQL_CONNECTION_PORT"`
	User                      string `default:"postgres" envconfig:"CLOUDSQL_USER"`
	Password                  string `default:"example" envconfig:"CLOUDSQL_PASSWORD"`
	DbName                    string `default:"postgres" envconfig:"CLOUDSQL_DB"`
	CloudSqlPrefix            string `default:"" envconfig:"CLOUDSQL_PREFIX"`
	TableGormPrefix           string `default:"" envconfig:"TABLEGORM_PREFIX"`
	GormLogging               bool   `default:"false" envconfig:"GORM_LOGGING"`
	GormMaxOpenConns          int    `default:"20" envconfig:"GORM_MAX_OPEN_CONNS"`
	GormMaxIdleDuration       int    `default:"0" envconfig:"GORM_MAX_IDLE_DURATION"`
	GormMaxConnectionLifeTime int    `default:"0" envconfig:"GORM_MAX_LIFETIME_DURATION"`
}

func LoadConfig() {
	utils.LoadGenericConfig(&Config)
}

func OpenConnection() *gorm.DB {

	ctxLogger := utils.BuildLogger()

	dbURI := fmt.Sprintf("host=%s%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		Config.CloudSqlPrefix, Config.Host,
		Config.Port, Config.User,
		Config.Password, Config.DbName)
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   Config.TableGormPrefix,
			SingularTable: true,
		}}

	var err error

	DbConnection, err := gorm.Open(postgres.Open(dbURI), config)
	if err != nil {
		panic(err)
	}

	sqlDB, err := DbConnection.DB()
	if err != nil {
		ctxLogger.Errorf(connectionFailedMsg, err)
		panic(err)
	}

	if Config.GormLogging {
		DbConnection.Config.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}

	sqlDB.SetMaxOpenConns(Config.GormMaxOpenConns)
	sqlDB.SetMaxIdleConns(Config.GormMaxOpenConns - 1)
	sqlDB.SetConnMaxIdleTime(time.Duration(Config.GormMaxIdleDuration) * time.Second)
	sqlDB.SetConnMaxLifetime(time.Duration(Config.GormMaxConnectionLifeTime) * time.Second)

	err = sqlDB.Ping()

	if err != nil {
		ctxLogger.Errorf(connectionFailedMsg, err)
	} else {
		ctxLogger.Info("postgres-gorm connection successfully established")
	}

	if err = DbConnection.AutoMigrate(&user.User{}); err != nil {
		ctxLogger.Errorf("postgres-gorm migration failed: %s", err)
		return nil
	}

	return DbConnection
}
