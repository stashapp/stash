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

export class PerformersCriterion extends Criterion<IOptionType, ILabeledId[]> {
  public type: CriterionType = "performers";
  public parameterName: string = "performers";
  public modifier = CriterionModifier.Equals;
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];
}

export class PerformersCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performers");
  public value: CriterionType = "performers";
}
