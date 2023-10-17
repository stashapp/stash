import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption, StringCriterion } from "./criterion";

export const valueToCode = (value?: string | null) => {
  if (!value) {
    return undefined;
  }
};

class CriterionOptionFactory {
  static create(type: "audio_codec" | "video_codec"): CriterionOption {
    return new CriterionOption({
      messageID: "codecs",
      type,
      modifierOptions: [
        CriterionModifier.Includes,
        CriterionModifier.Excludes,
        CriterionModifier.IsNull,
        CriterionModifier.NotNull,
      ],
      defaultModifier: CriterionModifier.Includes,
      makeCriterion: () => new CodecsCriterion(type),
    });
  }
}
export class CodecsCriterion extends StringCriterion {
  constructor(type: "audio_codec" | "video_codec") {
    const criterionOption = CriterionOptionFactory.create(type);
    super(criterionOption);
  }

  protected toCriterionInput() {
    const value = valueToCode(this.value) ?? "";

    return {
      value,
      modifier: this.modifier,
    };
  }
}
