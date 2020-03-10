import { CriterionModifier } from "../../../core/generated-graphql";
import { ILabeledId } from "../types";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
} from "./criterion";

interface IOptionType {
  id: string;
  name?: string;
  image_path?: string;
}

export class MoviesCriterion extends Criterion<IOptionType, ILabeledId[]> {
  public type: CriterionType = "movies";
  public parameterName: string = "movies";
  public modifier = CriterionModifier.Includes;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];
}

export class MoviesCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("movies");
  public value: CriterionType = "movies";
}
