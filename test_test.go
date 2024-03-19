package main

import "testing"


func TestBasic(t *testing.T) {
    n, err := MakeNetwork()
    if err != nil {
        t.Errorf("create network failed: %s", err)
    }

    client, err := MakeClient("client-1")
    if err != nil {
        t.Errorf("create client failed: %s", err)
    }

    ser, err := MakeServer("server-99")
    if err != nil {
        t.Errorf("create server failed: %s", err)
    }

    n.AddServer(ser)
    n.AddClient(client)
    n.Connect("client-1", "server-99")

    {
        reply := ""
        echoMsg := "hi, this is a msg."
        client.Call("service.echo", "hi, this is a msg.", &reply)
        if reply != echoMsg {
            t.Errorf("rpc service.echo failed.")
        }
    }

    {
        reply := 0
        echoMsg := "hi, this is a msg."
        client.Call("service.add", 5, 10, &reply)
        if reply != 15 {
            t.Errorf("rpc service.add failed.")
        }
    }
}
