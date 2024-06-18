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

class BaseTagsCriterionOption extends CriterionOption {
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
      inputType,
      makeCriterion: () => new TagsCriterion(this),
    });
  }
}

export const TagsCriterionOption = new BaseTagsCriterionOption(
  "tags",
  "tags",
  defaultModifierOptions
);

export const SceneTagsCriterionOption = new BaseTagsCriterionOption(
  "scene_tags",
  "scene_tags",
  defaultModifierOptions
);

export const PerformerTagsCriterionOption = new BaseTagsCriterionOption(
  "performer_tags",
  "performer_tags",
  withoutEqualsModifierOptions
);

// TODO - this requires using a nested studios_filter which needs to be added separately
// export const StudioTagsCriterionOption = new BaseTagsCriterionOption(
//   "studio_tags",
//   "studio_tags",
//   withoutEqualsModifierOptions
// );

export const ParentTagsCriterionOption = new BaseTagsCriterionOption(
  "parent_tags",
  "parents",
  withoutEqualsModifierOptions
);

export const ChildTagsCriterionOption = new BaseTagsCriterionOption(
  "sub_tags",
  "children",
  withoutEqualsModifierOptions
);

export class TagsCriterion extends IHierarchicalLabeledIdCriterion {}
