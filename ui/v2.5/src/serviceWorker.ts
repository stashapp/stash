import { clientsClaim } from "workbox-core";
import { registerRoute, Route } from "workbox-routing";
import { NetworkOnly, StrategyHandler } from "workbox-strategies";
import { ExpirationPlugin } from "workbox-expiration";
import type { ManifestEntry } from "workbox-build";
import { precacheAndRoute, cleanupOutdatedCaches } from "workbox-precaching";
import { Strategy } from "workbox-strategies";

class CacheNetworkRace extends Strategy {
  _handle(
    request: Request,
    handler: StrategyHandler
  ): Promise<Response | undefined> {
    const fetchAndCachePutDone = handler.fetchAndCachePut(request);
    const cacheMatchDone = handler.cacheMatch(request);

    return new Promise((resolve, reject) => {
      fetchAndCachePutDone.then(resolve);
      cacheMatchDone.then((response) => response && resolve(response));

      // Reject if both network and cache error or find no response.
      Promise.allSettled([fetchAndCachePutDone, cacheMatchDone]).then(
        (results) => {
          const [fetchAndCachePutResult, cacheMatchResult] = results;
          if (
            fetchAndCachePutResult.status === "rejected" &&
            !("value" in cacheMatchResult)
          ) {
            reject(fetchAndCachePutResult.reason);
          }
        }
      );
    });
  }
}

const cacheNetworkRaceExpiration = (
  cacheName: string,
  entries: number,
  maxAgeDays: number
) => {
  return new CacheNetworkRace({
    cacheName: cacheName,
    plugins: [
      new ExpirationPlugin({
        maxAgeSeconds: maxAgeDays * 24 * 60 * 60, // Days
        maxEntries: entries,
        purgeOnQuotaError: true,
        matchOptions: {
          ignoreSearch: true,
        },
      }),
    ],
  });
};

// Give TypeScript the correct global.
declare let self: ServiceWorkerGlobalScope;

const manifest = self.__WB_MANIFEST as Array<ManifestEntry>;

cleanupOutdatedCaches();
precacheAndRoute(manifest);

// handle gzipped files using runtime cache
registerRoute(
  new Route(({ request, sameOrigin }) => {
    return sameOrigin && new RegExp(".(js|json|css|svg|md)$").test(request.url);
  }, new CacheNetworkRace())
);

// scene screenshots
registerRoute(
  new Route(({ request, sameOrigin }) => {
    return (
      sameOrigin && new RegExp(`\/scene\/.*\/screenshot`).test(request.url)
    );
  }, cacheNetworkRaceExpiration("scene-screenshots", 100, 30))
);

// scene vtt files
registerRoute(
  new Route(({ request, sameOrigin }) => {
    return (
      sameOrigin &&
      new RegExp(`(\/scene\/[0-9]*((_thumbs\.vtt$)|(\/vtt\/chapter)))`).test(
        request.url
      )
    );
  }, cacheNetworkRaceExpiration("scene-vtt", 100, 30))
);

// scene sprites
registerRoute(
  new Route(({ request, sameOrigin }) => {
    return (
      sameOrigin && new RegExp(`\/scene\/[0-9]*_sprite\.jpg$`).test(request.url)
    );
  }, cacheNetworkRaceExpiration("scene-sprites", 100, 30))
);

// image thumbnails; studio images
registerRoute(
  new Route(({ request, sameOrigin }) => {
    return (
      sameOrigin &&
      new RegExp(`(\/image\/.*\/thumbnail)|(\/studio\/.*\/image)`).test(
        request.url
      )
    );
  }, cacheNetworkRaceExpiration("image-thumbnail", 100, 30))
);

// https://web.dev/sw-range-requests/
// https://bugzilla.mozilla.org/show_bug.cgi?id=1465074
// https://wpt.fyi/results/fetch/range/sw.https.window.html?label=master&label=experimental&aligned&view=subtest
// setDefaultHandler(new NetworkOnly());
// use custom default handler to fix firefox range requests
registerRoute(
  new Route(({ request }) => {
    if (
      request.headers.has("range") &&
      navigator.userAgent.search("Firefox") !== -1
    ) {
      return false;
    }
    return true;
  }, new NetworkOnly())
);

// this immediately takes over the page with the new service worker
self.skipWaiting();
clientsClaim();
