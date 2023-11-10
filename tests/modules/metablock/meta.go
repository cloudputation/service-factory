package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	Service Service `hcl:"service,block"`
}

type Service struct {
	Name string    `hcl:"name,label"`
	Meta cty.Value `hcl:"meta,attr"`
}

func main() {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL([]byte(`
		service "example" {
		  meta = {
		    key1 = "value1"
		    key2 = 42
		  }
		}
	`), "example.hcl")
	if diags.HasErrors() {
		fmt.Println("Error parsing HCL:", diags)
		return
	}

	var config Config
	diags = gohcl.DecodeBody(f.Body, nil, &config)
	if diags.HasErrors() {
		fmt.Println("Error decoding HCL:", diags)
		return
	}

	fmt.Printf("Service Name: %s\n", config.Service.Name)

	metaMap := config.Service.Meta.AsValueMap()
	for k, v := range metaMap {
		if v.Type() == cty.String {
			fmt.Printf("Meta key: %s, value: %s\n", k, v.AsString())
		} else if v.Type() == cty.Number {
			val, _ := v.AsBigFloat().Int64()
			fmt.Printf("Meta key: %s, value: %d\n", k, val)
		}
	}
}
