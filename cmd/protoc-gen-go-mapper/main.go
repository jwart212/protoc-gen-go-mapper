package main

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/jwart212/protoc-gen-go-mapper/internal/config"
	"github.com/jwart212/protoc-gen-go-mapper/internal/plugin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// Read CodeGeneratorRequest from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	var req pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Load configuration from parameter or default
	configPath := "mapper.yaml"
	if req.Parameter != nil {
		params := strings.Split(*req.Parameter, ",")
		for _, param := range params {
			if strings.HasSuffix(param, ".yaml") || strings.HasSuffix(param, ".yml") {
				configPath = param
				break
			}
		}
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		// If config file not found, use defaults
		cfg = &config.Config{
			Version:   "v1",
			Database:  "sqlc",
			DBPackage: "", // Will be derived from proto package
			Package: config.Package{
				Proto: "internal/gen",
				DB:    "internal/postgres",
			},
		}
	}

	// Create plugin
	p := plugin.New(cfg)

	// Generate code for each proto file (only those explicitly requested)
	var resp pluginpb.CodeGeneratorResponse
	// Indicate support for proto3 optional fields
	resp.SupportedFeatures = proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))
	for _, fileName := range req.FileToGenerate {
		// Find the file descriptor for this file
		var fileProto *descriptorpb.FileDescriptorProto
		for _, fp := range req.ProtoFile {
			if *fp.Name == fileName {
				fileProto = fp
				break
			}
		}
		if fileProto == nil {
			io.WriteString(os.Stderr, "File not found: "+fileName)
			os.Exit(1)
		}

		genReq := &plugin.GenerateRequest{
			FileProto: fileProto,
		}

		var buf bytes.Buffer
		err := p.Generate(genReq, &buf)
		if err != nil {
			io.WriteString(os.Stderr, err.Error())
			os.Exit(1)
		}

		// Determine output path based on go_package option
		outputPath := *fileProto.Name + "_mapper.pb.go"
		// Get the base filename from the proto file
		baseName := *fileProto.Name
		if idx := strings.LastIndex(baseName, "/"); idx != -1 {
			baseName = baseName[idx+1:]
		}
		outputPath = baseName + "_mapper.pb.go"

		// Add generated file to response
		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(outputPath),
			Content: proto.String(buf.String()),
		})
	}

	// Write response to stdout
	respData, err := proto.Marshal(&resp)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	os.Stdout.Write(respData)
}
