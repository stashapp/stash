/**
 * Query Complexity Scaling Benchmark
 * 
 * Tests how latency scales as query complexity increases across 8 tiers:
 *   Tier 0: Trivial    — COUNT only
 *   Tier 1: Minimal    — single field lookup, no joins  
 *   Tier 2: Simple     — list with 1 filter, basic sort
 *   Tier 3: Moderate   — 2-3 filters, join-based sort
 *   Tier 4: Medium     — 4+ filters, text search, join
 *   Tier 5: Complex    — many filters, computed sorts, multiple joins
 *   Tier 6: Heavy      — all applicable filters stacked
 *   Tier 7: Deep GQL   — deeply nested response (GQL advantage test)
 *   Tier 8: Wide GQL   — maximum field selection (GQL field selection advantage)
 * 
 * Also tests:
 *   - Page size scaling (1, 10, 25, 50, 100, 200)
 *   - Response size vs latency correlation
 *   - Concurrent scaling (1, 2, 5, 10, 20 parallel requests)
 *
 * Usage: node benchmarks/complexity-benchmark.mjs [--iterations 50]
 */

import { writeFileSync } from "fs";

const args = process.argv.slice(2);
function getArg(name, def) {
  const i = args.indexOf(`--${name}`);
  return i >= 0 && args[i + 1] ? args[i + 1] : def;
}

const ORIG = getArg("original", "http://localhost:9998");
const REW = getArg("rewrite", "http://localhost:9999");
const ITERS = parseInt(getArg("iterations", "50"), 10);
const WARMUP = 5;

// ── Helpers ─────────────────────────────────────────────────

async function timedFetch(url, options = {}) {
  const start = performance.now();
  try {
    const res = await fetch(url, options);
    const text = await res.text();
    return { ms: performance.now() - start, ok: res.ok, size: text.length, status: res.status };
  } catch (e) {
    return { ms: performance.now() - start, ok: false, size: 0, error: e.message };
  }
}

function gql(query) {
  return {
    url: `${ORIG}/graphql`,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ query }),
  };
}

function rest(path, method = "GET", body = null) {
  const c = { url: `${REW}${path}`, method };
  if (body) {
    c.headers = { "Content-Type": "application/json" };
    c.body = JSON.stringify(body);
  }
  return c;
}

function stats(times) {
  if (!times.length) return null;
  const sorted = [...times].sort((a, b) => a - b);
  const avg = times.reduce((a, b) => a + b) / times.length;
  return {
    avg: +avg.toFixed(2),
    min: +sorted[0].toFixed(2),
    max: +sorted[sorted.length - 1].toFixed(2),
    p50: +sorted[Math.floor(times.length * 0.5)].toFixed(2),
    p95: +sorted[Math.floor(times.length * 0.95)].toFixed(2),
    p99: +sorted[Math.min(sorted.length - 1, Math.floor(times.length * 0.99))].toFixed(2),
  };
}

function fmtMs(ms) {
  if (ms == null) return "N/A";
  return ms < 1 ? `${(ms * 1000).toFixed(0)}µs` : `${ms.toFixed(1)}ms`;
}

async function bench(config, iterations = ITERS, warmup = WARMUP) {
  const times = [];
  const sizes = [];
  for (let i = 0; i < warmup; i++) await timedFetch(config.url, config);
  for (let i = 0; i < iterations; i++) {
    const r = await timedFetch(config.url, config);
    if (r.ok) { times.push(r.ms); sizes.push(r.size); }
  }
  return { stats: stats(times), avgSize: sizes.length ? Math.round(sizes.reduce((a, b) => a + b) / sizes.length) : 0, successCount: times.length };
}

async function concurrentBench(config, concurrency, totalReqs) {
  const times = [];
  let done = 0;
  const start = performance.now();
  async function worker() {
    while (done < totalReqs) {
      done++;
      const r = await timedFetch(config.url, config);
      if (r.ok) times.push(r.ms);
    }
  }
  await Promise.all(Array.from({ length: concurrency }, () => worker()));
  const wall = performance.now() - start;
  return { stats: stats(times), throughput: +(times.length / wall * 1000).toFixed(1), wall };
}

// ── Complexity Tiers ────────────────────────────────────────

const complexityTiers = [
  {
    tier: 0,
    name: "COUNT only",
    description: "Simplest possible query — just count scenes",
    original: gql(`query { findScenes(filter: { per_page: 0 }) { count } }`),
    rewrite: rest("/api/scenes?perPage=0"),
  },
  {
    tier: 1,
    name: "Single entity by ID",
    description: "Fetch one scene by ID, minimal fields (GQL advantage: field selection)",
    original: gql(`query { findScene(id: 1) { id title date } }`),
    rewrite: rest("/api/scenes/1"),
  },
  {
    tier: 2,
    name: "List + 1 filter + simple sort",
    description: "25 scenes sorted by date desc",
    original: gql(`query { findScenes(filter: { per_page: 25, sort: "date", direction: DESC }) { count scenes { id title date rating100 } } }`),
    rewrite: rest("/api/scenes?perPage=25&sort=date&direction=desc"),
  },
  {
    tier: 3,
    name: "List + text search + sort",
    description: "Text search 'a' + sort by title",
    original: gql(`query { findScenes(filter: { q: "a", per_page: 25, sort: "title", direction: ASC }) { count scenes { id title date rating100 organized } } }`),
    rewrite: rest("/api/scenes?q=a&perPage=25&sort=title&direction=asc"),
  },
  {
    tier: 4,
    name: "List + join-based sort (duration)",
    description: "Sort by duration requires joining files table",
    original: gql(`query { findScenes(filter: { per_page: 25, sort: "duration", direction: DESC }) { count scenes { id title files { duration } } } }`),
    rewrite: rest("/api/scenes?perPage=25&sort=duration&direction=desc"),
  },
  {
    tier: 5,
    name: "List + join-based sort (filesize)",
    description: "Sort by file_size requires joining files table",
    original: gql(`query { findScenes(filter: { per_page: 25, sort: "filesize", direction: DESC }) { count scenes { id title files { size } } } }`),
    rewrite: rest("/api/scenes?perPage=25&sort=file_size&direction=desc"),
  },
  {
    tier: 6,
    name: "List + computed filter (duration criterion)",
    description: "Filter scenes with duration > 600s — requires file join + aggregate",
    original: gql(`query { findScenes(scene_filter: { duration: { value: 600, modifier: GREATER_THAN } }, filter: { per_page: 25, sort: "duration", direction: DESC }) { count scenes { id title files { duration } } } }`),
    rewrite: rest("/api/scenes/find", "POST", {
      findFilter: { page: 1, perPage: 25, sort: "duration", direction: "desc" },
      objectFilter: { durationCriterion: { value: 600, modifier: "greaterThan" } }
    }),
  },
  {
    tier: 7,
    name: "List + 2 computed filters (duration + resolution)",
    description: "Duration > 600s AND resolution >= 1080 — two file-join filters",
    original: gql(`query { findScenes(scene_filter: { duration: { value: 600, modifier: GREATER_THAN }, resolution: { value: FULL_HD, modifier: GREATER_THAN } }, filter: { per_page: 25 }) { count scenes { id title files { duration width height } } } }`),
    rewrite: rest("/api/scenes/find", "POST", {
      findFilter: { page: 1, perPage: 25 },
      objectFilter: {
        durationCriterion: { value: 600, modifier: "greaterThan" },
        resolutionCriterion: { value: 1080, modifier: "greaterThan" }
      }
    }),
  },
  {
    tier: 8,
    name: "List + 3 filters + text search + sort",
    description: "Duration > 300, resolution >= 720, organized=false, search 'e', sort by filesize",
    original: gql(`query { findScenes(scene_filter: { duration: { value: 300, modifier: GREATER_THAN }, resolution: { value: STANDARD_HD, modifier: GREATER_THAN }, organized: false }, filter: { q: "e", per_page: 25, sort: "filesize", direction: DESC }) { count scenes { id title date rating100 files { path size duration width height video_codec } } } }`),
    rewrite: rest("/api/scenes/find", "POST", {
      findFilter: { q: "e", page: 1, perPage: 25, sort: "file_size", direction: "desc" },
      objectFilter: {
        durationCriterion: { value: 300, modifier: "greaterThan" },
        resolutionCriterion: { value: 720, modifier: "greaterThan" },
        organized: false
      }
    }),
  },
  {
    tier: 9,
    name: "List + 5 filters stacked",
    description: "Duration > 300, resolution >= 720, organized=false, framerate >= 24, bitrate > 5000, sort by duration",
    original: gql(`query { findScenes(scene_filter: { duration: { value: 300, modifier: GREATER_THAN }, resolution: { value: STANDARD_HD, modifier: GREATER_THAN }, organized: false, framerate: { value: 24, modifier: GREATER_THAN }, bitrate: { value: 5000, modifier: GREATER_THAN } }, filter: { per_page: 25, sort: "duration", direction: DESC }) { count scenes { id title date files { path size duration width height frame_rate bit_rate video_codec audio_codec } } } }`),
    rewrite: rest("/api/scenes/find", "POST", {
      findFilter: { page: 1, perPage: 25, sort: "duration", direction: "desc" },
      objectFilter: {
        durationCriterion: { value: 300, modifier: "greaterThan" },
        resolutionCriterion: { value: 720, modifier: "greaterThan" },
        organized: false,
        frameRateCriterion: { value: 24, modifier: "greaterThan" },
        bitrateInterval: { value: 5000, modifier: "greaterThan" }
      }
    }),
  },
  {
    tier: 10,
    name: "Maximum filters (7 criteria)",
    description: "All applicable file+scene filters combined",
    original: gql(`query { findScenes(scene_filter: { duration: { value: 60, modifier: GREATER_THAN }, resolution: { value: STANDARD_HD, modifier: GREATER_THAN }, organized: false, framerate: { value: 20, modifier: GREATER_THAN }, bitrate: { value: 1000, modifier: GREATER_THAN }, video_codec: { value: "h264", modifier: EQUALS }, file_count: { value: 1, modifier: EQUALS } }, filter: { per_page: 25, sort: "filesize", direction: DESC }) { count scenes { id title date files { path size duration width height frame_rate bit_rate video_codec audio_codec } } } }`),
    rewrite: rest("/api/scenes/find", "POST", {
      findFilter: { page: 1, perPage: 25, sort: "file_size", direction: "desc" },
      objectFilter: {
        durationCriterion: { value: 60, modifier: "greaterThan" },
        resolutionCriterion: { value: 720, modifier: "greaterThan" },
        organized: false,
        frameRateCriterion: { value: 20, modifier: "greaterThan" },
        bitrateInterval: { value: 1000, modifier: "greaterThan" },
        videoCodecCriterion: { value: "h264", modifier: "equals" },
        fileCountCriterion: { value: 1, modifier: "equals" }
      }
    }),
  },
];

// ── GQL-specific depth tests ────────────────────────────────

const gqlDepthTests = [
  {
    name: "Minimal fields (3 fields)",
    query: gql(`query { findScenes(filter: { per_page: 25 }) { count scenes { id title date } } }`),
    rewrite: rest("/api/scenes?perPage=25"),
  },
  {
    name: "Medium fields (10 fields)",
    query: gql(`query { findScenes(filter: { per_page: 25 }) { count scenes { id title date rating100 organized o_counter play_count play_duration resume_time interactive } } }`),
    rewrite: rest("/api/scenes?perPage=25"),
  },
  {
    name: "With files (1 join)",
    query: gql(`query { findScenes(filter: { per_page: 25 }) { count scenes { id title date files { path size duration width height video_codec audio_codec frame_rate bit_rate } } } }`),
    rewrite: rest("/api/scenes?perPage=25"),
  },
  {
    name: "With files + fingerprints (2 joins)",
    query: gql(`query { findScenes(filter: { per_page: 25 }) { count scenes { id title files { path size duration fingerprints { type value } } } } }`),
    rewrite: rest("/api/scenes?perPage=25"),
  },
  {
    name: "Full kitchen sink response",
    query: gql(`query { findScenes(filter: { per_page: 25 }) { count scenes { id title code details director date rating100 organized o_counter play_count play_duration resume_time interactive last_played_at urls stash_ids { endpoint stash_id } files { id path basename size duration width height video_codec audio_codec frame_rate bit_rate fingerprints { type value } } scene_markers { id title seconds primary_tag { id name } tags { id name } } paths { screenshot } } } }`),
    rewrite: rest("/api/scenes?perPage=25"),
  },
];

// ── Page size scaling ───────────────────────────────────────

const pageSizes = [1, 5, 10, 25, 50, 100, 200];

// ── Concurrent scaling ──────────────────────────────────────

const concurrencyLevels = [1, 2, 5, 10, 20];

// ── Main runner ─────────────────────────────────────────────

async function main() {
  console.log("╔══════════════════════════════════════════════════════════════════════╗");
  console.log("║     QUERY COMPLEXITY SCALING BENCHMARK: Original vs Rewrite          ║");
  console.log("╠══════════════════════════════════════════════════════════════════════╣");
  console.log(`║  Iterations: ${ITERS}   Warmup: ${WARMUP}                                          ║`);
  console.log("╚══════════════════════════════════════════════════════════════════════╝\n");

  // Connectivity
  let origUp = false, rewUp = false;
  try { const r = await timedFetch(`${ORIG}/graphql`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ query: "{ stats { scene_count } }" })}); origUp = r.ok; } catch {}
  try { const r = await timedFetch(`${REW}/api/system/stats`); rewUp = r.ok; } catch {}
  console.log(`  Original: ${origUp ? "✅" : "❌"}  Rewrite: ${rewUp ? "✅" : "❌"}\n`);

  if (!origUp || !rewUp) { console.log("Both servers must be running."); process.exit(1); }

  const allResults = {};

  // ═══════════════════════════════════════════════════════════
  // SECTION 1: Complexity Tier Scaling
  // ═══════════════════════════════════════════════════════════
  console.log("\n" + "═".repeat(72));
  console.log("  SECTION 1: QUERY COMPLEXITY TIER SCALING");
  console.log("  Testing how latency grows as queries get more complex");
  console.log("═".repeat(72));

  const tierResults = [];
  for (const tier of complexityTiers) {
    process.stdout.write(`\n  Tier ${tier.tier}: ${tier.name}\n    ${tier.description}\n`);

    // Run original first, then rewrite (sequential to avoid interference)
    const origResult = await bench(tier.original);
    const rewResult = await bench(tier.rewrite);

    const oAvg = origResult.stats?.avg;
    const rAvg = rewResult.stats?.avg;
    const ratio = oAvg && rAvg ? (rAvg / oAvg).toFixed(1) : "N/A";

    console.log(`    Original: avg=${fmtMs(oAvg)}  p50=${fmtMs(origResult.stats?.p50)}  p95=${fmtMs(origResult.stats?.p95)}  (${origResult.avgSize}B)`);
    console.log(`    Rewrite:  avg=${fmtMs(rAvg)}  p50=${fmtMs(rewResult.stats?.p50)}  p95=${fmtMs(rewResult.stats?.p95)}  (${rewResult.avgSize}B)`);
    console.log(`    Ratio: Rewrite is ${ratio}x slower`);

    tierResults.push({
      tier: tier.tier,
      name: tier.name,
      original: origResult.stats,
      rewrite: rewResult.stats,
      origSize: origResult.avgSize,
      rewSize: rewResult.avgSize,
      ratio: +ratio || null,
    });
  }
  allResults.tierScaling = tierResults;

  // Print scaling summary
  console.log("\n  ┌──────┬─────────────────────────────────┬───────────┬───────────┬────────┐");
  console.log("  │ Tier │ Description                     │ Orig avg  │ Rew avg   │ Ratio  │");
  console.log("  ├──────┼─────────────────────────────────┼───────────┼───────────┼────────┤");
  for (const t of tierResults) {
    console.log(`  │ ${String(t.tier).padEnd(4)} │ ${t.name.padEnd(31).slice(0, 31)} │ ${fmtMs(t.original?.avg).padStart(9)} │ ${fmtMs(t.rewrite?.avg).padStart(9)} │ ${(t.ratio ? t.ratio + "x" : "N/A").padStart(6)} │`);
  }
  console.log("  └──────┴─────────────────────────────────┴───────────┴───────────┴────────┘");

  // Compute scaling trend
  const origTierAvgs = tierResults.filter(t => t.original).map(t => t.original.avg);
  const rewTierAvgs = tierResults.filter(t => t.rewrite).map(t => t.rewrite.avg);
  if (origTierAvgs.length > 1 && rewTierAvgs.length > 1) {
    const origSlope = (origTierAvgs[origTierAvgs.length - 1] - origTierAvgs[0]) / (origTierAvgs.length - 1);
    const rewSlope = (rewTierAvgs[rewTierAvgs.length - 1] - rewTierAvgs[0]) / (rewTierAvgs.length - 1);
    console.log(`\n  Scaling gradient (latency increase per tier):`);
    console.log(`    Original: +${origSlope.toFixed(2)}ms/tier`);
    console.log(`    Rewrite:  +${rewSlope.toFixed(2)}ms/tier`);
    console.log(`    ${rewSlope > origSlope ? "⚠️  Rewrite scales WORSE" : "✅ Rewrite scales BETTER"} with complexity (${(rewSlope / origSlope).toFixed(1)}x gradient)`);
  }

  // ═══════════════════════════════════════════════════════════
  // SECTION 2: GraphQL Field Selection Advantage
  // ═══════════════════════════════════════════════════════════
  console.log("\n\n" + "═".repeat(72));
  console.log("  SECTION 2: GRAPHQL FIELD SELECTION vs REST FIXED RESPONSE");
  console.log("  GQL can request minimal fields; REST always returns full DTO");
  console.log("═".repeat(72));

  const fieldResults = [];
  for (const test of gqlDepthTests) {
    process.stdout.write(`\n  ▸ ${test.name}\n`);
    const origResult = await bench(test.query);
    const rewResult = await bench(test.rewrite);

    console.log(`    Original: avg=${fmtMs(origResult.stats?.avg)}  (${origResult.avgSize}B)`);
    console.log(`    Rewrite:  avg=${fmtMs(rewResult.stats?.avg)}  (${rewResult.avgSize}B)`);
    console.log(`    Size ratio: Rewrite response is ${(rewResult.avgSize / origResult.avgSize).toFixed(1)}x larger`);

    fieldResults.push({
      name: test.name,
      original: origResult.stats,
      rewrite: rewResult.stats,
      origSize: origResult.avgSize,
      rewSize: rewResult.avgSize,
    });
  }
  allResults.fieldSelection = fieldResults;

  // ═══════════════════════════════════════════════════════════
  // SECTION 3: Page Size Scaling
  // ═══════════════════════════════════════════════════════════
  console.log("\n\n" + "═".repeat(72));
  console.log("  SECTION 3: PAGE SIZE SCALING");
  console.log("  How latency grows with result set size");
  console.log("═".repeat(72));

  const pageSizeResults = [];
  for (const ps of pageSizes) {
    process.stdout.write(`\n  ▸ ${ps} items per page\n`);
    const origResult = await bench(gql(`query { findScenes(filter: { per_page: ${ps}, sort: "date", direction: DESC }) { count scenes { id title date rating100 files { path size duration width height } } } }`));
    const rewResult = await bench(rest(`/api/scenes?perPage=${ps}&sort=date&direction=desc`));

    console.log(`    Original: avg=${fmtMs(origResult.stats?.avg)}  (${origResult.avgSize}B)`);
    console.log(`    Rewrite:  avg=${fmtMs(rewResult.stats?.avg)}  (${rewResult.avgSize}B)`);

    pageSizeResults.push({
      pageSize: ps,
      original: origResult.stats,
      rewrite: rewResult.stats,
      origSize: origResult.avgSize,
      rewSize: rewResult.avgSize,
    });
  }
  allResults.pageSizeScaling = pageSizeResults;

  // Print page size scaling table
  console.log("\n  ┌──────────┬───────────┬───────────┬────────┬──────────┬──────────┐");
  console.log("  │ PageSize │ Orig avg  │ Rew avg   │ Ratio  │ Orig sz  │ Rew sz   │");
  console.log("  ├──────────┼───────────┼───────────┼────────┼──────────┼──────────┤");
  for (const p of pageSizeResults) {
    const ratio = p.original?.avg && p.rewrite?.avg ? (p.rewrite.avg / p.original.avg).toFixed(1) + "x" : "N/A";
    console.log(`  │ ${String(p.pageSize).padEnd(8)} │ ${fmtMs(p.original?.avg).padStart(9)} │ ${fmtMs(p.rewrite?.avg).padStart(9)} │ ${ratio.padStart(6)} │ ${(p.origSize / 1024).toFixed(1).padStart(6)}KB │ ${(p.rewSize / 1024).toFixed(1).padStart(6)}KB │`);
  }
  console.log("  └──────────┴───────────┴───────────┴────────┴──────────┴──────────┘");

  // Compute page size scaling
  if (pageSizeResults.length > 1) {
    const first = pageSizeResults[0];
    const last = pageSizeResults[pageSizeResults.length - 1];
    const origGrowth = last.original?.avg / first.original?.avg;
    const rewGrowth = last.rewrite?.avg / first.rewrite?.avg;
    console.log(`\n  Page size scaling (1 → ${last.pageSize} items):`);
    console.log(`    Original latency grew ${origGrowth?.toFixed(1)}x`);
    console.log(`    Rewrite latency grew ${rewGrowth?.toFixed(1)}x`);
    console.log(`    ${rewGrowth > origGrowth ? "⚠️  Rewrite scales WORSE" : "✅ Rewrite scales BETTER"} with page size`);
  }

  // ═══════════════════════════════════════════════════════════
  // SECTION 4: Concurrent Scaling
  // ═══════════════════════════════════════════════════════════
  console.log("\n\n" + "═".repeat(72));
  console.log("  SECTION 4: CONCURRENT LOAD SCALING");
  console.log("  How throughput and latency change under concurrent load");
  console.log("═".repeat(72));

  // Simple query under load (count — shows raw framework overhead)
  console.log("\n  Query: COUNT (minimal — shows framework overhead)\n");
  const concCountResults = [];
  for (const c of concurrencyLevels) {
    const origR = await concurrentBench(gql(`query { findScenes(filter: { per_page: 0 }) { count } }`), c, 100);
    const rewR = await concurrentBench(rest("/api/scenes?perPage=0"), c, 100);
    console.log(`    ${String(c).padStart(2)} users:  Original: avg=${fmtMs(origR.stats?.avg)} ${origR.throughput} req/s  |  Rewrite: avg=${fmtMs(rewR.stats?.avg)} ${rewR.throughput} req/s`);
    concCountResults.push({ concurrency: c, original: { ...origR.stats, throughput: origR.throughput }, rewrite: { ...rewR.stats, throughput: rewR.throughput } });
  }

  // Complex query under load
  console.log("\n  Query: Complex (5 filters + sort + search)\n");
  const concComplexResults = [];
  for (const c of concurrencyLevels) {
    const origR = await concurrentBench(
      gql(`query { findScenes(scene_filter: { duration: { value: 300, modifier: GREATER_THAN }, resolution: { value: STANDARD_HD, modifier: GREATER_THAN }, organized: false, framerate: { value: 24, modifier: GREATER_THAN }, bitrate: { value: 5000, modifier: GREATER_THAN } }, filter: { per_page: 25, sort: "duration", direction: DESC }) { count scenes { id title files { duration } } } }`),
      c, 100);
    const rewR = await concurrentBench(
      rest("/api/scenes/find", "POST", {
        findFilter: { page: 1, perPage: 25, sort: "duration", direction: "desc" },
        objectFilter: { durationCriterion: { value: 300, modifier: "greaterThan" }, resolutionCriterion: { value: 720, modifier: "greaterThan" }, organized: false, frameRateCriterion: { value: 24, modifier: "greaterThan" }, bitrateInterval: { value: 5000, modifier: "greaterThan" } }
      }),
      c, 100);
    console.log(`    ${String(c).padStart(2)} users:  Original: avg=${fmtMs(origR.stats?.avg)} ${origR.throughput} req/s  |  Rewrite: avg=${fmtMs(rewR.stats?.avg)} ${rewR.throughput} req/s`);
    concComplexResults.push({ concurrency: c, original: { ...origR.stats, throughput: origR.throughput }, rewrite: { ...rewR.stats, throughput: rewR.throughput } });
  }
  allResults.concurrentCount = concCountResults;
  allResults.concurrentComplex = concComplexResults;

  // ═══════════════════════════════════════════════════════════
  // SECTION 5: GQL-Only Features (Rewrite cannot match)
  // ═══════════════════════════════════════════════════════════
  console.log("\n\n" + "═".repeat(72));
  console.log("  SECTION 5: GRAPHQL-EXCLUSIVE QUERY PATTERNS");
  console.log("  Queries only possible with GraphQL (no REST equivalent)");
  console.log("═".repeat(72));

  const gqlOnlyTests = [
    {
      name: "Boolean composition (OR)",
      description: "Scenes with duration > 1800 OR organized=true",
      query: gql(`query { findScenes(scene_filter: { OR: { duration: { value: 1800, modifier: GREATER_THAN }, organized: true } }, filter: { per_page: 25 }) { count scenes { id title } } }`),
    },
    {
      name: "Nested AND + OR",
      description: "(duration>600 AND resolution>=720) OR organized=true",
      query: gql(`query { findScenes(scene_filter: { AND: { duration: { value: 600, modifier: GREATER_THAN }, resolution: { value: STANDARD_HD, modifier: GREATER_THAN } }, OR: { organized: true } }, filter: { per_page: 25 }) { count scenes { id title } } }`),
    },
    {
      name: "NOT filter",
      description: "Scenes NOT organized",
      query: gql(`query { findScenes(scene_filter: { NOT: { organized: true } }, filter: { per_page: 25 }) { count scenes { id title } } }`),
    },
    {
      name: "Aggregated response (count + duration + filesize)",
      description: "Aggregate totals returned with query",
      query: gql(`query { findScenes(filter: { per_page: 0 }) { count duration filesize } }`),
    },
  ];

  for (const test of gqlOnlyTests) {
    process.stdout.write(`\n  ▸ ${test.name}\n    ${test.description}\n`);
    const result = await bench(test.query);
    console.log(`    Original: avg=${fmtMs(result.stats?.avg)}  p95=${fmtMs(result.stats?.p95)}  (${result.successCount}/${ITERS} ok)`);
    console.log(`    Rewrite:  N/A (no REST equivalent)`);
  }

  // ═══════════════════════════════════════════════════════════
  // SECTION 6: Cold vs Warm (Cache Effect)
  // ═══════════════════════════════════════════════════════════
  console.log("\n\n" + "═".repeat(72));
  console.log("  SECTION 6: COLD START vs WARM (Cache Impact)");
  console.log("  First request after 2s wait vs subsequent requests");
  console.log("═".repeat(72));

  const cacheTests = [
    { name: "Scene list (25)", orig: gql(`query { findScenes(filter: { per_page: 25 }) { count scenes { id title } } }`), rew: rest("/api/scenes?perPage=25") },
    { name: "Scene by ID", orig: gql(`query { findScene(id: 1) { id title } }`), rew: rest("/api/scenes/1") },
  ];

  for (const test of cacheTests) {
    console.log(`\n  ▸ ${test.name}`);
    
    // Cold: wait 2s (cache expires at 1s), then one request
    await new Promise(r => setTimeout(r, 2000));
    const origCold = await timedFetch(test.orig.url, test.orig);
    await new Promise(r => setTimeout(r, 2000));
    const rewCold = await timedFetch(test.rew.url, test.rew);

    // Warm: 10 rapid requests
    const origWarm = [];
    for (let i = 0; i < 10; i++) { const r = await timedFetch(test.orig.url, test.orig); if (r.ok) origWarm.push(r.ms); }
    const rewWarm = [];
    for (let i = 0; i < 10; i++) { const r = await timedFetch(test.rew.url, test.rew); if (r.ok) rewWarm.push(r.ms); }

    console.log(`    Original: cold=${fmtMs(origCold.ms)}  warm avg=${fmtMs(stats(origWarm)?.avg)}`);
    console.log(`    Rewrite:  cold=${fmtMs(rewCold.ms)}  warm avg=${fmtMs(stats(rewWarm)?.avg)}`);
  }

  // ═══════════════════════════════════════════════════════════
  // Final Analysis
  // ═══════════════════════════════════════════════════════════
  console.log("\n\n" + "═".repeat(72));
  console.log("  FINAL ANALYSIS");
  console.log("═".repeat(72));

  // Compute overall scaling characteristics
  const ratios = tierResults.filter(t => t.ratio).map(t => t.ratio);
  const avgRatio = ratios.reduce((a, b) => a + b) / ratios.length;
  const minRatio = Math.min(...ratios);
  const maxRatio = Math.max(...ratios);
  const simpleTierRatios = tierResults.filter(t => t.tier <= 3 && t.ratio).map(t => t.ratio);
  const complexTierRatios = tierResults.filter(t => t.tier >= 7 && t.ratio).map(t => t.ratio);
  const simpleAvgRatio = simpleTierRatios.length ? simpleTierRatios.reduce((a, b) => a + b) / simpleTierRatios.length : 0;
  const complexAvgRatio = complexTierRatios.length ? complexTierRatios.reduce((a, b) => a + b) / complexTierRatios.length : 0;

  console.log(`\n  Rewrite/Original latency ratio across all tiers:`);
  console.log(`    Average: ${avgRatio.toFixed(1)}x`);
  console.log(`    Range: ${minRatio.toFixed(1)}x – ${maxRatio.toFixed(1)}x`);
  console.log(`    Simple queries (tier 0-3): ${simpleAvgRatio.toFixed(1)}x`);
  console.log(`    Complex queries (tier 7+): ${complexAvgRatio.toFixed(1)}x`);

  if (complexAvgRatio < simpleAvgRatio) {
    console.log(`\n  ✅ GAP NARROWS WITH COMPLEXITY: Complex queries have ${simpleAvgRatio.toFixed(1)}x → ${complexAvgRatio.toFixed(1)}x ratio`);
    console.log("     PostgreSQL's query planner potentially handles complex queries well");
  } else if (complexAvgRatio > simpleAvgRatio * 1.3) {
    console.log(`\n  ⚠️  GAP WIDENS WITH COMPLEXITY: Simple ${simpleAvgRatio.toFixed(1)}x → Complex ${complexAvgRatio.toFixed(1)}x`);
    console.log("     Consider GraphQL for complex query patterns");
  } else {
    console.log(`\n  ➡️  GAP IS ROUGHLY CONSTANT across complexity levels`);
    console.log("     The difference is mostly fixed overhead, not query-plan related");
  }

  // Save results
  writeFileSync("benchmarks/complexity-results.json", JSON.stringify(allResults, null, 2));
  console.log("\n  📁 Results saved to benchmarks/complexity-results.json");
  console.log("\n✅ Complexity scaling benchmark complete.\n");
}

main().catch(console.error);
