# Plugins

Stash supports the running external tasks via plugins. Plugins are implemented by calling an external binary.

> **⚠️ Note:** Plugin support is still experimental and is likely to change.

# Adding plugins

By default, Stash looks for plugin configurations in the `plugins` sub-directory of the directory where the stash `config.yml` is read. This will either be the `$HOME/.stash` directory or the current working directory.

Plugins are added by adding configuration yaml files (format: `pluginName.yml`) to the `plugins` directory.

Loaded plugins can be viewed in the Plugins page of the Settings. After plugins are added, removed or edited while stash is running, they can be reloaded by clicking `Reload Plugins` button.

# Using plugins

Plugins provide tasks which can be run from the Tasks page. 

> **⚠️ Note:** It is currently only possible to run one task at a time. No queuing is currently implemented.

# Plugin configuration file format

The basic structure of a plugin configuration file is as follows:

```
name: <plugin name>
description: <optional description of the plugin>
version: <optional version tag>
url: <optional url>
exec:
  - <binary name>
  - <other args...>
interface: [interface type]
errLog: [one of none trace, debug, info, warning, error]
tasks:
  - ...
```

## Plugin process execution

The `exec` field is a list with the first element being the binary that will be executed, and the subsequent elements are the arguments passed. The execution process will search the path for the binary, then will attempt to find the program in the same directory as the plugin configuration file. The `exe` extension is not necessary on Windows systems. 

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

## Plugin interfaces

The `interface` field currently accepts one of two possible values: `rpc` and `raw`. It defaults to `raw` if not provided.

Plugins may log to the stash server by writing to stderr. By default, data written to stderr will be logged by stash at the `error` level. This default behaviour can be changed by setting the `errLog` field.

Plugins can log for specific levels or log progress by prefixing the output string with special control characters. See `pkg/plugin/common/log` for how this is done in go.

### RPC interface

The RPC interface uses JSON-RPC to communicate with the plugin process. A golang plugin utilising the RPC interface is available in the stash source code under `pkg/plugin/examples/gorpc`. RPC plugins are expected to provide an interface that fulfils the `RPCRunner` interface in `pkg/plugin/common`.

RPC plugins are expected to accept requests asynchronously.

When stopping an RPC plugin task, the stash server sends a stop request to the plugin and relies on the plugin to stop itself.

### Raw interface

Raw interface plugins are not required to conform to any particular interface. The stash server will send the plugin input to the plugin process via its stdin stream, encoded as JSON. Raw interface plugins are not required to read the input.

The stash server reads stdout for the plugin's output. If the output can be decoded as a JSON representation of the plugin output data structure then it will do so. If not, it will treat the entire stdout string as the plugin's output.

When stopping a raw plugin task, the stash server kills the spawned process without warning or signals.

## Plugin input

Plugins may accept an input from the stash server. This input is encoded according to the interface, and has the following structure (presented here in JSON format):
```
{
    "server_connection": {
        "Scheme": "http",
        "Port": 9999,
        "SessionCookie": {
            "Name":"session",
            "Value":"cookie-value",
            "Path":"",
            "Domain":"",
            "Expires":"0001-01-01T00:00:00Z",
            "RawExpires":"",
            "MaxAge":0,
            "Secure":false,
            "HttpOnly":false,
            "SameSite":0,
            "Raw":"",
            "Unparsed":null
        },
        "Dir": <path to stash config directory>,
        "PluginDir": <path to plugin config directory>,
    },
    "args": {
        "argKey": "argValue"
    }
}
```

The `server_connection` field contains all the information needed for a plugin to access the parent stash server.

## Plugin output

Plugin output is expected in the following structure (presented here as JSON format):

```
{
    "error": <optional error string>
    "output": <anything>
}
```

The `error` field is logged in stash at the `error` log level if present. The `output` is written at the `debug` log level.

## Task configuration

Tasks are configured using the following structure:

```
tasks:
  - name: <operation name>
    description: <optional description>
    defaultArgs:
      argKey: argValue
    execArgs:
      - <arg to add to the exec line>
```

A plugin configuration may contain multiple tasks. 

The `defaultArgs` field is used to add inputs to the plugin input sent to the plugin.

The `execArgs` field allows adding extra parameters to the execution arguments for this task.
