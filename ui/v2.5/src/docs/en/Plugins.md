# Plugins

Stash supports the running tasks via plugins. Plugins can be implemented using embedded Javascript, or by calling an external binary.

Stash also supports triggering of plugin hooks from specific stash operations.

> **⚠️ Note:** Plugin support is still experimental and is likely to change.

# Adding plugins

By default, Stash looks for plugin configurations in the `plugins` sub-directory of the directory where the stash `config.yml` is read. This will either be the `$HOME/.stash` directory or the current working directory.

Plugins are added by adding configuration yaml files (format: `pluginName.yml`) to the `plugins` directory.

Loaded plugins can be viewed in the Plugins page of the Settings. After plugins are added, removed or edited while stash is running, they can be reloaded by clicking `Reload Plugins` button.

# Using plugins

Plugins provide tasks which can be run from the Tasks page. 

# Creating plugins

See [External Plugins](/help/ExternalPlugins.md) for details for making external plugins.

See [Embedded Plugins](/help/EmbeddedPlugins.md) for details for making embedded plugins.

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

The `server_connection` field contains all the information needed for a plugin to access the parent stash server, if necessary.

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
```

A plugin configuration may contain multiple tasks. 

The `defaultArgs` field is used to add inputs to the plugin input sent to the plugin.

## Hook configuration

Stash supports executing plugin operations via triggering of a hook during a stash operation.

Hooks are configured using a similar structure to tasks:

```
hooks:
  - name: <operation name>
    description: <optional description>
    triggeredBy:
      - <trigger types>...
    defaultArgs:
      argKey: argValue
```

**Note:** it is possible for hooks to trigger eachother or themselves if they perform mutations. For safety, hooks will not be triggered if they have already been triggered in the context of the operation. Stash uses cookies to track this context, so it's important for plugins to send cookies when performing operations.

### Trigger types

Trigger types use the following format:
`<object type>.<operation>.<hook type>`

For example, a post-hook on a scene create operation will be `Scene.Create.Post`.

The following object types are supported:
* `Scene`
* `SceneMarker`
* `Image`
* `Gallery`
* `Movie`
* `Performer`
* `Studio`
* `Tag`

The following operations are supported:
* `Create`
* `Update`
* `Destroy`
* `Merge` (for `Tag` only)

Currently, only `Post` hook types are supported. These are executed after the operation has completed and the transaction is committed.

### Hook input

Plugin tasks triggered by a hook include an argument named `hookContext` in the `args` object structure. The `hookContext` is structured as follows:

```
{
    "id": <object id>,
    "type": <trigger type>,
    "input": <operation input>,
    "inputFields": <fields included in input>
}
```

The `input` field contains the JSON graphql input passed to the original operation. This will differ between operations. For hooks triggered by operations in a scan or clean, the input will be nil. `inputFields` is populated in update operations to indicate which fields were passed to the operation, to differentiate between missing and empty fields.

For example, here is the `args` values for a Scene update operation:

```
{
    "hookContext": {
        "type":"Scene.Update.Post",
        "id":45,
        "input":{
            "clientMutationId":null,
            "id":"45",
            "title":null,
            "details":null,
            "url":null,
            "date":null,
            "rating":null,
            "organized":null,
            "studio_id":null,
            "gallery_ids":null,
            "performer_ids":null,
            "movies":null,
            "tag_ids":["21"],
            "cover_image":null,
            "stash_ids":null
        },
        "inputFields":[
            "tag_ids",
            "id"
        ]
    }
}
```
