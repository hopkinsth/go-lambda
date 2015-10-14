package main

import "github.com/hopkinsth/lambda-phage"

func main() {
	ch := lambda.Listen()
	for r := range ch {
		r.Succeed("yay")
	}
}
