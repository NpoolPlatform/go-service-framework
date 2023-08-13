package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/hashicorp/consul/api"
)

type db struct {
	db      *sql.DB
	address string
}

var (
	mu        = sync.Mutex{}
	mysqlConn *db
)

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"

	checkDuration = time.Second * 10
	pingTimeout   = time.Second * 5
)

func init() {
	go heartbeat()
}

func GetConn() (conn *sql.DB, err error) {
	mu.Lock()
	if mysqlConn != nil {
		conn = mysqlConn.db
		mu.Unlock()
		return
	}
	mu.Unlock()

	dsn, err := getMySQLConfig()
	if err != nil {
		logger.Sugar().Warnf("call getMySQLConfig error: %v", err)
		return nil, err
	}

	conn, err = open("mysql", dsn)
	if err != nil {
		logger.Sugar().Warnf("call open error: %v", err)
		return nil, err
	}

	return
}

func getApolloConfig() (*api.AgentService, error) {
	return config.PeekService(constant.MysqlServiceName)
}

func getMySQLConfig() (string, error) {
	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	myServiceName := config.GetStringValueWithNameSpace("", config.KeyHostname)
	dbname := config.GetStringValueWithNameSpace(myServiceName, keyDBName)

	svc, err := getApolloConfig()
	if err != nil {
		logger.Sugar().Warnf("call getApolloConfig error: %v", err)
		return "", err
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true&multiStatements=true",
		username, password,
		svc.Address,
		svc.Port,
		dbname,
	), nil
}

func open(driverName, dataSourceName string) (conn *sql.DB, err error) {
	// double lock check
	mu.Lock()
	if mysqlConn != nil && mysqlConn.address == dataSourceName {
		conn = mysqlConn.db
		mu.Unlock()
		return
	}

	logger.Sugar().Infof("Reopen database %v: %v", driverName, dataSourceName)
	conn, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		mu.Unlock()
		logger.Sugar().Warnf("call Open error: %v", err)
		return nil, err
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(2)
	conn.SetMaxIdleConns(1)

	// maybe should close
	if mysqlConn != nil {
		mysqlConn.db.Close()
	}

	mysqlConn = &db{db: conn, address: dataSourceName}
	mu.Unlock()

	return conn, nil
}

func heartbeat() {
	for range time.After(checkDuration) {
		mu.Lock()
		if mysqlConn == nil {
			mu.Unlock()
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
		err := mysqlConn.db.PingContext(ctx)
		mu.Unlock()
		cancel()

		if err == nil {
			continue
		}
		logger.Sugar().Warnf("call ping mysql error: %v try to create new conn", err)
		if err != nil && strings.Contains(err.Error(), "Too many connections") {
			continue
		}

		dsn, err := getMySQLConfig()
		if err != nil {
			logger.Sugar().Warnf("call getMySQLConfig error: %v", err)
			continue
		}

		_, err = open("mysql", dsn)
		if err != nil {
			logger.Sugar().Warnf("call open error: %v", err)
		}
	}
}
