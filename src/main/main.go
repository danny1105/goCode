package main

import (
	"encoding/json"
	"fmt"
	"flag"
	//"strconv"
	"log"
	"runtime"
	"time"
	"net/http"
	"math/rand"

	//"github.com/andreagrandi/go-amqp-example/contracts"
	"github.com/streadway/amqp"	
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

var per = Person{
		ID: "2",
		Firstname: "Ashish",
		Lastname:  "Tiwari",
		Age: 24,
	}	
// var people []Person

var (
	amqpURI = flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func init() {
	flag.Parse()
	initAmqp()
}

var Conn *amqp.Connection
var ch *amqp.Channel

func initAmqp() {
	var err error

	Conn, err = amqp.Dial(*amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err = Conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"test-exchange", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // noWait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare the Exchange")
}


//CreatePersonEndpoint for entering Person Info
func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Data added to Redis")
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
	//server.Start()
	//queue := redismq.CreateQueue("localhost", "6379", "", 9, "example")

	//defining router
	router := mux.NewRouter()
	//people = append(people, Person{ID: "1", Firstname: "Ashish", Lastname: "Tiwari", Age: 24})
	//defining Endpoints
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	log.Println("Starting publisher...")

	// Publish messages
	publishMessages(10000)

	// Close Channel
	defer ch.Close()

	// Close Connection
	defer Conn.Close()

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
		//go write(queue)
		return nil
}

func publishMessages(messages int) {
	for i := 0; i < messages; i++ {
		//t := strconv.Itoa(keyval)
		
		payload, err := json.Marshal(per)
		failOnError(err, "Failed to marshal JSON")
		
		//fmt.Println(t)
		err = ch.Publish(
			"go-test-exchange", // exchange
			"go-test-key",      // routing key
			false,              // mandatory
			false,              // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Transient,
				ContentType:  "application/json",
				Body:         payload,
				Timestamp:    time.Now(),
			})

		failOnError(err, "Failed to Publish on RabbitMQ")
	}
}	



