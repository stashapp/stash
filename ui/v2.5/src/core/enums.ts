import {
  ImageLightboxDisplayMode,
  ImageLightboxScrollMode,
} from "../core/generated-graphql";

export const imageLightboxDisplayModeIntlMap = new Map<
  ImageLightboxDisplayMode,
  string
>([
  [ImageLightboxDisplayMode.Original, "dialogs.lightbox.display_mode.original"],
  [
    ImageLightboxDisplayMode.FitXy,
    "dialogs.lightbox.display_mode.fit_to_screen",
  ],
  [
    ImageLightboxDisplayMode.FitX,
    "dialogs.lightbox.display_mode.fit_horizontally",
  ],
]);

export const imageLightboxScrollModeIntlMap = new Map<
  ImageLightboxScrollMode,
  string
>([
  [ImageLightboxScrollMode.Zoom, "dialogs.lightbox.scroll_mode.zoom"],
  [ImageLightboxScrollMode.PanY, "dialogs.lightbox.scroll_mode.pan_y"],
]);
