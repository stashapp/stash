import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, encodeILabeledId } from "../types";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
  ILabeledIdCriterion,
} from "./criterion";

interface IOptionType {
  id: string;
  name?: string;
  image_path?: string;
}

export class MoviesCriterion extends ILabeledIdCriterion {
  public type: CriterionType = "movies";
  public parameterName: string = "movies";
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

export class MoviesCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("movies");
  public value: CriterionType = "movies";
}
