import * as GQL from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import { Criterion, CriterionOption, ILabeledIdCriterion } from "./criterion";

const galleriesCriterionOption = new CriterionOption("galleries", "galleries");

export class GalleriesCriterion extends ILabeledIdCriterion {
  public modifier = GQL.CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(GQL.CriterionModifier.IncludesAll),
    Criterion.getModifierOption(GQL.CriterionModifier.Includes),
    Criterion.getModifierOption(GQL.CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];

  constructor() {
    super(galleriesCriterionOption);
  }

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
  }
}
