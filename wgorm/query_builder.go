package wgorm

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type WGorm struct {
	*gorm.DB
}

func NewWGorm(db *gorm.DB) *WGorm {
	return &WGorm{DB: db}
}

func (q *WGorm) Gorm() *gorm.DB {
	return q.DB
}

// IfWhere if conditions is true, where will add
func (q *WGorm) IfWhere(conditions bool, query interface{}, arg ...interface{}) *WGorm {
	if !conditions || q.Error != nil {
		return q
	}
	if arg == nil || len(arg) == 0 {
		q.DB = q.Where(query)
	} else {
		q.DB = q.Where(query, arg)
	}
	return q
}

// NzWhere - NonzeroWhere
func (q *WGorm) NzWhere(field string, arg interface{}) *WGorm {
	if q.Error != nil {
		return q
	}
	//最后调用的都是gorm本身的函数,而gorm本身就会将自身复制一遍,因此不需要在此函数内调用getInstance
	switch v := arg.(type) {
	case string:
		q.IfWhere(v != "", field, arg)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		q.IfWhere(v != 0, field, arg)
	case float32, float64:
		q.IfWhere(v != 0, field, arg)
	case time.Time:
		q.IfWhere(!v.IsZero(), field, arg)
	case decimal.Decimal:
		q.IfWhere(!v.IsZero(), field, arg)
	case bool:
		q.IfWhere(v, field, arg)
	default:
		q.Error = fmt.Errorf("ERR wgorm AddWhere TypeOf \"" + reflect.TypeOf(arg).String() + "\" Not Suppose")
		return q
	}
	q.DB = q.Where(field, arg)
	return q
}

// NzRadius 范围查询
func (q *WGorm) NzRadius(field string, start, end interface{}) *WGorm {
	if q.Error != nil {
		return q
	}
	strStart := fmt.Sprintf("%s >= ?", field)
	strEnd := fmt.Sprintf("%s <= ?", field)
	if reflect.TypeOf(start) != reflect.TypeOf(end) {
		q.Error = fmt.Errorf("ERR: wgorm AddRadius - The type of start and end are not euqal")
	}
	q = q.NzWhere(strStart, start)
	q = q.NzWhere(strEnd, end)
	return q
}

func (q *WGorm) NzLimit(page, pageSize int) *WGorm {
	if q.Error != nil {
		return q
	}
	if pageSize < 0 {
		q.Error = fmt.Errorf("ERR: wgorm AddPageLimit -pagesize Less than 0")
	}
	if page <= 0 {
		q.Error = fmt.Errorf("ERR: wgorm AddPageLimit -page Less than 0")
	}
	q.DB = q.Offset((page - 1) * pageSize).Limit(pageSize)
	return q
}

func (q *WGorm) Count(count *int64) *WGorm {
	if q.Error != nil {
		return q
	}
	q.DB.Count(count)
	return q
}

func (q *WGorm) ToSQL() string {
	stmt := q.Session(&gorm.Session{DryRun: true, SkipDefaultTransaction: true}).Find(nil).Statement
	return q.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
}
