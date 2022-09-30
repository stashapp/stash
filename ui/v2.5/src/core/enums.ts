import {
  ImageLightboxDisplayMode,
  ImageLightboxScrollMode,
  RatingSystem,
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

export const ratingSystemIntlMap = new Map<RatingSystem, string>([
  [
    RatingSystem.FivePointFiveStar,
    "config.ui.editing.rating_system.options.five_pointfive_stars",
  ],
  [RatingSystem.FiveStar, "config.ui.editing.rating_system.options.five_stars"],
  [
    RatingSystem.FivePointTwoFiveStar,
    "config.ui.editing.rating_system.options.five_pointtwofive_stars",
  ],
  [RatingSystem.TenStar, "config.ui.editing.rating_system.options.ten_stars"],
  [
    RatingSystem.TenPointFiveStar,
    "config.ui.editing.rating_system.options.ten_pointfive_stars",
  ],
  [
    RatingSystem.TenPointTwoFiveStar,
    "config.ui.editing.rating_system.options.ten_pointtwofive_stars",
  ],
  [
    RatingSystem.TenPointDecimal,
    "config.ui.editing.rating_system.options.ten_point_decimal",
  ],
  [RatingSystem.None, "config.ui.editing.rating_system.options.none"],
]);
