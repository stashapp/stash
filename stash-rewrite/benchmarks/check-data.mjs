const q = `{ findScenes(filter:{per_page:3,sort:"file_size",direction:DESC}){count scenes{id title date rating100 o_counter play_count organized interactive resume_time play_duration files{path duration width height size video_codec audio_codec frame_rate bit_rate fingerprints{type value}}}} }`;
const r = await fetch('http://localhost:9998/graphql', { method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify({query:q}) });
const d = await r.json();
console.log(JSON.stringify(d, null, 2).slice(0, 3000));

// Check rewrite too
const r2 = await fetch('http://localhost:9999/api/scenes?perPage=3&sort=file_size&direction=desc');
const d2 = await r2.json();
console.log("\n--- REWRITE ---");
console.log(JSON.stringify(d2.items?.[0], null, 2));
