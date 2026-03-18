package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

var (
	flags         flag.FlagSet
	enableMocks   = flags.Bool("mocks", true, "generate mock files")
)

func main() {
	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	opts.Run(func(gen *protogen.Plugin) error {
		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}

			if err := processFile(gen, file); err != nil {
				return err
			}
		}

		return nil
	})
}

func processFile(gen *protogen.Plugin, file *protogen.File) error {
	if file == nil {
		return fmt.Errorf("nil file")
	}

	if len(file.Services) == 0 {
		return nil
	}

	if strings.Contains(file.GeneratedFilenamePrefix, "wire") {
		return nil
	}

	if *enableMocks {
		if err := generateMocks(file); err != nil {
			return fmt.Errorf("generate mocks for %s: %w", file.Desc.Path(), err)
		}
	}

	return nil
}

func generateMocks(file *protogen.File) error {
	baseFileName := filepath.Base(file.GeneratedFilenamePrefix)
	sourceFilePath := filepath.Join("gen", "go", file.GeneratedFilenamePrefix+"_grpc.pb.go")

	if _, err := os.Stat(sourceFilePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source file does not exist: %s", sourceFilePath)
		}
		return fmt.Errorf("stat source file: %w", err)
	}

	mockDir := filepath.Join(filepath.Dir(sourceFilePath), "mock")
	if err := os.MkdirAll(mockDir, 0o755); err != nil {
		return fmt.Errorf("create mock directory: %w", err)
	}

	mockFilePath := filepath.Join(mockDir, baseFileName+".pb.go")
	packageName := "mock_" + string(file.GoPackageName)

	cmd := exec.Command(
		"mockgen",
		"-source="+sourceFilePath,
		"-destination="+mockFilePath,
		"-package="+packageName,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run mockgen: %w", err)
	}

	return nil
}