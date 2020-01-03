import { CriterionModifier } from "../../../core/generated-graphql";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class ResolutionCriterion extends Criterion<string, string> { // TODO <string, string[]>
  public type: CriterionType = "resolution";
  public parameterName: string = "resolution";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: string[] = ["240p", "480p", "720p", "1080p", "4k"];
  public value: string = "";
}

export class ResolutionCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("resolution");
  public value: CriterionType = "resolution";
}
