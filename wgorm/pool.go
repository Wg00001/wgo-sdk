package wgorm

import (
	"errors"
	"github.com/wg00001/wgo-sdk/wg_pool"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Option struct {
	DSN                     string
	Driver                  string
	MaxConnection           int
	MaxIdleConnection       int
	TimeOut                 time.Duration
	ConnPoolCleanupInterval time.Duration
	SlowThreshold           time.Duration
}

type Pool struct {
	mu     sync.Mutex
	conn   []*WGorm
	option Option
}

var _ wg_pool.Pool[WGorm] = (*Pool)(nil)
var _ wg_pool.Option = (*Option)(nil)

func (p *Pool) Init(option wg_pool.Option) error {
	var ok bool
	p.option, ok = option.(Option)
	if !ok {
		return errors.New("Init Fail: option type fail")
	}
	for i := 0; i < p.option.MaxIdleConnection; i++ {
		cur, err := p.Open()
		if err != nil {
			return err
		}
		p.conn = append(p.conn, &cur)
	}
	return nil
}

func (p *Pool) Get() (WGorm, error) {
	return WGorm{}, nil
}

func (p *Pool) CloseAll() {

}

func (p *Pool) Len() int {
	return 0
}

func (p *Pool) Open() (WGorm, error) {
	var dialect gorm.Dialector
	switch p.option.Driver {
	case "mysql":
		dialect = mysql.Open(p.option.DSN)
	case "postgres":
		dialect = postgres.Open(p.option.DSN)
	case "sqlite":
		dialect = sqlite.Open(p.option.DSN)
	case "sqlserver":
		dialect = sqlserver.Open(p.option.DSN)
	case "clickhouse":
		dialect = clickhouse.Open(p.option.DSN)
	default:
		return WGorm{}, errors.New("Driver not support: " + p.option.Driver)
	}
	db, err := gorm.Open(dialect)
	if err != nil {
		return WGorm{}, err
	}
	return WGorm{db}, err
}
