package myorm

import (
	"reflect"
	"strconv"
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

func value(v reflect.Value) string {

	switch v.Type().Name() {
	case "string":
		return v.String()
	case "int8":
		fallthrough
	case "int16":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		fallthrough
	case "int":
		return strconv.Itoa(int(v.Int()))
	case "float32":
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case "float64":
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case "bool":
		if v.Bool() {
			return "1"
		} else {
			return "0"
		}
	}

	return ""
}

func fields(typ reflect.Type, val reflect.Value) []string {
	var (
		fields []string
	)
	for i := 0; i < val.NumField(); i++ {

		x := true
		if v, ok := typ.Field(i).Tag.Lookup("myorm"); ok {
			if strings.Contains(v, "primary") {
				x = false
			}
		}

		if x && !reflect.DeepEqual(val.Field(i).Interface(), reflect.Zero(val.Field(i).Type()).Interface()) {
			v := value(val.Field(i))
			fields = append(fields, typ.Field(i).Name+"='"+v+"'")
		}
	}

	return fields
}
