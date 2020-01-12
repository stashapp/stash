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

export class DvdsCriterion extends Criterion<IOptionType, ILabeledId[]> {
  public type: CriterionType = "dvds";
  public parameterName: string = "dvds";
  public modifier = CriterionModifier.Includes;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];
}

export class DvdsCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("dvds");
  public value: CriterionType = "dvds";
}
