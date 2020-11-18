package myorm

import (
	"fmt"
	"reflect"
	"strings"
)

func (ev env) InsertOne(in interface{}) error {
	conn, e := ev.db.GetConn()
	if e != nil {
		return e
	}
	defer ev.db.PutConn(conn)

	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	lot := set(t, v)

	dml := "INSERT INTO `" + t.Name() + "` SET " + strings.Join(lot, ",")

	_, e = conn.Exec(dml)

	return e
}

func set(t reflect.Type, v reflect.Value) []string {

	var lot []string

	for i := 0; i < t.NumField(); i++ {

		if !reflect.DeepEqual(v.Field(i).Interface(), reflect.Zero(v.Field(i).Type()).Interface()) {
			val := value(v.Field(i))
			lot = append(lot, t.Field(i).Name+"='"+val+"'")
		}
	}

	return lot
}

func (ev env) InsertAll(arr interface{}) error {
	conn, e := ev.db.GetConn()
	if e != nil {
		return e
	}
	defer ev.db.PutConn(conn)

	s := reflect.ValueOf(arr)
	for i := 0; i < s.Len(); i++ {

		t := reflect.TypeOf(s.Index(i).Interface())
		v := reflect.ValueOf(s.Index(i).Interface())

		lot := set(t, v)

		dml := "INSERT INTO `" + t.Name() + "` SET " + strings.Join(lot, ",")

		fmt.Println(dml)

		_, e = conn.Exec(dml)
		if e != nil {
			return e
		}
	}

	return nil
}
