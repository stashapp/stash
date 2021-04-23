import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, encodeILabeledId } from "../types";
import { Criterion, CriterionOption, ILabeledIdCriterion } from "./criterion";

interface IOptionType {
  id: string;
  name?: string;
  image_path?: string;
}

export class MoviesCriterionOption extends CriterionOption {
  constructor() {
    super("movies", "movies");
  }
}

export class MoviesCriterion extends ILabeledIdCriterion {
  public modifier = CriterionModifier.Includes;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];
  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];

  constructor() {
    super(new MoviesCriterionOption());
  }

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
  }
}
