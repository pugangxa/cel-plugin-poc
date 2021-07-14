package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"plugin"
	"regexp"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

const (
	pluginsDir = "./plugins"
	celLibName = "CustomLib"
)

func listFiles(dir, pattern string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			return nil, err
		}
		if matched {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

func loadPlugins(env *cel.Env) (*cel.Env, error) {
	if _, err := os.Stat(pluginsDir); err != nil {
		return nil, err
	}

	plugins, err := listFiles(pluginsDir, `.*_plugin.so`)
	if err != nil {
		return nil, err
	}

	for _, celPlugin := range plugins {
		plug, err := plugin.Open(path.Join(pluginsDir, celPlugin.Name()))
		if err != nil {
			fmt.Printf("failed to open plugin %s: %v\n", celPlugin.Name(), err)
			continue
		}
		celLibSymbol, err := plug.Lookup(celLibName)
		if err != nil {
			fmt.Printf("plugin %s does not export symbol \"%s\"\n",
				celPlugin.Name(), celLibName)
			continue
		}
		celLib, ok := celLibSymbol.(cel.Library)
		if !ok {
			fmt.Printf("Symbol %s (from %s) does not implement cel.Library interface\n",
				celLibName, celPlugin.Name())
			continue
		}
		env, _ = env.Extend(cel.Lib(celLib))
	}
	return env, nil
}

func main() {
	d := cel.Declarations(decls.NewVar("name", decls.String))

	// Create a program environment configured with the standard library of CEL functions and macros
	env, err := cel.NewEnv(
		d,
	)
	if err != nil {
		log.Fatalln("Couldn't create a program env with standard library of CEL functions and loaded plugins")
	}
	env, err = loadPlugins(env)
	if err != nil {
		log.Fatalln("Error when loading plugins.")
	}
	// no plugin, evaluateVarExpression := "Hello world!"
	// just prefix plugin, evaluateVarExpression := `custom.AddPrefix("Hello world!", name)`
	evaluateVarExpression := `custom.AddSuffix(custom.AddPrefix("Hello world!", name), " Done.")`

	ast, iss := env.Compile(evaluateVarExpression)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, _ := env.Program(ast)
	out, _, _ := prg.Eval(map[string]interface{}{
		"name": "CEL, ",
	})
	fmt.Println(out)
}
