package lambda

import "os"
import "io"
import "encoding/json"
import "fmt"

type getReq chan<- Request

func (g getReq) deliver() {
}

type streams struct {
	in    io.Reader
	out   io.Writer
	inch  chan *Request
	outch chan *Response
	cfg   *Config
}

func openStreams(cfg *Config) *streams {
	s := &streams{
		in:    os.Stdin,
		out:   os.Stdout,
		inch:  make(chan *Request, 10),
		outch: make(chan *Response, 10),
		cfg:   cfg,
	}

	return s
}

func (s *streams) begin() {
	go s.listener()
	go s.responder()
}

func (s *streams) listener() {
	d := json.NewDecoder(s.in)

	for d.More() {
		var req Request
		err := d.Decode(&req)
		if err != nil {
			fmt.Println(err) //some invalid json
		}

		req.ipc = s
		req.response = Response{
			RequestId: req.Context.AwsRequestID,
		}

		s.inch <- &req
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
