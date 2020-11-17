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
	Create(interface{}) error

	InsertOne(interface{}) error
	InsertAll(interface{}) error

	FindByID(interface{}, interface{}) error
	FindOne(interface{}, map[string]interface{}) error
	Find(interface{}, []string) Where
	load(reflect.Type, reflect.Value, string, []string) error

	UpdateByID(interface{}) (bool, error)
	UpdateAnd(interface{}, map[string]interface{}) (bool, error)
	UpdateOr(interface{}, map[string]interface{}) (bool, error)
}

var _ Handler = (*env) (nil)

var (
	ErrInvalidID   = errors.New("myorm: invalid or nil id")
	ErrNoPrimaryID = errors.New("myorm: no primary key found")
	ErrInvalidTYPE = errors.New("myorm: invalid struct type")
	ErrRequiredPTR = errors.New("myorm: required pointer to struct in FindByID")
	ErrInvalidPTR  = errors.New("myorm: required struct, pointer found in Find")
	ErrEmptyMap    = errors.New("myorm: empty map is send in Find")
	ErrEmptyWhere  = errors.New("myorm: empty where string is send in Find")
)

func Dial(dsn string, pool int) Handler {
	return env{
		db: mysqldriver.NewDB(dsn, pool, -1),
	}
}

func (ev env) Create(in interface{}) error {

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
		ddl = "ALTER TABLE `" + t.Name() + "` ADD INDEX (" + strings.Join(idx, ",") + ");"
		fmt.Println(ddl)
		_, e = conn.Exec(ddl)
	}

	return e
}

func joinMap(m map[string]interface{}, k string) string {
	var where []string
	for k, v := range m {
		where = append(where, k+" = '"+cast(v)+"'")
	}
	return strings.Join(where, " "+k+" ")
}

func cast(id interface{}) string {
	var idx string
	switch id := id.(type) {
	case int:
		idx = strconv.Itoa(id)
	case string:
		idx = id
	}
	return idx
}