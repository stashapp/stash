import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption, IHierarchicalLabeledIdCriterion } from "./criterion";

export class TagsCriterion extends IHierarchicalLabeledIdCriterion {}

export const TagsCriterionOption = new CriterionOption({
  messageID: "tags",
  type: "tags",
  parameterName: "tags",
  modifierOptions: [
    CriterionModifier.IncludesAll,
    CriterionModifier.Includes,
    CriterionModifier.Equals,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  defaultModifier: CriterionModifier.IncludesAll,
});
export const SceneTagsCriterionOption = new CriterionOption({
  messageID: "sceneTags",
  type: "sceneTags",
  parameterName: "scene_tags",
  modifierOptions: [
    CriterionModifier.IncludesAll,
    CriterionModifier.Includes,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  defaultModifier: CriterionModifier.IncludesAll,
});
export const PerformerTagsCriterionOption = new CriterionOption({
  messageID: "performerTags",
  type: "performerTags",
  parameterName: "performer_tags",
  modifierOptions: [
    CriterionModifier.IncludesAll,
    CriterionModifier.Includes,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  defaultModifier: CriterionModifier.IncludesAll,
});
export const ParentTagsCriterionOption = new CriterionOption({
  messageID: "parent_tags",
  type: "parentTags",
  parameterName: "parents",
  modifierOptions: [
    CriterionModifier.IncludesAll,
    CriterionModifier.Includes,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  defaultModifier: CriterionModifier.IncludesAll,
});
export const ChildTagsCriterionOption = new CriterionOption({
  messageID: "sub_tags",
  type: "childTags",
  parameterName: "children",
  modifierOptions: [
    CriterionModifier.IncludesAll,
    CriterionModifier.Includes,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  defaultModifier: CriterionModifier.IncludesAll,
});
