package main

import (
	"errors"
	"fmt"
	"github.com/gogo/googleapis/google/api"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity/command"
	"log"
)

type Generator struct {
	*generator.Generator
	generator.PluginImports
	overwrite bool
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Name() string {
	return "http-client"
}

func (g *Generator) Init(gg *generator.Generator) {
	g.Generator = gg
}

func (g *Generator) Overwrite() {
	g.overwrite = true
}

func (g *Generator) Generate(file *generator.FileDescriptor) {
	//proto3 := gogoproto.IsProto3(file.FileDescriptorProto)
	g.PluginImports = generator.NewPluginImports(g.Generator)

	httpPkg := g.NewImport("net/http")
	jsonPkg := g.NewImport("encoding/json")
	//fmtPkg := g.NewImport("fmt")
	ctxPkg := g.NewImport("context")

	for _, svc := range file.Service {
		svcName := *svc.Name + "HttpClient"

		// type HelloHttpClient struct {
		//   client net_http.Client
		//   baseURL string
		// }
		g.P("type ", svcName, " struct {")
		g.In()
		g.P("client ", httpPkg.Use(), ".Client")
		g.P("baseURL string")
		g.Out()
		g.P("}")

		// func NewHelloHttpClient() {
		//   return &HelloHttpClient{
		//     client: net_http.Client{},
		//	   baseURL: baseURL,
		//   }
		// }
		g.P("func New", svcName, "(baseURL string) *", svcName, "{")
		g.In()
		g.P("return &", svcName, "{")
		g.In()
		g.P("client: ", httpPkg.Use(), ".Client{},")
		g.P("baseURL: baseURL,")
		g.Out()
		g.P("}")
		g.Out()
		g.P("}")

		for _, method := range svc.Method {
			// func (c *GreeterHttpClient) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
			//   out := new(HelloReply)
			// }
			in := g.TypeName(g.ObjectNamed(method.GetInputType()))
			out := g.TypeName(g.ObjectNamed(method.GetOutputType()))
			endpoint, err := GetMethodType(method)
			if err != nil {
				log.Fatal(err)
			}
			g.P("func (c *", svcName, ") ", method.Name, "(ctx ", ctxPkg.Use(), ".Context, in *", in, ") (*", out, ", error) {")
			g.In()
			g.P("out := new(", out, ")")
			g.P(fmt.Sprintf(`request, err := %s.NewRequest(%s.%s, c.baseURL+"%s", nil)`, httpPkg.Use(), httpPkg.Use(), endpoint.Method, endpoint.Path))
			g.checkErr()
			g.P("response, err := c.client.Do(request)")
			g.checkErr()
			g.P("err = ", jsonPkg.Use(), ".NewDecoder(response.Body).Decode(out)")
			g.checkErr()
			g.P("return out, nil")
			g.Out()
			g.P("}")
		}
	}
}
func (g *Generator) checkErr() {
	g.P("if err != nil {")
	g.In()
	g.P("return nil, err")
	g.Out()
	g.P("}")
}

func GetMethodType(m *descriptor.MethodDescriptorProto) (Endpoint, error) {
	if !proto.HasExtension(m.Options, api.E_Http) {
		return Endpoint{}, errors.New("no extension")
	}
	ext, err := proto.GetExtension(m.Options, api.E_Http)
	if err != nil {
		return Endpoint{}, err
	}
	rule, ok := ext.(*api.HttpRule)
	if !ok {
		return Endpoint{}, errors.New("no http rule")
	}
	return GetEndpoint(rule), err
}

type Endpoint struct {
	Method string
	Path   string
}

func GetEndpoint(opts *api.HttpRule) Endpoint {
	if opts == nil {
		return Endpoint{}
	}
	switch opt := opts.GetPattern().(type) {
	case *api.HttpRule_Get:
		return Endpoint{
			Method: "MethodGet",
			Path:   opt.Get,
		}
	case *api.HttpRule_Put:
		return Endpoint{
			Method: "MethodPut",
			Path:   opt.Put,
		}
	case *api.HttpRule_Post:
		return Endpoint{
			Method: "MethodPost",
			Path:   opt.Post,
		}
	case *api.HttpRule_Delete:
		return Endpoint{
			Method: "MethodDelete",
			Path:   opt.Delete,
		}
	case *api.HttpRule_Patch:
		return Endpoint{
			Method: "MethodPatch",
			Path:   opt.Patch,
		}
	}
	return Endpoint{}
}

func init() {
	generator.RegisterPlugin(NewGenerator())
}

func main() {
	req := command.Read()
	p := NewGenerator()
	p.Overwrite()
	resp := command.GeneratePlugin(req, p, "_httpclient.gen.go")
	command.Write(resp)
}
