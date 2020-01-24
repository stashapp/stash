import * as GQL from "src/core/generated-graphql";
import { ILabeledId } from "../types";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class TagsCriterion extends Criterion<GQL.Tag, ILabeledId[]> {
  public type: CriterionType;
  public parameterName: string;
  public modifier = GQL.CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(GQL.CriterionModifier.IncludesAll),
    Criterion.getModifierOption(GQL.CriterionModifier.Includes),
    Criterion.getModifierOption(GQL.CriterionModifier.Excludes)
  ];
  public options: GQL.Tag[] = [];
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
