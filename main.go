package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity/command"
	"google.golang.org/genproto/googleapis/api/annotations"
	"reflect"
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

func (p *Generator) Overwrite() {
	p.overwrite = true
}

func (g *Generator) Generate(file *generator.FileDescriptor) {
	//proto3 := gogoproto.IsProto3(file.FileDescriptorProto)
	g.PluginImports = generator.NewPluginImports(g.Generator)

	httpPkg := g.NewImport("net/http")
	//fmtPkg := g.NewImport("fmt")
	ctxPkg := g.NewImport("context")

	for _, svc := range file.Service {
		svcName := *svc.Name + "HttpClient"

		// type HelloHttpClient struct {
		//   client net_http.Client
		// }
		g.P("type ", svcName, " struct {")
		g.In()
		g.P("client ", httpPkg.Use(), ".Client")
		g.Out()
		g.P("}")

		// func NewHelloHttpClient() {
		//   return &HelloHttpClient{
		//     client: net_http.Client{},
		//   }
		// }
		g.P("func New", svcName, "() *", svcName, "{")
		g.In()
		g.P("return &", svcName, "{")
		g.In()
		g.P("client: ", httpPkg.Use(), ".Client{},")
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

			g.P("// ", g.GetMethodType(method))

			g.P("func (c *", svcName, ") ", method.Name, "(ctx ", ctxPkg.Use(), ".Context, in *", in, ") (*", out, ", error) {")
			g.In()
			g.P("out := new(", out, ")")
			g.P("request, err := ", httpPkg.Use(), ".NewRequest()")
			g.P("return nil, nil")
			g.Out()
			g.P("}")
		}
	}
}

func (g *Generator) GetMethodType(m *descriptor.MethodDescriptorProto) string {
	ext, _ := proto.GetExtension(m.Options, &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.MethodOptions)(nil),
		ExtensionType: (*annotations.HttpRule)(nil),
		Field:         72295728,
		Name:          "google.api.http",
		Tag:           "bytes,72295728,opt,name=http",
		Filename:      "google/api/annotations.proto",
	})

	rule := ext.(*annotations.HttpRule)

	return fmt.Sprintf("%+v, %+v\n", string(rule.XXX_unrecognized), reflect.TypeOf(rule.Pattern))
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
