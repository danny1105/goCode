package main

import (
	"fmt"
	"encoding/json"
	"flag"
	"net/http"
	"log"
	"strconv"
	//"io"
	"io/ioutil"
	
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

var currentPersonId int

func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	HandleError(err)
	return c
}

//main	
func main() {
	fmt.Println("Hello Ashish")
	flag.Parse()
	
	router := mux.NewRouter()
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/post", CreatePost).Methods("POST")
	

	//defining port
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
	
	currentPersonId += 1
	
	p.Id = currentPersonId
	
	c := RedisConnect()
	defer c.Close()
	
	b, err := json.Marshal(p)
	HandleError(err)
	// Save JSON blob to Redis
	reply, err := c.Do("SET", "post:" + strconv.Itoa(p.Id), b)
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
			Body: []byte(strconv.Itoa(p.Id)),
		})

		//log.Printf(" [x] Sent %s", body)
		HandleError(err)
}