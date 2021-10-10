package mysql

import (
	"fmt"
	"time"

	"golang.org/x/xerrors"

	"database/sql" //nolint
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/go-sql-driver/mysql" //nolint

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/mysql/const" //nolint
)

type Client struct {
	Driver *entsql.Driver
}

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"
)

func NewMysqlClient() (*Client, error) {
	service, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		return nil, xerrors.Errorf("Fail to query mysql service: %v", err)
	}

	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	myServiceName := config.GetStringValueWithNameSpace("", config.KeyHostname)
	dbname := config.GetStringValueWithNameSpace(myServiceName, keyDBName)

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=True", username, password, service.Address, dbname))
	if err != nil {
		return nil, xerrors.Errorf("Fail to initialize sql driver: %v", err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	return &Client{
		Driver: entsql.OpenDB("mysql", db),
	}, nil
}
