package main

import (
	"fmt"
	"revision/cli-yaml/types"
	"revision/cli-yaml/utils"

	"github.com/golang/glog"
)

func main() {
	file := utils.ParseFlag("file", "The yaml definition for the operation")
	f, err := utils.ReadFile(*file)
	if err != nil {
		glog.Fatalf("Error reading YAML file definition: %v", err)
	}

	var definition types.Definition
	utils.GetFileData(f, &definition)

	url := utils.ConstructUrl(definition)

	fmt.Printf("Making http call to %s\n", url)

	res, err := utils.MakeHttpCall(definition.Details.Method, url, definition.Details.Headers)
	if err != nil {
		glog.Fatalf("Error making http call to URL: %s", url)
	}
	body, err := utils.GetBodyFromHTTPResponse(res)
	if err != nil {
		glog.Fatal("Error getting body from HTTP request")
	}

	if err := utils.HandleResponse(body, definition); err != nil {
		glog.Fatalf("Error occured handling the response: %v", err)
	}

}
