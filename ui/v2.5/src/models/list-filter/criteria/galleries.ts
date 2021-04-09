import * as GQL from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
  ILabeledIdCriterion,
} from "./criterion";

export class GalleriesCriterion extends ILabeledIdCriterion {
  public type: CriterionType = "galleries";
  public parameterName: string = "galleries";
  public modifier = GQL.CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(GQL.CriterionModifier.IncludesAll),
    Criterion.getModifierOption(GQL.CriterionModifier.Includes),
    Criterion.getModifierOption(GQL.CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
  }
}

export class GalleriesCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("galleries");
  public value: CriterionType = "galleries";
}
