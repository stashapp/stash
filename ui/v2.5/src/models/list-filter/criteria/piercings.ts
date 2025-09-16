import { StringCriterion, StringCriterionOption } from "./criterion";

export const PiercingsCriterionOption = new StringCriterionOption({
  messageID: "piercings",
  type: "piercings",
  makeCriterion: () => new PiercingsCriterion(),
});

export class PiercingsCriterion extends StringCriterion {
  constructor() {
    super(PiercingsCriterionOption);
  }
}
