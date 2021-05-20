package types

type Definition struct {
	Title string
	Desc  string
	Details  struct {
		Endpoint string
		Headers  []Headers
		Path string
		Method string
	}
	Response struct {
		Type string
		File string
		Folder string
		Permission int
		Port int
	}
}

type Headers struct {
	Header string
	Key string
	Value string
}