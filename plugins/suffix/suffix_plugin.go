package main

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type suffixType struct{}

func (suffixType) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Declarations(
			decls.NewFunction("custom.AddSuffix",
				decls.NewOverload("add_suffix",
					[]*exprpb.Type{decls.String, decls.String},
					decls.String),
			),
		),
	}
}

func (suffixType) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "custom.AddSuffix",
				Binary:   callInStringStringOutString(addSuffixImpl),
			},
		),
	}
}

func callInStringStringOutString(fn func(string, string) (string, error)) functions.BinaryOp {
	return func(arg1, arg2 ref.Val) ref.Val {
		argVal1, ok := arg1.(types.String)
		if !ok {
			return types.MaybeNoSuchOverloadErr(arg1)
		}
		argVal2, ok := arg2.(types.String)
		if !ok {
			return types.MaybeNoSuchOverloadErr(arg2)
		}
		out, err := fn(string(argVal1), string(argVal2))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.String(out)
	}
}

var CustomLib suffixType
