# Interface Options

## Language

Setting the language affects the formatting of numbers and dates.

## SFW Content Mode

SFW Content Mode is used to indicate that the content being managed is _not_ adult content. 

When SFW Content Mode is enabled, the following changes are made to the UI:
- default performer images are changed to less adult-oriented images
- certain adult-specific metadata fields are hidden (e.g. performer genital fields)
- `O`-Counter is replaced with `Like`-counter

## Scene/Marker Wall Preview Type

The Scene Wall and Marker pages display scene preview videos (mp4) by default. This can be changed to animated image (webp) or static image. 

> **⚠️ Note:** scene/marker preview videos must be generated to see them in the applicable wall page if Video preview type is selected. Likewise, if Animated Image is selected, then Image Previews must be generated.

## Show Studios as text

By default, a scene's studio will be shown as an image overlay. Checking this option changes this to display studios as a text name instead.

## Scene Player options

By default, scene videos do not automatically start when navigating to the scenes page. Checking the "Auto-start video" option changes this to auto play scene videos.

The maximum loop duration option allows looping of shorter videos. Set this value to the maximum scene duration that scene videos should loop. Setting this to 0 disables this functionality.

### Activity tracking

The "Track Activity" option allows tracking of scene play count and duration, and sets the resume point when a scene video is not finished.

The "Minimum Play Percent" gives the minimum proportion of a video that must be played before the play count of the scene is incremented.

By default, when a scene has a resume point, the scene player will automatically seek to this point when the scene is played. Setting "Always start video from beginning" to true disables this behaviour.

## Custom CSS

The stash UI can be customised using custom CSS. See [here](https://docs.stashapp.cc/themes/custom-css-snippets/) for a community-curated set of CSS snippets to customise your UI. 

There is also a [collection of community-created themes](https://docs.stashapp.cc/themes/list/#browse-themes) available.

## Custom Javascript

Stash supports the injection of custom javascript to assist with theming or adding additional functionality. Be aware that bad Javascript could break the UI or worse.

## Custom Locales

The localisation strings can be customised. The master list of default (en-GB) locale strings can be found [here](https://github.com/stashapp/stash/blob/develop/ui/v2.5/src/locales/en-GB.json). The custom locale format is the same as this json file.

For example, to override the `actions.add_directory` label (which is `Add Directory` by default), you would have the following in the custom locale:

```
{
  "actions": {
    "add_directory": "Some other description"
  }
}
```

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
