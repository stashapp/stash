import * as GQL from "../../../core/generated-graphql";
import { ILabeledId } from "../types";
import {
  Criterion,
  CriterionModifier,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class TagsCriterion extends Criterion<GQL.AllTagsForFilterAllTags, ILabeledId[]> {
  public type: CriterionType;
  public parameterName: string;
  public modifier = CriterionModifier.Equals;
  public options: GQL.AllTagsForFilterAllTags[] = [];
  public value: ILabeledId[] = [];

  constructor(type: "tags" | "sceneTags") {
    super();
    this.type = type;
    this.parameterName = type;
    if (type === "sceneTags") {
      this.parameterName = "scene_tags";
    }
  }
}

export class TagsCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("tags");
  public value: CriterionType = "tags";
}

export class SceneTagsCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("sceneTags");
  public value: CriterionType = "sceneTags";
}
