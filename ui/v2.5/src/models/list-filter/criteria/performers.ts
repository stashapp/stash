import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import { Criterion, CriterionOption, ILabeledIdCriterion } from "./criterion";

export const PerformersCriterionOption = new CriterionOption(
  "performers",
  "performers"
);

export class PerformersCriterion extends ILabeledIdCriterion {
  public modifier = CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.IncludesAll),
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];

  constructor() {
    super(PerformersCriterionOption);
  }

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
  }
}
