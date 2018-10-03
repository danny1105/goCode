package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"log"
	"runtime"
	"net/http"
	"math/rand"

	"github.com/adjust/redismq"	
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

// newPool returns a pointer to a redis.Pool
	var pool = newPool()
// get a connection from the pool (redis.Conn)
	var conn = pool.Get()
// use defer to close the connection when the function completes

	var keyval = rand.Intn(100)
	
//Person Information
type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Age       int    `json:"age,omitempty"`
}

// var people []Person

//CreatePersonEndpoint for entering Person Info
func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	//var person Person
	//_ = json.NewDecoder(req.Body).Decode(&person)
	//person.ID = params["id"]
	//people = append(people, person)
	//json.NewEncoder(w).Encode(people)
	//const objectPrefix string = "person:"
	}
	

func main() {

	defer conn.Close()
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })
	// call Redis PING command to test connectivity
	err := ping(conn)
	if err != nil {
		fmt.Println(err)
	}
	err = setStruct(conn)
	if err != nil {
		fmt.Println(err)
	}
	
	// err = setStruct(conn)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	
	runtime.GOMAXPROCS(5)
	//server := redismq.NewServer("localhost", "6379", "", 9, "9999")
	server.Start()
	//queue := redismq.CreateQueue("localhost", "6379", "", 9, "example")

	//defining router
	router := mux.NewRouter()
	//people = append(people, Person{ID: "1", Firstname: "Ashish", Lastname: "Tiwari", Age: 24})
	//people = append(people, Person{ID: "2", Firstname: "Ayush", Lastname: "Goyal", Age: 24})
	//defining Endpoints
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")

	//defining port
	log.Fatal(http.ListenAndServe(":3000", router))
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// ping tests connectivity for redis (PONG should be returned)
func ping(c redis.Conn) error {
	// Send PING command to Redis
	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

func setStruct(c redis.Conn) error {

	const objectPrefix string = "person:"
	//keyval := rand.Intn(100)

	per := Person{
		ID: "2",
		Firstname: "Ashish",
		Lastname:  "Tiwari",
		Age: 24,
	}
	
	fmt.Println("In Redis SET")

	// serialize User object to JSON
	json, err := json.Marshal(per)
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", keyval, json)
	if err != nil {
		return err
		}
		//convert random string to key to put in the queue
		go write(queue)
		return nil
}

func write(queue *redismq.Queue) {
	t := strconv.Itoa(keyval)
	fmt.Println(t)
	for {
		queue.Put(t)
	}
}


