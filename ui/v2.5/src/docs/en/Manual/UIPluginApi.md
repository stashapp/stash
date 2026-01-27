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
- `FontAwesomeBrands`
- `Mousetrap`
- `MousetrapPause`
- `ReactSelect`

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

This namespace provides access to the `NavUtils` , `StashService` and `InteractiveUtils` namespaces. It also provides access to the `loadComponents` method.

#### `PluginApi.utils.loadComponents`

Due to code splitting, some components may not be loaded and available when a plugin page is rendered. `loadComponents` loads all of the components that a plugin page may require.

In general, `PluginApi.hooks.useLoadComponents` hook should be used instead.

| Parameter | Type | Description |
|-----------|------|-------------|
| `components` | `Promise[]` | The list of components to load. These values should come from the `PluginApi.loadableComponents` namespace. |

Returns a `Promise<void>` that resolves when all of the components have been loaded.

#### `PluginApi.utils.InteractiveUtils`
This namespace provides access to `interactiveClientProvider` and `getPlayer`
 - `getPlayer` returns the current `videojs` player object
 - `interactiveClientProvider` takes `IInteractiveClientProvider` which allows a developer to hook into the lifecycle of funscripts.
```ts
  export interface IDeviceSettings {
  connectionKey: string;
  scriptOffset: number;
  estimatedServerTimeOffset?: number;
  useStashHostedFunscript?: boolean;
  [key: string]: unknown;
}

export interface IInteractiveClientProviderOptions {
  handyKey: string;
  scriptOffset: number;
  defaultClientProvider?: IInteractiveClientProvider;
  stashConfig?: GQL.ConfigDataFragment;
}
export interface IInteractiveClientProvider {
  (options: IInteractiveClientProviderOptions): IInteractiveClient;
}

/**
 * Interface that is used for InteractiveProvider
 */
export interface IInteractiveClient {
  connect(): Promise<void>;
  handyKey: string;
  uploadScript: (funscriptPath: string, apiKey?: string) => Promise<void>;
  sync(): Promise<number>;
  configure(config: Partial<IDeviceSettings>): Promise<void>;
  play(position: number): Promise<void>;
  pause(): Promise<void>;
  ensurePlaying(position: number): Promise<void>;
  setLooping(looping: boolean): Promise<void>;
  readonly connected: boolean;
  readonly playing: boolean;
}

```
##### Example
For instance say I wanted to add extra logging when `IInteractiveClient.connect()` is called.
In my plugin you would install your own client provider as seen below

```ts
InteractiveUtils.interactiveClientProvider = (
  opts
) => {
  if (!opts.defaultClientProvider) {
    throw new Error('invalid setup');
  }

  const client = opts.defaultClientProvider(opts);
  const connect = client.connect;
  client.connect = async () => {
      console.log('patching connect method');
      return connect.call(client);
    };
   
  return client;
};

```


### `hooks`

This namespace provides access to the following core utility hooks:

- `useGalleryLightbox`
- `useLightbox`
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

This component also includes coarse-grained entries for every lazily loaded import in the stock UI. If a component is not available in `components` when the page loads, it can be loaded using the coarse-grained entry. For example, `PerformerCard` can be loaded using `loadableComponents.Performers`.

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

Registers a replacement function for a component. The provided function will be called with the arguments passed to the original render function, plus the next render function as the last argument. Replacement functions will be called in the order that they are registered. If a replacement function does not call the next render function then the following replacement functions will not be called or applied.

| Parameter | Type | Description |
|-----------|------|-------------|
| `component` | `string` | The name of the component to patch. |
| `fn` | `Function` | The replacement function. It accepts the same arguments as the original render function, plus the next render function, and is expected to return the replacement component. |

Returns `void`.

#### `PluginApi.patch.after`

Registers an after function. An after function is called after the render function of the component. It accepts the arguments passed to the original render function, plus the result of the original render function. It is expected to return the rendered component.

| Parameter | Type | Description |
|-----------|------|-------------|
| `component` | `string` | The name of the component to patch. |
| `fn` | `Function` | The after function. It accepts the same arguments as the original render function, plus the result of the original render function, and is expected to return the rendered component. |

Returns `void`.

#### Patchable components and functions

- `AlertModal`
- `App`
- `BackgroundImage`
- `BooleanSetting`
- `ChangeButtonSetting`
- `CompressedPerformerDetailsPanel`
- `ConstantSetting`
- `CountrySelect`
- `CustomFieldInput`
- `CustomFields`
- `CustomFieldsInput`
- `DateInput`
- `DetailImage`
- `ExternalLinkButtons`
- `ExternalLinksButton`
- `FolderSelect`
- `FrontPage`
- `GalleryCard`
- `GalleryCard.Details`
- `GalleryCard.Image`
- `GalleryCard.Overlays`
- `GalleryCard.Popovers`
- `GalleryIDSelect`
- `GallerySelect`
- `GallerySelect.sort`
- `GroupIDSelect`
- `GroupSelect`
- `GroupSelect.sort`
- `HeaderImage`
- `HoverPopover`
- `Icon`
- `ImageDetailPanel`
- `ImageInput`
- `LightboxLink`
- `LoadingIndicator`
- `MainNavBar.MenuItems`
- `MainNavBar.UtilityItems`
- `ModalSetting`
- `NumberSetting`
- `Pagination`
- `PaginationIndex`
- `PerformerAppearsWithPanel`
- `PerformerCard`
- `PerformerCard.Details`
- `PerformerCard.Image`
- `PerformerCard.Overlays`
- `PerformerCard.Popovers`
- `PerformerCard.Title`
- `PerformerDetailsPanel`
- `PerformerDetailsPanel.DetailGroup`
- `PerformerGalleriesPanel`
- `PerformerGroupsPanel`
- `PerformerHeaderImage`
- `PerformerIDSelect`
- `PerformerImagesPanel`
- `PerformerPage`
- `PerformerScenesPanel`
- `PerformerSelect`
- `PerformerSelect.sort`
- `PluginRoutes`
- `PluginSettings`
- `RatingNumber`
- `RatingStars`
- `RatingSystem`
- `SceneCard`
- `SceneCard.Details`
- `SceneCard.Image`
- `SceneCard.Overlays`
- `SceneCard.Popovers`
- `SceneFileInfoPanel`
- `SceneIDSelect`
- `ScenePage`
- `ScenePage.TabContent`
- `ScenePage.Tabs`
- `ScenePlayer`
- `SceneSelect`
- `SceneSelect.sort`
- `SelectSetting`
- `Setting`
- `SettingGroup`
- `SettingModal`
- `StringListSetting`
- `StringSetting`
- `StudioIDSelect`
- `StudioSelect`
- `StudioSelect.sort`
- `SweatDrops`
- `TabTitleCounter`
- `TagCard`
- `TagCard.Details`
- `TagCard.Image`
- `TagCard.Overlays`
- `TagCard.Popovers`
- `TagCard.Title`
- `TagIDSelect`
- `TagLink`
- `TagSelect`
- `TagSelect.sort`
- `TruncatedText`

### `PluginApi.Event`

Allows plugins to listen for Stash's events.

```js
PluginApi.Event.addEventListener("stash:location", (e) => console.log("Page Changed", e.detail.data.location.pathname))
```


