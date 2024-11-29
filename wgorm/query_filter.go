package wgorm

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type WGorm struct {
	*gorm.DB
}

func InitWGorm(db *gorm.DB) *WGorm {
	return &WGorm{DB: db}
}

func (q *WGorm) AddWhere(field string, arg interface{}) *WGorm {
	if q.Error != nil {
		return q
	}
	//最后调用的都是gorm本身的函数,而gorm本身就会将自身复制一遍,因此不需要在此函数内调用getInstance
	switch v := arg.(type) {
	case string:
		if v == "" || v == "ALL" {
			return q
		}
	case int:
		if v == 0 {
			return q
		}
	case int64:
		if v == 0 {
			return q
		}
	case int8:
		if v == 0 {
			return q
		}
	default:
		q.Error = fmt.Errorf("ERR wgorm AddWhere TypeOf \"" + reflect.TypeOf(arg).String() + "\" Not Suppose")
		return q
	}
	q.DB = q.Where(field, arg)
	return q
}

// AddRadius 范围查询
func (q *WGorm) AddRadius(field string, start, end interface{}) *WGorm {
	strStart := fmt.Sprintf("%s >= ?", field)
	strEnd := fmt.Sprintf("%s <= ?", field)
	if reflect.TypeOf(start) != reflect.TypeOf(end) {
		q.Error = fmt.Errorf("ERR: wgorm AddRadius - The type of start and end are not euqal")
	}
	switch v := start.(type) {
	case int, int64, int8:
		if v != 0 {
			q.DB = q.DB.Where(strStart, start)
		}
		if endInt, ok := end.(int); ok && endInt != 0 {
			q.DB = q.DB.Where(strEnd, end)
		}
	case float64, float32:
		if v != 0.0 {
			q.DB = q.DB.Where(strStart, start)
		}
		if endFloat, ok := end.(float64); ok && endFloat != 0.0 {
			q.DB = q.DB.Where(strEnd, end)
		}
	case time.Time:
		if !v.IsZero() {
			q.DB = q.DB.Where(strStart, start)
		}
		if endTime, ok := end.(time.Time); ok && !endTime.IsZero() {
			q.DB = q.DB.Where(strEnd, end)
		}
	case string:
		if v != "" {
			q.DB = q.DB.Where(strStart, start)
		}
		if endStr, ok := end.(string); ok && endStr != "" {
			q.DB = q.DB.Where(strEnd, end)
		}
	default:
		q.Error = fmt.Errorf("ERR: wgorm AddRadius - type not suppose: %T", start)
	}
	return q
}

func (q *WGorm) AddLimit(page, pageSize int) *WGorm {
	if q.Error != nil {
		return q
	}
	if pageSize <= 0 {
		q.Error = fmt.Errorf("ERR: wgorm AddPageLimit -pagesize Less than 0")
	}
	if page <= 0 {
		q.Error = fmt.Errorf("ERR: wgorm AddPageLimit -page Less than 0")
	}
	q.DB = q.Offset((page - 1) * pageSize).Limit(pageSize)
	return q
}

// AddTimeRadius 用于避免时区的影响
func (q *WGorm) AddTimeRadius(field string, gte, lte time.Time) *WGorm {
	if q.Error != nil {
		return q
	}
	if gte.IsZero() && lte.IsZero() {
		return q
	}
	if !gte.IsZero() {
		q.DB = q.Where(fmt.Sprintf("%s >= ?", field), gte.Format(time.DateOnly)+" 00:00:00")
	}
	if !lte.IsZero() {
		q.DB = q.Where(fmt.Sprintf("%s <= ?", field), lte.Format(time.DateOnly)+" 23:59:59")
	}
	return q
}
func (q *WGorm) AddDateRadius(field string, gte, lte time.Time) *WGorm {
	if q.Error != nil {
		return q
	}
	if gte.IsZero() && lte.IsZero() {
		return q
	}
	if !gte.IsZero() {
		q.DB = q.Where(fmt.Sprintf("%s >= ?", field), gte.Format(time.DateOnly))
	}
	if !lte.IsZero() {
		q.DB = q.Where(fmt.Sprintf("%s <= ?", field), lte.Format(time.DateOnly))
	}
	return q
}
