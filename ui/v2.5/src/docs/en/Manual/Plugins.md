# Plugins

Stash supports plugins that can do the following:

- perform custom tasks when triggered by the user from the Tasks page
- perform custom tasks when triggered from specific events
- add custom CSS to the UI
- add custom JavaScript to the UI

Plugin tasks can be implemented using embedded Javascript, or by calling an external binary.

> **⚠️ Note:** Plugin support is still experimental and is likely to change.

## Managing Plugins

Plugins can be installed and managed from the `Settings > Plugins` page. 

Plugins are installed using the `Available Plugins` section. This section allows configuring sources from which to install plugins. The `Community (stable)` source is configured by default. This source contains plugins for the current _stable_ version of stash.

These are the plugin sources maintained by the stashapp organisation:

| Name | Source URL | Recommended Local Path | Notes |
|------|-----------|------------------------|-------|
| Community (stable) | `https://stashapp.github.io/CommunityScripts/stable/index.yml` | `stable` | For the current stable version of stash. |
| Community (develop) | `https://stashapp.github.io/CommunityScripts/develop/index.yml` | `develop` | For the develop version of stash. |

Installed plugins can be updated or uninstalled from the `Installed Plugins` section.

### Source URLs

The source URL must return a yaml file containing all the available packages for the source. An example source yaml file looks like the following:

```
- id: <package id>
  name: <package name>
  version: <version>
  date: <date>
  requires:
  - <ids of packages required by this package (optional)>
  - ...
  path: <path to package zip file>
  sha256: <sha256 of zip>
  metadata:
    <optional key/value pairs for extra information>
- ...
```

Path can be a relative path to the zip file or an external URL.

## Adding plugins manually

By default, Stash looks for plugin configurations in the `plugins` sub-directory of the directory where the stash `config.yml` is read. This will either be the `$HOME/.stash` directory or the current working directory.

Plugins are added by adding configuration yaml files (format: `pluginName.yml`) to the `plugins` directory.

Loaded plugins can be viewed in the Plugins page of the Settings. After plugins are added, removed or edited while stash is running, they can be reloaded by clicking `Reload Plugins` button.

## Using plugins

Plugins provide tasks which can be run from the Tasks page. 

## Creating plugins

### Plugin configuration file format

The basic structure of a plugin configuration file is as follows:

```yaml
name: <plugin name> 
# optional list of dependencies to be included
# "#" is is part of the config - do not remove
# requires: <plugin ID>
description: <optional description of the plugin>
version: <optional version tag>
url: <optional url>

ui:
  # optional list of css files to include in the UI
  css:
    - <path to css file>

  # optional list of js files to include in the UI
  javascript:
    - <path to javascript file>

  # optional list of plugin IDs to load prior to this plugin
  requires:
    - <plugin ID>

  # optional list of assets 
  assets:
    urlPrefix: fsLocation
    ...

  # content-security policy overrides
  csp:
    script-src:
      - http://alloweddomain.com
    
    style-src:
      - http://alloweddomain.com
    
    connect-src:
      - http://alloweddomain.com

# map of setting names to be displayed in the plugins page in the UI
settings:
  # internal name
  foo:
  # name to display in the UI
  displayName: Foo
  # type of the attribute to show in the UI
  # can be BOOLEAN, NUMBER, or STRING
  type: BOOLEAN

# the following are used for plugin tasks only
exec:
  - ...
interface: [interface type]
errLog: [one of none trace, debug, info, warning, error]
tasks:
  - ...
```

The `name`, `description`, `version` and `url` fields are displayed on the plugins page.

`# requires` will make the plugin manager select plugins matching the specified IDs to be automatically installed as dependencies. Only works with plugins within the same index.

The `exec`, `interface`, `errLog` and `tasks` fields are used only for plugins with tasks.

The `settings` field is used to display plugin settings on the plugins page. Plugin settings can also be set using the graphql mutation `configurePlugin` - the settings set this way do _not_ need to be specified in the `settings` field unless they are to be displayed in the stock plugin settings UI.

### UI Configuration

The `css` and `javascript` field values may be relative paths to the plugin configuration file, or
may be full external URLs.

The `requires` field is a list of plugin IDs which must have their javascript/css files loaded
before this plugins javascript/css files.

The `assets` field is a map of URL prefixes to filesystem paths relative to the plugin configuration file.
Assets are mounted to the `/plugin/{pluginID}/assets` path. 

As an example, for a plugin with id `foo` with the following `assets` value:
```
assets:
  foo: bar
  /: .
```
The following URLs will be mapped to these locations:
`/plugin/foo/assets/foo/file.txt` -> `{pluginDir}/bar/file.txt`
`/plugin/foo/assets/file.txt` -> `{pluginDir}/file.txt`
`/plugin/foo/assets/bar/file.txt` -> `{pluginDir}/bar/file.txt` (via the `/` entry)

Mappings that try to go outside of the directory containing the plugin configuration file will be
ignored.

The `csp` field contains overrides to the content security policies. The URLs in `script-src`,
`style-src` and `connect-src` will be added to the applicable content security policy.

See [External Plugins](/help/ExternalPlugins.md) for details for making plugins with external tasks.

See [Embedded Plugins](/help/EmbeddedPlugins.md) for details for making plugins with embedded tasks.

### Plugin task input

Plugin tasks may accept an input from the stash server. This input is encoded according to the interface, and has the following structure (presented here in JSON format):
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

### Plugin task output

Plugin task output is expected in the following structure (presented here as JSON format):

```
{
    "error": <optional error string>
    "output": <anything>
}
```

The `error` field is logged in stash at the `error` log level if present. The `output` is written at the `debug` log level.

### Task configuration

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

### Hook configuration

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

#### Trigger types

Trigger types use the following format: `<object type>.<operation>.<hook type>`

For example, a post-hook on a scene create operation will be `Scene.Create.Post`.

The following object types are supported:

* `Scene`
* `SceneMarker`
* `Image`
* `Gallery`
* `Group`
* `Performer`
* `Studio`
* `Tag`

The following operations are supported:

* `Create`
* `Update`
* `Destroy`
* `Merge` (for `Tag` only)

Currently, only `Post` hook types are supported. These are executed after the operation has completed and the transaction is committed.

#### Hook input

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
            "groups":null,
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
