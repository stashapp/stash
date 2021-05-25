import { CriterionOption, ILabeledIdCriterion } from "./criterion";

export class TagsCriterion extends ILabeledIdCriterion {
  constructor(type: CriterionOption) {
    super(type, true);
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
