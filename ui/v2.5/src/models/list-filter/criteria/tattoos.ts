import { StringCriterion, StringCriterionOption } from "./criterion";

export const TattoosCriterionOption = new StringCriterionOption({
  messageID: "tattoos",
  type: "tattoos",
  makeCriterion: () => new TattoosCriterion(),
});

export class TattoosCriterion extends StringCriterion {
  constructor() {
    super(TattoosCriterionOption);
  }
}
