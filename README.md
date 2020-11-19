# myorm - Mysql ORM
Simple MySQL ORM which wrap around pubnative mysql driver. This orm support simple CRUD operations and doesn't support Inter-Table complex relations yet.

### Install
`go get -v github.com/reoxey/myorm`

### Import
```go
import "github.com/reoxey/myorm"
```

### Dependency
This package uses `github.com/pubnative/mysqldriver-go` under the hood.

### Examples
Here are a few CRUD operations.

#### Create Table
```go
type User struct {
	Id      int    `myorm:"primary,auto_increment"`
	Name    string `myorm:"notnull,unique"`
	Age     int    `myorm:"notnull,index"`
	Address string `myorm:"notnull,default='IN'"`
	Code    int    `myorm:"index"`
}

func main() {

	db := myorm.Dial("user:****@tcp(127.0.0.1:3306)/orm", 10)

	e := db.Create(User{})

	if e != nil {
		log.Fatal(e)
	}
}
```

#### INSERT
```go
func main(){
    e = db.InsertOne(User{1, "One", 12, "US", 5})
    if e != nil {
        log.Println(e)
    }
    
    e = db.InsertAll(
        []User{
            {Name: "Ok", Age: 20},
            {Name: "KO", Age: 30},
        },
    )
    if e != nil {
        log.Println(e)
    }
}
```

#### SELECT
```go
func main(){
    var u User
    
    e = db.FindByID(&u, 2)
    if e != nil {
        log.Println(e)
    }
    
    fmt.Printf("%+v \n\n", u)
    
    var us User
    
    x, e := db.Find(us, nil).All()
    if e != nil {
        log.Println(e)
    }
    
    for _, v := range x {
        fmt.Printf("%+v \n", v)
    }
    
    var uo User
    e = db.FindOne(&uo, map[string]interface{}{"Age":15})
    if e != nil {
        log.Println(e)
    }
    
    fmt.Printf("\n\n %+v \n", uo)
    
    
    z, e := db.Find(us, nil).ByEqAnd(map[string]interface{}{"Address":"IN"})
    if e != nil {
        log.Println(e)
    }
    
    for _, v := range z {
        fmt.Printf("%+v \n", v)
    }
}
```

#### Update
```go
func main(){
	uu := User{
		Id: 3,
		Name: "BoB",
	}
	x, e := db.UpdateByID(uu)
	
	if e != nil {
		log.Fatal(e)
	} else if !x {
		fmt.Println("\n\n Failed")
	} else {
		fmt.Println("\n\n Super!")
	}
}
```

#### Delete
```go
func main(){
	ud := User{
		Id: 6,
	}
	x, e := db.DeleteByID(ud)

	if e != nil {
		log.Fatal(e)
	} else if !x {
		fmt.Println("\n\n Failed")
	} else {
		fmt.Println("\n\n Super!")
	}
}
```