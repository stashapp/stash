# Building

From the base stash source directory:
```
go build -tags=plugin_example -o plugin_goraw.exe ./pkg/plugin/examples/goraw/...
go build -tags=plugin_example -o plugin_gorpc.exe ./pkg/plugin/examples/gorpc/...
```

Place the resulting binaries together with the yml files in the `plugins` subdirectory of your stash directory.