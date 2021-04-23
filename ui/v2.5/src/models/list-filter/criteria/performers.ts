import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import {
  Criterion,
  CriterionOption,
  ILabeledIdCriterion,
} from "./criterion";

export class PerformersCriterionOption extends CriterionOption {
  constructor() {
    super("performers", "performers");
  }
}

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
    super(new PerformersCriterionOption());
  }

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
  }
}
