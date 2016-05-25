package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sg3des/rattle"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/rattle/", http.StripPrefix("/rattle/", http.FileServer(http.Dir("../"))))

	//set debug mode
	rattle.Debug = true

	//bind controllers and get handler
	wshandle := rattle.SetControllers(&Main{})
	http.Handle("/ws", wshandle)

	println("web server listen on 127.0.0.1:8080")
	http.ListenAndServe("127.0.0.1:8080", nil)
}

//Main controller, into fields parse JSON requests
//!in real project controller will be located in another package
type Main struct {
	Text string
}

//Index is method of controller Main on request takes incoming message and possible return answer message
func (c *Main) Index(r *rattle.Message) *rattle.Message {
	//return answer - insert data to field with id description
	return r.NewMessage("=#description", []byte(`Rattle is tiny websocket double-sided RPC framework, designed for create dynamic web applications`))
}

//JSON method
func (c *Main) JSON(r *rattle.Message) *rattle.Message {
	data, err := json.Marshal(c)
	if err != nil {
		return r.NewMessage("+#errors", []byte("failed parse JSON request, error: "+err.Error()))
	}
	//call "test.RecieveJSON frontend function and send to it JSON data"
	return r.NewMessage("test.RecieveJSON", data)
}

//RAW method
func (c *Main) RAW(r *rattle.Message) *rattle.Message {
	//call "test.RecieveRAW frontend function and send to raw data"
	return r.NewMessage("test.RecieveRAW", []byte(c.Text))
}

//Timer is example of periodic send data, note the that function does not return anything
func (c *Main) Timer(r *rattle.Message) {
	for {
		t := time.Now().Local().Format("2006.01.02 15:04:05")
		if err := r.NewMessage("=#timer", []byte(t)).Send(); err != nil {
			//if err then connection is closed
			return
		}

		time.Sleep(time.Second)
	}
}
