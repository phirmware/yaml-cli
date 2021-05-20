package main

import (
	"fmt"
	"revision/cli-yaml/types"
	"revision/cli-yaml/utils"
)

func main() {
	file := utils.ParseFlag("file", "The yaml definition for the operation")
	fmt.Println(*file)
	f, err := utils.ReadFile(*file)
	if err != nil {
		panic(err)
	}
	var definition types.Definition
	utils.GetFileData(f, &definition)

	url := utils.ConstructUrl(definition)

	fmt.Printf("Making http call to %s\n", url)

	res, err := utils.MakeHttpCall(definition.Details.Method, url, definition.Details.Headers)
	if err != nil {
		panic(err)
	}
	body, err := utils.GetBodyFromHTTPResponse(res)
	if err != nil {
		panic(err)
	}

	if err := utils.HandleResponse(body, definition); err != nil {
		panic(err)
	}

}
