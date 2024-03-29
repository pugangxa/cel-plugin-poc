# CEL-PLUGIN-POC to load cel extension as plugin

This is to allow cel extensions to be loaded dynamically. User need implement their own functions as cel Library and use these functions in the cel expression, the evaluator will load it dynamically and use it for evaluation.

## Getting started

Make plugins under `plugins` directory:
```
cd plugins
make prefix_plugin
make suffix_plugin
```
There can be no plugin, 1 plugin or 2 plugins.

Run the cel evaluation:
```
go run cel_eval.go
```

For cases of no plugin, output should be `Hello world!`, with just prefix plugin, output should be `CEL, Hello world!`, with both prefix and suffix plugin, output should be `CEL, Hello world! Done.`.

## Some notes:

- User must implement their own functions as CEL library.

- To integrate with the cel-eval controller, can watch the plugins directory or watch a configmap which will have the `.so` file list updated during user registration stage, and then force the controller to load the plugins found.

- Consider controller scale out case, the `.so` files should be located on shared storage which can be a volume to mount.
