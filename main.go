package main

import (
	"fmt"
	"encoding/json"
	"flag"
	"net/http"
	"log"
	"strconv"
	//"io"
	"math/rand"
	"io/ioutil"
	
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/bdwilliams/go-jsonify/jsonify"
	_ "github.com/go-sql-driver/mysql"
    "database/sql"
)

//database instance
var appdatabase *sql.DB

func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	HandleError(err)
	return c
}

//main	
func main() {
	var err error
	fmt.Println("Hello Ashish")
	flag.Parse()
	
	appdatabase, err = sql.Open("mysql", "root:ashish@/gocode")
       HandleError(err)
       err = appdatabase.Ping()
       if err != nil {
              fmt.Println(err.Error())
       }
	
	router := mux.NewRouter()
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/post", CreatePost).Methods("POST")
	router.HandleFunc("/get", GetData).Methods("GET")

	//defining port
	defer appdatabase.Close()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Give us some seed data
//func init() {
//	CreatePerson(Person{
//		First: "Ashish",
//		Last: "Tiwari",
//		Text: "Hello everyone! This is Ashish.",
//	})
//	
//	CreatePerson(Person{
//		First: "Ayush",
//		Last: "Goyal",
//		Text: "Hello everyone! This is Ayush.",
//	})
//}



//Person Information
type Person struct {
	Id	   int	  `json:"id,omitempty"`
	First  string `json:"first,omitempty"`
	Last   string `json:"last,omitempty"`
	Text   string `json:"text,omitempty"`
	
}

//Index 
func Index(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<h1 style=\"font-family: Helvetica;\">Hello, welcome to my assignment</h1>")
}


//CreatePost
func CreatePost(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
        panic(err)
		}
		log.Println(string(body))
		var person Person
		err = json.Unmarshal(body, &person)
		if err != nil {
			panic(err)
		}
		
		// Save JSON to Post struct
		//if err := json.Unmarshal(body, &person); err != nil {
		//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    //    w.WriteHeader(422)
	    //    if err := json.NewEncoder(w).Encode(err); err != nil {
	    //            panic(err)
	    //    }
		//}
		
		CreatePerson(person)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "<h3 style=\"font-family: Helvetica;\">Successful POST</h3>")
		
}

//CeatePerson
func CreatePerson(p Person) {
	
	//currentPersonId += 1
	
	//p.Id = currentPersonId
	
	c := RedisConnect()
	defer c.Close()
	
	b, err := json.Marshal(p)
	HandleError(err)
	
	keyval := rand.Intn(100)
	
	// Save JSON blob to Redis
	reply, err := c.Do("SET", "person:" + strconv.Itoa(keyval), b)
	HandleError(err)
	fmt.Println("GET ", reply)
	
	//publishing
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

		HandleError(err)

		//defer conn.Close()
		
		ch, err := conn.Channel()

		HandleError(err)

		defer ch.Close()
		
		q, err := ch.QueueDeclare(

			"redis-assign", // name
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait (wait time for processing)
			nil, // arguments
		)
		HandleError(err)
		err = ch.Publish(

			"", // exchange
			q.Name, // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte("person:"+strconv.Itoa(keyval)),
		})

		//log.Printf(" [x] Sent %s", body)
		HandleError(err)
}

func GetData(w http.ResponseWriter, r *http.Request) {

      rows, err := appdatabase.Query("SELECT * FROM goinfo")
	  if err!= nil {
	  panic(err.Error())
	  }
	  
	  defer rows.Close()
	  
	  json.NewEncoder(w).Encode(jsonify.Jsonify(rows))
}