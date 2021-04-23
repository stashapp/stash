import { CriterionModifier } from "src/core/generated-graphql";
import { ILabeledId, IOptionType, encodeILabeledId } from "../types";
import {
  Criterion,
  CriterionOption,
  CriterionType,
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
  public type: CriterionType = "parent_studios";
  public parameterName: string = "parents";

  constructor() {
    super(new ParentStudiosCriterionOption());
  }
}
