/* eslint-disable consistent-return, default-case */
import { CriterionModifier } from "src/core/generated-graphql";
import {
  Criterion,
  CriterionType,
  StringCriterion,
  NumberCriterion,
  DurationCriterion,
  MandatoryStringCriterion,
} from "./criterion";
import { OrganizedCriterion } from "./organized";
import { FavoriteCriterion } from "./favorite";
import { HasMarkersCriterion } from "./has-markers";
import {
  PerformerIsMissingCriterion,
  SceneIsMissingCriterion,
  GalleryIsMissingCriterion,
  TagIsMissingCriterion,
  StudioIsMissingCriterion,
  MovieIsMissingCriterion,
  ImageIsMissingCriterion,
} from "./is-missing";
import { NoneCriterion } from "./none";
import { PerformersCriterion } from "./performers";
import { RatingCriterion } from "./rating";
import { AverageResolutionCriterion, ResolutionCriterion } from "./resolution";
import { StudiosCriterion, ParentStudiosCriterion } from "./studios";
import { TagsCriterion } from "./tags";
import { GenderCriterion } from "./gender";
import { MoviesCriterion } from "./movies";
import { GalleriesCriterion } from "./galleries";

export function makeCriteria(type: CriterionType = "none") {
  switch (type) {
    case "none":
      return new NoneCriterion();
    case "path":
      return new MandatoryStringCriterion(type, type);
    case "rating":
      return new RatingCriterion();
    case "organized":
      return new OrganizedCriterion();
    case "o_counter":
    case "scene_count":
    case "marker_count":
    case "image_count":
    case "gallery_count":
    case "performer_count":
      return new NumberCriterion(type, type);
    case "resolution":
      return new ResolutionCriterion();
    case "average_resolution":
      return new AverageResolutionCriterion();
    case "duration":
      return new DurationCriterion(type, type);
    case "favorite":
      return new FavoriteCriterion();
    case "hasMarkers":
      return new HasMarkersCriterion();
    case "sceneIsMissing":
      return new SceneIsMissingCriterion();
    case "imageIsMissing":
      return new ImageIsMissingCriterion();
    case "performerIsMissing":
      return new PerformerIsMissingCriterion();
    case "galleryIsMissing":
      return new GalleryIsMissingCriterion();
    case "tagIsMissing":
      return new TagIsMissingCriterion();
    case "studioIsMissing":
      return new StudioIsMissingCriterion();
    case "movieIsMissing":
      return new MovieIsMissingCriterion();
    case "tags":
      return new TagsCriterion("tags");
    case "sceneTags":
      return new TagsCriterion("sceneTags");
    case "performerTags":
      return new TagsCriterion("performerTags");
    case "performers":
      return new PerformersCriterion();
    case "studios":
      return new StudiosCriterion();
    case "parent_studios":
      return new ParentStudiosCriterion();
    case "movies":
      return new MoviesCriterion();
    case "galleries":
      return new GalleriesCriterion();
    case "birth_year":
      return new NumberCriterion(type, type);
    case "age": {
      const ret = new NumberCriterion(type, type);
      // null/not null doesn't make sense for these criteria
      ret.modifierOptions = [
        Criterion.getModifierOption(CriterionModifier.Equals),
        Criterion.getModifierOption(CriterionModifier.NotEquals),
        Criterion.getModifierOption(CriterionModifier.GreaterThan),
        Criterion.getModifierOption(CriterionModifier.LessThan),
      ];
      return ret;
    }
    case "gender":
      return new GenderCriterion();
    case "ethnicity":
    case "country":
    case "eye_color":
    case "height":
    case "measurements":
    case "fake_tits":
    case "career_length":
    case "tattoos":
    case "piercings":
    case "aliases":
      return new StringCriterion(type, type);
  }
}
