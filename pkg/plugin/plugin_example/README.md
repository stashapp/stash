# Building

From the base stash source directory:
```
go build -tags=plugin_example ./pkg/plugin/plugin_example/...
```

Place the resulting binary together with `example.yml` in the `plugins` subdirectory of your stash directory.