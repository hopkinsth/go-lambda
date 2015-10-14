package lambda

import "encoding/json"

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

	ResponseData chan interface{} `json:"-"`
	response     Response
}

type Response struct {
	RequestId string      `json:"requestId"`
	Error     *string     `json:"error"`
	Data      interface{} `json:"data"`
}

func Listen() <-chan Request {
	ch := make(chan Request)
	gStream.add(ch)
	return ch
}

func init() {
	gStream.begin()
}
