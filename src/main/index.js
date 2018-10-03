const express = require('express');
// const exphbs = require('express-handlebars');
// const path = require('path');
// const bodyParser = require('body-parser');
// const methodOverride = require('method-override');
const redis = require('redis');
const Step = require('step');

//set port
const port = 3000;

//init app
const app = express();

// config file
var config = require('config');

var amqp = require('amqplib/callback_api');

amqp.connect('amqp://localhost', function(err, conn) {
  conn.createChannel(function(err, ch) {
    var ex = 'logs';

    ch.assertExchange(ex, 'fanout', {durable: false});

    ch.assertQueue('', {exclusive: true}, function(err, q) {
      console.log(" [*] Waiting for messages in %s. To exit press CTRL+C", q.queue);
      ch.bindQueue(q.queue, ex, '');

      ch.consume(q.queue, function(msg) {
        console.log(" [x] %s", msg.content.toString());
      }, {noAck: true});
    });
  });
});


// mysql client connect
var mysql      = require('mysql');
var connection = mysql.createConnection({
	host        	  : 'localhost',
	port        	  : '',
	database    	  : '',
	user        	  : '',
	password    	  : '',
	insecureAuth	  : true,
	multipleStatements: true
});

function handleDisconnect() {
  connection.on('error', function(err) {  
    if (!err.fatal) {
      return;
    }

    if (err.code !== 'PROTOCOL_CONNECTION_LOST') {
      throw err;
    }

    connection = mysql.createConnection(connection.config);
    handleDisconnect();
    connection.connect();
  });
}
handleDisconnect();

// redis client connect
client = redis.createClient();
client.on("error", function (err) {
    console.log("Error " + err);
});

var timeOut = 10;
function dumpMysql() {

	Step(
		
		function checkWorkingVars() {
			client.multi()
				.exists('data_tmp')
				.exists('data')
				.exec(this);
		},
		
		function setWorkingVars(err, replies) {
			existsTmp = replies[0];
			existsData = replies[1];
			
			if(!existsData && !existsTmp) {
				setTimeout(dumpMysql, timeOut*1000);
			} else {
				if(existsData && !existsTmp) {
					client.rename('data', 'data_tmp');
				}
				client.smembers('data_tmp', this);
			}
		},
		
		function readData(err, data) {
			for(var attr in data) {
				connection.query('SELECT 1'); // db ping
				
				var query  = 'INSERT IGNORE INTO data (log) ';
					query += 'VALUES ("'+data[attr]+'")';
				connection.query(query);				
			}
			
				
		},
				
		function cleanUp(err, r) {
			client.del('data_tmp', this);
		},
		
		function restart(err, r) {
			setTimeout(dumpMysql, timeOut*1000);
		}
		
	);

};
setTimeout(dumpMysql, 1000);

// stdout log started
console.log('DUMP Mysql running...');

app.get('/person/{id}', function(req, res, next){
    res.render('query');
});

// app.delete('/peer/delete/:id', function(req, res, next){
//     client.del(req.params.id);
//     res.redirect('/');
// })

app.listen(port, function(){
    console.log('server started on port '+port)
});