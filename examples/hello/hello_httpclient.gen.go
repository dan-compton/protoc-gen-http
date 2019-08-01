// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hello.proto

package hello

import (
	context "context"
	encoding_json "encoding/json"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
	net_http "net/http"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type GreeterHttpClient struct {
	client  net_http.Client
	baseURL string
}

func NewGreeterHttpClient(baseURL string) *GreeterHttpClient {
	return &GreeterHttpClient{
		client:  net_http.Client{},
		baseURL: baseURL,
	}
}
func (c *GreeterHttpClient) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	out := new(HelloReply)
	request, err := net_http.NewRequest("GET", c.baseURL+"/hello", nil)
	if err != nil {
		return nil, err
	}
	q := request.URL.Query()
	q.Add("name", in.Name)
	request.URL.RawQuery = q.Encode()
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	err = encoding_json.NewDecoder(response.Body).Decode(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
