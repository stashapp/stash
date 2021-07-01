import { ResolutionEnum } from "src/core/generated-graphql";
import { CriterionType } from "../types";
import { CriterionOption, StringCriterion } from "./criterion";

abstract class AbstractResolutionCriterion extends StringCriterion {
  protected toCriterionInput(): ResolutionEnum | undefined {
    switch (this.value) {
      case "144p":
        return ResolutionEnum.VeryLow;
      case "240p":
        return ResolutionEnum.Low;
      case "360p":
        return ResolutionEnum.R360P;
      case "480p":
        return ResolutionEnum.Standard;
      case "540p":
        return ResolutionEnum.WebHd;
      case "720p":
        return ResolutionEnum.StandardHd;
      case "1080p":
        return ResolutionEnum.FullHd;
      case "1440p":
        return ResolutionEnum.QuadHd;
      case "1920p":
        return ResolutionEnum.VrHd;
      case "4k":
        return ResolutionEnum.FourK;
      case "5k":
        return ResolutionEnum.FiveK;
      case "6k":
        return ResolutionEnum.SixK;
      case "8k":
        return ResolutionEnum.EightK;
      // no default
    }
  }
}

class ResolutionCriterionOptionType extends CriterionOption {
  constructor(value: CriterionType) {
    super({
      messageID: value,
      type: value,
      parameterName: value,
      options: [
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
      ],
    });
  }
}

export const ResolutionCriterionOption = new ResolutionCriterionOptionType(
  "resolution"
);
export class ResolutionCriterion extends AbstractResolutionCriterion {
  constructor() {
    super(ResolutionCriterionOption);
  }
}

export const AverageResolutionCriterionOption = new ResolutionCriterionOptionType(
  "average_resolution"
);

export class AverageResolutionCriterion extends AbstractResolutionCriterion {
  constructor() {
    super(AverageResolutionCriterionOption);
  }
}
