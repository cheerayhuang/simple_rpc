package main

import "log"
import "reflect"
import "sync"


type Network struct {
    mu sync.Mutex

    servers map[any]*Server
    clients map[any]*Client
    conns map[*Client]*Server

    reqCh chan ReqMsg
    done  chan struct{}
}


type Server struct {
    mu sync.Mutex

    name any

    services map[any]*Service

}

type Service struct {

}

type Client struct {
    name any

    reqCh chan ReqMsg
    done chan struct {}
}

type ReqMsg struct {
    clinetName any

    serviceMethod string
    argsType reflect.Type
    args []byte

    replyChan chan ReplyMsg
}

type ReplyMsg struct {
    ok bool
    reply []byte
}

func main() {
    log.Println("hello world!")
}
