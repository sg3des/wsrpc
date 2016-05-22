package rattle

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

var (
	conn *websocket.Conn
	err  error
	addr = "127.0.0.1:8088"
)

// ****
// fake controllers for tests

type FakeController struct {
	Name string
}

func (c *FakeController) FakeMethod(r *Message) *Message {
	fmt.Println("recieve message:", c)
	r.NewMessage("tovoid0").Send()
	return r.NewMessage("tovoid1")
}

func (c *FakeController) FakeEmptyMethod(r *Message) {}

// TESTS

func init() {
	Debug = true
	log.SetFlags(log.Lshortfile)

	wshandle := SetControllers(
		&FakeController{},
	)
	http.Handle("/ws", wshandle)

	go func() {
		err = http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
	}()

	//so the server had go up
	time.Sleep(300 * time.Millisecond)
}

func TestSetControllers(t *testing.T) {
	// controllers already set in init, just check the correctness of this

	if len(Controllers) != 1 {
		t.Fatal("failed set controllers, length of Controllers map is incorrect")
	}

	if conInterface, ok := Controllers["FakeController"]; ok {
		controller := reflect.ValueOf(conInterface)
		if !controller.IsValid() {
			t.Error("failed set controllers, incorrect reflect of controller interface")
		}
		if !controller.MethodByName("FakeEmptyMethod").IsValid() {
			t.Error("failed set controllers, required method not found")
		}
		if !controller.MethodByName("FakeMethod").IsValid() {
			t.Error("failed set controllers, required method not found")
		}
	} else {
		t.Error("failed set controllers, incorrect determine name of controller")
	}
}

func TestRequest(t *testing.T) {
	conn, err = websocket.Dial("ws://"+addr+"/ws", "", "http://"+addr)
	if err != nil {
		t.Error(err)
	}

	go fakeReciever()

	msg := &Message{From: []byte("test.From"), To: []byte("FakeController.FakeMethod"), Data: []byte(`{"Name":"testname"}`)}
	if _, err := conn.Write(msg.Bytes()); err != nil {
		t.Error(err)
	}
}

func fakeReciever() {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		bmsg := scanner.Bytes()
		fmt.Println("msg for frontend:", string(bmsg))
	}
}

func TestParsemsg(t *testing.T) {
	incorrectmsgs := []string{"\n", "fakeController.fakeMethod {}\n", "{}\n", "test.From test.toMethod\n", "test.From fakeController.fakeMethod"}

	for _, smsg := range incorrectmsgs {
		_, err := Parsemsg([]byte(smsg))
		if err == nil {
			t.Error("failed parse msg: '" + smsg + "' must be error")
		}
	}

	correctmsg := []byte("test.From FakeController.FakeMethod {\"name\":\"value\"}\n")

	msg, err := Parsemsg(correctmsg)
	if err != nil {
		t.Error(err)
	}

	rpcmethod, err := splitRPC(msg.To)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(rpcmethod.Join(), msg.To) {
		t.Error("failed inverse transformation controller and method")
	}

	if !bytes.Equal(msg.Bytes(), correctmsg) {
		t.Error("failed convert msg to bytes")
	}

	//disable debug, because otherwise there will be a warning of failed determine caller function name - it`s ok, as this not use controller-method architecture
	Debug = false
	newmsg := msg.NewMessage("test.To")
	if !bytes.Equal(newmsg.To, []byte(`test.To`)) {
		t.Error("failed create new message field To fill incorrect")
	}
	Debug = true

	if !bytes.Equal(newmsg.Data, []byte(`{}`)) {
		t.Error("failed create new message field Data fill incorrect")
	}
}
