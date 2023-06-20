import {
  VideoCodecCriterionInput,
  CriterionModifier,
} from "src/core/generated-graphql";
import { stringToVideoCodec, videoCodecStrings } from "src/utils/videoCodec";
import { CriterionType } from "../types";
import { CriterionOption, StringCriterion } from "./criterion";

abstract class AbstractVideoCodecCriterion extends StringCriterion {
  protected toCriterionInput(): VideoCodecCriterionInput | undefined {
    const value = stringToVideoCodec(this.value);

    if (value !== undefined) {
      return {
        value,
        modifier: this.modifier,
      };
    }
  }
}

class VideoCodecCriterionOptionType extends CriterionOption {
  constructor(value: CriterionType) {
    super({
      messageID: value,
      type: value,
      parameterName: value,
      modifierOptions: [CriterionModifier.Equals, CriterionModifier.NotEquals],
      options: videoCodecStrings,
    });
  }
}

export const VideoCodecCriterionOption = new VideoCodecCriterionOptionType(
  "video_codec"
);

export class VideoCodecCriterion extends AbstractVideoCodecCriterion {
  constructor() {
    super(VideoCodecCriterionOption);
  }
}
