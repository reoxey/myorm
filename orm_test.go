package myorm_test

import (
	"testing"

	"myorm"
)

type Order struct {
	Id    int     `myorm:"primary,auto_increment"`
	Hash  string  `myorm:"notnull,index"`
	Price float32 `myorm:"notnull"`
	Qty   int     `myorm:"notnull,default='1'"`
}

func TestCreate(t *testing.T) {

	db := myorm.Dial("work:work@tcp(127.0.0.1:3306)/orm", 10)

	e := db.Create(Order{})

	if e != nil {
		t.Errorf("TestCreate failed: %v", e)
	}
}
