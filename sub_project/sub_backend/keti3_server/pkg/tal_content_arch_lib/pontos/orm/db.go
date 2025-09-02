package orm

import (
	"database/sql"
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

const (
	Master = "master"
)

var MasterLostError = errors.New("master config lost")

type MysqlConfig struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	User                string        `mapstructure:"user"`
	Password            string        `mapstructure:"password"`
	Database            string        `mapstructure:"database"`
	Charset             string        `mapstructure:"charset"`
	Timeout             time.Duration `mapstructure:"timeout"`
	ReadTimeout         time.Duration `mapstructure:"read_timeout"`
	WriteTimeout        time.Duration `mapstructure:"write_timeout"`
	Loc                 string        `mapstructure:"loc"`
	SetMaxIdleConns     int           `mapstructure:"set_max_idle_conns"`
	SetMaxOpenConns     int           `mapstructure:"set_max_open_conns"`
	SetConnMaxLifetime  time.Duration `mapstructure:"set_conn_max_lifetime"`
	IsCutMasterRelation bool          `mapstructure:"is_cut_master_relation"` // 是否加入主从配置
}

func (m *MysqlConfig) Init() *MysqlConfig {
	if m.Port == 0 {
		m.Port = 3306
	}
	if m.Charset == "" {
		m.Charset = "utf8"
	}
	if m.Timeout == 0 {
		m.Timeout = 10 * time.Second
	}
	if m.ReadTimeout == 0 {
		m.ReadTimeout = 20 * time.Second
	}
	if m.WriteTimeout == 0 {
		m.WriteTimeout = 20 * time.Second
	}
	return m
}

func NewConfig() *MysqlConfig {
	return &MysqlConfig{
		Port:         3306,
		Charset:      "utf8",
		Timeout:      10 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
}

func newDialector(config *MysqlConfig) (*sql.DB, error) {
	opts := NewDBOption(config.Database, config.User, config.Password, config.Host)
	opts.Set(
		SetCharset(config.Charset),
		SetParseTime(true),
		SetLoc(config.Loc),
		SetTimeout(config.Timeout),
		SetReadTimeout(config.ReadTimeout),
		SetWriteTimeout(config.WriteTimeout),
	).Port(config.Port)
	sqlDB, err := opts.Open(true)
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.SetMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.SetMaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.SetConnMaxLifetime)
	return sqlDB, nil
}

func NewOrm(config *MysqlConfig, log logger.Interface) (*gorm.DB, error) {
	sqlDB, err := newDialector(config)
	if err != nil {
		return nil, err
	}

	orm, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger:          log,
		CreateBatchSize: 200,
	})
	if err != nil {
		return nil, err
	}

	if err = orm.Use(NewDefaultPlugin(config.Host, config.Database, config.User)); err != nil {
		return nil, err
	}
	return orm, nil
}

func NewMLOrm(configs map[string]*MysqlConfig, log logger.Interface) (*gorm.DB, error) {
	masterConfig, ok := configs[Master]
	if !ok {
		return nil, MasterLostError
	}
	sqlDB, err := newDialector(masterConfig.Init())
	if err != nil {
		return nil, err
	}
	orm, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger:          log,
		CreateBatchSize: 200,
	})

	if err != nil {
		return nil, err
	}

	if err = orm.Use(NewDefaultPlugin(masterConfig.Host, masterConfig.Database, masterConfig.User)); err != nil {
		return nil, err
	}

	dialectorSli := []gorm.Dialector{}
	for name, config := range configs {
		if name == Master || config.IsCutMasterRelation {
			continue
		}
		if sqlDB, err := newDialector(config.Init()); err != nil {
			return nil, err
		} else {
			dialectorSli = append(dialectorSli, mysql.New(mysql.Config{
				Conn: sqlDB,
			}))
		}
	}
	if err = orm.Use(dbresolver.Register(
		dbresolver.Config{
			Replicas: dialectorSli,
			Policy:   dbresolver.RandomPolicy{},
		})); err != nil {
		return nil, err
	}

	return orm, nil
}
