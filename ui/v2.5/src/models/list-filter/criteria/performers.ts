import { CriterionModifier } from "src/core/generated-graphql";
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

export class PerformersCriterion extends Criterion<IOptionType, ILabeledId[]> {
  public type: CriterionType = "performers";
  public parameterName: string = "performers";
  public modifier = CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.IncludesAll),
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];
}

export class PerformersCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performers");
  public value: CriterionType = "performers";
}
