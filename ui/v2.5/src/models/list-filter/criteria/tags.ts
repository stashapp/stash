import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

export class TagsCriterion extends ILabeledIdCriterion {}

export const TagsCriterionOption = new ILabeledIdCriterionOption(
  "tags",
  "tags",
  "tags",
  true
);
export const SceneTagsCriterionOption = new ILabeledIdCriterionOption(
  "sceneTags",
  "sceneTags",
  "scene_tags",
  true
);
export const PerformerTagsCriterionOption = new ILabeledIdCriterionOption(
  "performerTags",
  "performerTags",
  "performer_tags",
  true
);
