package headers

// Header foo
type Header map[string]string

// JSONHeader foo
var JSONHeader Header = Header{
	"Content-Type":                "application/json",
	"Access-Control-Allow-Origin": "*",
}
