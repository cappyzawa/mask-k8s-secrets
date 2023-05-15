package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	sensitiveMask = "***"
)

type cli struct {
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
}

func (c *cli) run(args []string) int {
	input, err := io.ReadAll(c.inStream)
	if err != nil {
		fmt.Fprintf(c.errStream, "failed to read input: %s\n", err.Error())
		return 1
	}
	yamls := bytes.Split(input, []byte("\n---\n"))
	for i, y := range yamls {
		if i > 0 {
			fmt.Fprintf(c.outStream, "---\n")
		}
		var obj unstructured.Unstructured
		if err := yaml.Unmarshal(y, &obj.Object); err != nil {
			fmt.Fprintf(c.errStream, "failed to unmarshal yaml: %s\n", err.Error())
			return 1
		}
		if gvk := obj.GroupVersionKind(); gvk.Kind == "Secret" && gvk.Version == "v1" {
			// mask secret data with sensitiveMask
			if data, ok, err := unstructured.NestedMap(obj.Object, "data"); err == nil && ok {
				for k := range data {
					data[k] = sensitiveMask
				}
				if err := unstructured.SetNestedMap(obj.Object, data, "data"); err != nil {
					fmt.Fprintf(c.errStream, "failed to set nested map: %s\n", err.Error())
				}
			}
		}
		if err := yaml.NewEncoder(c.outStream).Encode(&obj.Object); err != nil {
			fmt.Fprintf(c.errStream, "failed to encode yaml: %s\n", err.Error())
			return 1
		}
	}

	return 0
}

func main() {
	c := &cli{inStream: os.Stdin, outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(c.run(os.Args))
}
