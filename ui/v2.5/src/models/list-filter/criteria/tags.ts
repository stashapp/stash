import * as GQL from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import { Criterion, CriterionOption, ILabeledIdCriterion } from "./criterion";

export class TagsCriterion extends ILabeledIdCriterion {
  public modifier = GQL.CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(GQL.CriterionModifier.IncludesAll),
    Criterion.getModifierOption(GQL.CriterionModifier.Includes),
    Criterion.getModifierOption(GQL.CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
  }
}

export const TagsCriterionOption = new CriterionOption("tags", "tags");
export const SceneTagsCriterionOption = new CriterionOption(
  "sceneTags",
  "sceneTags",
  "scene_tags"
);
export const PerformerTagsCriterionOption = new CriterionOption(
  "performerTags",
  "performerTags",
  "performer_tags"
);
