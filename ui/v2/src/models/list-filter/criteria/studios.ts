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

export class StudiosCriterion extends Criterion<IOptionType, ILabeledId[]> {
  public type: CriterionType = "studios";
  public parameterName: string = "studios";
  public modifier = CriterionModifier.Includes;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];
}

export class StudiosCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("studios");
  public value: CriterionType = "studios";
}
