# MINIMESS

The goal of this project is for me to learn making an async messaging platform using go concurrency, websockets, and kafka. This project should prepare components that can run on a single laptop. Scaling is not the top priority here.

## Day 1 achievement

### Server related

 we first setup a server which listens on port 8080 for ws connections

* the server will send "success"
* the server will send timestamps every few seconds
* close connection when client is inactive

### Client related

* will connect to the server
* will print messages it receives from the server

## Day 2 goals

### Server related

* support channels
* create an API to subscribe to a topic

the client will listen to topic
/subscribe/id:`<uid>`/topic:`<topic>`

clients that are newly subscribed to a topic will only receive new messages

todo:
    have one go routine per connection
    decouple writer, reader, pinger, ponger
    communicate via go channels
    create a shared struct between the separate goroutines
    topics updates will be stored in mem and will be dropped after 5 minutes

### misc goals

* try out domain model by [Kat](https://github.com/katzien/go-structure-examples)
