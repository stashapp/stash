import { ScraperSourceInput } from "src/core/generated-graphql";

export const STASH_BOX_PREFIX = "stashbox:";
export const SCRAPER_PREFIX = "scraper:";

export interface IScraperSource {
  id: string;
  stashboxEndpoint?: string;
  sourceInput: ScraperSourceInput;
  displayName: string;
  supportQuery?: boolean;
  supportFragment?: boolean;
}
