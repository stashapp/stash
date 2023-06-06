import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionType } from "../types";
import { CriterionOption, IHierarchicalLabeledIdCriterion } from "./criterion";

export class TagsCriterion extends IHierarchicalLabeledIdCriterion {}

const tagsModifierOptions = [
  CriterionModifier.Includes,
  CriterionModifier.IncludesAll,
  CriterionModifier.Equals,
];

const withoutEqualsModifierOptions = [
  CriterionModifier.Includes,
  CriterionModifier.IncludesAll,
];

class tagsCriterionOption extends CriterionOption {
  constructor(
    messageID: string,
    value: CriterionType,
    parameterName: string,
    modifierOptions: CriterionModifier[]
  ) {
    let defaultModifier = CriterionModifier.IncludesAll;

    super({
      messageID,
      type: value,
      parameterName,
      modifierOptions,
      defaultModifier,
    });
  }
}

export const TagsCriterionOption = new tagsCriterionOption(
  "tags",
  "tags",
  "tags",
  tagsModifierOptions
);
export const SceneTagsCriterionOption = new tagsCriterionOption(
  "sceneTags",
  "sceneTags",
  "scene_tags",
  tagsModifierOptions
);
export const PerformerTagsCriterionOption = new tagsCriterionOption(
  "performerTags",
  "performerTags",
  "performer_tags",
  withoutEqualsModifierOptions
);
export const ParentTagsCriterionOption = new tagsCriterionOption(
  "parent_tags",
  "parentTags",
  "parents",
  withoutEqualsModifierOptions
);
export const ChildTagsCriterionOption = new tagsCriterionOption(
  "sub_tags",
  "childTags",
  "children",
  withoutEqualsModifierOptions
);
