/**
 * Benchmark harness comparing Original Stash (GraphQL) vs Rewrite (REST).
 * Usage: node benchmarks/run.mjs [--original http://localhost:9998] [--rewrite http://localhost:9999]
 */

const args = process.argv.slice(2);
function getArg(name, defaultVal) {
  const idx = args.indexOf(`--${name}`);
  return idx >= 0 && args[idx + 1] ? args[idx + 1] : defaultVal;
}

const ORIGINAL_URL = getArg("original", "http://localhost:9998");
const REWRITE_URL = getArg("rewrite", "http://localhost:9999");
const ITERATIONS = parseInt(getArg("iterations", "50"), 10);
const WARMUP = 5;

// ── Scenario definitions ──────────────────────────────────────────────

const scenarios = [
  {
    name: "List scenes (page 1, 25 items)",
    original: {
      url: `${ORIGINAL_URL}/graphql`,
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { findScenes(filter: { page: 1, per_page: 25, sort: "date", direction: DESC }) { count scenes { id title date } } }`,
      }),
    },
    rewrite: {
      url: `${REWRITE_URL}/api/scenes?page=1&perPage=25&sort=date&direction=desc`,
      method: "GET",
    },
  },
  {
    name: "Get scene by ID (id=1)",
    original: {
      url: `${ORIGINAL_URL}/graphql`,
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { findScene(id: 1) { id title date details } }`,
      }),
    },
    rewrite: {
      url: `${REWRITE_URL}/api/scenes/1`,
      method: "GET",
    },
  },
  {
    name: "List performers (page 1, 25 items)",
    original: {
      url: `${ORIGINAL_URL}/graphql`,
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { findPerformers(filter: { page: 1, per_page: 25, sort: "name", direction: ASC }) { count performers { id name } } }`,
      }),
    },
    rewrite: {
      url: `${REWRITE_URL}/api/performers?page=1&perPage=25&sort=name&direction=asc`,
      method: "GET",
    },
  },
  {
    name: "System stats",
    original: {
      url: `${ORIGINAL_URL}/graphql`,
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { stats { scene_count image_count gallery_count performer_count studio_count tag_count } }`,
      }),
    },
    rewrite: {
      url: `${REWRITE_URL}/api/system/stats`,
      method: "GET",
    },
  },
  {
    name: "List tags (all)",
    original: {
      url: `${ORIGINAL_URL}/graphql`,
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: `query { findTags(filter: { page: 1, per_page: 100, sort: "name", direction: ASC }) { count tags { id name } } }`,
      }),
    },
    rewrite: {
      url: `${REWRITE_URL}/api/tags?page=1&perPage=100&sort=name&direction=asc`,
      method: "GET",
    },
  },
];

// ── Benchmark runner ──────────────────────────────────────────────────

async function runRequest(config) {
  const start = performance.now();
  try {
    const res = await fetch(config.url, {
      method: config.method,
      headers: config.headers,
      body: config.body,
    });
    const elapsed = performance.now() - start;
    const ok = res.ok;
    return { elapsed, ok, status: res.status };
  } catch {
    return { elapsed: performance.now() - start, ok: false, status: 0 };
  }
}

function percentile(sorted, pct) {
  const idx = Math.ceil((pct / 100) * sorted.length) - 1;
  return sorted[Math.max(0, idx)];
}

function formatMs(ms) {
  return ms < 1 ? `${(ms * 1000).toFixed(0)}µs` : `${ms.toFixed(1)}ms`;
}

async function benchmarkScenario(scenario) {
  const results = { original: [], rewrite: [] };

  // Warmup
  for (let i = 0; i < WARMUP; i++) {
    await runRequest(scenario.original);
    await runRequest(scenario.rewrite);
  }

  // Actual runs
  for (let i = 0; i < ITERATIONS; i++) {
    const orig = await runRequest(scenario.original);
    const rew = await runRequest(scenario.rewrite);
    if (orig.ok) results.original.push(orig.elapsed);
    if (rew.ok) results.rewrite.push(rew.elapsed);
  }

  return results;
}

function reportResults(name, results) {
  console.log(`\n━━━ ${name} ━━━`);

  for (const [label, times] of [["Original (Go/SQLite)", results.original], ["Rewrite  (.NET/PG) ", results.rewrite]]) {
    if (times.length === 0) {
      console.log(`  ${label}: ❌ No successful responses`);
      continue;
    }
    const sorted = [...times].sort((a, b) => a - b);
    const avg = times.reduce((s, t) => s + t, 0) / times.length;
    const p50 = percentile(sorted, 50);
    const p95 = percentile(sorted, 95);
    const p99 = percentile(sorted, 99);
    const rps = 1000 / avg;
    console.log(
      `  ${label}: avg=${formatMs(avg)}  p50=${formatMs(p50)}  p95=${formatMs(p95)}  p99=${formatMs(p99)}  ~${rps.toFixed(0)} req/s  (${times.length}/${ITERATIONS} ok)`
    );
  }

  // Comparison
  if (results.original.length > 0 && results.rewrite.length > 0) {
    const avgOrig = results.original.reduce((s, t) => s + t, 0) / results.original.length;
    const avgRew = results.rewrite.reduce((s, t) => s + t, 0) / results.rewrite.length;
    const ratio = avgOrig / avgRew;
    const winner = ratio > 1 ? "Rewrite" : "Original";
    const factor = ratio > 1 ? ratio : 1 / ratio;
    console.log(`  → ${winner} is ${factor.toFixed(2)}x faster`);
  }
}

// ── Main ──────────────────────────────────────────────────────────────

async function main() {
  console.log("╔══════════════════════════════════════════════════╗");
  console.log("║   Stash Benchmark: Original vs Rewrite          ║");
  console.log("╠══════════════════════════════════════════════════╣");
  console.log(`║ Original: ${ORIGINAL_URL.padEnd(38)}║`);
  console.log(`║ Rewrite:  ${REWRITE_URL.padEnd(38)}║`);
  console.log(`║ Iterations: ${String(ITERATIONS).padEnd(36)}║`);
  console.log(`║ Warmup: ${String(WARMUP).padEnd(40)}║`);
  console.log("╚══════════════════════════════════════════════════╝");

  // Check connectivity
  console.log("\nChecking connectivity...");
  let origAlive = false, rewAlive = false;
  try { const r = await fetch(`${ORIGINAL_URL}/graphql`, { method: "POST", headers: { "Content-Type": "application/json" }, body: '{"query":"{ stats { scene_count } }"}' }); origAlive = r.ok; } catch {}
  try { const r = await fetch(`${REWRITE_URL}/api/system/stats`); rewAlive = r.ok; } catch {}

  console.log(`  Original: ${origAlive ? "✅ Connected" : "❌ Not reachable"}`);
  console.log(`  Rewrite:  ${rewAlive ? "✅ Connected" : "❌ Not reachable"}`);

  if (!origAlive && !rewAlive) {
    console.log("\n⚠️  Neither server is running. Start at least one to benchmark.");
    process.exit(1);
  }

  for (const scenario of scenarios) {
    const results = await benchmarkScenario(scenario);
    reportResults(scenario.name, results);
  }

  console.log("\n✅ Benchmark complete.");
}

main().catch(console.error);
