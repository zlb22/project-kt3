package pontos

import (
	"keti3/pkg/tal_content_arch_lib/pontos/bootstrap"
	"keti3/pkg/tal_content_arch_lib/pontos/logger"
	"keti3/pkg/tal_content_arch_lib/pontos/orm"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

const (
	Master        = "master"
	DefaultSqlMap = "mysql_map.master"
	MapDefKey     = "default"
)

type MysqlClient struct {
	Orm *gorm.DB
}

func newMysqlClient(config *orm.MysqlConfig, log *logger.ZLogger) (*MysqlClient, error) {
	ormObj, err := orm.NewOrm(config, orm.NewSQLLogger(log))
	if err != nil {
		return nil, err
	}
	return &MysqlClient{
		Orm: ormObj,
	}, nil
}

type Model struct {
	*orm.Model
}

func NewModel(tableName string) *Model {
	return &Model{
		Model: &orm.Model{
			TableName: tableName,
			DB:        Mysql.Orm,
		},
	}
}

func NewModelByName(tableName, name string) *Model {
	ormObj := MysqlMap[name]
	return &Model{
		Model: &orm.Model{
			TableName: tableName,
			DB:        ormObj.Orm,
		},
	}
}

// InitMysql 使用关键字 [mysql] 初始化MySQL配置
func InitMysql(log *logger.ZLogger) error {
	if !Config.IsSet(sqlSection) {
		return nil
	}

	conf := orm.NewConfig()
	err := Config.UnmarshalKey(sqlSection, conf)
	if err != nil {
		return err
	}

	Mysql, err = newMysqlClient(conf, log)
	if err != nil {
		return err
	}
	return nil
}

// InitMysqlMap 使用关键字 [mysql_map] 初始化多实例的MySQL配置，支持主从
// 主从时 pontos.MysqlMap["a"], pontos.MysqlMap["b"]
// [mysql_map.a.master] [mysql_map.a.slave]
// [mysql_map.b.master] [mysql_map.b.slave] [mysql_map.b.slave2]
//
// 非主从时 pontos.MysqlMap["a"], pontos.MysqlMap["b"]
// [mysql_map.a]
// [mysql_map.b]
func InitMysqlMap(log *logger.ZLogger) error {
	if !Config.IsSet(sqlMapSection) {
		return nil
	}

	// 兼容之前默认用法, [mysql_map.master] [mysql_map.slave] => 直接使用 pontos.Mysql 或者 pontos.MysqlMap["default"]
	if Config.IsSet(DefaultSqlMap) {
		configMap, err := unmarshalDefMap()
		if err != nil {
			return err
		}

		if _, ok := configMap[Master]; ok {
			if ormObj, err := orm.NewMLOrm(configMap, orm.NewSQLLogger(log)); err != nil {
				return err
			} else {
				Mysql = &MysqlClient{ormObj}
				MysqlMap[MapDefKey] = &MysqlClient{ormObj}
			}
		}

		return nil
	}

	var isMs bool
	mapCfg := Config.GetStringMap(sqlMapSection)
	//判断是主从多实例还是单纯多实例
	for _, val := range mapCfg {
		m := cast.ToStringMap(val)
		if _, ok := m[Master]; ok {
			isMs = true
			break
		}
	}

	//非主从多实例
	if !isMs {
		return pureMysqlMap(log)
	}

	//多实例主从处理
	groupMap := make(map[string]map[string]*orm.MysqlConfig, 2)
	if err := Config.UnmarshalKey(sqlMapSection, &groupMap); err != nil {
		return err
	}

	for key, value := range groupMap {
		ormObj, err := orm.NewMLOrm(value, orm.NewSQLLogger(log))
		if err != nil {
			return err
		}

		MysqlMap[key] = &MysqlClient{ormObj}

		//[mysql_map.default.master] => pontos.Mysql 或者 pontos.MysqlMap["default"]
		if key == MapDefKey {
			Mysql = &MysqlClient{ormObj}
		}
	}
	return nil
}

func unmarshalDefMap() (map[string]*orm.MysqlConfig, error) {
	configMap := make(map[string]*orm.MysqlConfig, 2)
	if err := Config.UnmarshalKey(sqlMapSection, &configMap); err != nil {
		return nil, err
	}

	return configMap, nil
}

func pureMysqlMap(log *logger.ZLogger) error {
	configMap, err := unmarshalDefMap()
	if err != nil {
		return err
	}

	for name, config := range configMap {
		if cli, err := newMysqlClient(config.Init(), log); err != nil {
			return err
		} else {
			MysqlMap[name] = cli
		}
	}

	return nil
}

func CloseMysql() bootstrap.AfterServerFunc {
	return func() {
		if Config.IsSet(sqlSection) {
			sqlDB, _ := Mysql.Orm.DB()
			_ = sqlDB.Close()
		}
	}
}

func CloseMysqlMap() bootstrap.AfterServerFunc {
	return func() {
		if Config.IsSet(sqlMapSection) {
			for _, cli := range MysqlMap {
				sqlDB, _ := cli.Orm.DB()
				_ = sqlDB.Close()
			}
		}
	}
}
