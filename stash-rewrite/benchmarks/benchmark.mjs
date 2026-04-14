/**
 * Comprehensive side-by-side benchmark: Original Stash (Go/SQLite/GraphQL) vs Rewrite (.NET/PostgreSQL/REST).
 *
 * This script assumes:
 *   - Original stash is running on port 9998 with matching library paths
 *   - Rewrite stash is running on port 9999 with matching library paths
 *   - Both have completed a scan so they have the same content in their databases
 *
 * Usage:
 *   node benchmarks/benchmark.mjs [--original http://localhost:9998] [--rewrite http://localhost:9999] [--iterations 100]
 */

import { writeFileSync } from "fs";

const args = process.argv.slice(2);
function getArg(name, defaultVal) {
  const idx = args.indexOf(`--${name}`);
  return idx >= 0 && args[idx + 1] ? args[idx + 1] : defaultVal;
}

const ORIGINAL_URL = getArg("original", "http://localhost:9998");
const REWRITE_URL = getArg("rewrite", "http://localhost:9999");
const ITERATIONS = parseInt(getArg("iterations", "100"), 10);
const WARMUP = 10;
const CONCURRENT_USERS = [1, 5, 10]; // for throughput tests

// ── Helpers ───────────────────────────────────────────────────────────

async function timedFetch(config) {
  const start = performance.now();
  try {
    const res = await fetch(config.url, {
      method: config.method || "GET",
      headers: config.headers,
      body: config.body,
    });
    const elapsed = performance.now() - start;
    const bodySize = (await res.text()).length;
    return { elapsed, ok: res.ok, status: res.status, bodySize };
  } catch (e) {
    return { elapsed: performance.now() - start, ok: false, status: 0, bodySize: 0, error: e.message };
  }
}

function percentile(sorted, pct) {
  const idx = Math.ceil((pct / 100) * sorted.length) - 1;
  return sorted[Math.max(0, idx)];
}

function stats(times) {
  if (times.length === 0) return null;
  const sorted = [...times].sort((a, b) => a - b);
  const avg = times.reduce((s, t) => s + t, 0) / times.length;
  return {
    count: times.length,
    avg: +avg.toFixed(2),
    min: +sorted[0].toFixed(2),
    max: +sorted[sorted.length - 1].toFixed(2),
    p50: +percentile(sorted, 50).toFixed(2),
    p95: +percentile(sorted, 95).toFixed(2),
    p99: +percentile(sorted, 99).toFixed(2),
    rps: +(1000 / avg).toFixed(1),
  };
}

function fmtMs(ms) {
  if (ms == null) return "N/A";
  return ms < 1 ? `${(ms * 1000).toFixed(0)}µs` : `${ms.toFixed(1)}ms`;
}

function printRow(label, s) {
  if (!s) { console.log(`  ${label}: ❌ No successful responses`); return; }
  console.log(`  ${label}: avg=${fmtMs(s.avg)}  p50=${fmtMs(s.p50)}  p95=${fmtMs(s.p95)}  p99=${fmtMs(s.p99)}  min=${fmtMs(s.min)}  max=${fmtMs(s.max)}  ~${s.rps} req/s  (${s.count}/${ITERATIONS})`);
}

function printWinner(origStats, rewStats, name) {
  if (!origStats || !rewStats) return;
  const ratio = origStats.avg / rewStats.avg;
  const winner = ratio > 1 ? "Rewrite" : "Original";
  const factor = ratio > 1 ? ratio : 1 / ratio;
  const pct = ((factor - 1) * 100).toFixed(1);
  console.log(`  → ${winner} is ${factor.toFixed(2)}x faster (${pct}% improvement on avg latency)`);
}

// ── GraphQL helpers for original stash ────────────────────────────────

function gql(query) {
  return {
    url: `${ORIGINAL_URL}/graphql`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ query }),
  };
}

// ── Rewrite REST helpers ──────────────────────────────────────────────

function rest(path, method = "GET", body = null) {
  const config = { url: `${REWRITE_URL}${path}`, method };
  if (body) {
    config.headers = { "Content-Type": "application/json" };
    config.body = JSON.stringify(body);
  }
  return config;
}

// ── Scenario Categories ───────────────────────────────────────────────

const categories = [
  {
    category: "📋 BASIC READS — Single Entity Fetch",
    scenarios: [
      {
        name: "Get scene by ID (id=1)",
        original: gql(`query { findScene(id: 1) { id title date details rating100 o_counter organized resume_time play_count play_duration files { path size duration video_codec audio_codec width height frame_rate bit_rate } } }`),
        rewrite: rest("/api/scenes/1"),
      },
      {
        name: "Get performer by ID (id=1)",
        original: gql(`query { findPerformer(id: "1") { id name disambiguation gender birthdate ethnicity country hair_color eye_color height_cm weight measurements fake_tits career_length tattoos piercings favorite scene_count image_count } }`),
        rewrite: rest("/api/performers/1"),
      },
      {
        name: "Get tag by ID (id=1)",
        original: gql(`query { findTag(id: "1") { id name description scene_count image_count } }`),
        rewrite: rest("/api/tags/1"),
      },
    ],
  },
  {
    category: "📑 PAGINATED LISTS — Various Page Sizes",
    scenarios: [
      {
        name: "List scenes (page 1, 25 items)",
        original: gql(`query { findScenes(filter: { page: 1, per_page: 25, sort: "date", direction: DESC }) { count scenes { id title date rating100 } } }`),
        rewrite: rest("/api/scenes?page=1&perPage=25&sort=date&direction=desc"),
      },
      {
        name: "List scenes (page 1, 100 items)",
        original: gql(`query { findScenes(filter: { page: 1, per_page: 100, sort: "date", direction: DESC }) { count scenes { id title date rating100 } } }`),
        rewrite: rest("/api/scenes?page=1&perPage=100&sort=date&direction=desc"),
      },
      {
        name: "List performers (page 1, 25 items)",
        original: gql(`query { findPerformers(filter: { page: 1, per_page: 25, sort: "name", direction: ASC }) { count performers { id name gender } } }`),
        rewrite: rest("/api/performers?page=1&perPage=25&sort=name&direction=asc"),
      },
      {
        name: "List tags (page 1, 100 items)",
        original: gql(`query { findTags(filter: { page: 1, per_page: 100, sort: "name", direction: ASC }) { count tags { id name } } }`),
        rewrite: rest("/api/tags?page=1&perPage=100&sort=name&direction=asc"),
      },
      {
        name: "List studios (page 1, 25 items)",
        original: gql(`query { findStudios(filter: { page: 1, per_page: 25, sort: "name", direction: ASC }) { count studios { id name } } }`),
        rewrite: rest("/api/studios?page=1&perPage=25&sort=name&direction=asc"),
      },
      {
        name: "List images (page 1, 25 items)",
        original: gql(`query { findImages(filter: { page: 1, per_page: 25, sort: "path", direction: ASC }) { count images { id title } } }`),
        rewrite: rest("/api/images?page=1&perPage=25&sort=path&direction=asc"),
      },
    ],
  },
  {
    category: "🔍 FILTERED QUERIES — Complex Searches",
    scenarios: [
      {
        name: "Scenes filtered by rating >= 60",
        original: gql(`query { findScenes(scene_filter: { rating100: { value: 60, modifier: GREATER_THAN } }, filter: { page: 1, per_page: 25 }) { count scenes { id title rating100 } } }`),
        rewrite: rest("/api/scenes/find", "POST", { criteria: [{ type: "rating100", value: { value: 60, modifier: "greaterThan" } }], filter: { page: 1, perPage: 25 } }),
      },
      {
        name: "Scenes text search (query 'a')",
        original: gql(`query { findScenes(filter: { q: "a", page: 1, per_page: 25 }) { count scenes { id title } } }`),
        rewrite: rest("/api/scenes?q=a&page=1&perPage=25"),
      },
      {
        name: "Scenes sorted by file size desc",
        original: gql(`query { findScenes(filter: { page: 1, per_page: 25, sort: "file_size", direction: DESC }) { count scenes { id title } } }`),
        rewrite: rest("/api/scenes?page=1&perPage=25&sort=fileSize&direction=desc"),
      },
      {
        name: "Scenes sorted by duration asc",
        original: gql(`query { findScenes(filter: { page: 1, per_page: 25, sort: "duration", direction: ASC }) { count scenes { id title } } }`),
        rewrite: rest("/api/scenes?page=1&perPage=25&sort=duration&direction=asc"),
      },
    ],
  },
  {
    category: "📊 AGGREGATION — Stats & Counts",
    scenarios: [
      {
        name: "System stats (all entity counts)",
        original: gql(`query { stats { scene_count image_count gallery_count performer_count studio_count tag_count total_file_size total_play_duration } }`),
        rewrite: rest("/api/system/stats"),
      },
      {
        name: "Scene count only",
        original: gql(`query { findScenes(filter: { per_page: 0 }) { count } }`),
        rewrite: rest("/api/scenes?perPage=0"),
      },
    ],
  },
  {
    category: "⚡ RAPID-FIRE — Sequential Rapid Requests",
    scenarios: [
      {
        name: "20 sequential scene fetches (id 1-20)",
        isCustom: true,
        run: async () => {
          const origTimes = [];
          const rewTimes = [];
          // Run all original fetches first, then all rewrite fetches
          for (let i = 1; i <= 20; i++) {
            const o = await timedFetch(gql(`query { findScene(id: ${i}) { id title } }`));
            if (o.ok) origTimes.push(o.elapsed);
          }
          for (let i = 1; i <= 20; i++) {
            const r = await timedFetch(rest(`/api/scenes/${i}`));
            if (r.ok) rewTimes.push(r.elapsed);
          }
          return { original: origTimes, rewrite: rewTimes };
        },
      },
    ],
  },
];

// ── Concurrent throughput test ────────────────────────────────────────

async function concurrentBench(config, concurrency, totalRequests) {
  const times = [];
  let completed = 0;
  const start = performance.now();

  async function worker() {
    while (completed < totalRequests) {
      completed++;
      const r = await timedFetch(config);
      if (r.ok) times.push(r.elapsed);
    }
  }

  await Promise.all(Array.from({ length: concurrency }, () => worker()));
  const wallTime = performance.now() - start;
  return { times, wallTime, throughput: (times.length / wallTime) * 1000 };
}

// ── Main benchmark runner ─────────────────────────────────────────────

async function main() {
  console.log("╔════════════════════════════════════════════════════════════════╗");
  console.log("║    COMPREHENSIVE STASH BENCHMARK: Original vs Rewrite         ║");
  console.log("╠════════════════════════════════════════════════════════════════╣");
  console.log(`║  Original (Go/SQLite/GQL):  ${ORIGINAL_URL.padEnd(35)}║`);
  console.log(`║  Rewrite  (.NET/PG/REST):   ${REWRITE_URL.padEnd(35)}║`);
  console.log(`║  Iterations per scenario:   ${String(ITERATIONS).padEnd(35)}║`);
  console.log(`║  Warmup iterations:         ${String(WARMUP).padEnd(35)}║`);
  console.log(`║  Concurrent user levels:    ${CONCURRENT_USERS.join(", ").padEnd(35)}║`);
  console.log("╚════════════════════════════════════════════════════════════════╝\n");

  // Connectivity check
  console.log("🔌 Checking connectivity...");
  let origAlive = false, rewAlive = false;
  let origStats = null, rewStats = null;
  try { const r = await timedFetch(gql(`{ stats { scene_count } }`)); origAlive = r.ok; if (r.ok) origStats = JSON.parse((await (await fetch(gql(`{ stats { scene_count image_count performer_count studio_count tag_count } }`).url, { method: "POST", headers: { "Content-Type": "application/json" }, body: gql(`{ stats { scene_count image_count performer_count studio_count tag_count } }`).body })).text()).data?.stats); } catch {}
  try { const r = await timedFetch(rest("/api/system/stats")); rewAlive = r.ok; if (r.ok) rewStats = await (await fetch(rest("/api/system/stats").url)).json(); } catch {}

  console.log(`  Original: ${origAlive ? "✅ Connected" : "❌ Not reachable"}`);
  console.log(`  Rewrite:  ${rewAlive ? "✅ Connected" : "❌ Not reachable"}`);

  if (origStats) console.log(`  Original DB: ${origStats.scene_count} scenes, ${origStats.performer_count} performers, ${origStats.tag_count} tags`);
  if (rewStats) console.log(`  Rewrite DB:  ${rewStats.sceneCount} scenes, ${rewStats.performerCount} performers, ${rewStats.tagCount} tags`);

  if (!origAlive && !rewAlive) {
    console.log("\n⚠️  Neither server is running. Start at least one.\n");
    process.exit(1);
  }
  console.log("");

  const allResults = [];

  // ── Run categorized scenarios ──────────────────────────────────────
  for (const cat of categories) {
    console.log(`\n${"═".repeat(70)}`);
    console.log(`  ${cat.category}`);
    console.log(`${"═".repeat(70)}`);

    for (const scenario of cat.scenarios) {
      console.log(`\n  ▸ ${scenario.name}`);

      let origTimes = [], rewTimes = [];

      if (scenario.isCustom) {
        const result = await scenario.run();
        origTimes = result.original;
        rewTimes = result.rewrite;
      } else {
        // Run ALL original iterations first, then ALL rewrite iterations
        // (avoids simultaneous load on shared PC/disk which would skew results)
        if (origAlive) {
          for (let i = 0; i < WARMUP; i++) await timedFetch(scenario.original);
          for (let i = 0; i < ITERATIONS; i++) {
            const r = await timedFetch(scenario.original);
            if (r.ok) origTimes.push(r.elapsed);
          }
        }
        if (rewAlive) {
          for (let i = 0; i < WARMUP; i++) await timedFetch(scenario.rewrite);
          for (let i = 0; i < ITERATIONS; i++) {
            const r = await timedFetch(scenario.rewrite);
            if (r.ok) rewTimes.push(r.elapsed);
          }
        }
      }

      const oStats = stats(origTimes);
      const rStats = stats(rewTimes);

      printRow("Original (Go/SQLite) ", oStats);
      printRow("Rewrite  (.NET/PG)   ", rStats);
      printWinner(oStats, rStats, scenario.name);

      allResults.push({ category: cat.category, name: scenario.name, original: oStats, rewrite: rStats });
    }
  }

  // ── Concurrent throughput tests ────────────────────────────────────
  console.log(`\n${"═".repeat(70)}`);
  console.log(`  🏋️ CONCURRENT THROUGHPUT — List Scenes under Load`);
  console.log(`${"═".repeat(70)}`);

  const totalReqs = 100;
  for (const concurrency of CONCURRENT_USERS) {
    console.log(`\n  ▸ ${concurrency} concurrent users, ${totalReqs} total requests`);

    let origResult = null, rewResult = null;

    // Run original first, then rewrite (no simultaneous load on shared PC)
    if (origAlive) {
      const origC = await concurrentBench(gql(`query { findScenes(filter: { page: 1, per_page: 25 }) { count scenes { id title } } }`), concurrency, totalReqs);
      const oS = stats(origC.times);
      origResult = { ...oS, throughput: origC.throughput };
      console.log(`    Original: avg=${fmtMs(oS?.avg)}  p95=${fmtMs(oS?.p95)}  throughput=${origC.throughput.toFixed(1)} req/s  wall=${fmtMs(origC.wallTime)}`);
    }

    if (rewAlive) {
      const rewC = await concurrentBench(rest("/api/scenes?page=1&perPage=25"), concurrency, totalReqs);
      const rS = stats(rewC.times);
      rewResult = { ...rS, throughput: rewC.throughput };
      console.log(`    Rewrite:  avg=${fmtMs(rS?.avg)}  p95=${fmtMs(rS?.p95)}  throughput=${rewC.throughput.toFixed(1)} req/s  wall=${fmtMs(rewC.wallTime)}`);
    }

    allResults.push({ category: "Concurrent Throughput", name: `${concurrency} users`, original: origResult, rewrite: rewResult });
  }

  // ── Response size comparison ────────────────────────────────────────
  console.log(`\n${"═".repeat(70)}`);
  console.log(`  📦 RESPONSE SIZE COMPARISON`);
  console.log(`${"═".repeat(70)}`);

  const sizeTests = [
    { name: "Scene list (25)", original: gql(`query { findScenes(filter: { page: 1, per_page: 25 }) { count scenes { id title date details rating100 o_counter organized play_count play_duration files { path size duration video_codec audio_codec width height frame_rate bit_rate } } } }`), rewrite: rest("/api/scenes?page=1&perPage=25") },
    { name: "Single scene", original: gql(`query { findScene(id: 1) { id title date details rating100 o_counter organized play_count play_duration files { path size duration video_codec audio_codec width height frame_rate bit_rate } } }`), rewrite: rest("/api/scenes/1") },
    { name: "Stats", original: gql(`query { stats { scene_count image_count gallery_count performer_count studio_count tag_count } }`), rewrite: rest("/api/system/stats") },
  ];

  for (const t of sizeTests) {
    let origSize = "N/A", rewSize = "N/A";
    if (origAlive) { const r = await timedFetch(t.original); origSize = r.ok ? `${(r.bodySize / 1024).toFixed(1)} KB` : "err"; }
    if (rewAlive) { const r = await timedFetch(t.rewrite); rewSize = r.ok ? `${(r.bodySize / 1024).toFixed(1)} KB` : "err"; }
    console.log(`  ${t.name.padEnd(20)} Original: ${origSize.padEnd(12)} Rewrite: ${rewSize}`);
  }

  // ── Summary table ──────────────────────────────────────────────────
  console.log(`\n${"═".repeat(70)}`);
  console.log("  📊 SUMMARY TABLE");
  console.log(`${"═".repeat(70)}\n`);

  console.log("  Scenario".padEnd(45) + "Original avg".padEnd(15) + "Rewrite avg".padEnd(15) + "Winner");
  console.log("  " + "─".repeat(68));

  let origWins = 0, rewWins = 0;
  for (const r of allResults.filter(r => r.original && r.rewrite)) {
    const origAvg = r.original.avg;
    const rewAvg = r.rewrite.avg;
    const winner = origAvg < rewAvg ? "Original" : "Rewrite";
    const factor = origAvg < rewAvg ? (rewAvg / origAvg) : (origAvg / rewAvg);
    if (origAvg < rewAvg) origWins++; else rewWins++;
    console.log(`  ${r.name.padEnd(43)} ${fmtMs(origAvg).padEnd(15)}${fmtMs(rewAvg).padEnd(15)}${winner} (${factor.toFixed(1)}x)`);
  }

  if (origWins + rewWins > 0) {
    console.log(`\n  Overall: Original won ${origWins} scenarios, Rewrite won ${rewWins} scenarios`);
  }

  // Save results JSON
  const outputPath = "benchmarks/results.json";
  writeFileSync(outputPath, JSON.stringify({ timestamp: new Date().toISOString(), config: { originalUrl: ORIGINAL_URL, rewriteUrl: REWRITE_URL, iterations: ITERATIONS, warmup: WARMUP }, results: allResults }, null, 2));
  console.log(`\n  📁 Results saved to ${outputPath}`);
  console.log("\n✅ Benchmark complete.\n");
}

main().catch(console.error);
