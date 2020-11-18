package myorm

import (
	"reflect"
	"strings"
)

func (ev env) UpdateByID(in interface{}) (bool, error) {

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

	var (
		fields []string
		pri    string
		idx    string
	)
	for i := 0; i < val.NumField(); i++ {

		x := true
		if v, ok := typ.Field(i).Tag.Lookup("myorm"); ok {
			if strings.Contains(v, "primary") {
				pri = typ.Field(i).Name
				x = false
				idx = cast(val.Field(i).Interface())
			}
		}

		if x && !reflect.DeepEqual(val.Field(i).Interface(), reflect.Zero(val.Field(i).Type()).Interface()) {
			v := value(val.Field(i))
			fields = append(fields, typ.Field(i).Name+"='"+v+"'")
		}
	}
	if idx == "" {
		return false, ErrInvalidID
	}
	if pri == "" {
		return false, ErrNoPrimaryID
	}

	s := "UPDATE `" + typ.Name() + "` SET " + strings.Join(fields, ",") + " WHERE " + pri + " = " + idx

	x, e := conn.Exec(s)

	if e != nil {
		return false, e
	}

	return x.AffectedRows != 0, nil
}

func (ev env) UpdateAnd(in interface{}, and map[string]interface{}) (bool, error) {

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

	s := "UPDATE `" + typ.Name() + "` SET " + strings.Join(f, ",") + " WHERE " + joinMap(and, "AND")

	x, e := conn.Exec(s)

	if e != nil {
		return false, e
	}

	return x.AffectedRows != 0, nil
}

func (ev env) UpdateOr(in interface{}, or map[string]interface{}) (bool, error) {

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

	s := "UPDATE `" + typ.Name() + "` SET " + strings.Join(f, ",") + " WHERE " + joinMap(or, "OR")

	x, e := conn.Exec(s)

	if e != nil {
		return false, e
	}

	return x.AffectedRows != 0, nil
}
