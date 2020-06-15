import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class StudiosCriterion extends Criterion {
  public type: CriterionType = "studios";
  public parameterName: string = "studios";
  public modifier = CriterionModifier.Includes;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
  }
}

export class StudiosCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("studios");
  public value: CriterionType = "studios";
}

export class ParentStudiosCriterion extends StudiosCriterion {
  public type: CriterionType = "parent_studios";
  public parameterName: string = "parents";
}

export class ParentStudiosCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("parent_studios");
  public value: CriterionType = "parent_studios";
}
