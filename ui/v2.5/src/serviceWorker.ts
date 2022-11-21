import { clientsClaim } from "workbox-core";
import { registerRoute, Route, setDefaultHandler } from "workbox-routing";
import { NetworkFirst, StaleWhileRevalidate } from "workbox-strategies";
import type { ManifestEntry } from "workbox-build";
import { precacheAndRoute, cleanupOutdatedCaches } from "workbox-precaching";

// Give TypeScript the correct global.
declare let self: ServiceWorkerGlobalScope;

const manifest = self.__WB_MANIFEST as Array<ManifestEntry>;

precacheAndRoute(manifest);

const gzRoute = new Route(({ request, sameOrigin }) => {
  return new RegExp(".(js|json|css|svg|md)$").test(request.url);
}, new StaleWhileRevalidate());
registerRoute(gzRoute);

setDefaultHandler(new NetworkFirst());
cleanupOutdatedCaches();

// this is necessary, since the new service worker will keep on skipWaiting state
// and then, caches will not be cleared since it is not activated
self.skipWaiting();
clientsClaim();
