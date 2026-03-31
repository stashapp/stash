import { shouldPolyfill as shouldPolyfillCanonicalLocales } from "@formatjs/intl-getcanonicallocales/should-polyfill";
import { shouldPolyfill as shouldPolyfillLocale } from "@formatjs/intl-locale/should-polyfill";
import { shouldPolyfill as shouldPolyfillNumberformat } from "@formatjs/intl-numberformat/should-polyfill";
import { shouldPolyfill as shouldPolyfillPluralRules } from "@formatjs/intl-pluralrules/should-polyfill";

// needed for older safari versions
import "event-target-polyfill";

async function checkPolyfills() {
  if (shouldPolyfillCanonicalLocales()) {
    await import("@formatjs/intl-getcanonicallocales/polyfill");
  }
  if (shouldPolyfillLocale()) {
    await import("@formatjs/intl-locale/polyfill");
  }
  if (shouldPolyfillNumberformat()) {
    await import("@formatjs/intl-numberformat/polyfill");
    await import("@formatjs/intl-numberformat/locale-data/en");
    await import("@formatjs/intl-numberformat/locale-data/en-GB");
  }
  if (shouldPolyfillPluralRules()) {
    await import("@formatjs/intl-pluralrules/polyfill");
    await import("@formatjs/intl-pluralrules/locale-data/en");
  }
}

export const initPolyfills = async () => {
  await checkPolyfills();
};
