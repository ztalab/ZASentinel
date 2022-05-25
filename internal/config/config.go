package config

import (
	"encoding/base64"
	"github.com/ztalab/ZASentinel/pkg/influxdb"
	"github.com/ztalab/ZASentinel/pkg/logger"
	"github.com/ztalab/ZASentinel/pkg/util/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/koding/multiconfig"
)

var (
	// C Global configuration (Must Load first, otherwise the configuration will not be available)
	C    = new(Config)
	once sync.Once
	Is   = new(I)
)

// I ...
type I struct {
	Metrics    *influxdb.Metrics
	HttpClient *http.Client
}

// MustLoad load config
func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}
		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
}

func ParseConfigByEnv() error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	C.Common.Hostname = hostname
	if v := os.Getenv("PODIP"); v != "" {
		C.Common.PodIP = v
	}
	if v := os.Getenv("CONTRO_HOST"); v != "" {
		C.Common.ControHost = v
	}
	if v := os.Getenv("LOG_HOOK_ENABLED"); v == "true" {
		C.Log.EnableHook = true
	}
	if v := os.Getenv("LOG_REDIS_ADDR"); v != "" {
		C.LogRedisHook.Addr = v
	}
	if v := os.Getenv("LOG_REDIS_KEY"); v != "" {
		C.LogRedisHook.Key = v
	}
	if C.Certificate.CertPem == "" {
		if v := os.Getenv("CERT_PEM"); v != "" {
			cv, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return err
			}
			C.Certificate.CertPem = string(cv)
		} else {
			cert, err := ioutil.ReadFile(C.Certificate.CertPemPath)
			if err != nil {
				logger.Errorf("can not open the `%s`, err is %+v", C.Certificate.CertPemPath, err)
			}
			C.Certificate.CertPem = string(cert)
		}
	}
	if C.Certificate.KeyPem == "" {
		if v := os.Getenv("KEY_PEM"); v != "" {
			cv, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return err
			}
			C.Certificate.KeyPem = string(cv)
		} else {
			cert, err := ioutil.ReadFile(C.Certificate.KeyPemPath)
			if err != nil {
				logger.Errorf("can not open the `%s`, err is %+v", C.Certificate.KeyPemPath, err)
			}
			C.Certificate.KeyPem = string(cert)
		}
	}
	if C.Certificate.CaPem == "" {
		if v := os.Getenv("CA_PEM"); v != "" {
			cv, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return err
			}
			C.Certificate.CaPem = string(cv)
		} else {
			cert, err := ioutil.ReadFile(C.Certificate.CaPemPath)
			if err != nil {
				logger.Errorf("can not open the `%s`, err is %+v", C.Certificate.CaPemPath, err)
			}
			C.Certificate.CaPem = string(cert)
		}
	}
	// influxdb
	if v := os.Getenv("INFLUXDB_ENABLED"); v == "true" {
		C.Influxdb.Enabled = true
	}
	if v := os.Getenv("INFLUXDB_ADDRESS"); v != "" {
		C.Influxdb.Address = v
	}
	if v := os.Getenv("INFLUXDB_PORT"); v != "" {
		cv, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		C.Influxdb.Port = cv
	}
	if v := os.Getenv("INFLUXDB_USERNAME"); v != "" {
		C.Influxdb.Username = v
	}
	if v := os.Getenv("INFLUXDB_PASSWORD"); v != "" {
		C.Influxdb.Password = v
	}
	if v := os.Getenv("INFLUXDB_DATABASE"); v != "" {
		C.Influxdb.Database = v
	}
	if v := os.Getenv("INFLUXDB_PRECISION"); v != "" {
		C.Influxdb.Precision = v
	}
	return nil
}

func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

type Config struct {
	RunMode      string
	PrintConfig  bool
	Common       Common
	Machine      Machine
	Log          Log
	LogRedisHook LogRedisHook
	Certificate  Certificate
	Influxdb     Influxdb
}

func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

func (c *Config) IsReleaseMode() bool {
	return c.RunMode == "release"
}

type LogHook string

func (h LogHook) IsRedis() bool {
	return h == "redis"
}

type Log struct {
	Level         int
	Format        string
	Output        string
	OutputFile    string
	EnableHook    bool
	HookLevels    []string
	Hook          LogHook
	HookMaxThread int
	HookMaxBuffer int
}

type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

type LogRedisHook struct {
	Addr string
	Key  string
}

// Common Configuration parameters
type Common struct {
	UniqueID   string
	AppName    string
	Hostname   string
	PodIP      string
	ControHost string
}

// Certificate certificate
type Certificate struct {
	CertPem string
	CaPem   string
	KeyPem  string

	CertPemPath string
	CaPemPath   string
	KeyPemPath  string
}

type Influxdb struct {
	Enabled             bool
	Address             string
	Port                int
	Username            string
	Password            string
	Database            string
	Precision           string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	FlushTime           int
	FlushSize           int
}

// Machine
type Machine struct {
	MachineId string
	Cookie    string
}

const lockPath = "./machine.lock"

func (a *Machine) SetMachineId(macid string) {
	C.Machine.MachineId = macid
	a.MachineId = macid
}

func (a *Machine) SetCookie(cookie string) {
	C.Machine.Cookie = cookie
	a.Cookie = cookie
}

func (a *Machine) Write() error {
	b, err := json.Marshal(a)
	if err != nil {
		return err
	}
	// Write lock to filesystem to indicate an existing running daemon.
	err = os.WriteFile(lockPath, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (a *Machine) Read() (*Machine, error) {
	in, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, err
	}
	var result Machine
	err = json.Unmarshal(in, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
