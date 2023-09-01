import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption, IHierarchicalLabeledIdCriterion } from "./criterion";
import { CriterionType } from "../types";

const defaultModifierOptions = [
  CriterionModifier.IncludesAll,
  CriterionModifier.Includes,
  CriterionModifier.Equals,
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
];

const withoutEqualsModifierOptions = [
  CriterionModifier.IncludesAll,
  CriterionModifier.Includes,
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
];

const defaultModifier = CriterionModifier.IncludesAll;
const inputType = "tags";

export class TagsCriterionOptionClass extends CriterionOption {
  constructor(
    messageID: string,
    type: CriterionType,
    modifierOptions: CriterionModifier[]
  ) {
    super({
      messageID,
      type,
      modifierOptions,
      defaultModifier,
      makeCriterion: () => new TagsCriterion(this),
      inputType,
    });
  }
}

export const TagsCriterionOption = new TagsCriterionOptionClass(
  "tags",
  "tags",
  defaultModifierOptions
);

export const SceneTagsCriterionOption = new TagsCriterionOptionClass(
  "scene_tags",
  "scene_tags",
  defaultModifierOptions
);

export const PerformerTagsCriterionOption = new TagsCriterionOptionClass(
  "performer_tags",
  "performer_tags",
  withoutEqualsModifierOptions
);

export const ParentTagsCriterionOption = new TagsCriterionOptionClass(
  "parent_tags",
  "parents",
  withoutEqualsModifierOptions
);

export const ChildTagsCriterionOption = new TagsCriterionOptionClass(
  "sub_tags",
  "children",
  withoutEqualsModifierOptions
);

export class TagsCriterion extends IHierarchicalLabeledIdCriterion {}
