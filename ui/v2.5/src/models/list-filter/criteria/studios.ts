import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import { Criterion, CriterionOption, ILabeledIdCriterion } from "./criterion";

abstract class AbstractStudiosCriterion extends ILabeledIdCriterion {
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

export const StudiosCriterionOption = new CriterionOption("studios", "studios");
export class StudiosCriterion extends AbstractStudiosCriterion {
  constructor() {
    super(StudiosCriterionOption);
  }
}

export const ParentStudiosCriterionOption = new CriterionOption(
  "parent_studios",
  "parent_studios"
);
export class ParentStudiosCriterion extends AbstractStudiosCriterion {
  constructor() {
    super(ParentStudiosCriterionOption);
  }
}
