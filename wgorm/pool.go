package wgorm

import (
	"gorm.io/gorm"
	"sync"
)

type wgormPool struct {
	pool sync.Map
}

var Pool wgormPool

func InitWGormPool() {

}

func RegisterWGorm(wg WGorm) error {
	return nil
}

func Use() WGorm {
	return WGorm{}
}

func UseGormDB() *gorm.DB {
	return nil
}
