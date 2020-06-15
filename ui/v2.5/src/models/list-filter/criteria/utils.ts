/* eslint-disable consistent-return, default-case */
import { CriterionModifier } from "src/core/generated-graphql";
import {
  Criterion,
  CriterionType,
  StringCriterion,
  NumberCriterion,
  DurationCriterion,
} from "./criterion";
import { FavoriteCriterion } from "./favorite";
import { HasMarkersCriterion } from "./has-markers";
import {
  PerformerIsMissingCriterion,
  SceneIsMissingCriterion,
} from "./is-missing";
import { NoneCriterion } from "./none";
import { PerformersCriterion } from "./performers";
import { RatingCriterion } from "./rating";
import { ResolutionCriterion } from "./resolution";
import { StudiosCriterion, ParentStudiosCriterion } from "./studios";
import { TagsCriterion } from "./tags";
import { GenderCriterion } from "./gender";
import { MoviesCriterion } from "./movies";

export function makeCriteria(type: CriterionType = "none") {
  switch (type) {
    case "none":
      return new NoneCriterion();
    case "rating":
      return new RatingCriterion();
    case "o_counter":
      return new NumberCriterion(type, type);
    case "resolution":
      return new ResolutionCriterion();
    case "duration":
      return new DurationCriterion(type, type);
    case "favorite":
      return new FavoriteCriterion();
    case "hasMarkers":
      return new HasMarkersCriterion();
    case "sceneIsMissing":
      return new SceneIsMissingCriterion();
    case "performerIsMissing":
      return new PerformerIsMissingCriterion();
    case "tags":
      return new TagsCriterion("tags");
    case "sceneTags":
      return new TagsCriterion("sceneTags");
    case "performers":
      return new PerformersCriterion();
    case "studios":
      return new StudiosCriterion();
    case "parent_studios":
      return new ParentStudiosCriterion();
    case "movies":
      return new MoviesCriterion();
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
