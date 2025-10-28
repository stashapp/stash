# External Plugin Tasks

External plugin tasks are executed by running an external binary.

## Plugin interfaces

Stash communicates with external plugin tasks using an interface. Stash currently supports RPC and raw interface types.

### RPC interface

The RPC interface uses JSON-RPC to communicate with the plugin process. A golang plugin utilising the RPC interface is available in the stash source code under `pkg/plugin/examples/gorpc`. RPC plugins are expected to provide an interface that fulfils the `RPCRunner` interface in `pkg/plugin/common`.

RPC plugins are expected to accept requests asynchronously.

When stopping an RPC plugin task, the stash server sends a stop request to the plugin and relies on the plugin to stop itself.

### Raw interface

Raw interface plugins are not required to conform to any particular interface. The stash server will send the plugin input to the plugin process via its stdin stream, encoded as JSON. Raw interface plugins are not required to read the input.

The stash server reads stdout for the plugin's output. If the output can be decoded as a JSON representation of the plugin output data structure then it will do so. If not, it will treat the entire stdout string as the plugin's output.

When stopping a raw plugin task, the stash server kills the spawned process without warning or signals.

## Logging

External plugins may log to the stash server by writing to stderr. By default, data written to stderr will be logged by stash at the `error` level. This default behaviour can be changed by setting the `errLog` field in the plugin configuration file.

Plugins can log for specific levels or log progress by prefixing the output string with special control characters. See `pkg/plugin/common/log` for how this is done in go.

## Plugin configuration file format

### exec

For external plugin tasks, the `exec` field is a list with the first element being the binary that will be executed, and the subsequent elements are the arguments passed. The execution process will search the path for the binary, then will attempt to find the program in the same directory as the plugin configuration file. The `exe` extension is not necessary on Windows systems. 

> **⚠️ Note:** The plugin execution process sets the current working directory to that of the stash process.

Arguments can include the plugin's directory with the special string `{pluginDir}`. 

For example, if the plugin executable `my_plugin` is placed in the `plugins` subdirectory and requires arguments `foo` and `bar`, then the `exec` part of the configuration would look like the following:

```
exec:
  - my_plugin
  - foo
  - bar
```

Another example might use a python script to execute the plugin. Assuming the python script `foo.py` is placed in the same directory as the plugin config file, the `exec` fragment would look like the following:

```
exec:
  - python
  - {pluginDir}/foo.py
```

### interface

For external plugin tasks, the `interface` field must be set to one of the following values:
* `rpc`
* `raw`

See the `Plugin interfaces` section above for details on these interface types.

The `interface` field defaults to `raw` if not provided.

### errLog

The `errLog` field tells stash what the default log level should be when the plugin outputs to stderr without encoding a log level. It defaults to the `error` level if no provided. This field is not necessary if the plugin outputs logging with the appropriate encoding. See the `Logging` section above for details.

## Task configuration

In addition to the standard task configuration, external tasks may be configured with an optional `execArgs` field to add extra parameters to the execution arguments for the task.

For example:

```
tasks:
  - name: <operation name>
    description: <optional description>
    execArgs:
      - <arg to add to the exec line>
```
