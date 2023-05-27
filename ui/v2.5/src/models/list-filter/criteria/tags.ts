import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption, IHierarchicalLabeledIdCriterion } from "./criterion";

// TODO: should be modifierOptions for all CriterionOptions below but
// the Equals modifier only works for the "tags" and "sceneTags" criterion.
// Replace all modifierOptions below when the Equals modifier works
// everywhere and add the NotEquals modifier here when that works.
const modifierOptions = [
  CriterionModifier.IncludesAll,
  CriterionModifier.Includes,
  CriterionModifier.Equals,
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
];

const defaultModifier = CriterionModifier.IncludesAll;

export const TagsCriterionOption = new CriterionOption({
  messageID: "tags",
  type: "tags",
  parameterName: "tags",
  modifierOptions,
  defaultModifier,
});
export const SceneTagsCriterionOption = new CriterionOption({
  messageID: "sceneTags",
  type: "sceneTags",
  parameterName: "scene_tags",
  modifierOptions,
  defaultModifier,
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
  defaultModifier,
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
  defaultModifier,
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
  defaultModifier,
});

export class TagsCriterion extends IHierarchicalLabeledIdCriterion {}
