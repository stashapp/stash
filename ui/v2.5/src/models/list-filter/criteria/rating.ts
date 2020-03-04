import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class RatingCriterion extends Criterion {
  public type: CriterionType = "rating";
  public parameterName: string = "rating";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Equals),
    Criterion.getModifierOption(CriterionModifier.NotEquals),
    Criterion.getModifierOption(CriterionModifier.GreaterThan),
    Criterion.getModifierOption(CriterionModifier.LessThan),
    Criterion.getModifierOption(CriterionModifier.IsNull),
    Criterion.getModifierOption(CriterionModifier.NotNull)
  ];
  public options: number[] = [1, 2, 3, 4, 5];
  public value: number = 0;
}

export class RatingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("rating");
  public value: CriterionType = "rating";
}
