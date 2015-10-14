package lambda

import "encoding/json"

var cfg *Config
var ipc *streams

type LambdaContext struct {
	AwsRequestID             string `json:"awsRequestId"`
	FunctionName             string `json:"functionName"`
	FunctionVersion          string `json:"functionVersion"`
	InvokeID                 string `json:"invokeid"`
	IsDefaultFunctionVersion bool   `json:"isDefaultFunctionVersion"`
	LogGroupName             string `json:"logGroupName"`
	LogStreamName            string `json:"logStreamName"`
	MemoryLimitInMB          string `json:"memoryLimitInMB"`
}

type Request struct {
	// custom event fields
	Event json.RawMessage `json:"event"`
	// default context object
	Context *LambdaContext `json:"context"`

	response Response `json:"-"`
	ipc      *streams `json:"-"`
}

func (r *Request) Succeed(data interface{}) {
	r.response.Data = data
	r.ipc.outch <- &r.response
}

func (r *Request) Error(err error) {
	r.response.Error = err.Error()
	r.ipc.outch <- &r.response
}

type Response struct {
	RequestId string      `json:"requestId"`
	Error     string      `json:"error"`
	Data      interface{} `json:"data"`
}

type Config struct {
}

//
func Setup(icfg *Config) {
	if cfg == nil {
		if icfg == nil {
			icfg = &Config{}
		}

		cfg = icfg
	}

	if ipc == nil {
		ipc = openStreams(cfg)
		ipc.begin()
	}
}

func Listen() <-chan *Request {
	Setup(nil)
	return ipc.inch
}
