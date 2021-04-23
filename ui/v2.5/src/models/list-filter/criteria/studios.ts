import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import {
  Criterion,
  CriterionOption,
  ILabeledIdCriterion,
} from "./criterion";

abstract class AbstractStudiosCriterion extends ILabeledIdCriterion {
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

export class StudiosCriterionOption extends CriterionOption {
  constructor() {
    super("studios", "studios");
  }
}
export class StudiosCriterion extends AbstractStudiosCriterion {
  constructor() {
    super(new StudiosCriterionOption());
  }
}

export class ParentStudiosCriterionOption extends CriterionOption {
  constructor() {
    super("parent_studios", "parent_studios");
  }
}
export class ParentStudiosCriterion extends AbstractStudiosCriterion {
  constructor() {
    super(new ParentStudiosCriterionOption());
  }
}
