# UI Plugin API

The `PluginApi` object is a global object in the `window` object.

`PluginApi` is considered experimental and is subject to change without notice. This documentation covers only the plugin-specific API. It does not necessarily cover the core UI API. Information on these methods should be referenced in the UI source code.

An example using various aspects of `PluginApi` may be found in the source code under `pkg/plugin/examples/react-component`.

## Properties

### `React`

An instance of the React library.

### `ReactDOM`

An instance of the ReactDOM library.

### `GQL`

This namespace contains the generated graphql client interface. This is a low-level interface. In many cases, `StashService` should be used instead.

### `libraries`

`libraries` provides access to the following UI libraries:
- `ReactRouterDOM`
- `Bootstrap`
- `Apollo`
- `Intl`
- `FontAwesomeRegular`
- `FontAwesomeSolid`
- `Mousetrap`
- `MousetrapPause`

### `register`

This namespace contains methods used to register page routes and components.

#### `PluginApi.register.route`

Registers a route in the React Router.

| Parameter | Type | Description |
|-----------|------|-------------|
| `path` | `string` | The path to register. This should generally use the `/plugin/` prefix. |
| `component` | `React.FC` | A React function component that will be rendered when the route is loaded. |

Returns `void`.

#### `PluginApi.register.component`

Registers a component to be used by plugins. The component will be available in the `components` namespace.

| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | `string` | The name of the component to register. This should be unique and should ideally be prefixed with `plugin-`. |
| `component` | `React.FC` | A React function component. |

Returns `void`.

### `components`

This namespace contains all of the components available to plugins. These include a selection of core components and components registered using `PluginApi.register.component`.

### `utils`

This namespace provides access to the `NavUtils` and `StashService` namespaces. It also provides access to the `loadComponents` method.

#### `PluginApi.utils.loadComponents`

Due to code splitting, some components may not be loaded and available when a plugin page is rendered. `loadComponents` loads all of the components that a plugin page may require.

In general, `PluginApi.hooks.useLoadComponents` hook should be used instead.

| Parameter | Type | Description |
|-----------|------|-------------|
| `components` | `Promise[]` | The list of components to load. These values should come from the `PluginApi.loadableComponents` namespace. |

Returns a `Promise<void>` that resolves when all of the components have been loaded.

### `hooks`

This namespace provides access to the following core utility hooks:
- `useSpriteInfo`
- `useToast`

It also provides plugin-specific hooks.

#### `PluginApi.hooks.useLoadComponents`

This is a hook used to load components, using the `PluginApi.utils.loadComponents` method.

| Parameter | Type | Description |
|-----------|------|-------------|
| `components` | `Promise[]` | The list of components to load. These values should come from the `PluginApi.loadableComponents` namespace. |

Returns a `boolean` which will be `true` if the components are loading.

### `loadableComponents`

This namespace contains all of the components that may need to be loaded using the `loadComponents` method. Components are added to this namespace as needed. Please make a development request if a required component is not in this namespace.

### `patch`

This namespace provides methods to patch components to change their behaviour.

#### `PluginApi.patch.before`

Registers a before function. A before function is called prior to calling a component's render function. It accepts the same parameters as the component's render function, and is expected to return a list of new arguments that will be passed to the render.

| Parameter | Type | Description |
|-----------|------|-------------|
| `component` | `string` | The name of the component to patch. |
| `fn` | `Function` | The before function. It accepts the same arguments as the component render function and is expected to return a list of arguments to pass to the render function. |

Returns `void`.

#### `PluginApi.patch.instead`

Registers a replacement function for a component. The provided function will be called with the arguments passed to the original render function, plus the original render function as the last argument. An error will be thrown if the component already has a replacement function registered.

| Parameter | Type | Description |
|-----------|------|-------------|
| `component` | `string` | The name of the component to patch. |
| `fn` | `Function` | The replacement function. It accepts the same arguments as the original render function, plus the original render function, and is expected to return the replacement component. |

Returns `void`.

#### `PluginApi.patch.after`

Registers an after function. An after function is called after the render function of the component. It accepts the arguments passed to the original render function, plus the result of the original render function. It is expected to return the rendered component.

| Parameter | Type | Description |
|-----------|------|-------------|
| `component` | `string` | The name of the component to patch. |
| `fn` | `Function` | The after function. It accepts the same arguments as the original render function, plus the result of the original render function, and is expected to return the rendered component. |

Returns `void`.

#### Patchable components and functions

- `CountrySelect`
- `DateInput`
- `FolderSelect`
- `GalleryIDSelect`
- `GallerySelect`
- `GallerySelect.sort`
- `Icon`
- `MovieIDSelect`
- `MovieSelect`
- `MovieSelect.sort`
- `PerformerIDSelect`
- `PerformerSelect`
- `PerformerSelect.sort`
- `PluginRoutes`
- `SceneCard`
- `SceneCard.Details`
- `SceneCard.Image`
- `SceneCard.Overlays`
- `SceneCard.Popovers`
- `SceneIDSelect`
- `SceneSelect`
- `SceneSelect.sort`
- `Setting`
- `StudioIDSelect`
- `StudioSelect`
- `StudioSelect.sort`
- `TagIDSelect`
- `TagSelect`
- `TagSelect.sort`

### `PluginApi.Event`

Allows plugins to listen for Stash's events.

```js
PluginApi.Event.addEventListener("stash:location", (e) => console.log("Page Changed", e.detail.data.location.pathname))
```
