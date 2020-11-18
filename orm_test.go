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

func TestInsertOne(t *testing.T) {

	db := myorm.Dial("work:work@tcp(127.0.0.1:3306)/orm", 10)

	o := Order{
		Hash:  "One",
		Price: 50.63,
		Qty:   10,
	}

	e := db.InsertOne(o)
	if e != nil {
		t.Errorf("TestInsertOne failed: %v", e)
	}
}

func TestInsertAll(t *testing.T) {

	db := myorm.Dial("work:work@tcp(127.0.0.1:3306)/orm", 10)

	o := []Order{
		{Hash: "Two", Price: 3.5},
		{Hash: "Three", Price: 5.36, Qty: 3},
		{Hash: "Four", Price: 2.3},
		{Hash: "Five", Price: 1.2, Qty: 2},
	}

	e := db.InsertAll(o)
	if e != nil {
		t.Errorf("TestInsertAll failed: %v", e)
	}
}

func TestFindByID(t *testing.T) {
	db := myorm.Dial("work:work@tcp(127.0.0.1:3306)/orm", 10)

	var o Order

	e := db.FindByID(&o, 1)

	if e != nil {
		t.Errorf("TestFindByID failed: %v", e)
		return
	}

	if o.Qty != 10 {
		t.Errorf("TestFindByID value mismatched")
	}
}

func TestFindOne(t *testing.T) {
	db := myorm.Dial("work:work@tcp(127.0.0.1:3306)/orm", 10)

	var o Order

	e := db.FindOne(&o, map[string]interface{}{"Qty": 3})

	if e != nil {
		t.Errorf("TestFindOne failed: %v", e)
		return
	}

	if o.Hash != "Three" {
		t.Errorf("TestFindOne value mismatched")
	}
}

func TestFind(t *testing.T) {
	db := myorm.Dial("work:work@tcp(127.0.0.1:3306)/orm", 10)

	z, e := db.Find(Order{}, nil).ByEqAnd(map[string]interface{}{"Qty": 1, "Hash": "Two"})
	if e != nil {
		t.Errorf("TestFind ByEqAnd failed: %v", e)
		return
	}

	if len(z) != 1 {
		t.Errorf("TestFind ByEqAnd value mismatched")
		return
	}

	z, e = db.Find(Order{}, nil).ByEqOr(map[string]interface{}{"Qty": 1, "Hash": "Five"})
	if e != nil {
		t.Errorf("TestFind ByEqOr failed: %v", e)
		return
	}

	if len(z) != 3 {
		t.Errorf("TestFind ByEqOr value mismatched")
	}
}

func TestUpdateByID(t *testing.T) {
	db := myorm.Dial("work:work@tcp(127.0.0.1:3306)/orm", 10)

	o := Order{
		Id:   2,
		Hash: "Twotwo",
		Qty:  22,
	}

	x, e := db.UpdateByID(o)

	if e != nil {
		t.Errorf("TestUpdateByID failed: %v", e)
	} else if !x {
		t.Errorf("TestUpdateByID values not updated")
	}
}

func TestDeleteByID(t *testing.T) {
	db := myorm.Dial("work:work@tcp(127.0.0.1:3306)/orm", 10)

	oi := Order{
		Id:    6,
		Hash:  "Sex",
		Price: 6.9,
	}
	db.InsertOne(oi)

	o := Order{
		Id: 6,
	}

	x, e := db.DeleteByID(o)

	if e != nil {
		t.Errorf("TestDeleteByID failed: %v", e)
	} else if !x {
		t.Errorf("TestDeleteByID not deleted")
	}
}