package utils

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"revision/cli-yaml/types"
	"runtime"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

func ParseFlag(name string, description string) *string {
	result := flag.String(name, "", description)
	flag.Parse()
	return result
}

func ReadFile(file string) (*os.File, error) {
	result, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetFileData(file *os.File, target interface{}) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, target)
}

func constructHeader(headers types.Headers) string {
	return headers.Key + " " + headers.Value
}

func ConstructUrl(details types.Definition) string {
	endpoint := details.Details.Endpoint
	path := details.Details.Path

	return endpoint + path
}

func appendHeaders(req *http.Request, headers []types.Headers) {
	for _, header := range headers {
		h := constructHeader(header)
		req.Header.Add(header.Header, h)
	}
}

func MakeHttpCall(method string, url string, headers []types.Headers) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		appendHeaders(req, headers)
	}
	client := &http.Client{}

	return client.Do(req)
}

func GetBodyFromHTTPResponse(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func HandleResponse(data []byte, definition types.Definition) error {
	switch definition.Response.Type {
	case "json":
		return handleJsonResponse(data, definition)
	case "browser":
		return handleBrowserResponse(data, definition)
	}
	return errors.New("Unsupported response type")
}

func handleJsonResponse(data []byte, d types.Definition) error {
	filePath := d.Response.Folder + d.Response.File + "." + d.Response.Type
	fmt.Printf("Handling response for file path %s\n", filePath)
	return ioutil.WriteFile(filePath, data, os.FileMode(d.Response.Permission))
}

func handleBrowserResponse(data []byte, d types.Definition) error {
	var port int
	if (d.Response.Port == 0) {
		port = 8080
	} else {
		port = d.Response.Port
	}

	strport := strconv.Itoa(port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(data))
	})
	openbrowser("http://localhost:" + strport)
	// use the termination of the main thread to stop go routine execution
	go http.ListenAndServe(":" + strport, nil)
	time.Sleep(2 * time.Second)
	return nil
}

func openbrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return errors.New("unsupported platform")
	}
}
