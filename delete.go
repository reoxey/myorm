package myorm

import (
	"fmt"
	"reflect"
	"strings"
)

func (ev env) DeleteByID(in interface{}) (bool, error) {

	typ := reflect.TypeOf(in)
	val := reflect.ValueOf(in)

	if val.Kind() == reflect.Ptr {
		return false, ErrInvalidPTR
	}

	conn, e := ev.db.GetConn()
	if e != nil {
		return false, e
	}
	defer ev.db.PutConn(conn)

	var pri, idx string
	for i := 0; i < val.NumField(); i++ {

		if v, ok := typ.Field(i).Tag.Lookup("myorm"); ok {
			if strings.Contains(v, "primary") {
				pri = typ.Field(i).Name
				idx = cast(val.Field(i).Interface())
				break
			}
		}
	}
	if idx == "" {
		return false, ErrInvalidID
	}
	if pri == "" {
		return false, ErrNoPrimaryID
	}

	s := "DELETE FROM `" + typ.Name() + "` WHERE " + pri + " = " + idx

	fmt.Println(s)

	x, e := conn.Exec(s)

	if e != nil {
		return false, e
	}

	return x.AffectedRows != 0, nil
}

func (ev env) DeleteAnd(in interface{}) (bool, error) {

	typ := reflect.TypeOf(in)
	val := reflect.ValueOf(in)

	if val.Kind() == reflect.Ptr {
		return false, ErrInvalidPTR
	}

	conn, e := ev.db.GetConn()
	if e != nil {
		return false, e
	}
	defer ev.db.PutConn(conn)

	f := fields(typ, val)

	s := "DELETE FROM `" + typ.Name() + "` WHERE " + strings.Join(f, " AND ")

	fmt.Println(s)

	x, e := conn.Exec(s)

	if e != nil {
		return false, e
	}

	return x.AffectedRows != 0, nil
}

func (ev env) DeleteOr(in interface{}) (bool, error) {

	typ := reflect.TypeOf(in)
	val := reflect.ValueOf(in)

	if val.Kind() == reflect.Ptr {
		return false, ErrInvalidPTR
	}

	conn, e := ev.db.GetConn()
	if e != nil {
		return false, e
	}
	defer ev.db.PutConn(conn)

	f := fields(typ, val)

	s := "DELETE FROM `" + typ.Name() + "` WHERE " + strings.Join(f, " OR ")

	fmt.Println(s)

	x, e := conn.Exec(s)

	if e != nil {
		return false, e
	}

	return x.AffectedRows != 0, nil
}
