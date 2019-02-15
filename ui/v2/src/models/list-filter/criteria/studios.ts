import { ILabeledId } from "../types";
import {
  Criterion,
  CriterionModifier,
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
  public modifier = CriterionModifier.Equals;
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];
}

export class StudiosCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("studios");
  public value: CriterionType = "studios";
}
