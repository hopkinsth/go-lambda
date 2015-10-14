package lambda

import "os"
import "io"
import "encoding/json"
import "fmt"

var gStream = openStreams()

type streams struct {
	in    io.Reader
	out   io.Writer
	inch  chan Request
	outch chan Response
	chs   chan (chan<- Request)
}

func openStreams() *streams {
	s := &streams{
		in:    os.Stdin,
		out:   os.Stdout,
		inch:  make(chan Request, 10),
		outch: make(chan Response, 10),
		chs:   make(chan (chan<- Request), 10),
	}

	return s
}

func (s *streams) begin() {
	go s.listener()
	go s.responder()
	go s.repeater()
}

func (s *streams) add(ch chan<- Request) {
	s.chs <- ch
}

func (s *streams) repeater() {
	for req := range s.inch {
		r := req
		tot := len(s.chs)
		//all := make([]chan<- Request, tot)

		for i := 0; i < tot; i += 1 {
			ch := <-s.chs
			ch <- r
			s.chs <- ch
		}
	}
}

func (s *streams) listener() {
	d := json.NewDecoder(s.in)

	for d.More() {
		var req Request
		err := d.Decode(&req)
		if err != nil {
			fmt.Println(err) //some invalid json
		}

		req.ResponseData = make(chan interface{})
		req.response = Response{
			RequestId: req.Context.AwsRequestID,
		}

		s.inch <- req
	}
}

func (s *streams) responder() {
	for res := range s.outch {
		o, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Fprintf(s.out, "%s", o)
		}
	}
}
