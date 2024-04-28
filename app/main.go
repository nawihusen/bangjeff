package main

import (
	"bangjeff/config"
	"strconv"

	// "time"

	_DeliveryHTTPArticle "bangjeff/bangjeff/delivery/http"
	_RepoMySQL "bangjeff/bangjeff/repository/mysql"
	_Usecase "bangjeff/bangjeff/usecase"

	"database/sql"
	"flag"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// CLI options parse
	configFile := flag.String("c", "config.yaml", "Config file")
	flag.Parse()

	// Config file
	config.ReadConfig(*configFile)

	// Set log level
	switch viper.GetString("server.log_level") {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// Initialize database
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", viper.GetString("mysql.user"), viper.GetString("mysql.password"), viper.GetString("mysql.host"), viper.GetString("mysql.port"), viper.GetString("mysql.database"))
	val := url.Values{}
	val.Add("multiStatements", "true")
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Migrate database if any new schema
	driver, err := mysql.WithInstance(dbConn, &mysql.Config{})
	if err == nil {
		mig, err := migrate.NewWithDatabaseInstance(viper.GetString("mysql.path_migrate"), viper.GetString("mysql.database"), driver)
		log.Info(viper.GetString("mysql.path_migrate"))
		if err == nil {
			err = mig.Up()
			if err != nil {
				if err == migrate.ErrNoChange {
					log.Debug("No database migration")
				} else {
					log.Error(err)
				}
			} else {
				log.Info("Migrate database success")
			}
			version, dirty, err := mig.Version()
			if err != nil && err != migrate.ErrNilVersion {
				log.Error(err)
			}
			log.Debug("Current DB version: " + strconv.FormatUint(uint64(version), 10) + "; Dirty: " + strconv.FormatBool(dirty))
		} else {
			log.Warn(err)
		}
	} else {
		log.Warn(err)
	}

	// Register repository & usecase
	repoMySQLUser := _RepoMySQL.NewMySQLUserRepository(dbConn)
	usecaseUser := _Usecase.NewUserUsecase(repoMySQLUser)

	// Initialize HTTP web framework
	app := fiber.New(fiber.Config{
		Prefork:       viper.GetBool("server.prefork"),
		StrictRouting: viper.GetBool("server.strict_routing"),
		CaseSensitive: viper.GetBool("server.case_sensitive"),
		BodyLimit:     viper.GetInt("server.body_limit"),
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("middleware.allows_origin"),
	}))

	// HTTP routing
	app.Get(viper.GetString("server.base_path")+"/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	_DeliveryHTTPArticle.RouterAPI(app, usecaseUser)

	// go func() {
	if err := app.Listen(":" + viper.GetString("server.port")); err != nil {
		log.Fatal(err)
	}
	// }()

	// Wait for interrupt signal to gracefully shutdown the server
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// <-quit
	// log.Info("Gracefully shutdown")
	// app.Shutdown()
}
