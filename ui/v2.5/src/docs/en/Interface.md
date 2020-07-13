# Interface Options

## Language

Setting the language affects the formatting of numbers and dates.

## Scene/Marker Wall Preview Type

The Scene Wall and Marker pages display scene preview videos by default. This can be changed to animated image (webp) or static image. 

> **⚠️ Note:** scene/marker preview videos must be generated to see them in the applicable wall page if Video preview type is selected. Likewise, if Animated Image is selected, then Image Previews must be generated.

## Show Studios as text

By default, a scene's studio will be shown as an image overlay. Checking this option changes this to display studios as a text name instead.

## Scene Player options

By default, scene videos do not automatically start when navigating to the scenes page. Checking the "Auto-start video" option changes this to auto play scene videos.

The maximum loop duration option allows looping of shorter videos. Set this value to the maximum scene duration that scene videos should loop. Setting this to 0 disables this functionality.

## Custom CSS

The stash UI can be customised using custom CSS. See [here](https://github.com/stashapp/stash/wiki/Custom-CSS-snippets) for a community-curated set of CSS snippets to customise your UI. 

[Stash Plex Theme](https://github.com/stashapp/stash/wiki/Stash-Plex-Theme) is a community created theme inspired by the popular Plex interface.


## Custom served folders

It is possible to expose specific folders to the UI. This configuration is performed manually in the `config.yml` file only.

Custom served content is exposed via the `/custom` URL path prefix.

For example, in the `config.yml` file:
```
custom_served_folders:
  /: D:\stash\static
  /foo: D:\bar
```

With the above configuration, a request for `/custom/foo/bar.png` would return `D:\bar\bar.png`. The `/` entry matches anything that is not otherwise mapped by the other entries. For example, `/custom/baz/xyz.png` would return `D:\stash\static\baz\xyz.png`.

Applications for this include using static images in custom css, like the Plex theme. For example, using the following config:
```yml
custom_served_folders:
  /: <stash folder>\custom
```

The `background.png` and `noise.png` files can be placed in the `custom` folder, then in the custom css, the `./background.png` and `./noise.png` strings can be replaced with `/custom/background.png` and `/custom/noise.png` respectively.

Other applications are to add custom UIs to stash, accessible via `/custom`.