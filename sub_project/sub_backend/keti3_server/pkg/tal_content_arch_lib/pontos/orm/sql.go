package orm

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	defaultDriver = "mysql"
	defaultPort   = 3306
	cDSNFormat    = "%s%s=%s&"
)

//Setting is the signature of setting function
type Setting func(string) string

func boolSetting(source, param string, ok bool) string {
	return fmt.Sprintf(cDSNFormat, source, param, strconv.FormatBool(ok))
}

func timeSetting(source, param string, t time.Duration) string {
	//make sure 1ms<=t<24h
	if t < time.Millisecond || t >= 24*time.Hour {
		return source
	}
	return fmt.Sprintf(cDSNFormat, source, param, t)
}

func stringSetting(source, param, value string) string {
	if "" == value {
		return source
	}
	return fmt.Sprintf(cDSNFormat, source, param, value)
}

//SetCharset Sets the charset used for client-server interaction
func SetCharset(v string) Setting {
	return func(source string) string {
		return stringSetting(source, "charset", v)
	}
}

//SetLoc Sets the location for time.Time values (when using parseTime=true). "Local" sets the system's location. See time.LoadLocation for details.
func SetLoc(v string) Setting {
	return func(source string) string {
		return stringSetting(source, "loc", v)
	}
}

//SetParseTime parseTime=true changes the output type of DATE and DATETIME values to time.Time instead of []byte / string
func SetParseTime(ok bool) Setting {
	return func(source string) string {
		return boolSetting(source, "parseTime", ok)
	}
}

//SetTimeout Driver side connection timeout. The value must be a decimal number with an unit suffix ( "ms", "s", "m", "h" ), such as "30s", "0.5m" or "1m30s". To set a server side timeout, use the parameter wait_timeout.
func SetTimeout(timeout time.Duration) Setting {
	return func(source string) string {
		return timeSetting(source, "timeout", timeout)
	}
}

//SetReadTimeout I/O read timeout. The value must be a decimal number with an unit suffix ( "ms", "s", "m", "h" ), such as "30s", "0.5m" or "1m30s".
func SetReadTimeout(timeout time.Duration) Setting {
	return func(source string) string {
		return timeSetting(source, "readTimeout", timeout)
	}
}

//SetWriteTimeout I/O write timeout. The value must be a decimal number with an unit suffix ( "ms", "s", "m", "h" ), such as "30s", "0.5m" or "1m30s".
func SetWriteTimeout(timeout time.Duration) Setting {
	return func(source string) string {
		return timeSetting(source, "writeTimeout", timeout)
	}
}

//Option stands for a series of options for creating a DB
type Option struct {
	driver   string
	dbName   string
	user     string
	password string
	host     string
	port     int
	settings []Setting
}

//NewDBOption returns an Option
func NewDBOption(dbName, user, password, host string) *Option {
	return &Option{
		dbName:   dbName,
		user:     user,
		password: password,
		host:     host,
		port:     defaultPort,
		driver:   defaultDriver,
	}
}

//Port sets the server port,default 3306
func (o *Option) Port(port int) *Option {
	o.port = port
	return o
}

//Driver sets the driver, default mysql
func (o *Option) Driver(driver string) *Option {
	o.driver = driver
	return o
}

//Set receives a series of Set*-like functions
func (o *Option) Set(sets ...Setting) *Option {
	o.settings = append(o.settings, sets...)
	return o
}

//Open is used for creating a *sql.DB
//Use it at the last
func (o *Option) Open(ping bool) (*sql.DB, error) {
	db, err := open(o)
	if nil != err {
		return nil, err
	}
	if ping {
		err = db.Ping()
	}
	return db, err
}

func concatDSN(settings []Setting) string {
	s := ""
	for _, f := range settings {
		s = f(s)
	}
	return strings.TrimRight(s, "&")
}

func realDSN(info *Option) string {
	format := "%s:%s@tcp(%s:%d)/%s?%s"
	return strings.TrimRight(fmt.Sprintf(format, info.user, info.password, info.host, info.port, info.dbName, concatDSN(info.settings)), "?")
}

func open(o *Option) (*sql.DB, error) {
	return sql.Open(o.driver, realDSN(o))
}
