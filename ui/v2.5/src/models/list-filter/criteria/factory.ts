/* eslint-disable consistent-return, default-case */
import {
  StringCriterion,
  NumberCriterion,
  DurationCriterion,
  NumberCriterionOption,
  MandatoryStringCriterionOption,
  MandatoryNumberCriterionOption,
  StringCriterionOption,
  ILabeledIdCriterion,
  BooleanCriterion,
  BooleanCriterionOption,
} from "./criterion";
import { OrganizedCriterion } from "./organized";
import { FavoriteCriterion, PerformerFavoriteCriterion } from "./favorite";
import { HasMarkersCriterion } from "./has-markers";
import {
  PerformerIsMissingCriterionOption,
  ImageIsMissingCriterionOption,
  TagIsMissingCriterionOption,
  SceneIsMissingCriterionOption,
  IsMissingCriterion,
  GalleryIsMissingCriterionOption,
  StudioIsMissingCriterionOption,
  MovieIsMissingCriterionOption,
} from "./is-missing";
import { NoneCriterion } from "./none";
import { PerformersCriterion } from "./performers";
import { AverageResolutionCriterion, ResolutionCriterion } from "./resolution";
import { StudiosCriterion, ParentStudiosCriterion } from "./studios";
import {
  ChildTagsCriterionOption,
  ParentTagsCriterionOption,
  PerformerTagsCriterionOption,
  SceneTagsCriterionOption,
  TagsCriterion,
  TagsCriterionOption,
} from "./tags";
import { GenderCriterion } from "./gender";
import { MoviesCriterionOption } from "./movies";
import { GalleriesCriterion } from "./galleries";
import { CriterionType } from "../types";
import { InteractiveCriterion } from "./interactive";
import { RatingCriterionOption } from "./rating";
import { DuplicatedCriterion, PhashCriterionOption } from "./phash";

export function makeCriteria(type: CriterionType = "none") {
  switch (type) {
    case "none":
      return new NoneCriterion();
    case "name":
    case "path":
    case "checksum":
      return new StringCriterion(
        new MandatoryStringCriterionOption(type, type)
      );
    case "oshash":
      return new StringCriterion(
        new MandatoryStringCriterionOption("media_info.hash", type, type)
      );
    case "rating":
      return new NumberCriterion(RatingCriterionOption);
    case "organized":
      return new OrganizedCriterion();
    case "o_counter":
    case "interactive_speed":
    case "scene_count":
    case "marker_count":
    case "image_count":
    case "gallery_count":
    case "performer_count":
    case "performer_age":
    case "tag_count":
      return new NumberCriterion(
        new MandatoryNumberCriterionOption(type, type)
      );
    case "resolution":
      return new ResolutionCriterion();
    case "average_resolution":
      return new AverageResolutionCriterion();
    case "duration":
      return new DurationCriterion(new NumberCriterionOption(type, type));
    case "favorite":
      return new FavoriteCriterion();
    case "hasMarkers":
      return new HasMarkersCriterion();
    case "sceneIsMissing":
      return new IsMissingCriterion(SceneIsMissingCriterionOption);
    case "imageIsMissing":
      return new IsMissingCriterion(ImageIsMissingCriterionOption);
    case "performerIsMissing":
      return new IsMissingCriterion(PerformerIsMissingCriterionOption);
    case "galleryIsMissing":
      return new IsMissingCriterion(GalleryIsMissingCriterionOption);
    case "tagIsMissing":
      return new IsMissingCriterion(TagIsMissingCriterionOption);
    case "studioIsMissing":
      return new IsMissingCriterion(StudioIsMissingCriterionOption);
    case "movieIsMissing":
      return new IsMissingCriterion(MovieIsMissingCriterionOption);
    case "tags":
      return new TagsCriterion(TagsCriterionOption);
    case "sceneTags":
      return new TagsCriterion(SceneTagsCriterionOption);
    case "performerTags":
      return new TagsCriterion(PerformerTagsCriterionOption);
    case "parentTags":
      return new TagsCriterion(ParentTagsCriterionOption);
    case "childTags":
      return new TagsCriterion(ChildTagsCriterionOption);
    case "performers":
      return new PerformersCriterion();
    case "performer_favorite":
      return new PerformerFavoriteCriterion();
    case "studios":
      return new StudiosCriterion();
    case "parent_studios":
      return new ParentStudiosCriterion();
    case "movies":
      return new ILabeledIdCriterion(MoviesCriterionOption);
    case "galleries":
      return new GalleriesCriterion();
    case "birth_year":
    case "death_year":
    case "weight":
      return new NumberCriterion(new NumberCriterionOption(type, type));
    case "age":
      return new NumberCriterion(
        new MandatoryNumberCriterionOption(type, type)
      );
    case "gender":
      return new GenderCriterion();
    case "sceneChecksum":
    case "galleryChecksum":
      return new StringCriterion(
        new StringCriterionOption("checksum", type, "checksum")
      );
    case "phash":
      return new StringCriterion(PhashCriterionOption);
    case "duplicated":
      return new DuplicatedCriterion();
    case "ethnicity":
    case "country":
    case "hair_color":
    case "eye_color":
    case "height":
    case "measurements":
    case "fake_tits":
    case "career_length":
    case "tattoos":
    case "piercings":
    case "aliases":
    case "url":
    case "stash_id":
    case "details":
    case "title":
    case "director":
    case "synopsis":
      return new StringCriterion(new StringCriterionOption(type, type));
    case "interactive":
      return new InteractiveCriterion();
    case "parent_tag_count":
      return new NumberCriterion(
        new MandatoryNumberCriterionOption(
          "parent_tag_count",
          "parent_tag_count",
          "parent_count"
        )
      );
    case "child_tag_count":
      return new NumberCriterion(
        new MandatoryNumberCriterionOption(
          "sub_tag_count",
          "child_tag_count",
          "child_count"
        )
      );
    case "ignore_auto_tag":
      return new BooleanCriterion(new BooleanCriterionOption(type, type));
  }
}
