//const express = require('express');
const redis = require('redis');
const mysql = require('mysql');
const amqplib = require('amqplib/callback_api');

//const app = express();
const client = redis.createClient();
//var msgtosql;

amqplib.connect('amqp://localhost',function(err, conn) {
        if(err){
            throw err;
        }
        conn.createChannel(function(err, ch) {
            if(err){
                throw err;
            }
            var q = 'redis-assign';
            ch.assertQueue(q, {durable: false});
            ch.consume(q, function(msg) {
                console.log('message recieved');
                msgtosql = msg.content.toString();
                mysql_insert(msgtosql);
            }, {noAck: true});
        });
});



//create mysql db connection
var db = mysql.createConnection({
    host     : 'localhost',
    port     : '',
    user     : 'root',
    password : 'ashish',
    database : 'gocode'
});

db.connect((err) => {
    if(err){
        throw err;
    }
    console.log('Mysql connected...');
});

//app.get('/createdb', (req, res) => {
    // let sql = 'CREATE TABLE goinfo(id int AUTO_INCREMENT, first VARCHAR(50), last VARCHAR(50), text VARCHAR(255), PRIMARY KEY (id))';
    // db.query(sql, (err, result) => {
    //          if(err){
    //             throw err;
    //         }
    //         console.log(result);
function mysql_insert(data) {
    client.hgetall(data, (err, reply) => {
        if(err){
            throw err;
        }
        db.query('INSERT INTO gocode set ?', reply, (err, result) => {
            if(err){
             throw err;
           }
        console.log('MySQL table updated');
        client.del(data, (err, reply) => {
            if(err){
                throw err;
            }
            console.log('Data deleted');
        });
        res.send('Everything is OK!!!');
    });
    //  });
  });
}

//app.listen('8080', () => {
//    console.log('Server started on port 8080');
//})