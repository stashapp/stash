import * as GQL from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class TagsCriterion extends Criterion {
  public type: CriterionType;
  public parameterName: string;
  public modifier = GQL.CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(GQL.CriterionModifier.IncludesAll),
    Criterion.getModifierOption(GQL.CriterionModifier.Includes),
    Criterion.getModifierOption(GQL.CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];

  constructor(type: "tags" | "sceneTags" | "performerTags") {
    super();
    this.type = type;
    this.parameterName = type;
    if (type === "sceneTags") {
      this.parameterName = "scene_tags";
    }
    if (type === "performerTags") {
      this.parameterName = "performer_tags";
    }
  }

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
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

export class PerformerTagsCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performerTags");
  public value: CriterionType = "performerTags";
}
