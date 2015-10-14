package main

import "github.com/hopkinsth/lambda-phage"
import "fmt"

func main() {
	ch := lambda.Listen()
	for r := range ch {
		fmt.Println(r.Context.AwsRequestID)
		r.ResponseData <- nil
	}
}
