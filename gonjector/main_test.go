package main

import (
	"go/ast"
	"reflect"
	"testing"
)

func TestGenAstFromFile(t *testing.T) {
	type args struct {
		fnName string
	}
	tests := []struct {
		name   string
		fnName string
		want   []ast.Stmt
	}{
		{name: "test01",
			fnName: "string",
			want:   []ast.Stmt{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenAstFromFile(tt.fnName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenAstFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
