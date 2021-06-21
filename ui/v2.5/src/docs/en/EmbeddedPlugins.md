# Embedded Plugins

Embedded plugins are executed within the stash process using a scripting system.

## Supported script languages

Stash currently supports Javascript embedded plugins using [otto](https://github.com/robertkrimen/otto).

# Javascript plugins

## Plugin input

The input is provided to Javascript plugins using the `input` global variable, and is an object based on the structure provided in the `Plugin input` section of the [Plugins](/help/Plugins.md) page. Note that the `server_connection` field should not be necessary in most embedded plugins.

## Plugin output

The output of a Javascript plugin is derived from the evaluated value of the script. The output should conform to the structure provided in the `Plugin output` section of the [Plugins](/help/Plugins.md) page.

There are a number of ways to return the plugin output:

### Example #1
```
(function() {
    return {
        Output: "ok"
    };
})();
```

### Example #2
```
function main() {
    return {
        Output: "ok"
    };
}

main();
```

### Example #3
```
var output = {
    Output: "ok"
};

output;
```

## Logging

See the `Javascript API` section below on how to log with Javascript plugins.

# Plugin configuration file format

The basic structure of an embedded plugin configuration file is as follows:

```
name: <plugin name>
description: <optional description of the plugin>
version: <optional version tag>
url: <optional url>
exec:
  - <path to script>
interface: [interface type]
tasks:
  - ...
```

The `name`, `description`, `version` and `url` fields are displayed on the plugins page.

## exec

For embedded plugins, the `exec` field is a list with the first element being the path to the Javascript file that will be executed. It is expected that the path to the Javascript file is relative to the directory of the plugin configuration file.

## interface

For embedded plugins, the `interface` field must be set to one of the following values:
* `js`

# Javascript API

## Logging

Stash provides the following API for logging in Javascript plugins:

| Method | Description |
|--------|-------------|
| `log.Trace(<string>)` | Log with the `trace` log level. |
| `log.Debug(<string>)` | Log with the `debug` log level. |
| `log.Info(<string>)` | Log with the `info` log level. |
| `log.Warn(<string>)` | Log with the `warn` log level. |
| `log.Error(<string>)` | Log with the `error` log level. |
| `log.Progress(<float between 0 and 1>)` | Sets the progress of the plugin task, as a float, where `0` represents 0% and `1` represents 100%. |

## GQL

Stash provides the following API for communicating with stash using the graphql interface:

| Method | Description |
|--------|-------------|
| `gql.Do(<query/mutation string>, <variables object>)` | Executes a graphql query/mutation on the stash server. Returns an object in the same way as a graphql query does. |

### Example

```
// creates a tag
var mutation = "\
mutation tagCreate($input: TagCreateInput!) {\
  tagCreate(input: $input) {\
    id\
  }\
}";

var variables = {
    input: {
        'name': tagName
    }
};

result = gql.Do(mutation, variables);
log.Info("tag id = " + result.tagCreate.id);
```

## Utility functions

Stash provides the following API for utility functions:

| Method | Description |
|--------|-------------|
| `util.Sleep(<milliseconds>)` | Suspends the current thread for the specified duration. |
