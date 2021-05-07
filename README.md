# MINIMESS

The goal of this project is for me to learn making an async messaging platform using go concurrency, websockets, and kafka. This project should prepare components that can run on a single laptop. Scaling is not the top priority here.

## Phase 1 achievement

### Server related 1

 we first setup a server which listens on port 8080 for ws connections

* the server will send "success"
* the server will send timestamps every few seconds
* close connection when client is inactive

### Client related

* will connect to the server
* will print messages it receives from the server

## Phase 2 goals

### Server related 2

* support channels
* create an API to subscribe to a topic

the client will listen to topic
/subscribe/id:`<uid>`/topic:`<topic>`

clients that are newly subscribed to a topic will only receive new messages

done:
    use fanout to broadcast
    have one go routine per connection
    inter goroutine communicate via go channels
    create a shared struct between the separate goroutines
    decouple writer, reader, pinger, ponger
    topics mapping will be stored in mem and will be broadcased immediately - the server will not store messages for now

## Client related 2

make an interactive client that looks something like below

```bash
> sub afrozalm
> send zorfa What\'s up Zorfa
> send mala How ya doin\' Mala
> zorfa: I\'m fantastic. Learning some magic tricks.
> mala: Just finished RoW. My mind is blow into millions of pieces. BrandoSando is legend
> unsub afrozalm
```

## Future Goals

* Use thrift
* use kafka
* use nginx

### misc goals

* try out domain model by [Kat](https://github.com/katzien/go-structure-examples)
