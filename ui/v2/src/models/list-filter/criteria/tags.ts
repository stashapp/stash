import * as GQL from "../../../core/generated-graphql";
import { CriterionModifier } from "../../../core/generated-graphql";
import { ILabeledId, encodeILabeledId } from "../types";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class TagsCriterion extends Criterion<GQL.AllTagsForFilterAllTags, ILabeledId[]> {
  public type: CriterionType;
  public parameterName: string;
  public modifier = CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.IncludesAll),
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];
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

  public encodeValue() {
    return this.value.map((o) => { return encodeILabeledId(o); });
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
