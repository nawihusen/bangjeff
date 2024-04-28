package config

import (
	"bytes"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Config main structure
type Config struct {
	Server Server `yaml:"server"`
	MySQL  MySQL  `yaml:"mysql"`
}

// Server is server related config
type Server struct {
	// Port is the local machine TCP Port to bind the HTTP Server to
	Port string `yaml:"port"`

	// Prefork will spawn multiple Go processes listening on the same port
	Prefork bool `yaml:"prefork"`

	// StrictRouting
	// When enabled, the router treats /foo and /foo/ as different.
	// Otherwise, the router treats /foo and /foo/ as the same.
	StrictRouting bool `yaml:"strict_routing"`

	// CaseSensitive
	// When enabled, /Foo and /foo are different routes.
	// When disabled, /Foo and /foo are treated the same.
	CaseSensitive bool `yaml:"case_sensitive"`

	// BodyLimit
	// Sets the maximum allowed size for a request body, if the size exceeds
	// the configured limit, it sends 413 - Request Entity Too Large response.
	BodyLimit int `yaml:"body_limit"`

	// Concurrency maximum number of concurrent connections
	Concurrency int `yaml:"concurrency"`

	Timeout Timeout `yaml:"timeout"`

	// LogLevel is log level, available value: error, warning, info, debug
	LogLevel string `yaml:"log_level"`

	// GRPCPort is the local machine TCP port to bind the gRPC server to
	GRPCPort string `yaml:"grpc_port"`

	// BasePath is router base path
	BasePath string `yaml:"base_path"`
}

// Timeout is server timeout related config
type Timeout struct {
	// Read is the amount of time to wait until an HTTP server
	// read operation is cancelled
	Read time.Duration `yaml:"read"`

	// Write is the amount of time to wait until an HTTP server
	// write opperation is cancelled
	Write time.Duration `yaml:"write"`

	// Read is the amount of time to wait
	// until an IDLE HTTP session is closed
	Idle time.Duration `yaml:"idle"`
}

// MySQL is MySQL related config
type MySQL struct {
	// Host is the MySQL IP Address to connect to
	Host string `yaml:"host,omitempty"`

	// Port is the MySQL Port to connect to
	Port string `yaml:"port,omitempty"`

	// Database is MySQL database name
	Database string `yaml:"database"`

	// User is MySQL username
	User string `yaml:"user"`

	// Password is MySQL password
	Password string `yaml:"password"`

	// PathMigrate is directory for migration file
	PathMigrate string `yaml:"path_migrate"`
}

// Default config
var defaultConfig = &Config{
	Server: Server{
		Port:          "8000",
		Prefork:       false,
		StrictRouting: false,
		CaseSensitive: false,
		BodyLimit:     4 * 1024 * 1024,
		Concurrency:   256 * 1024,
		Timeout: Timeout{
			Read:  5,
			Write: 10,
			Idle:  120,
		},
		LogLevel: "debug",
		BasePath: "",
	},
	MySQL: MySQL{
		Host:        "localhost",
		Port:        "3306",
		Database:    "article",
		User:        "root",
		Password:    "1sampai8",
		PathMigrate: "file:./db/migration",
	},
}

func lookupEnv(parent string, rt reflect.Type, rv reflect.Value) {
	for i := 0; i < rt.NumField(); i++ {
		structField := rt.Field(i)
		tag := strings.Split(structField.Tag.Get("yaml"), ",")[0]
		if structField.Type.Kind() == reflect.Struct {
			lookupEnv(parent+strings.ToUpper(tag)+"_", structField.Type, rv.Field(i))
		} else {
			env := parent + strings.ToUpper(tag)
			value, exist := os.LookupEnv(env)
			if exist {
				log.Info(env + " = " + value)
				switch structField.Type.Kind().String() {
				case "string":
					rv.Field(i).SetString(value)
				case "bool":
					val, err := strconv.ParseBool(value)
					if err == nil {
						rv.Field(i).SetBool(val)
					}
				case "int", "int8", "int16", "int32", "int64":
					val, err := strconv.ParseInt(value, 10, 64)
					if err == nil {
						rv.Field(i).SetInt(val)
					}
				case "uint", "uint8", "uint16", "uint32", "uint64":
					val, err := strconv.ParseUint(value, 10, 64)
					if err == nil {
						rv.Field(i).SetUint(val)
					}
				case "float32", "float64":
					val, err := strconv.ParseFloat(value, 64)
					if err == nil {
						rv.Field(i).SetFloat(val)
					}
				case "slice":
					values := strings.Split(strings.ReplaceAll(value, " ", ""), ",")
					slice := reflect.MakeSlice(rt.Field(i).Type, len(values), len(values))
					for idx, val := range values {
						switch rt.Field(i).Type.String() {
						case "[]string":
							slice.Index(idx).Set(reflect.ValueOf(val))
						case "[]bool":
							v, err := strconv.ParseBool(val)
							if err == nil {
								slice.Index(idx).Set(reflect.ValueOf(v))
							}
						case "[]int", "[]int8", "[]int16", "[]int32", "[]int64":
							v, err := strconv.ParseInt(val, 10, 64)
							if err == nil {
								slice.Index(idx).Set(reflect.ValueOf(v))
							}
						case "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
							v, err := strconv.ParseUint(val, 10, 64)
							if err == nil {
								slice.Index(idx).Set(reflect.ValueOf(v))
							}
						case "[]float32", "[]float64":
							v, err := strconv.ParseFloat(val, 64)
							if err == nil {
								slice.Index(idx).Set(reflect.ValueOf(v))
							}
						}
					}
					rv.Field(i).Set(slice)
				}
			}
		}
	}
}

// Init function of config
func init() {
	config := *defaultConfig
	rt := reflect.TypeOf(&config).Elem()
	rv := reflect.ValueOf(&config).Elem()
	lookupEnv("", rt, rv)
	*defaultConfig = rv.Interface().(Config)
}

// ReadConfig is main function to read configuration file
func ReadConfig(configFile string) {
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")
	viper.AddConfigPath("/usr/local/etc/")
	viper.AddConfigPath(".")
	rt := reflect.TypeOf(defaultConfig).Elem()
	rv := reflect.ValueOf(defaultConfig).Elem()
	for i := 0; i < rt.NumField(); i++ {
		tag := strings.Split(rt.Field(i).Tag.Get("yaml"), ",")[0]
		name := rt.Field(i).Name
		viper.SetDefault(tag, rv.FieldByName(name).Interface())
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("Use default config")
			cfgYAML, err := yaml.Marshal(defaultConfig)
			if err != nil {
				log.Fatal(err)
			}
			err = viper.ReadConfig((bytes.NewBuffer(cfgYAML)))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Info("Use config file " + configFile)
	}

	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Errorf("Unable to marshal config to YAML: %v", err)
	}
	log.Info(string(bs))
}
