import {
  AudioCodecCriterionInput,
  CriterionModifier,
} from "src/core/generated-graphql";
import { stringToAudioCodec, audioCodecStrings } from "src/utils/audioCodec";
import { CriterionType } from "../types";
import { CriterionOption, StringCriterion } from "./criterion";

abstract class AbstractAudioCodecCriterion extends StringCriterion {
  protected toCriterionInput(): AudioCodecCriterionInput | undefined {
    const value = stringToAudioCodec(this.value);

    if (value !== undefined) {
      return {
        value,
        modifier: this.modifier,
      };
    }
  }
}

class AudioCodecCriterionOptionType extends CriterionOption {
  constructor(value: CriterionType) {
    super({
      messageID: value,
      type: value,
      parameterName: value,
      modifierOptions: [CriterionModifier.Equals, CriterionModifier.NotEquals],
      options: audioCodecStrings,
    });
  }
}

export const AudioCodecCriterionOption = new AudioCodecCriterionOptionType(
  "audio_codec"
);

export class AudioCodecCriterion extends AbstractAudioCodecCriterion {
  constructor() {
    super(AudioCodecCriterionOption);
  }
}
