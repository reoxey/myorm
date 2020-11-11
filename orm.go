package myorm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pubnative/mysqldriver-go"
)

type env struct{
	db *mysqldriver.DB
}

type Handler interface {
	Model(interface{}) error
	InsertOne(interface{}) error
	InsertAll(interface{}) error
	Find(interface{}) Where
}

var _ Handler = (*env) (nil)

var (
	ErrInvalidID = errors.New("myorm: invalid or nil id")
)

func Dial(dsn string, pool int) Handler {
	return env{
		db: mysqldriver.NewDB(dsn, pool, -1),
	}
}

func (ev env) Model(in interface{}) error {

	conn, e := ev.db.GetConn()
	if e != nil {
		return e
	}
	defer ev.db.PutConn(conn)

	t := reflect.TypeOf(in)

	var lot []string
	var idx []string

	for i := 0; i < t.NumField(); i++ {

		var l strings.Builder

		l.WriteString("`"+t.Field(i).Name+"` ")

		switch t.Field(i).Type.Name() {
		case "string": l.WriteString("TEXT ")
		case "int": l.WriteString("INT ")
		case "int8": l.WriteString("TINYINT ")
		case "int16": l.WriteString("SMALLINT ")
		case "int32": l.WriteString("MEDIUMINT ")
		case "int64": l.WriteString("BIGINT ")
		case "float32": l.WriteString("FLOAT ")
		case "float64": l.WriteString("DOUBLE ")
		case "bool": l.WriteString("TINYINT(1) ")
		default:
			l.WriteString("VARCHAR(100) ")
		}

		if x, ok := t.Field(i).Tag.Lookup("myorm"); ok {
			s := strings.Split(x, ",")
			for o := 0; o < len(s); o++ {
				switch s[o] {
				case "primary": l.WriteString("PRIMARY KEY ")
				case "notnull": l.WriteString("NOT NULL ")
				case "index": idx = append(idx, t.Field(i).Name)
				default:
					l.WriteString(strings.Replace(s[o],"=", " ",1)+" ")
				}
			}
		} else {
			l.WriteString("NULL ")
		}

		lot = append(lot, l.String())
	}

	ddl := "CREATE TABLE IF NOT EXISTS `"+t.Name()+"` ("+
		strings.Join(lot, ",")+
		") ENGINE=InnoDB DEFAULT CHARSET=utf8 ;"

	fmt.Println(ddl)

	x, e := conn.Exec(ddl)
	if e != nil {
		return e
	}

	if x.Warnings == 0 && len(idx) > 0 {
		ddl = "ALTER TABLE `"+t.Name()+"` ADD INDEX ("+strings.Join(idx, ",")+");"
		fmt.Println(ddl)
		_, e = conn.Exec(ddl)
	}

	return e
}

func  (ev env) InsertOne(in interface{}) error {
	conn, e := ev.db.GetConn()
	if e != nil {
		return e
	}
	defer ev.db.PutConn(conn)

	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	lot := set(t, v)

	dml := "INSERT INTO `"+t.Name()+"` SET "+strings.Join(lot, ",")

	fmt.Println(dml)

	_, e = conn.Exec(dml)

	return e
}

func set(t reflect.Type, v reflect.Value) []string {

	var lot []string

	for i := 0; i < t.NumField(); i++ {
		val := ""
		switch t.Field(i).Type.Name() {
		case "int": val = strconv.Itoa(int(v.Field(i).Int()))
		case "float32": val = strconv.FormatFloat(v.Field(i).Float(),'f',-1,32)
		case "float64": val = strconv.FormatFloat(v.Field(i).Float(),'f',-1,64)
		case "string": val = v.Field(i).String()
			if val == "" {
				continue
			}
		}
		lot = append(lot, t.Field(i).Name+"='"+val+"'")
	}

	return lot
}

func  (ev env) InsertAll(arr interface{}) error {
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

		dml := "INSERT INTO `"+t.Name()+"` SET "+strings.Join(lot, ",")

		fmt.Println(dml)

		_, e = conn.Exec(dml)
		if e != nil {
			return e
		}
	}

	return nil
}


func  (ev env) Find(in interface{}) Where {
	conn, _ := ev.db.GetConn()
	defer ev.db.PutConn(conn)

	w := condition{db: ev.db}
	w.store = reflect.TypeOf(in)

	return w
}

type condition struct {
	where string
	store reflect.Type
	db *mysqldriver.DB
}

type Where interface {
	ByID(interface{}) error
	//All() error

}

var _ Where = (*condition) (nil)

func  (w condition) ByID(id interface{}) error {
	conn, e := w.db.GetConn()
	if e != nil {
		return e
	}
	defer w.db.PutConn(conn)

	var idx string
	switch id := id.(type) {
	case int: idx = strconv.Itoa(id)
	case string: idx = id
	default:
		return ErrInvalidID
	}

	if idx == "" {
		return ErrInvalidID
	}



	return e
}