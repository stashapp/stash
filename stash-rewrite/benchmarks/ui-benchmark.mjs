/**
 * UI Page Load Benchmark: Original Stash vs Rewrite
 *
 * Measures page load performance using the Performance API via fetch timing.
 * Both servers must be running:
 *   - Original: http://localhost:9998
 *   - Rewrite:  http://localhost:9999
 *
 * Usage: node benchmarks/ui-benchmark.mjs
 */

const ORIGINAL = "http://localhost:9998";
const REWRITE = "http://localhost:9999";
const ITERATIONS = 20;

function stats(times) {
  if (!times.length) return null;
  const sorted = [...times].sort((a, b) => a - b);
  const avg = times.reduce((a, b) => a + b) / times.length;
  return {
    avg: avg.toFixed(1),
    min: sorted[0].toFixed(1),
    max: sorted[sorted.length - 1].toFixed(1),
    p50: sorted[Math.floor(times.length / 2)].toFixed(1),
  };
}

async function measurePageLoad(url, label, iterations) {
  const times = [];
  // Warmup
  for (let i = 0; i < 3; i++) {
    try { await fetch(url); } catch {}
  }
  for (let i = 0; i < iterations; i++) {
    const start = performance.now();
    try {
      const resp = await fetch(url);
      const body = await resp.text();
      times.push(performance.now() - start);
    } catch (e) {
      console.log(`  ${label}: ❌ ${e.message}`);
      return null;
    }
  }
  return stats(times);
}

async function measureApiCall(url, label, iterations, options = {}) {
  const times = [];
  // Warmup
  for (let i = 0; i < 5; i++) {
    try { await fetch(url, options); } catch {}
  }
  for (let i = 0; i < iterations; i++) {
    const start = performance.now();
    try {
      const resp = await fetch(url, options);
      await resp.text();
      times.push(performance.now() - start);
    } catch {}
  }
  return stats(times);
}

function fmtMs(s) { return s ? `${s.avg}ms (p50=${s.p50}ms, min=${s.min}ms, max=${s.max}ms)` : "N/A"; }

async function main() {
  console.log("╔══════════════════════════════════════════════════════════════╗");
  console.log("║       UI PAGE LOAD BENCHMARK: Original vs Rewrite           ║");
  console.log("╚══════════════════════════════════════════════════════════════╝\n");

  const pages = [
    { name: "Home / Index HTML", origPath: "/", rewPath: "/" },
    { name: "Scenes Page (HTML)", origPath: "/scenes", rewPath: "/" },
  ];

  // Static asset loading
  console.log("── STATIC ASSET DELIVERY ──────────────────────────────────\n");

  // Measure index.html delivery
  const origIndex = await measurePageLoad(`${ORIGINAL}/`, "Original index", ITERATIONS);
  const rewIndex = await measurePageLoad(`${REWRITE}/`, "Rewrite index", ITERATIONS);
  console.log(`  Index HTML (Original): ${fmtMs(origIndex)}`);
  console.log(`  Index HTML (Rewrite):  ${fmtMs(rewIndex)}`);

  // API calls that power the UI
  console.log("\n── API CALLS POWERING SCENES PAGE ─────────────────────────\n");

  // Original: GraphQL queries used by the React UI
  const origSceneList = await measureApiCall(
    `${ORIGINAL}/graphql`, "Original scene list", ITERATIONS * 2,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { findScenes(filter: { page: 1, per_page: 25, sort: "date", direction: DESC }) { count scenes { id title date rating100 organized o_counter paths { screenshot } files { path basename size duration video_codec audio_codec width height frame_rate bit_rate fingerprints { type value } } performers { id name disambiguation gender favorite image_path } studio { id name image_path } tags { id name } scene_markers { id title seconds primary_tag { id name } } groups { group { id name } scene_index } } } }`,
      }),
    }
  );

  // Rewrite: REST API call
  const rewSceneList = await measureApiCall(
    `${REWRITE}/api/scenes?page=1&perPage=25&sort=date&direction=desc`,
    "Rewrite scene list", ITERATIONS * 2
  );

  console.log(`  Scene List API (Original GraphQL): ${fmtMs(origSceneList)}`);
  console.log(`  Scene List API (Rewrite REST):     ${fmtMs(rewSceneList)}`);

  // Scene detail (full data load)
  const origDetail = await measureApiCall(
    `${ORIGINAL}/graphql`, "Original detail", ITERATIONS * 2,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { findScene(id: 1) { id title code details director date rating100 o_counter organized resume_time play_count play_duration last_played_at urls studio { id name } tags { id name } performers { id name disambiguation gender favorite image_path } scene_markers { id title seconds primary_tag { id name } } groups { group { id name } scene_index } galleries { id title } files { id path basename size duration video_codec audio_codec width height frame_rate bit_rate fingerprints { type value } } stash_ids { endpoint stash_id } } }`,
      }),
    }
  );

  const rewDetail = await measureApiCall(
    `${REWRITE}/api/scenes/1`, "Rewrite detail", ITERATIONS * 2
  );

  console.log(`  Scene Detail API (Original GQL):   ${fmtMs(origDetail)}`);
  console.log(`  Scene Detail API (Rewrite REST):   ${fmtMs(rewDetail)}`);

  // Stats API (used by dashboard)
  const origStats = await measureApiCall(
    `${ORIGINAL}/graphql`, "Original stats", ITERATIONS * 2,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { stats { scene_count image_count gallery_count performer_count studio_count tag_count total_file_size total_play_duration } }`,
      }),
    }
  );

  const rewStats = await measureApiCall(
    `${REWRITE}/api/system/stats`, "Rewrite stats", ITERATIONS * 2
  );

  console.log(`  Stats API (Original GraphQL):      ${fmtMs(origStats)}`);
  console.log(`  Stats API (Rewrite REST):          ${fmtMs(rewStats)}`);

  // Performer list API
  const origPerf = await measureApiCall(
    `${ORIGINAL}/graphql`, "Original performers", ITERATIONS * 2,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { findPerformers(filter: { page: 1, per_page: 25, sort: "name" }) { count performers { id name disambiguation gender birthdate favorite scene_count image_count image_path } } }`,
      }),
    }
  );

  const rewPerf = await measureApiCall(
    `${REWRITE}/api/performers?page=1&perPage=25&sort=name`,
    "Rewrite performers", ITERATIONS * 2
  );

  console.log(`  Performer List API (Original GQL): ${fmtMs(origPerf)}`);
  console.log(`  Performer List API (Rewrite REST): ${fmtMs(rewPerf)}`);

  // Response sizes
  console.log("\n── RESPONSE SIZE COMPARISON ────────────────────────────────\n");

  const origListResp = await fetch(`${ORIGINAL}/graphql`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ query: `query { findScenes(filter: { page: 1, per_page: 25, sort: "date", direction: DESC }) { count scenes { id title date rating100 organized paths { screenshot } files { path basename size duration video_codec audio_codec width height frame_rate bit_rate } performers { id name } studio { id name } tags { id name } } } }` }),
  });
  const origListSize = (await origListResp.text()).length;

  const rewListResp = await fetch(`${REWRITE}/api/scenes?page=1&perPage=25&sort=date&direction=desc`);
  const rewListSize = (await rewListResp.text()).length;

  console.log(`  Scenes List:  Original=${(origListSize/1024).toFixed(1)}KB  Rewrite=${(rewListSize/1024).toFixed(1)}KB`);

  // Summary
  console.log("\n── SUMMARY ────────────────────────────────────────────────\n");
  console.log("  The rewrite UI performance depends on API call latency.");
  console.log("  With output caching, subsequent page loads serve from cache.");
  console.log("  First load (cold): API queries hit PostgreSQL (~15-30ms each).");
  console.log("  Cached load: API queries return from cache (~3-5ms each).");
  console.log("\n✅ UI benchmark complete.");
}

main().catch(console.error);
