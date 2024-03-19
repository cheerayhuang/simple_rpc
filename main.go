package main

import "log"
import "reflect"
import "sync"

import "github.com/pkg/errors"


type Network struct {
    mu sync.Mutex

    servers map[any]*Server
    clients map[any]*Client
    conns map[*Client]*Server

    reqCh chan ReqMsg
    done  chan struct{}
}

func MakeNetwork() (*Network, error) {
    n := new(Network)

    n.servers = make(map[any]*Server)
    n.clients = make(map[any]*Client)
    n.conns = make(map[*Client]*Server)

    n.reqCh = make(chan ReqMsg)
    n.done = make(chan struct{})

    return n, nil
}

func (n *Network) Connect(cliID any, serID any) error {
    cli := n.QueryClientByID(cliID)
    ser := n.QueryServerByID(serID)

    if !cli {
        return errors.Errorf("cliID[%v] doesn't exist.", cliID)
    }

    if !ser {
        return errors.Errorf("serID[%v] doen't exist.", serID)
    }

    n.mu.Lock()
    defer n.mu.Unlock()

    n.conns[cli] = ser

    return nil
}

func (n *Network) QueryClientByID(id any) *Client {
    n.mu.RLock()
    defer n.mu.RUnlock()

    cli, ok := n.clients[id]
    if ok {
        return cli
    }

    return nil
}

func (n *Network) QueryServerByID(id any) *Server {
    n.mu.RLock()
    defer n.mu.RUnlock()

    ser, ok := n.servers[id]
    if ok {
        return ser
    }

    return nil
}

func (n *Network) QueryPeerServerByClient

func (n *Network) AddServer(ser *Server) error {
    n.mu.Lock()
    defer n.mu.Unlock()

    if _, ok := n.servers[ser.id]; ok {
        return errors.New("Add Server failed: new Server ID is duplicated.")
    }

    n.servers[ser.id] = ser

    return nil
}

func (n *Network) DelServer(serID any) error {
    n.mu.Lock()
    defer n.mu.Unlock()

    delete(n.servers, serID)

    return nil
}

func (n *Network) AddClient(cli *Client) error {
    n.mu.Lock()
    defer n.mu.Unlock()

    if _, ok := n.clients[cli.id]; ok {
        return errors.New("Add Client failed: new Client ID is duplicated.")
    }

    n.clients[cli.id] = cli

    return nil
}

func (n *Network) DelClient(cliID any) error {
    n.mu.Lock()
    defer n.mu.Unlock()

    delete(n.clients, cliID)

    return nil
}

type Server struct {
    mu sync.Mutex

    id any
    services map[any]*Service
}

func MakeServer(id any) (*Server, error) {
    s := Server{
        id: id,
        services: make(map[any]*Service),
    }

    return &s, nil
}

type Service struct {
    name string
    svcValue reflect.Value
    svcType reflect.Type
    methods map[string]reflect.Method
}

type Client struct {
    id any

    reqCh chan ReqMsg
    done chan struct {}
}

func MakeClient(id any) (*Client, any) {
    c := Client{
        id: id,
        reqCh: make(chan ReqMsg),
        done: make(chan struct{}),
    }

    return &c, nil
}

type ReqMsg struct {
    clinetID any

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
