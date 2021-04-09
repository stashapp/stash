import { CriterionModifier } from "src/core/generated-graphql";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
  StringCriterion,
} from "./criterion";

export class ResolutionCriterion extends StringCriterion {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];

  constructor() {
    super("resolution", "resolution", [
      "144p",
      "240p",
      "360p",
      "480p",
      "540p",
      "720p",
      "1080p",
      "1440p",
      "4k",
      "5k",
      "6k",
      "8k",
    ]);
  }
}

export class ResolutionCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("resolution");
  public value: CriterionType = "resolution";
}

export class AverageResolutionCriterion extends ResolutionCriterion {
  public type: CriterionType = "average_resolution";
  public parameterName: string = "average_resolution";
}

export class AverageResolutionCriterionOption extends ResolutionCriterionOption {
  public label: string = Criterion.getLabel("average_resolution");
  public value: CriterionType = "average_resolution";
}
