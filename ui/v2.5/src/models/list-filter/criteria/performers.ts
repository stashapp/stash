import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
  ILabeledIdCriterion,
} from "./criterion";

export class PerformersCriterion extends ILabeledIdCriterion {
  public type: CriterionType = "performers";
  public parameterName: string = "performers";
  public modifier = CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.IncludesAll),
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

export class PerformersCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performers");
  public value: CriterionType = "performers";
}
