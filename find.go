package myorm

import (
	"reflect"
	"strings"

	"github.com/pubnative/mysqldriver-go"
)

func (ev env) FindByID(in, id interface{}) error {

	typ := reflect.TypeOf(in)
	val := reflect.ValueOf(in)

	if val.Kind() != reflect.Ptr {
		return ErrRequiredPTR
	}

	idx := cast(id)
	if idx == "" {
		return ErrInvalidID
	}

	var (
		fields []string
		pri    string
	)
	for i := 0; i < reflect.Indirect(val).NumField(); i++ {

		fields = append(fields, typ.Elem().Field(i).Name)

		if v, ok := typ.Elem().Field(i).Tag.Lookup("myorm"); ok {
			if strings.Contains(v, "primary") {
				pri = typ.Elem().Field(i).Name
			}
		}
	}
	if pri == "" {
		return ErrNoPrimaryID
	}

	attr := strings.Join(fields, ",")

	s := "SELECT " + attr + " FROM `" + typ.Elem().Name() + "` WHERE " + pri + " = " + idx

	return ev.load(typ, val, s, fields)
}

func (ev env) FindOne(in interface{}, and map[string]interface{}) error {

	typ := reflect.TypeOf(in)
	val := reflect.ValueOf(in)

	if val.Kind() != reflect.Ptr {
		return ErrRequiredPTR
	}
	if len(and) == 0 {
		return ErrEmptyMap
	}

	var where []string
	for k, v := range and {
		where = append(where, k+" = '"+cast(v)+"'")
	}

	var fields []string
	for i := 0; i < reflect.Indirect(val).NumField(); i++ {
		fields = append(fields, typ.Elem().Field(i).Name)
	}

	attr := strings.Join(fields, ",")

	s := "SELECT " + attr + " FROM `" + typ.Elem().Name() + "` WHERE " + strings.Join(where, " AND ")

	return ev.load(typ, val, s, fields)
}

func (ev env) Find(in interface{}, attr []string) Where {
	conn, _ := ev.db.GetConn()
	defer ev.db.PutConn(conn)

	w := condition{db: ev.db}
	w.rType = reflect.TypeOf(in)
	w.rVal = reflect.ValueOf(in)

	w.attr = attr

	return w
}

type condition struct {
	attr   []string
	fields []string
	rType  reflect.Type
	rVal   reflect.Value
	db     *mysqldriver.DB
}

type Where interface {
	ByEqAnd(map[string]interface{}) ([]interface{}, error)
	ByEqOr(map[string]interface{}) ([]interface{}, error)
	ByWhere(string) ([]interface{}, error)
	All() ([]interface{}, error)
	load(rows *mysqldriver.Rows) ([]interface{}, error)
}

var _ Where = (*condition)(nil)

func (w condition) All() ([]interface{}, error) {

	if w.rVal.Kind() == reflect.Ptr {
		return nil, ErrInvalidPTR
	}

	conn, e := w.db.GetConn()
	if e != nil {
		return nil, e
	}
	defer w.db.PutConn(conn)

	attr := w.list().attributes()

	s := "SELECT " + attr + " FROM `" + w.rType.Name() + "`"

	row, e := conn.Query(s)
	if e != nil {
		return nil, e
	}

	return w.load(row)
}

func (w condition) ByWhere(where string) ([]interface{}, error) {

	if w.rVal.Kind() == reflect.Ptr {
		return nil, ErrInvalidPTR
	}

	conn, e := w.db.GetConn()
	if e != nil {
		return nil, e
	}
	defer w.db.PutConn(conn)

	if len(where) == 0 {
		return nil, ErrEmptyWhere
	}

	attr := w.list().attributes()

	s := "SELECT " + attr + " FROM `" + w.rType.Name() + "` WHERE " + where

	row, e := conn.Query(s)
	if e != nil {
		return nil, e
	}

	return w.load(row)
}

func (w condition) ByEqAnd(and map[string]interface{}) ([]interface{}, error) {

	if w.rVal.Kind() == reflect.Ptr {
		return nil, ErrInvalidPTR
	}

	conn, e := w.db.GetConn()
	if e != nil {
		return nil, e
	}
	defer w.db.PutConn(conn)

	if len(and) == 0 {
		return nil, ErrEmptyMap
	}

	attr := w.list().attributes()

	s := "SELECT " + attr + " FROM `" + w.rType.Name() + "` WHERE " + joinMap(and, "AND")

	row, e := conn.Query(s)
	if e != nil {
		return nil, e
	}

	return w.load(row)
}

func (w condition) ByEqOr(or map[string]interface{}) ([]interface{}, error) {

	if w.rVal.Kind() == reflect.Ptr {
		return nil, ErrInvalidPTR
	}

	conn, e := w.db.GetConn()
	if e != nil {
		return nil, e
	}
	defer w.db.PutConn(conn)

	if len(or) == 0 {
		return nil, ErrEmptyMap
	}

	attr := w.list().attributes()

	s := "SELECT " + attr + " FROM `" + w.rType.Name() + "` WHERE " + joinMap(or, "OR")

	row, e := conn.Query(s)
	if e != nil {
		return nil, e
	}

	return w.load(row)
}

func (w *condition) list() *condition {
	for i := 0; i < w.rType.NumField(); i++ {
		w.fields = append(w.fields, w.rType.Field(i).Name)
	}
	return w
}

func (w *condition) attributes() (a string) {
	if w.attr != nil {
		a = strings.Join(w.attr, ",")
		w.fields = w.attr
	} else {
		a = strings.Join(w.fields, ",")
	}
	return a
}

func (w condition) load(row *mysqldriver.Rows) ([]interface{}, error) {
	var stack []interface{}

	for row.Next() {
		v := reflect.New(w.rType).Elem()
		for i := 0; i < len(w.fields); i++ {

			x := w.rVal.FieldByName(w.fields[i])
			switch x.Type().Name() {
			case "string":
				v.FieldByName(w.fields[i]).SetString(row.String())
			case "int":
				fallthrough
			case "int8":
				fallthrough
			case "int16":
				fallthrough
			case "int32":
				fallthrough
			case "int64":
				v.FieldByName(w.fields[i]).SetInt(row.Int64())
			case "float32":
				fallthrough
			case "float64":
				v.FieldByName(w.fields[i]).SetFloat(row.Float64())
			case "bool":
				v.FieldByName(w.fields[i]).SetBool(row.Bool())
			default:
				return stack, ErrInvalidTYPE
			}
		}
		stack = append(stack, v)
	}

	return stack, nil
}

func (ev env) load(typ reflect.Type, val reflect.Value, s string, f []string) error {

	conn, e := ev.db.GetConn()
	if e != nil {
		return e
	}
	defer ev.db.PutConn(conn)

	row, e := conn.Query(s)
	if e != nil {
		return e
	}

	for row.Next() {
		for i := 0; i < len(f); i++ {

			t := typ.Elem().Field(i)
			v := val.Elem().Field(i)

			if f[i] == t.Name {

				switch t.Type.Name() {
				case "string":
					v.SetString(row.String())
				case "int":
					fallthrough
				case "int8":
					fallthrough
				case "int16":
					fallthrough
				case "int32":
					fallthrough
				case "int64":
					v.SetInt(row.Int64())
				case "float32":
					fallthrough
				case "float64":
					v.SetFloat(row.Float64())
				case "bool":
					v.SetBool(row.Bool())
				default:
					return ErrInvalidTYPE
				}
			}
		}
	}

	return nil
}