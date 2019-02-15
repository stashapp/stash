import {
  Criterion,
  CriterionModifier,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class RatingCriterion extends Criterion<number, number> { // TODO <number, number[]>
  public type: CriterionType = "rating";
  public parameterName: string = "rating";
  public modifier = CriterionModifier.Equals;
  public options: number[] = [1, 2, 3, 4, 5];
  public value: number = 0;
}

export class RatingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("rating");
  public value: CriterionType = "rating";
}
