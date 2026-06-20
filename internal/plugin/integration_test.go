package plugin

import (
	"bytes"
	"testing"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/internal/config"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestPluginIntegration(t *testing.T) {
	cfg := &config.Config{
		Version:  "v1",
		Database: "sqlc",
		Package: config.Package{
			Proto: "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/testdata/gen",
			DB:    "gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/testdata/db",
		},
	}

	p := New(cfg)

	fileProto := &descriptorpb.FileDescriptorProto{
		Name: proto.String("test.proto"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("User"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("id"),
						Number:   proto.Int32(1),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						JsonName: proto.String("id"),
					},
					{
						Name:     proto.String("name"),
						Number:   proto.Int32(2),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						JsonName: proto.String("name"),
					},
				},
			},
		},
	}

	req := &GenerateRequest{
		FileProto: fileProto,
	}

	var buf bytes.Buffer
	err := p.Generate(req, &buf)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	generated := buf.String()
	if len(generated) == 0 {
		t.Error("Generate() produced no output")
	}
}

func TestProtocPluginProtocol(t *testing.T) {
	// Create a minimal CodeGeneratorRequest
	req := &pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			{
				Name: proto.String("test.proto"),
				MessageType: []*descriptorpb.DescriptorProto{
					{
						Name: proto.String("Simple"),
						Field: []*descriptorpb.FieldDescriptorProto{
							{
								Name:     proto.String("id"),
								Number:   proto.Int32(1),
								Type:     descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
								Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								JsonName: proto.String("id"),
							},
						},
					},
				},
			},
		},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Verify request can be unmarshaled
	var unmarshaled pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if len(unmarshaled.ProtoFile) != 1 {
		t.Errorf("Expected 1 proto file, got %d", len(unmarshaled.ProtoFile))
	}
}
