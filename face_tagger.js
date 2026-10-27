/**
 * ============================================================================
 * Stash Unified Image Zoom Navigator & Native Face Identification Plugin
 * Compatible with: https://github.com/ucbbrla9011-jpg/stash232323.git
 * ============================================================================
 *
 * Features:
 * 1. Interactive Pan & Zoom Image Viewer:
 *    - Cursor-anchored smooth wheel zoom (0.1x to 12x), drag-to-pan, 1:1 pixel
 *      toggle (Z), fit-to-screen, and interactive bottom-right Minimap.
 *    - Keyboard navigation (Left/Right arrows, A/D, 1-5 ratings, O-counter,
 *      B to toggle boxes, M for manual face box drawing, Esc to close).
 *    - Bottom thumbnail filmstrip for rapid gallery/collection scrubbing.
 *
 * 2. Native Facial Identification & Zero-Dependency Fallback Engine:
 *    - Tier 1: Hardware/OS accelerated W3C Shape Detection API (`window.FaceDetector`)
 *      leveraging Apple Vision (macOS), Windows Media FaceAnalysis, or Android/ChromeOS.
 *    - Tier 2: Pure-Canvas Integral-Image YCbCr Skin + Haar-Cascade Bilateral
 *      Symmetry detector so face detection works out-of-the-box on every browser
 *      without Python/dlib/C++ dependencies.
 *    - 128-D L2-Normalized Face Embedding Engine (8x8 spatial gradient/luminance
 *      descriptor + 4x4 chrominance contrast + eye-aligned symmetry + 2D-DCT pHash).
 *
 * 3. Unsupervised Face Clustering & Smart Image Sorting Workspace (/plugins/face-tagger):
 *    - Groups similar faces across images via Cosine Similarity clustering with
 *      a live interactive threshold slider.
 *    - Matches detected faces against existing Stash Performer profile images.
 *    - One-click bulk Performer labeling via Stash GraphQL (`BulkImageUpdate` /
 *      `ImageUpdate`) and inline Performer creation (`PerformerCreate`).
 *    - Smart Image Sorter: sort & filter images by Unassigned Faces, Face Count,
 *      Cluster Group, Rating, O-Counter, Resolution, or Date.
 */

(function () {
  "use strict";

  const PLUGIN_ID = "stash-native-face-tagger";
  const STORAGE_KEY_FACE_CACHE = "stash_native_face_tagger_cache_v1";
  const STORAGE_KEY_PERFORMER_DB = "stash_native_face_tagger_performers_v1";

  // Prevent duplicate registration if loaded both via App.tsx and plugins/ folder
  if (window.__STASH_NATIVE_FACE_TAGGER_LOADED__) {
    return;
  }
  window.__STASH_NATIVE_FACE_TAGGER_LOADED__ = true;

  // -------------------------------------------------------------------------
  // Self-Contained CSS Injection (Zero External CSS Dependency Required)
  // -------------------------------------------------------------------------
  function injectStyles() {
    if (document.getElementById("stash-face-tagger-styles")) return;
    const style = document.createElement("style");
    style.id = "stash-face-tagger-styles";
    style.textContent = `
      .sft-workspace {
        padding: 1.25rem 1.5rem;
        color: #f8f9fa;
        min-height: calc(100vh - 4rem);
      }
      .sft-header-bar {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        justify-content: space-between;
        gap: 1rem;
        background: rgba(30, 34, 42, 0.92);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        padding: 0.85rem 1.15rem;
        margin-bottom: 1.15rem;
      }
      .sft-badge-engine {
        display: inline-flex;
        align-items: center;
        gap: 0.4rem;
        font-size: 0.75rem;
        font-weight: 600;
        padding: 0.25rem 0.6rem;
        border-radius: 999px;
        background: rgba(40, 167, 69, 0.2);
        color: #5dd879;
        border: 1px solid rgba(40, 167, 69, 0.45);
      }
      .sft-badge-engine.fallback {
        background: rgba(23, 162, 184, 0.2);
        color: #6edff6;
        border-color: rgba(23, 162, 184, 0.45);
      }
      .sft-controls-row {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 0.75rem;
        background: rgba(24, 27, 34, 0.85);
        border: 1px solid rgba(255, 255, 255, 0.08);
        border-radius: 8px;
        padding: 0.75rem 1rem;
        margin-bottom: 1rem;
      }
      .sft-control-group {
        display: flex;
        align-items: center;
        gap: 0.45rem;
        font-size: 0.85rem;
      }
      .sft-control-group label {
        margin: 0;
        color: #adb5bd;
        font-weight: 500;
        white-space: nowrap;
      }
      .sft-select, .sft-input {
        background: #15181e;
        color: #f8f9fa;
        border: 1px solid rgba(255, 255, 255, 0.18);
        border-radius: 5px;
        padding: 0.3rem 0.6rem;
        font-size: 0.85rem;
      }
      .sft-cluster-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
        gap: 1.1rem;
      }
      .sft-cluster-card {
        background: rgba(32, 36, 45, 0.94);
        border: 1px solid rgba(255, 255, 255, 0.12);
        border-radius: 10px;
        padding: 1rem;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        box-shadow: 0 4px 14px rgba(0, 0, 0, 0.28);
        transition: border-color 0.15s ease;
      }
      .sft-cluster-card:hover {
        border-color: rgba(102, 175, 233, 0.55);
      }
      .sft-cluster-card.labeled {
        border-color: rgba(40, 167, 69, 0.6);
        background: rgba(26, 42, 34, 0.85);
      }
      .sft-cluster-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 0.5rem;
      }
      .sft-faces-strip {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        max-height: 176px;
        overflow-y: auto;
        padding: 0.25rem;
        background: rgba(15, 17, 22, 0.6);
        border-radius: 6px;
      }
      .sft-face-thumb-wrap {
        position: relative;
        width: 72px;
        height: 72px;
        border-radius: 6px;
        overflow: hidden;
        cursor: pointer;
        border: 2px solid rgba(255, 255, 255, 0.18);
        transition: transform 0.12s ease, border-color 0.12s ease;
      }
      .sft-face-thumb-wrap:hover {
        transform: scale(1.06);
        border-color: #4dabf7;
      }
      .sft-face-thumb-wrap.excluded {
        opacity: 0.32;
        border-color: #dc3545;
        filter: grayscale(1);
      }
      .sft-face-thumb-wrap img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
      .sft-face-exclude-btn {
        position: absolute;
        top: 2px;
        right: 2px;
        width: 18px;
        height: 18px;
        font-size: 10px;
        line-height: 16px;
        text-align: center;
        border-radius: 50%;
        background: rgba(0, 0, 0, 0.72);
        color: #fff;
        border: none;
        cursor: pointer;
      }
      .sft-face-exclude-btn:hover {
        background: #dc3545;
      }
      .sft-image-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
        gap: 1rem;
      }
      .sft-image-card {
        background: rgba(30, 34, 42, 0.94);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 8px;
        overflow: hidden;
        display: flex;
        flex-direction: column;
        transition: transform 0.12s ease, border-color 0.12s ease;
      }
      .sft-image-card:hover {
        border-color: #4dabf7;
        transform: translateY(-2px);
      }
      .sft-image-card.selected {
        border-color: #339af0;
        box-shadow: 0 0 0 2px rgba(51, 154, 240, 0.45);
      }
      .sft-image-thumb-box {
        position: relative;
        width: 100%;
        height: 220px;
        background: #0d0f14;
        cursor: pointer;
        overflow: hidden;
      }
      .sft-image-thumb-box img {
        width: 100%;
        height: 100%;
        object-fit: contain;
      }
      .sft-card-badges {
        position: absolute;
        top: 6px;
        right: 6px;
        display: flex;
        flex-direction: column;
        gap: 4px;
        align-items: flex-end;
      }
      .sft-pill {
        font-size: 0.72rem;
        font-weight: 600;
        padding: 2px 7px;
        border-radius: 999px;
        background: rgba(15, 17, 23, 0.85);
        color: #fff;
        border: 1px solid rgba(255, 255, 255, 0.18);
      }
      .sft-pill.unassigned {
        background: rgba(255, 193, 7, 0.9);
        color: #111;
      }
      .sft-pill.assigned {
        background: rgba(40, 167, 69, 0.88);
        color: #fff;
      }

      /* Fullscreen Pan & Zoom Modal */
      .sft-modal-backdrop {
        position: fixed;
        inset: 0;
        z-index: 9999;
        background: rgba(8, 10, 14, 0.96);
        display: flex;
        flex-direction: column;
        user-select: none;
      }
      .sft-modal-topbar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.55rem 1rem;
        background: rgba(20, 23, 30, 0.95);
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        gap: 0.75rem;
        flex-wrap: wrap;
        z-index: 10;
      }
      .sft-modal-title {
        font-weight: 600;
        font-size: 0.95rem;
        color: #f8f9fa;
        max-width: 340px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      .sft-modal-toolbar {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        flex-wrap: wrap;
      }
      .sft-viewport {
        position: relative;
        flex: 1;
        overflow: hidden;
        cursor: grab;
        background: radial-gradient(circle at center, #181c24 0%, #0b0d12 100%);
      }
      .sft-viewport.panning {
        cursor: grabbing;
      }
      .sft-viewport.drawing {
        cursor: crosshair;
      }
      .sft-transform-layer {
        position: absolute;
        top: 0;
        left: 0;
        transform-origin: 0 0;
        will-change: transform;
      }
      .sft-main-img {
        display: block;
        max-width: none;
        max-height: none;
        pointer-events: none;
        user-select: none;
      }
      .sft-face-box {
        position: absolute;
        border: 2px solid #00e5ff;
        background: rgba(0, 229, 255, 0.08);
        border-radius: 4px;
        cursor: pointer;
        transition: background 0.12s ease, box-shadow 0.12s ease;
      }
      .sft-face-box:hover, .sft-face-box.active {
        background: rgba(0, 229, 255, 0.22);
        box-shadow: 0 0 0 2px rgba(0, 229, 255, 0.55), 0 4px 12px rgba(0, 0, 0, 0.5);
      }
      .sft-face-box.assigned {
        border-color: #40c057;
        background: rgba(64, 192, 87, 0.12);
      }
      .sft-face-box-label {
        position: absolute;
        left: -2px;
        bottom: calc(100% + 4px);
        background: rgba(13, 17, 23, 0.92);
        color: #fff;
        font-size: 11px;
        font-weight: 600;
        padding: 2px 6px;
        border-radius: 4px;
        white-space: nowrap;
        border: 1px solid rgba(255, 255, 255, 0.2);
        pointer-events: none;
      }
      .sft-landmark-dot {
        position: absolute;
        width: 5px;
        height: 5px;
        border-radius: 50%;
        background: #ffec99;
        transform: translate(-50%, -50%);
        pointer-events: none;
      }
      .sft-popover-card {
        position: absolute;
        top: calc(100% + 8px);
        left: 0;
        width: 280px;
        background: #1c2029;
        border: 1px solid rgba(255, 255, 255, 0.22);
        border-radius: 8px;
        padding: 0.7rem;
        z-index: 50;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.65);
        color: #f8f9fa;
        cursor: default;
      }
      .sft-minimap {
        position: absolute;
        right: 16px;
        bottom: 16px;
        width: 150px;
        height: 110px;
        background: rgba(12, 15, 20, 0.88);
        border: 1px solid rgba(255, 255, 255, 0.24);
        border-radius: 6px;
        overflow: hidden;
        z-index: 20;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
      }
      .sft-minimap img {
        width: 100%;
        height: 100%;
        object-fit: contain;
        opacity: 0.65;
      }
      .sft-minimap-rect {
        position: absolute;
        border: 1.5px solid #339af0;
        background: rgba(51, 154, 240, 0.22);
        pointer-events: none;
      }
      .sft-filmstrip {
        display: flex;
        align-items: center;
        gap: 0.45rem;
        padding: 0.45rem 0.85rem;
        background: rgba(16, 19, 25, 0.96);
        border-top: 1px solid rgba(255, 255, 255, 0.1);
        overflow-x: auto;
        height: 76px;
      }
      .sft-filmstrip-item {
        width: 56px;
        height: 56px;
        flex-shrink: 0;
        border-radius: 5px;
        overflow: hidden;
        border: 2px solid transparent;
        cursor: pointer;
        opacity: 0.6;
        transition: opacity 0.12s ease, border-color 0.12s ease;
      }
      .sft-filmstrip-item:hover {
        opacity: 0.9;
      }
      .sft-filmstrip-item.active {
        opacity: 1;
        border-color: #339af0;
      }
      .sft-filmstrip-item img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
      .sft-card-quick-btn {
        position: absolute;
        top: 6px;
        left: 6px;
        z-index: 5;
        background: rgba(15, 19, 26, 0.82);
        color: #74c0fc;
        border: 1px solid rgba(116, 192, 252, 0.4);
        border-radius: 5px;
        padding: 2px 7px;
        font-size: 0.72rem;
        font-weight: 600;
        cursor: pointer;
        backdrop-filter: blur(4px);
      }
      .sft-card-quick-btn:hover {
        background: #1c7ed6;
        color: #fff;
      }
    `;
    document.head.appendChild(style);
  }

  // -------------------------------------------------------------------------
  // Persistent Cache Helpers (LocalStorage for Face Embeddings & Prototypes)
  // -------------------------------------------------------------------------
  function loadFaceCache() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY_FACE_CACHE);
      return raw ? JSON.parse(raw) : {};
    } catch (e) {
      return {};
    }
  }

  function saveFaceCache(cache) {
    try {
      localStorage.setItem(STORAGE_KEY_FACE_CACHE, JSON.stringify(cache));
    } catch (e) {
      // If quota exceeded, trim oldest half of entries
      try {
        const keys = Object.keys(cache);
        const trimmed = {};
        keys.slice(Math.floor(keys.length / 2)).forEach((k) => {
          trimmed[k] = cache[k];
        });
        localStorage.setItem(STORAGE_KEY_FACE_CACHE, JSON.stringify(trimmed));
      } catch (_) {}
    }
  }

  function loadPerformerPrototypes() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY_PERFORMER_DB);
      return raw ? JSON.parse(raw) : {};
    } catch (e) {
      return {};
    }
  }

  function savePerformerPrototype(performerId, performerName, embedding) {
    if (!performerId || !embedding || !embedding.length) return;
    const db = loadPerformerPrototypes();
    if (!db[performerId]) {
      db[performerId] = {
        id: String(performerId),
        name: performerName || "",
        embeddings: [embedding],
        centroid: embedding.slice(),
      };
    } else {
      db[performerId].name = performerName || db[performerId].name;
      db[performerId].embeddings.push(embedding);
      if (db[performerId].embeddings.length > 10) {
        db[performerId].embeddings.shift();
      }
      db[performerId].centroid = computeCentroid(db[performerId].embeddings);
    }
    try {
      localStorage.setItem(STORAGE_KEY_PERFORMER_DB, JSON.stringify(db));
    } catch (_) {}
  }

  // -------------------------------------------------------------------------
  // Vector Math & 128-D Face Embedding + Cosine Similarity Clustering Engine
  // -------------------------------------------------------------------------
  function l2Normalize(vec) {
    let sumSq = 0;
    for (let i = 0; i < vec.length; i++) {
      sumSq += vec[i] * vec[i];
    }
    const norm = Math.sqrt(sumSq) || 1e-8;
    const out = new Array(vec.length);
    for (let i = 0; i < vec.length; i++) {
      out[i] = Number((vec[i] / norm).toFixed(5));
    }
    return out;
  }

  function cosineSimilarity(a, b) {
    if (!a || !b || a.length !== b.length) return 0;
    let dot = 0;
    let normA = 0;
    let normB = 0;
    for (let i = 0; i < a.length; i++) {
      dot += a[i] * b[i];
      normA += a[i] * a[i];
      normB += b[i] * b[i];
    }
    const denom = Math.sqrt(normA) * Math.sqrt(normB);
    return denom > 1e-8 ? dot / denom : 0;
  }

  function computeCentroid(embeddings) {
    if (!embeddings || !embeddings.length) return [];
    const dim = embeddings[0].length;
    const sum = new Array(dim).fill(0);
    for (const emb of embeddings) {
      for (let i = 0; i < dim; i++) {
        sum[i] += emb[i] || 0;
      }
    }
    return l2Normalize(sum);
  }

  /**
   * Extracts a 128-dimensional L2-normalized facial feature vector + dataURL thumbnail
   * from an HTMLImageElement given a face boundingBox { x, y, width, height } and
   * optional landmarks.
   */
  function extractFaceEmbeddingAndThumb(imgEl, box, landmarks) {
    const SIZE = 64;
    const canvas = document.createElement("canvas");
    canvas.width = SIZE;
    canvas.height = SIZE;
    const ctx = canvas.getContext("2d");

    const natW = imgEl.naturalWidth || imgEl.width || 1;
    const natH = imgEl.naturalHeight || imgEl.height || 1;

    // Slight padding around face box to preserve hairline and chin contour
    const padX = box.width * 0.08;
    const padY = box.height * 0.08;
    const sx = Math.max(0, box.x - padX);
    const sy = Math.max(0, box.y - padY);
    const sw = Math.min(natW - sx, box.width + padX * 2);
    const sh = Math.min(natH - sy, box.height + padY * 2);

    // Deskew rotation if left & right eye landmarks are available
    let angle = 0;
    if (Array.isArray(landmarks)) {
      const eyes = landmarks.filter((l) => l.type === "eye" && l.locations && l.locations.length);
      if (eyes.length >= 2) {
        const p1 = eyes[0].locations[0];
        const p2 = eyes[1].locations[0];
        const left = p1.x < p2.x ? p1 : p2;
        const right = p1.x < p2.x ? p2 : p1;
        angle = Math.atan2(right.y - left.y, right.x - left.x);
      }
    }

    ctx.save();
    if (Math.abs(angle) > 0.02 && Math.abs(angle) < 0.65) {
      ctx.translate(SIZE / 2, SIZE / 2);
      ctx.rotate(-angle);
      ctx.translate(-SIZE / 2, -SIZE / 2);
    }
    ctx.drawImage(imgEl, sx, sy, Math.max(1, sw), Math.max(1, sh), 0, 0, SIZE, SIZE);
    ctx.restore();

    // Generate compact thumbnail dataURL for the cluster review UI
    const thumbCanvas = document.createElement("canvas");
    thumbCanvas.width = 72;
    thumbCanvas.height = 72;
    const thumbCtx = thumbCanvas.getContext("2d");
    thumbCtx.drawImage(canvas, 0, 0, 72, 72);
    const thumbDataUrl = thumbCanvas.toDataURL("image/jpeg", 0.78);

    const imgData = ctx.getImageData(0, 0, SIZE, SIZE).data;
    const gray = new Float32Array(SIZE * SIZE);
    const cbArr = new Float32Array(SIZE * SIZE);
    const crArr = new Float32Array(SIZE * SIZE);

    let minY = 255;
    let maxY = 0;
    for (let i = 0, p = 0; i < imgData.length; i += 4, p++) {
      const r = imgData[i];
      const g = imgData[i + 1];
      const b = imgData[i + 2];
      const y = 0.299 * r + 0.587 * g + 0.114 * b;
      const cb = 128 - 0.168736 * r - 0.331264 * g + 0.5 * b;
      const cr = 128 + 0.5 * r - 0.418688 * g - 0.081312 * b;
      gray[p] = y;
      cbArr[p] = cb / 255.0;
      crArr[p] = cr / 255.0;
      if (y < minY) minY = y;
      if (y > maxY) maxY = y;
    }

    // Contrast stretch luminance for lighting invariance
    const rangeY = Math.max(10, maxY - minY);
    for (let p = 0; p < gray.length; p++) {
      gray[p] = (gray[p] - minY) / rangeY;
    }

    const features = [];

    // 1. 64 dimensions: 8x8 spatial blocks of mean luminance & gradient magnitude
    const cell = SIZE / 8; // 8px per cell
    for (let by = 0; by < 8; by++) {
      for (let bx = 0; bx < 8; bx++) {
        let sumLum = 0;
        let sumGrad = 0;
        let count = 0;
        for (let y = by * cell; y < (by + 1) * cell; y++) {
          for (let x = bx * cell; x < (bx + 1) * cell; x++) {
            const idx = y * SIZE + x;
            sumLum += gray[idx];
            if (x > 0 && x < SIZE - 1 && y > 0 && y < SIZE - 1) {
              const gx = gray[idx + 1] - gray[idx - 1];
              const gy = gray[idx + SIZE] - gray[idx - SIZE];
              sumGrad += Math.sqrt(gx * gx + gy * gy);
            }
            count++;
          }
        }
        features.push((sumLum / count) * 0.65 + (sumGrad / count) * 0.35);
      }
    }

    // 2. 32 dimensions: 4x4 spatial grid of Cb and Cr chrominance averages
    const cCell = SIZE / 4;
    for (let by = 0; by < 4; by++) {
      for (let bx = 0; bx < 4; bx++) {
        let sumCb = 0;
        let sumCr = 0;
        let count = 0;
        for (let y = by * cCell; y < (by + 1) * cCell; y++) {
          for (let x = bx * cCell; x < (bx + 1) * cCell; x++) {
            const idx = y * SIZE + x;
            sumCb += cbArr[idx];
            sumCr += crArr[idx];
            count++;
          }
        }
        features.push(sumCb / count);
        features.push(sumCr / count);
      }
    }

    // 3. 16 dimensions: Horizontal band bilateral symmetry & vertical gradient ratios
    for (let band = 0; band < 8; band++) {
      const yStart = band * 8;
      let symDiff = 0;
      let horizEdge = 0;
      let n = 0;
      for (let y = yStart; y < yStart + 8; y++) {
        for (let x = 4; x < SIZE / 2; x++) {
          const leftVal = gray[y * SIZE + x];
          const rightVal = gray[y * SIZE + (SIZE - 1 - x)];
          symDiff += Math.abs(leftVal - rightVal);
          if (y > 0 && y < SIZE - 1) {
            horizEdge += Math.abs(gray[(y + 1) * SIZE + x] - gray[(y - 1) * SIZE + x]);
          }
          n++;
        }
      }
      features.push(1.0 - symDiff / Math.max(1, n));
      features.push(horizEdge / Math.max(1, n));
    }

    // 4. 16 dimensions: 4x4 Low-Frequency 2D-DCT Perceptual Hash coefficients
    for (let u = 0; u < 4; u++) {
      for (let v = 0; v < 4; v++) {
        let sum = 0;
        for (let y = 0; y < 16; y++) {
          const sy = y * 4;
          const cosY = Math.cos(((2 * y + 1) * v * Math.PI) / 32);
          for (let x = 0; x < 16; x++) {
            const sx2 = x * 4;
            const cosX = Math.cos(((2 * x + 1) * u * Math.PI) / 32);
            sum += gray[sy * SIZE + sx2] * cosX * cosY;
          }
        }
        features.push(sum / 256.0);
      }
    }

    return {
      embedding: l2Normalize(features),
      thumbDataUrl,
    };
  }

  // -------------------------------------------------------------------------
  // Multi-Tier Face Detection Engine (Native window.FaceDetector + Canvas Fallback)
  // -------------------------------------------------------------------------
  let nativeFaceDetectorState = "unknown"; // "supported" | "fallback"
  let cachedNativeDetector = null;

  async function detectFacesInImageElement(imgEl) {
    const natW = imgEl.naturalWidth || imgEl.width || 1;
    const natH = imgEl.naturalHeight || imgEl.height || 1;

    // Tier 1: Try Native OS / Browser Shape Detection API (window.FaceDetector)
    if ("FaceDetector" in window && nativeFaceDetectorState !== "fallback") {
      try {
        if (!cachedNativeDetector) {
          cachedNativeDetector = new window.FaceDetector({
            fastMode: false,
            maxDetectedFaces: 20,
          });
        }
        const rawFaces = await cachedNativeDetector.detect(imgEl);
        nativeFaceDetectorState = "supported";
        if (Array.isArray(rawFaces) && rawFaces.length > 0) {
          return rawFaces.map((f, idx) => {
            const bb = f.boundingBox;
            const box = {
              x: Math.max(0, bb.x),
              y: Math.max(0, bb.y),
              width: Math.min(natW - bb.x, bb.width),
              height: Math.min(natH - bb.y, bb.height),
            };
            const { embedding, thumbDataUrl } = extractFaceEmbeddingAndThumb(
              imgEl,
              box,
              f.landmarks || []
            );
            return {
              faceIndex: idx,
              engine: "native",
              box,
              normBox: {
                x: box.x / natW,
                y: box.y / natH,
                width: box.width / natW,
                height: box.height / natH,
              },
              landmarks: f.landmarks || [],
              embedding,
              thumbDataUrl,
            };
          });
        }
      } catch (err) {
        nativeFaceDetectorState = "fallback";
      }
    } else {
      nativeFaceDetectorState = "fallback";
    }

    // Tier 2: Pure-Canvas Integral-Image YCbCr Skin + Viola-Jones Haar Symmetry Detector
    const fallbackBoxes = runCanvasHaarSkinFaceDetector(imgEl);
    return fallbackBoxes.map((item, idx) => {
      const { embedding, thumbDataUrl } = extractFaceEmbeddingAndThumb(
        imgEl,
        item.box,
        item.landmarks
      );
      return {
        faceIndex: idx,
        engine: "canvas-haar",
        box: item.box,
        normBox: {
          x: item.box.x / natW,
          y: item.box.y / natH,
          width: item.box.width / natW,
          height: item.box.height / natH,
        },
        landmarks: item.landmarks,
        embedding,
        thumbDataUrl,
      };
    });
  }

  /**
   * Pure-Canvas Integral-Image Face Detector (Zero External Dependencies)
   * Combines YCbCr skin-likelihood integral map, Viola-Jones Haar eye-bridge
   * contrast cascade, and bilateral facial symmetry scoring with NMS.
   */
  function runCanvasHaarSkinFaceDetector(imgEl) {
    const natW = imgEl.naturalWidth || imgEl.width || 1;
    const natH = imgEl.naturalHeight || imgEl.height || 1;
    const MAX_DIM = 280;
    const scale = Math.min(1, MAX_DIM / Math.max(natW, natH));
    const w = Math.max(32, Math.round(natW * scale));
    const h = Math.max(32, Math.round(natH * scale));

    const canvas = document.createElement("canvas");
    canvas.width = w;
    canvas.height = h;
    const ctx = canvas.getContext("2d");
    ctx.drawImage(imgEl, 0, 0, w, h);
    const data = ctx.getImageData(0, 0, w, h).data;

    // Integral images (size (w+1) * (h+1)) for luminance, skin mask, and edge energy
    const stride = w + 1;
    const intLum = new Float64Array(stride * (h + 1));
    const intSkin = new Float64Array(stride * (h + 1));
    const intEdge = new Float64Array(stride * (h + 1));
    const lumGrid = new Float32Array(w * h);

    for (let y = 0; y < h; y++) {
      for (let x = 0; x < w; x++) {
        const i = (y * w + x) * 4;
        const r = data[i];
        const g = data[i + 1];
        const b = data[i + 2];
        const lum = (0.299 * r + 0.587 * g + 0.114 * b) / 255.0;
        lumGrid[y * w + x] = lum;

        const cb = 128 - 0.168736 * r - 0.331264 * g + 0.5 * b;
        const cr = 128 + 0.5 * r - 0.418688 * g - 0.081312 * b;
        const isSkin = cb >= 75 && cb <= 135 && cr >= 130 && cr <= 180 && lum > 0.15 && lum < 0.94 ? 1.0 : 0.0;

        let edge = 0;
        if (x > 0 && y > 0) {
          edge =
            Math.abs(lum - lumGrid[y * w + (x - 1)]) +
            Math.abs(lum - lumGrid[(y - 1) * w + x]);
        }

        const idx = (y + 1) * stride + (x + 1);
        intLum[idx] = lum + intLum[idx - 1] + intLum[idx - stride] - intLum[idx - stride - 1];
        intSkin[idx] = isSkin + intSkin[idx - 1] + intSkin[idx - stride] - intSkin[idx - stride - 1];
        intEdge[idx] = edge + intEdge[idx - 1] + intEdge[idx - stride] - intEdge[idx - stride - 1];
      }
    }

    function rectMean(intArr, x0, y0, rw, rh) {
      const x1 = Math.max(0, Math.min(w, Math.round(x0)));
      const y1 = Math.max(0, Math.min(h, Math.round(y0)));
      const x2 = Math.max(0, Math.min(w, Math.round(x0 + rw)));
      const y2 = Math.max(0, Math.min(h, Math.round(y0 + rh)));
      const area = Math.max(1, (x2 - x1) * (y2 - y1));
      const sum =
        intArr[y2 * stride + x2] -
        intArr[y1 * stride + x2] -
        intArr[y2 * stride + x1] +
        intArr[y1 * stride + x1];
      return sum / area;
    }

    const candidates = [];
    const minWin = Math.max(24, Math.round(Math.min(w, h) * 0.12));
    const maxWin = Math.round(Math.min(w, h) * 0.68);

    for (let winW = minWin; winW <= maxWin; winW = Math.round(winW * 1.24)) {
      const winH = Math.round(winW * 1.2);
      const step = Math.max(4, Math.round(winW * 0.14));

      for (let y = 2; y + winH < h - 2; y += step) {
        for (let x = 2; x + winW < w - 2; x += step) {
          const skinDensity = rectMean(intSkin, x + winW * 0.15, y + winH * 0.15, winW * 0.7, winH * 0.75);
          if (skinDensity < 0.35) continue;

          const edgeDensity = rectMean(intEdge, x + winW * 0.1, y + winH * 0.15, winW * 0.8, winH * 0.7);
          if (edgeDensity < 0.035) continue;

          // Viola-Jones Haar Feature 1: Eye region is darker than upper cheek / nose bridge
          const eyeBandLum = rectMean(intLum, x + winW * 0.15, y + winH * 0.24, winW * 0.7, winH * 0.18);
          const cheekBandLum = rectMean(intLum, x + winW * 0.15, y + winH * 0.44, winW * 0.7, winH * 0.22);
          const foreheadLum = rectMean(intLum, x + winW * 0.2, y + winH * 0.08, winW * 0.6, winH * 0.14);

          const haarEyeCheek = cheekBandLum - eyeBandLum;
          const haarForeheadEye = foreheadLum - eyeBandLum;
          if (haarEyeCheek < -0.02 && haarForeheadEye < -0.02) continue;

          // Viola-Jones Haar Feature 2: Nasal bridge brighter than left/right eye sockets
          const leftEyeLum = rectMean(intLum, x + winW * 0.16, y + winH * 0.25, winW * 0.24, winH * 0.16);
          const bridgeLum = rectMean(intLum, x + winW * 0.4, y + winH * 0.25, winW * 0.2, winH * 0.16);
          const rightEyeLum = rectMean(intLum, x + winW * 0.6, y + winH * 0.25, winW * 0.24, winH * 0.16);
          const bridgeContrast = bridgeLum - 0.5 * (leftEyeLum + rightEyeLum);
          const eyeBalance = 1.0 - Math.abs(leftEyeLum - rightEyeLum);

          // Upper-centre compositional bias (faces are more commonly in the upper 75% of portrait/scene frames)
          const posBias = y < h * 0.7 ? 0.12 : 0.0;

          const score =
            skinDensity * 0.42 +
            Math.min(0.25, edgeDensity * 2.2) +
            Math.max(0, haarEyeCheek) * 0.8 +
            Math.max(0, bridgeContrast) * 0.6 +
            eyeBalance * 0.2 +
            posBias;

          if (score > 0.58) {
            candidates.push({
              x: x / scale,
              y: y / scale,
              width: winW / scale,
              height: winH / scale,
              score,
            });
          }
        }
      }
    }

    // Non-Maximum Suppression (NMS)
    candidates.sort((a, b) => b.score - a.score);
    const picked = [];
    for (const cand of candidates) {
      let overlaps = false;
      for (const p of picked) {
        const ix1 = Math.max(cand.x, p.x);
        const iy1 = Math.max(cand.y, p.y);
        const ix2 = Math.min(cand.x + cand.width, p.x + p.width);
        const iy2 = Math.min(cand.y + cand.height, p.y + p.height);
        const inter = Math.max(0, ix2 - ix1) * Math.max(0, iy2 - iy1);
        const union = cand.width * cand.height + p.width * p.height - inter;
        if (inter / Math.max(1, union) > 0.32) {
          overlaps = true;
          break;
        }
      }
      if (!overlaps) {
        picked.push(cand);
        if (picked.length >= 6) break;
      }
    }

    return picked.map((b) => ({
      box: { x: b.x, y: b.y, width: b.width, height: b.height },
      landmarks: [
        {
          type: "eye",
          locations: [{ x: b.x + b.width * 0.32, y: b.y + b.height * 0.34 }],
        },
        {
          type: "eye",
          locations: [{ x: b.x + b.width * 0.68, y: b.y + b.height * 0.34 }],
        },
        {
          type: "nose",
          locations: [{ x: b.x + b.width * 0.5, y: b.y + b.height * 0.55 }],
        },
        {
          type: "mouth",
          locations: [{ x: b.x + b.width * 0.5, y: b.y + b.height * 0.76 }],
        },
      ],
    }));
  }

  function loadImageAsync(src) {
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.crossOrigin = "anonymous";
      img.onload = () => resolve(img);
      img.onerror = (e) => reject(e);
      img.src = src;
    });
  }

  /**
   * Matches a face embedding against stored Performer face prototypes.
   */
  function matchPerformerPrototypes(embedding, performersList) {
    const db = loadPerformerPrototypes();
    const matches = [];
    for (const pid of Object.keys(db)) {
      const entry = db[pid];
      if (!entry.centroid || !entry.centroid.length) continue;
      const sim = cosineSimilarity(embedding, entry.centroid);
      const perfObj = (performersList || []).find((p) => String(p.id) === String(pid));
      matches.push({
        performerId: String(pid),
        name: perfObj ? perfObj.name : entry.name,
        similarity: sim,
      });
    }
    matches.sort((a, b) => b.similarity - a.similarity);
    return matches.slice(0, 3);
  }

  /**
   * Groups an array of face records into similarity clusters using
   * Leader-Follower Centroid Clustering.
   */
  function clusterDetectedFaces(allFaces, similarityThreshold, performersList) {
    const clusters = [];
    for (const face of allFaces) {
      if (!face.embedding || !face.embedding.length) continue;

      let bestCluster = null;
      let bestSim = -1;

      for (const cluster of clusters) {
        const sim = cosineSimilarity(face.embedding, cluster.centroid);
        if (sim >= similarityThreshold && sim > bestSim) {
          bestSim = sim;
          bestCluster = cluster;
        }
      }

      if (bestCluster) {
        bestCluster.faces.push({ ...face, similarityToCluster: bestSim });
        bestCluster.centroid = computeCentroid(bestCluster.faces.map((f) => f.embedding));
      } else {
        clusters.push({
          id: `cluster-${clusters.length + 1}`,
          centroid: face.embedding.slice(),
          faces: [{ ...face, similarityToCluster: 1.0 }],
        });
      }
    }

    // Sort largest clusters first and attach best performer suggestion
    clusters.sort((a, b) => b.faces.length - a.faces.length);
    clusters.forEach((c, i) => {
      c.id = `cluster-${i + 1}`;
      c.suggestions = matchPerformerPrototypes(c.centroid, performersList);
    });
    return clusters;
  }

  // -------------------------------------------------------------------------
  // Main Initialization with PluginApi
  // -------------------------------------------------------------------------
  function initPluginWithApi(PluginApi) {
    injectStyles();

    const React = PluginApi.React;
    const { useState, useEffect, useRef, useMemo, useCallback } = React;
    const GQL = PluginApi.GQL;
    const { Button, Badge, Form, ProgressBar, Spinner } = PluginApi.libraries.Bootstrap;
    const { NavLink, Link, useHistory } = PluginApi.libraries.ReactRouterDOM;
    const Mousetrap = PluginApi.libraries.Mousetrap;
    const FA = PluginApi.libraries.FontAwesomeSolid || {};

    // Helper to safely render icons
    function renderIcon(faIcon, fallbackText) {
      const IconComp = PluginApi.components && PluginApi.components.Icon;
      if (IconComp && faIcon) {
        return React.createElement(IconComp, { icon: faIcon });
      }
      return React.createElement("span", null, fallbackText || "•");
    }

    // GraphQL Mutation Helper using Stash's Apollo Client
    async function assignPerformerToImagesGQL(imageIds, performerId, performerName, faceEmbeddings) {
      const client = PluginApi.utils.StashService.getClient();
      const ids = Array.from(new Set(imageIds.map(String)));
      if (!ids.length || !performerId) return;

      await client.mutate({
        mutation: GQL.BulkImageUpdateDocument,
        variables: {
          input: {
            ids,
            performer_ids: {
              mode: GQL.BulkUpdateIdMode ? GQL.BulkUpdateIdMode.Add : "ADD",
              ids: [String(performerId)],
            },
          },
        },
        refetchQueries: ["FindImages", "FindImage"],
      });

      // Save prototype embeddings so future scans auto-suggest this performer
      if (Array.isArray(faceEmbeddings)) {
        faceEmbeddings.forEach((emb) => {
          savePerformerPrototype(performerId, performerName, emb);
        });
      }
    }

    async function createNewPerformerGQL(name) {
      const client = PluginApi.utils.StashService.getClient();
      const res = await client.mutate({
        mutation: GQL.PerformerCreateDocument,
        variables: {
          input: { name: name.trim() },
        },
        refetchQueries: ["FindPerformers"],
      });
      return res && res.data && res.data.performerCreate;
    }

    async function updateSingleImageRatingOrOrganized(imageId, patchInput) {
      const client = PluginApi.utils.StashService.getClient();
      await client.mutate({
        mutation: GQL.ImageUpdateDocument,
        variables: {
          input: {
            id: String(imageId),
            ...patchInput,
          },
        },
      });
    }

    // =======================================================================
    // Component 1: Interactive Pan, Zoom & Face Tagging Image Viewer Modal
    // =======================================================================
    const ZoomViewerModal = function ({
      images,
      initialIndex = 0,
      initialFocusBox = null,
      performers = [],
      onClose,
      onImageUpdated,
    }) {
      const [currentIndex, setCurrentIndex] = useState(initialIndex);
      const currentImage = images[currentIndex] || images[0];

      const viewportRef = useRef(null);
      const imgRef = useRef(null);

      const [imgNatural, setImgNatural] = useState({ width: 1, height: 1, loaded: false });
      const [transform, setTransform] = useState({ scale: 1, x: 0, y: 0 });
      const [isPanning, setIsPanning] = useState(false);
      const panStartRef = useRef({ x: 0, y: 0, origX: 0, origY: 0 });

      const [showBoxes, setShowBoxes] = useState(true);
      const [showFilmstrip, setShowFilmstrip] = useState(true);
      const [faces, setFaces] = useState([]);
      const [detecting, setDetecting] = useState(false);
      const [activeFaceIdx, setActiveFaceIdx] = useState(null);

      // Manual Face Box Drawing state
      const [drawMode, setDrawMode] = useState(false);
      const [dragBox, setDragBox] = useState(null);

      // Popover state for performer assignment
      const [searchPerf, setSearchPerf] = useState("");
      const [newPerfName, setNewPerfName] = useState("");
      const [savingAssign, setSavingAssign] = useState(false);
      const [statusBanner, setStatusBanner] = useState("");

      const fitImageToViewport = useCallback(
        (natW, natH, focusNormBox) => {
          if (!viewportRef.current) return;
          const rect = viewportRef.current.getBoundingClientRect();
          const vw = Math.max(200, rect.width);
          const vh = Math.max(200, rect.height);

          if (focusNormBox) {
            const boxW = focusNormBox.width * natW;
            const boxH = focusNormBox.height * natH;
            const targetScale = Math.min(4.5, Math.max(0.5, Math.min(vw / (boxW * 2.8), vh / (boxH * 2.8))));
            const centerX = (focusNormBox.x + focusNormBox.width / 2) * natW;
            const centerY = (focusNormBox.y + focusNormBox.height / 2) * natH;
            setTransform({
              scale: targetScale,
              x: vw / 2 - centerX * targetScale,
              y: vh / 2 - centerY * targetScale,
            });
            return;
          }

          const scaleX = (vw * 0.94) / Math.max(1, natW);
          const scaleY = (vh * 0.94) / Math.max(1, natH);
          const fitScale = Math.min(scaleX, scaleY, 1.5);
          setTransform({
            scale: fitScale,
            x: (vw - natW * fitScale) / 2,
            y: (vh - natH * fitScale) / 2,
          });
        },
        []
      );

      const toggleOneToOneZoom = useCallback(() => {
        if (!viewportRef.current || !imgNatural.loaded) return;
        const rect = viewportRef.current.getBoundingClientRect();
        if (Math.abs(transform.scale - 1.0) < 0.08) {
          fitImageToViewport(imgNatural.width, imgNatural.height, null);
        } else {
          setTransform({
            scale: 1.0,
            x: (rect.width - imgNatural.width) / 2,
            y: (rect.height - imgNatural.height) / 2,
          });
        }
      }, [transform.scale, imgNatural, fitImageToViewport]);

      // Run face detection when image loads
      const handleImageLoad = async (e) => {
        const el = e.currentTarget;
        const natW = el.naturalWidth || 1;
        const natH = el.naturalHeight || 1;
        setImgNatural({ width: natW, height: natH, loaded: true });
        fitImageToViewport(
          natW,
          natH,
          currentIndex === initialIndex ? initialFocusBox : null
        );

        const cache = loadFaceCache();
        const imgId = String(currentImage.id);
        if (cache[imgId] && Array.isArray(cache[imgId].faces)) {
          const hydrated = cache[imgId].faces.map((f, idx) => ({
            ...f,
            faceIndex: idx,
            box: {
              x: f.normBox.x * natW,
              y: f.normBox.y * natH,
              width: f.normBox.width * natW,
              height: f.normBox.height * natH,
            },
          }));
          setFaces(hydrated);
          return;
        }

        setDetecting(true);
        try {
          const detected = await detectFacesInImageElement(el);
          setFaces(detected);
          cache[imgId] = {
            updatedAt: Date.now(),
            faces: detected.map((d) => ({
              normBox: d.normBox,
              landmarks: d.landmarks,
              embedding: d.embedding,
              thumbDataUrl: d.thumbDataUrl,
              assignedPerformerId: d.assignedPerformerId || null,
              assignedPerformerName: d.assignedPerformerName || null,
            })),
          };
          saveFaceCache(cache);
        } catch (err) {
          console.warn("Face detection error on image:", err);
        } finally {
          setDetecting(false);
        }
      };

      // Reset state when switching images
      useEffect(() => {
        setActiveFaceIdx(null);
        setImgNatural({ width: 1, height: 1, loaded: false });
        setFaces([]);
      }, [currentIndex]);

      // Keyboard navigation & shortcuts inside ZoomViewerModal
      useEffect(() => {
        const handleKeyDown = (e) => {
          if (e.target && (e.target.tagName === "INPUT" || e.target.tagName === "SELECT")) {
            return;
          }
          if (e.key === "Escape") {
            e.preventDefault();
            onClose();
          } else if (e.key === "ArrowRight" || e.key === "d" || e.key === "D") {
            e.preventDefault();
            setCurrentIndex((i) => (i + 1) % images.length);
          } else if (e.key === "ArrowLeft" || e.key === "a" || e.key === "A") {
            e.preventDefault();
            setCurrentIndex((i) => (i - 1 + images.length) % images.length);
          } else if (e.key === "z" || e.key === "Z") {
            e.preventDefault();
            toggleOneToOneZoom();
          } else if (e.key === "r" || e.key === "R") {
            e.preventDefault();
            fitImageToViewport(imgNatural.width, imgNatural.height, null);
          } else if (e.key === "b" || e.key === "B") {
            e.preventDefault();
            setShowBoxes((v) => !v);
          } else if (e.key === "m" || e.key === "M") {
            e.preventDefault();
            setDrawMode((v) => !v);
          } else if (["1", "2", "3", "4", "5"].includes(e.key)) {
            const stars = parseInt(e.key, 10);
            updateSingleImageRatingOrOrganized(currentImage.id, { rating100: stars * 20 });
            setStatusBanner(`Rated ${stars} ★`);
            setTimeout(() => setStatusBanner(""), 1800);
          } else if (e.key === "o" || e.key === "O") {
            if (PluginApi.utils.StashService.mutateImageIncrementO) {
              PluginApi.utils.StashService.mutateImageIncrementO(currentImage.id);
              setStatusBanner("Incremented O-Counter");
              setTimeout(() => setStatusBanner(""), 1800);
            }
          }
        };
        window.addEventListener("keydown", handleKeyDown);
        return () => window.removeEventListener("keydown", handleKeyDown);
      }, [images.length, onClose, toggleOneToOneZoom, fitImageToViewport, imgNatural, currentImage]);

      // Cursor-anchored smooth mouse wheel zoom
      const handleWheel = (e) => {
        e.preventDefault();
        if (!viewportRef.current) return;
        const rect = viewportRef.current.getBoundingClientRect();
        const mouseX = e.clientX - rect.left;
        const mouseY = e.clientY - rect.top;

        const factor = e.deltaY < 0 ? 1.16 : 0.86;
        const nextScale = Math.min(Math.max(transform.scale * factor, 0.1), 12.0);
        const ratio = 1 - nextScale / transform.scale;

        setTransform({
          scale: nextScale,
          x: transform.x + (mouseX - transform.x) * ratio,
          y: transform.y + (mouseY - transform.y) * ratio,
        });
      };

      // Pointer Down for Pan or Manual Face Box Drawing
      const handleMouseDown = (e) => {
        if (e.button !== 0) return;
        if (e.target.closest(".sft-face-box") || e.target.closest(".sft-popover-card")) {
          return;
        }
        if (!viewportRef.current) return;
        const rect = viewportRef.current.getBoundingClientRect();

        if (drawMode || e.shiftKey) {
          const imgX = (e.clientX - rect.left - transform.x) / transform.scale;
          const imgY = (e.clientY - rect.top - transform.y) / transform.scale;
          setDragBox({ startX: imgX, startY: imgY, x: imgX, y: imgY, width: 0, height: 0 });
          return;
        }

        setActiveFaceIdx(null);
        setIsPanning(true);
        panStartRef.current = {
          x: e.clientX,
          y: e.clientY,
          origX: transform.x,
          origY: transform.y,
        };
      };

      const handleMouseMove = (e) => {
        if (dragBox && viewportRef.current) {
          const rect = viewportRef.current.getBoundingClientRect();
          const curX = (e.clientX - rect.left - transform.x) / transform.scale;
          const curY = (e.clientY - rect.top - transform.y) / transform.scale;
          setDragBox({
            ...dragBox,
            x: Math.min(dragBox.startX, curX),
            y: Math.min(dragBox.startY, curY),
            width: Math.abs(curX - dragBox.startX),
            height: Math.abs(curY - dragBox.startY),
          });
          return;
        }

        if (!isPanning) return;
        const dx = e.clientX - panStartRef.current.x;
        const dy = e.clientY - panStartRef.current.y;
        setTransform((prev) => ({
          ...prev,
          x: panStartRef.current.origX + dx,
          y: panStartRef.current.origY + dy,
        }));
      };

      const handleMouseUp = () => {
        if (dragBox) {
          if (dragBox.width > 16 && dragBox.height > 16 && imgRef.current) {
            const box = {
              x: Math.max(0, dragBox.x),
              y: Math.max(0, dragBox.y),
              width: Math.min(imgNatural.width - dragBox.x, dragBox.width),
              height: Math.min(imgNatural.height - dragBox.y, dragBox.height),
            };
            const { embedding, thumbDataUrl } = extractFaceEmbeddingAndThumb(imgRef.current, box, []);
            const newFace = {
              faceIndex: faces.length,
              engine: "manual",
              box,
              normBox: {
                x: box.x / imgNatural.width,
                y: box.y / imgNatural.height,
                width: box.width / imgNatural.width,
                height: box.height / imgNatural.height,
              },
              landmarks: [],
              embedding,
              thumbDataUrl,
            };
            const updatedFaces = [...faces, newFace];
            setFaces(updatedFaces);
            setActiveFaceIdx(updatedFaces.length - 1);

            const cache = loadFaceCache();
            cache[String(currentImage.id)] = {
              updatedAt: Date.now(),
              faces: updatedFaces.map((d) => ({
                normBox: d.normBox,
                landmarks: d.landmarks,
                embedding: d.embedding,
                thumbDataUrl: d.thumbDataUrl,
                assignedPerformerId: d.assignedPerformerId || null,
                assignedPerformerName: d.assignedPerformerName || null,
              })),
            };
            saveFaceCache(cache);
          }
          setDragBox(null);
          setDrawMode(false);
          return;
        }
        setIsPanning(false);
      };

      // Assign Performer to a specific detected face
      const handleAssignPerformerToFace = async (faceIdx, performer) => {
        if (!performer || !currentImage) return;
        setSavingAssign(true);
        try {
          const targetFace = faces[faceIdx];
          await assignPerformerToImagesGQL(
            [currentImage.id],
            performer.id,
            performer.name,
            targetFace && targetFace.embedding ? [targetFace.embedding] : []
          );

          const updated = faces.map((f, i) =>
            i === faceIdx
              ? {
                  ...f,
                  assignedPerformerId: String(performer.id),
                  assignedPerformerName: performer.name,
                }
              : f
          );
          setFaces(updated);

          const cache = loadFaceCache();
          if (cache[String(currentImage.id)]) {
            cache[String(currentImage.id)].faces = updated.map((d) => ({
              normBox: d.normBox,
              landmarks: d.landmarks,
              embedding: d.embedding,
              thumbDataUrl: d.thumbDataUrl,
              assignedPerformerId: d.assignedPerformerId || null,
              assignedPerformerName: d.assignedPerformerName || null,
            }));
            saveFaceCache(cache);
          }

          setStatusBanner(`Assigned ${performer.name}`);
          setTimeout(() => setStatusBanner(""), 2200);
          setActiveFaceIdx(null);
          if (onImageUpdated) onImageUpdated();
        } catch (err) {
          console.error("Failed to assign performer:", err);
        } finally {
          setSavingAssign(false);
        }
      };

      const handleCreateAndAssign = async (faceIdx) => {
        if (!newPerfName.trim()) return;
        setSavingAssign(true);
        try {
          const created = await createNewPerformerGQL(newPerfName.trim());
          if (created) {
            setNewPerfName("");
            await handleAssignPerformerToFace(faceIdx, created);
          }
        } catch (err) {
          console.error("Error creating performer:", err);
          setSavingAssign(false);
        }
      };

      if (!currentImage) return null;
      const imgSrc =
        (currentImage.paths && (currentImage.paths.image || currentImage.paths.preview || currentImage.paths.thumbnail)) ||
        "";
      const imgTitle =
        currentImage.title ||
        (currentImage.visual_files &&
          currentImage.visual_files[0] &&
          currentImage.visual_files[0].path) ||
        `Image #${currentImage.id}`;

      const filteredPerformers = performers
        .filter((p) =>
          !searchPerf.trim()
            ? true
            : (p.name || "").toLowerCase().includes(searchPerf.trim().toLowerCase())
        )
        .slice(0, 12);

      // Minimap viewport calculation
      let minimapStyle = { left: "0%", top: "0%", width: "100%", height: "100%" };
      if (viewportRef.current && imgNatural.loaded) {
        const rect = viewportRef.current.getBoundingClientRect();
        const totalW = imgNatural.width * transform.scale;
        const totalH = imgNatural.height * transform.scale;
        const leftPct = Math.max(0, Math.min(100, (-transform.x / Math.max(1, totalW)) * 100));
        const topPct = Math.max(0, Math.min(100, (-transform.y / Math.max(1, totalH)) * 100));
        const wPct = Math.min(100 - leftPct, (rect.width / Math.max(1, totalW)) * 100);
        const hPct = Math.min(100 - topPct, (rect.height / Math.max(1, totalH)) * 100);
        minimapStyle = {
          left: `${leftPct}%`,
          top: `${topPct}%`,
          width: `${Math.max(6, wPct)}%`,
          height: `${Math.max(6, hPct)}%`,
        };
      }

      return React.createElement(
        "div",
        { className: "sft-modal-backdrop" },
        // Top Toolbar
        React.createElement(
          "div",
          { className: "sft-modal-topbar" },
          React.createElement(
            "div",
            { className: "d-flex align-items-center", style: { gap: "0.75rem" } },
            React.createElement(
              "span",
              { className: "sft-modal-title", title: imgTitle },
              `${currentIndex + 1} / ${images.length} — ${imgTitle}`
            ),
            React.createElement(
              "span",
              {
                className: `sft-badge-engine ${
                  nativeFaceDetectorState === "supported" ? "" : "fallback"
                }`,
              },
              nativeFaceDetectorState === "supported"
                ? "Native OS FaceDetector"
                : "Canvas Haar+Symmetry Engine"
            ),
            detecting &&
              React.createElement(
                "span",
                { className: "text-info small" },
                "Scanning faces…"
              ),
            statusBanner &&
              React.createElement(
                Badge,
                { variant: "success", className: "px-2 py-1" },
                statusBanner
              )
          ),
          React.createElement(
            "div",
            { className: "sft-modal-toolbar" },
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "outline-light",
                onClick: () => setCurrentIndex((i) => (i - 1 + images.length) % images.length),
                title: "Previous Image (Left Arrow / A)",
              },
              renderIcon(FA.faArrowLeft, "◀"),
              " Prev"
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "outline-light",
                onClick: () => setCurrentIndex((i) => (i + 1) % images.length),
                title: "Next Image (Right Arrow / D)",
              },
              "Next ",
              renderIcon(FA.faArrowRight, "▶")
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "secondary",
                onClick: () =>
                  setTransform((t) => ({ ...t, scale: Math.max(0.1, t.scale * 0.8) })),
                title: "Zoom Out",
              },
              renderIcon(FA.faSearchMinus, "−")
            ),
            React.createElement(
              "span",
              {
                style: {
                  minWidth: "54px",
                  textAlign: "center",
                  fontSize: "0.82rem",
                  fontWeight: 600,
                },
              },
              `${Math.round(transform.scale * 100)}%`
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "secondary",
                onClick: () =>
                  setTransform((t) => ({ ...t, scale: Math.min(12.0, t.scale * 1.25) })),
                title: "Zoom In",
              },
              renderIcon(FA.faSearchPlus, "+")
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "secondary",
                onClick: toggleOneToOneZoom,
                title: "Toggle 1:1 100% Pixel View (Z)",
              },
              "1:1"
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "secondary",
                onClick: () => fitImageToViewport(imgNatural.width, imgNatural.height, null),
                title: "Fit to Screen (R)",
              },
              "Fit"
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: showBoxes ? "info" : "outline-secondary",
                onClick: () => setShowBoxes((v) => !v),
                title: "Toggle Face Bounding Boxes (B)",
              },
              `Faces (${faces.length})`
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: drawMode ? "warning" : "outline-info",
                onClick: () => setDrawMode((v) => !v),
                title: "Draw Manual Face Box (M or Shift+Drag)",
              },
              drawMode ? "Drag on Face…" : "+ Draw Face"
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: showFilmstrip ? "secondary" : "outline-secondary",
                onClick: () => setShowFilmstrip((v) => !v),
                title: "Toggle Bottom Filmstrip",
              },
              "Filmstrip"
            ),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "danger",
                onClick: onClose,
                title: "Close Viewer (Esc)",
              },
              renderIcon(FA.faTimes, "✕"),
              " Close"
            )
          )
        ),

        // Main Interactive Pan & Zoom Viewport
        React.createElement(
          "div",
          {
            ref: viewportRef,
            className: `sft-viewport ${isPanning ? "panning" : ""} ${
              drawMode ? "drawing" : ""
            }`,
            onWheel: handleWheel,
            onMouseDown: handleMouseDown,
            onMouseMove: handleMouseMove,
            onMouseUp: handleMouseUp,
            onDoubleClick: toggleOneToOneZoom,
          },
          React.createElement(
            "div",
            {
              className: "sft-transform-layer",
              style: {
                transform: `translate3d(${transform.x}px, ${transform.y}px, 0) scale(${transform.scale})`,
              },
            },
            React.createElement("img", {
              ref: imgRef,
              src: imgSrc,
              alt: imgTitle,
              className: "sft-main-img",
              crossOrigin: "anonymous",
              onLoad: handleImageLoad,
            }),

            // Active Drag Box Preview when drawing a manual face box
            dragBox &&
              React.createElement("div", {
                style: {
                  position: "absolute",
                  left: `${dragBox.x}px`,
                  top: `${dragBox.y}px`,
                  width: `${dragBox.width}px`,
                  height: `${dragBox.height}px`,
                  border: "2px dashed #fcc419",
                  background: "rgba(252, 196, 25, 0.18)",
                  pointerEvents: "none",
                },
              }),

            // Face Bounding Boxes Overlay
            showBoxes &&
              imgNatural.loaded &&
              faces.map((face, idx) => {
                const suggestions = matchPerformerPrototypes(face.embedding, performers);
                const topMatch = suggestions[0];
                const existingPerfNames = (currentImage.performers || []).map((p) => p.name);
                const labelText =
                  face.assignedPerformerName ||
                  (topMatch && topMatch.similarity >= 0.72
                    ? `${topMatch.name} (${Math.round(topMatch.similarity * 100)}%)`
                    : existingPerfNames[idx] || `Face #${idx + 1} (Click to tag)`);

                const isAssigned = Boolean(
                  face.assignedPerformerId || existingPerfNames[idx]
                );

                return React.createElement(
                  "div",
                  {
                    key: `face-${idx}`,
                    className: `sft-face-box ${isAssigned ? "assigned" : ""} ${
                      activeFaceIdx === idx ? "active" : ""
                    }`,
                    style: {
                      left: `${face.box.x}px`,
                      top: `${face.box.y}px`,
                      width: `${face.box.width}px`,
                      height: `${face.box.height}px`,
                      borderWidth: `${Math.max(1.5, 2.2 / transform.scale)}px`,
                    },
                    onClick: (e) => {
                      e.stopPropagation();
                      setActiveFaceIdx(activeFaceIdx === idx ? null : idx);
                    },
                  },
                  React.createElement(
                    "div",
                    {
                      className: "sft-face-box-label",
                      style: {
                        transform: `scale(${Math.min(1.4, Math.max(0.45, 1 / transform.scale))})`,
                        transformOrigin: "left bottom",
                      },
                    },
                    labelText
                  ),

                  // Landmark dots
                  Array.isArray(face.landmarks) &&
                    face.landmarks.map((lm, lmi) =>
                      (lm.locations || []).map((loc, loci) =>
                        React.createElement("div", {
                          key: `lm-${lmi}-${loci}`,
                          className: "sft-landmark-dot",
                          style: {
                            left: `${loc.x - face.box.x}px`,
                            top: `${loc.y - face.box.y}px`,
                          },
                        })
                      )
                    ),

                  // Click-to-Assign Performer Popover
                  activeFaceIdx === idx &&
                    React.createElement(
                      "div",
                      {
                        className: "sft-popover-card",
                        style: {
                          transform: `scale(${1 / transform.scale})`,
                          transformOrigin: "top left",
                        },
                        onClick: (e) => e.stopPropagation(),
                      },
                      React.createElement(
                        "div",
                        { className: "d-flex justify-content-between align-items-center mb-2" },
                        React.createElement(
                          "strong",
                          { style: { fontSize: "0.85rem" } },
                          `Tag Face #${idx + 1}`
                        ),
                        React.createElement(
                          "button",
                          {
                            className: "btn btn-sm btn-link text-light p-0",
                            onClick: () => setActiveFaceIdx(null),
                          },
                          "✕"
                        )
                      ),
                      suggestions.length > 0 &&
                        React.createElement(
                          "div",
                          { className: "mb-2" },
                          React.createElement(
                            "div",
                            { className: "text-muted", style: { fontSize: "0.75rem" } },
                            "Suggested Matches:"
                          ),
                          suggestions.map((sug) =>
                            React.createElement(
                              Button,
                              {
                                key: sug.performerId,
                                size: "sm",
                                variant: "outline-info",
                                className: "mr-1 mt-1",
                                disabled: savingAssign,
                                onClick: () =>
                                  handleAssignPerformerToFace(idx, {
                                    id: sug.performerId,
                                    name: sug.name,
                                  }),
                              },
                              `${sug.name} (${Math.round(sug.similarity * 100)}%)`
                            )
                          )
                        ),
                      React.createElement("input", {
                        type: "text",
                        className: "sft-input w-100 mb-2",
                        placeholder: "Search existing performers…",
                        value: searchPerf,
                        onChange: (e) => setSearchPerf(e.target.value),
                      }),
                      React.createElement(
                        "div",
                        {
                          style: {
                            maxHeight: "130px",
                            overflowY: "auto",
                            display: "flex",
                            flexDirection: "column",
                            gap: "4px",
                          },
                        },
                        filteredPerformers.map((perf) =>
                          React.createElement(
                            "button",
                            {
                              key: perf.id,
                              type: "button",
                              className:
                                "btn btn-sm btn-dark text-left d-flex justify-content-between align-items-center",
                              disabled: savingAssign,
                              onClick: () => handleAssignPerformerToFace(idx, perf),
                            },
                            React.createElement("span", null, perf.name),
                            React.createElement(
                              "span",
                              { className: "text-info small" },
                              "Assign"
                            )
                          )
                        )
                      ),
                      React.createElement(
                        "div",
                        { className: "d-flex mt-2", style: { gap: "4px" } },
                        React.createElement("input", {
                          type: "text",
                          className: "sft-input flex-grow-1",
                          placeholder: "Or new performer name…",
                          value: newPerfName,
                          onChange: (e) => setNewPerfName(e.target.value),
                        }),
                        React.createElement(
                          Button,
                          {
                            size: "sm",
                            variant: "success",
                            disabled: savingAssign || !newPerfName.trim(),
                            onClick: () => handleCreateAndAssign(idx),
                          },
                          "+ Add"
                        )
                      )
                    )
                );
              })
          ),

          // Interactive Bottom-Right Minimap
          imgNatural.loaded &&
            React.createElement(
              "div",
              { className: "sft-minimap", title: "Viewport Minimap" },
              React.createElement("img", { src: imgSrc, alt: "minimap" }),
              React.createElement("div", {
                className: "sft-minimap-rect",
                style: minimapStyle,
              })
            )
        ),

        // Bottom Thumbnail Filmstrip for Rapid Image Scrubbing
        showFilmstrip &&
          React.createElement(
            "div",
            { className: "sft-filmstrip" },
            images.map((imgItem, idx) => {
              const thumbUrl =
                (imgItem.paths && (imgItem.paths.thumbnail || imgItem.paths.preview || imgItem.paths.image)) ||
                "";
              return React.createElement(
                "div",
                {
                  key: imgItem.id,
                  className: `sft-filmstrip-item ${idx === currentIndex ? "active" : ""}`,
                  onClick: () => setCurrentIndex(idx),
                  title: imgItem.title || `Image #${imgItem.id}`,
                },
                React.createElement("img", {
                  src: thumbUrl,
                  alt: `thumb-${imgItem.id}`,
                  loading: "lazy",
                })
              );
            })
          )
      );
    };

    // =======================================================================
    // Component 2: Full Workspace Route (/plugins/face-tagger)
    // =======================================================================
    const FaceTaggerWorkspace = function () {
      const [activeTab, setActiveTab] = useState("clusters"); // "clusters" | "sorter"
      const [perPage, setPerPage] = useState(40);
      const [page, setPage] = useState(1);
      const [searchQuery, setSearchQuery] = useState("");
      const [filterMode, setFilterMode] = useState("all"); // "all" | "untagged" | "unorganized"
      const [sortMode, setSortMode] = useState("unassigned_first");
      const [similarityThreshold, setSimilarityThreshold] = useState(0.72);

      // Scanning & Clustering state
      const [scanning, setScanning] = useState(false);
      const [scanProgress, setScanProgress] = useState({ done: 0, total: 0 });
      const [faceMap, setFaceMap] = useState(() => loadFaceCache());
      const [excludedFaceKeys, setExcludedFaceKeys] = useState({});
      const [selectedClusterPerformer, setSelectedClusterPerformer] = useState({});
      const [newClusterPerfName, setNewClusterPerfName] = useState({});
      const [labeledClusters, setLabeledClusters] = useState({});
      const [toastMsg, setToastMsg] = useState("");

      // Multi-select state in Image Sorter tab
      const [selectedImageIds, setSelectedImageIds] = useState({});
      const [batchPerformerId, setBatchPerformerId] = useState("");

      // Modal Viewer state
      const [viewerState, setViewerState] = useState({
        open: false,
        index: 0,
        focusBox: null,
      });

      // Query Images from Stash GraphQL
      const {
        data: imagesData,
        loading: imagesLoading,
        refetch: refetchImages,
      } = GQL.useFindImagesQuery({
        variables: {
          filter: {
            q: searchQuery,
            page: page,
            per_page: perPage,
            sort: "created_at",
            direction: GQL.SortDirectionEnum ? GQL.SortDirectionEnum.Desc : "DESC",
          },
        },
        fetchPolicy: "cache-and-network",
      });

      // Query Performers from Stash GraphQL
      const { data: performersData, refetch: refetchPerformers } =
        GQL.useFindPerformersQuery({
          variables: {
            filter: {
              per_page: -1,
              sort: "name",
              direction: GQL.SortDirectionEnum ? GQL.SortDirectionEnum.Asc : "ASC",
            },
          },
        });

      const performers = useMemo(
        () =>
          (performersData &&
            performersData.findPerformers &&
            performersData.findPerformers.performers) ||
          [],
        [performersData]
      );

      const rawImages = useMemo(() => {
        const list =
          (imagesData && imagesData.findImages && imagesData.findImages.images) || [];
        if (filterMode === "untagged") {
          return list.filter((img) => !img.performers || img.performers.length === 0);
        }
        if (filterMode === "unorganized") {
          return list.filter((img) => !img.organized);
        }
        return list;
      }, [imagesData, filterMode]);

      // Batch scan all currently loaded images for faces
      const runBatchFaceScan = async (forceRescan = false) => {
        if (!rawImages.length || scanning) return;
        setScanning(true);
        const updatedCache = { ...loadFaceCache() };
        setScanProgress({ done: 0, total: rawImages.length });

        for (let i = 0; i < rawImages.length; i++) {
          const imgObj = rawImages[i];
          const idStr = String(imgObj.id);
          if (!forceRescan && updatedCache[idStr] && Array.isArray(updatedCache[idStr].faces)) {
            setScanProgress({ done: i + 1, total: rawImages.length });
            continue;
          }

          const src =
            (imgObj.paths && (imgObj.paths.image || imgObj.paths.preview || imgObj.paths.thumbnail)) ||
            "";
          if (src) {
            try {
              const htmlImg = await loadImageAsync(src);
              const detected = await detectFacesInImageElement(htmlImg);
              updatedCache[idStr] = {
                updatedAt: Date.now(),
                faces: detected.map((d) => ({
                  normBox: d.normBox,
                  landmarks: d.landmarks,
                  embedding: d.embedding,
                  thumbDataUrl: d.thumbDataUrl,
                  assignedPerformerId: null,
                  assignedPerformerName: null,
                })),
              };
            } catch (err) {
              console.warn(`Failed scanning image ${idStr}:`, err);
            }
          }
          setScanProgress({ done: i + 1, total: rawImages.length });
        }

        saveFaceCache(updatedCache);
        setFaceMap({ ...updatedCache });
        setScanning(false);
        setToastMsg(`Scanned ${rawImages.length} images and updated face clusters.`);
        setTimeout(() => setToastMsg(""), 3500);
      };

      // Sync Performer avatars into the reference prototype database
      const syncPerformerAvatars = async () => {
        if (!performers.length || scanning) return;
        setScanning(true);
        const withImages = performers.filter((p) => p.image_path);
        setScanProgress({ done: 0, total: withImages.length });
        let syncedCount = 0;

        for (let i = 0; i < withImages.length; i++) {
          const perf = withImages[i];
          try {
            const htmlImg = await loadImageAsync(perf.image_path);
            const detected = await detectFacesInImageElement(htmlImg);
            if (detected.length > 0 && detected[0].embedding) {
              savePerformerPrototype(perf.id, perf.name, detected[0].embedding);
              syncedCount++;
            }
          } catch (_) {}
          setScanProgress({ done: i + 1, total: withImages.length });
        }

        setScanning(false);
        setFaceMap({ ...loadFaceCache() });
        setToastMsg(`Indexed ${syncedCount} performer reference faces.`);
        setTimeout(() => setToastMsg(""), 3500);
      };

      // Build flat list of all detected faces across current images & cluster them
      const clusters = useMemo(() => {
        const allFaces = [];
        rawImages.forEach((imgObj, imgIdx) => {
          const entry = faceMap[String(imgObj.id)];
          if (!entry || !Array.isArray(entry.faces)) return;
          entry.faces.forEach((f, fIdx) => {
            allFaces.push({
              key: `${imgObj.id}_${fIdx}`,
              imageId: String(imgObj.id),
              imageIndex: imgIdx,
              imageTitle: imgObj.title || `Image #${imgObj.id}`,
              existingPerformers: imgObj.performers || [],
              faceIndex: fIdx,
              normBox: f.normBox,
              embedding: f.embedding,
              thumbDataUrl: f.thumbDataUrl,
              assignedPerformerId: f.assignedPerformerId,
              assignedPerformerName: f.assignedPerformerName,
            });
          });
        });
        return clusterDetectedFaces(allFaces, similarityThreshold, performers);
      }, [rawImages, faceMap, similarityThreshold, performers]);

      // Map imageId -> clusterId for sorting
      const imageClusterMap = useMemo(() => {
        const map = {};
        clusters.forEach((c) => {
          c.faces.forEach((f) => {
            if (!map[f.imageId]) map[f.imageId] = c.id;
          });
        });
        return map;
      }, [clusters]);

      // Sorted images for the Smart Image Sorter tab
      const sortedImages = useMemo(() => {
        const arr = [...rawImages];
        arr.sort((a, b) => {
          const facesA = (faceMap[String(a.id)] && faceMap[String(a.id)].faces) || [];
          const facesB = (faceMap[String(b.id)] && faceMap[String(b.id)].faces) || [];
          const perfA = (a.performers && a.performers.length) || 0;
          const perfB = (b.performers && b.performers.length) || 0;

          if (sortMode === "unassigned_first") {
            const unA = facesA.length > 0 && perfA === 0 ? 1 : 0;
            const unB = facesB.length > 0 && perfB === 0 ? 1 : 0;
            if (unB !== unA) return unB - unA;
            return facesB.length - facesA.length;
          }
          if (sortMode === "faces_desc") return facesB.length - facesA.length;
          if (sortMode === "faces_asc") return facesA.length - facesB.length;
          if (sortMode === "cluster") {
            return String(imageClusterMap[String(a.id)] || "zzz").localeCompare(
              String(imageClusterMap[String(b.id)] || "zzz")
            );
          }
          if (sortMode === "rating_desc") return (b.rating100 || 0) - (a.rating100 || 0);
          if (sortMode === "o_counter_desc") return (b.o_counter || 0) - (a.o_counter || 0);
          if (sortMode === "resolution_desc") {
            const resA =
              a.visual_files && a.visual_files[0]
                ? (a.visual_files[0].width || 0) * (a.visual_files[0].height || 0)
                : 0;
            const resB =
              b.visual_files && b.visual_files[0]
                ? (b.visual_files[0].width || 0) * (b.visual_files[0].height || 0)
                : 0;
            return resB - resA;
          }
          return String(a.title || a.id).localeCompare(String(b.title || b.id));
        });
        return arr;
      }, [rawImages, faceMap, sortMode, imageClusterMap]);

      // Apply a Performer label to an entire Face Cluster
      const handleApplyClusterPerformer = async (cluster) => {
        let perfId = selectedClusterPerformer[cluster.id];
        let perfName = "";

        if (!perfId && cluster.suggestions && cluster.suggestions[0]) {
          perfId = cluster.suggestions[0].performerId;
        }

        // Or create a new performer if typed
        const customName = (newClusterPerfName[cluster.id] || "").trim();
        if (customName) {
          const created = await createNewPerformerGQL(customName);
          if (created) {
            perfId = created.id;
            perfName = created.name;
            await refetchPerformers();
          }
        } else if (perfId) {
          const found = performers.find((p) => String(p.id) === String(perfId));
          perfName = found ? found.name : "";
        }

        if (!perfId) return;

        const activeFaces = cluster.faces.filter((f) => !excludedFaceKeys[f.key]);
        const targetImageIds = Array.from(new Set(activeFaces.map((f) => f.imageId)));
        const embeddings = activeFaces.map((f) => f.embedding);

        await assignPerformerToImagesGQL(targetImageIds, perfId, perfName, embeddings);
        setLabeledClusters((prev) => ({ ...prev, [cluster.id]: perfName || `Performer #${perfId}` }));
        setNewClusterPerfName((prev) => ({ ...prev, [cluster.id]: "" }));
        setToastMsg(
          `Assigned "${perfName || perfId}" to ${targetImageIds.length} image(s) in ${cluster.id}.`
        );
        setTimeout(() => setToastMsg(""), 3500);
        refetchImages();
      };

      // Batch assign performer to selected images in Sorter tab
      const handleBatchAssignSelectedImages = async () => {
        const ids = Object.keys(selectedImageIds).filter((k) => selectedImageIds[k]);
        if (!ids.length || !batchPerformerId) return;
        const found = performers.find((p) => String(p.id) === String(batchPerformerId));
        await assignPerformerToImagesGQL(ids, batchPerformerId, found ? found.name : "", []);
        setSelectedImageIds({});
        setToastMsg(`Assigned ${found ? found.name : "Performer"} to ${ids.length} selected images.`);
        setTimeout(() => setToastMsg(""), 3500);
        refetchImages();
      };

      const totalDetectedFaces = useMemo(
        () => clusters.reduce((acc, c) => acc + c.faces.length, 0),
        [clusters]
      );

      return React.createElement(
        "div",
        { className: "sft-workspace container-fluid" },
        // Header Bar
        React.createElement(
          "div",
          { className: "sft-header-bar" },
          React.createElement(
            "div",
            null,
            React.createElement(
              "h4",
              { className: "mb-1 d-flex align-items-center", style: { gap: "0.6rem" } },
              "Image Zoom Navigator & Native Facial Tagger",
              React.createElement(
                "span",
                {
                  className: `sft-badge-engine ${
                    nativeFaceDetectorState === "supported" ? "" : "fallback"
                  }`,
                },
                nativeFaceDetectorState === "supported"
                  ? "Native OS FaceDetector Active"
                  : "Canvas Haar+Symmetry Engine Active"
              )
            ),
            React.createElement(
              "div",
              { className: "text-muted small" },
              `${rawImages.length} images loaded • ${totalDetectedFaces} faces detected across ${clusters.length} similarity groups`
            )
          ),
          React.createElement(
            "div",
            { className: "d-flex align-items-center", style: { gap: "0.5rem" } },
            React.createElement(
              Button,
              {
                variant: activeTab === "clusters" ? "primary" : "outline-light",
                size: "sm",
                onClick: () => setActiveTab("clusters"),
              },
              `Face Clusters (${clusters.length})`
            ),
            React.createElement(
              Button,
              {
                variant: activeTab === "sorter" ? "primary" : "outline-light",
                size: "sm",
                onClick: () => setActiveTab("sorter"),
              },
              `Smart Image Sorter (${rawImages.length})`
            ),
            React.createElement(
              Button,
              {
                variant: "success",
                size: "sm",
                disabled: scanning || !rawImages.length,
                onClick: () => runBatchFaceScan(false),
              },
              scanning ? "Scanning…" : "Scan & Group Faces"
            ),
            React.createElement(
              Button,
              {
                variant: "outline-info",
                size: "sm",
                disabled: scanning,
                onClick: syncPerformerAvatars,
                title: "Build reference face embeddings from existing Stash Performer profile images",
              },
              "Sync Performer Face DB"
            )
          )
        ),

        // Progress or Toast Notification
        scanning &&
          React.createElement(
            "div",
            { className: "mb-3" },
            React.createElement(ProgressBar, {
              animated: true,
              now:
                scanProgress.total > 0
                  ? Math.round((scanProgress.done / scanProgress.total) * 100)
                  : 100,
              label: `${scanProgress.done} / ${scanProgress.total}`,
            })
          ),
        toastMsg &&
          React.createElement(
            "div",
            { className: "alert alert-success py-2 px-3 mb-3" },
            toastMsg
          ),

        // Filter & Clustering Controls Bar
        React.createElement(
          "div",
          { className: "sft-controls-row" },
          React.createElement(
            "div",
            { className: "sft-control-group" },
            React.createElement("label", null, "Search:"),
            React.createElement("input", {
              type: "text",
              className: "sft-input",
              placeholder: "Filter images…",
              value: searchQuery,
              onChange: (e) => {
                setSearchQuery(e.target.value);
                setPage(1);
              },
            })
          ),
          React.createElement(
            "div",
            { className: "sft-control-group" },
            React.createElement("label", null, "Filter:"),
            React.createElement(
              "select",
              {
                className: "sft-select",
                value: filterMode,
                onChange: (e) => setFilterMode(e.target.value),
              },
              React.createElement("option", { value: "all" }, "All Images"),
              React.createElement("option", { value: "untagged" }, "No Performers Linked"),
              React.createElement("option", { value: "unorganized" }, "Unorganized Only")
            )
          ),
          React.createElement(
            "div",
            { className: "sft-control-group" },
            React.createElement("label", null, "Batch Size:"),
            React.createElement(
              "select",
              {
                className: "sft-select",
                value: perPage,
                onChange: (e) => {
                  setPerPage(Number(e.target.value));
                  setPage(1);
                },
              },
              React.createElement("option", { value: 24 }, "24 Images"),
              React.createElement("option", { value: 40 }, "40 Images"),
              React.createElement("option", { value: 80 }, "80 Images"),
              React.createElement("option", { value: 150 }, "150 Images")
            )
          ),
          React.createElement(
            "div",
            { className: "sft-control-group" },
            React.createElement("label", null, "Page:"),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "secondary",
                disabled: page <= 1,
                onClick: () => setPage((p) => Math.max(1, p - 1)),
              },
              "◀"
            ),
            React.createElement("span", { className: "px-2" }, `Page ${page}`),
            React.createElement(
              Button,
              {
                size: "sm",
                variant: "secondary",
                onClick: () => setPage((p) => p + 1),
              },
              "▶"
            )
          ),
          activeTab === "clusters"
            ? React.createElement(
                "div",
                { className: "sft-control-group ml-auto" },
                React.createElement(
                  "label",
                  null,
                  `Grouping Similarity (${Math.round(similarityThreshold * 100)}%):`
                ),
                React.createElement("input", {
                  type: "range",
                  min: "0.50",
                  max: "0.94",
                  step: "0.02",
                  value: similarityThreshold,
                  onChange: (e) => setSimilarityThreshold( parseFloat(e.target.value) ),
                })
              )
            : React.createElement(
                "div",
                { className: "sft-control-group ml-auto" },
                React.createElement("label", null, "Sort Images By:"),
                React.createElement(
                  "select",
                  {
                    className: "sft-select",
                    value: sortMode,
                    onChange: (e) => setSortMode(e.target.value),
                  },
                  React.createElement(
                    "option",
                    { value: "unassigned_first" },
                    "Unassigned Faces First"
                  ),
                  React.createElement(
                    "option",
                    { value: "faces_desc" },
                    "Face Count (Most → Least)"
                  ),
                  React.createElement(
                    "option",
                    { value: "faces_asc" },
                    "Face Count (Least → Most)"
                  ),
                  React.createElement(
                    "option",
                    { value: "cluster" },
                    "Facial Similarity Cluster"
                  ),
                  React.createElement(
                    "option",
                    { value: "rating_desc" },
                    "Rating (High → Low)"
                  ),
                  React.createElement(
                    "option",
                    { value: "o_counter_desc" },
                    "O-Counter (High → Low)"
                  ),
                  React.createElement(
                    "option",
                    { value: "resolution_desc" },
                    "Resolution (High → Low)"
                  )
                )
              )
        ),

        // Main Content Area
        imagesLoading
          ? React.createElement(
              "div",
              { className: "text-center py-5" },
              React.createElement(Spinner, { animation: "border", variant: "info" }),
              React.createElement("div", { className: "mt-2" }, "Loading images from Stash…")
            )
          : activeTab === "clusters"
          ? // TAB 1: FACE CLUSTERS
            clusters.length === 0
            ? React.createElement(
                "div",
                {
                  className: "text-center py-5",
                  style: {
                    background: "rgba(25, 29, 38, 0.7)",
                    borderRadius: "10px",
                    border: "1px dashed rgba(255,255,255,0.2)",
                  },
                },
                React.createElement(
                  "h5",
                  null,
                  "No Face Clusters Indexed for This Page Yet"
                ),
                React.createElement(
                  "p",
                  { className: "text-muted mb-3" },
                  "Click 'Scan & Group Faces' above to detect faces across the current images and group similar people together automatically."
                ),
                React.createElement(
                  Button,
                  { variant: "success", onClick: () => runBatchFaceScan(true) },
                  "Scan & Group Faces Now"
                )
              )
            : React.createElement(
                "div",
                { className: "sft-cluster-grid" },
                clusters.map((cluster, cIdx) => {
                  const topSug = cluster.suggestions && cluster.suggestions[0];
                  const isLabeled = Boolean(labeledClusters[cluster.id]);
                  const activeFacesCount = cluster.faces.filter(
                    (f) => !excludedFaceKeys[f.key]
                  ).length;

                  return React.createElement(
                    "div",
                    {
                      key: cluster.id,
                      className: `sft-cluster-card ${isLabeled ? "labeled" : ""}`,
                    },
                    React.createElement(
                      "div",
                      { className: "sft-cluster-header" },
                      React.createElement(
                        "div",
                        null,
                        React.createElement(
                          "strong",
                          null,
                          labeledClusters[cluster.id]
                            ? `✓ ${labeledClusters[cluster.id]}`
                            : `Person Group #${cIdx + 1}`
                        ),
                        React.createElement(
                          "span",
                          { className: "text-muted small ml-2" },
                          `(${activeFacesCount} face${activeFacesCount === 1 ? "" : "s"})`
                        )
                      ),
                      topSug &&
                        topSug.similarity >= 0.65 &&
                        React.createElement(
                          Badge,
                          { variant: "info" },
                          `Match: ${topSug.name} (${Math.round(topSug.similarity * 100)}%)`
                        )
                    ),

                    // Cropped Face Thumbnails Strip
                    React.createElement(
                      "div",
                      { className: "sft-faces-strip" },
                      cluster.faces.map((face) => {
                        const isExcluded = Boolean(excludedFaceKeys[face.key]);
                        return React.createElement(
                          "div",
                          {
                            key: face.key,
                            className: `sft-face-thumb-wrap ${isExcluded ? "excluded" : ""}`,
                            title: `${face.imageTitle} — Click to inspect & zoom face`,
                            onClick: () => {
                              const idxInSorted = sortedImages.findIndex(
                                (im) => String(im.id) === String(face.imageId)
                              );
                              setViewerState({
                                open: true,
                                index: idxInSorted >= 0 ? idxInSorted : 0,
                                focusBox: face.normBox,
                              });
                            },
                          },
                          React.createElement("img", {
                            src: face.thumbDataUrl,
                            alt: face.imageTitle,
                          }),
                          React.createElement(
                            "button",
                            {
                              type: "button",
                              className: "sft-face-exclude-btn",
                              title: isExcluded
                                ? "Include face in cluster"
                                : "Exclude outlier face from cluster",
                              onClick: (e) => {
                                e.stopPropagation();
                                setExcludedFaceKeys((prev) => ({
                                  ...prev,
                                  [face.key]: !prev[face.key],
                                }));
                              },
                            },
                            isExcluded ? "+" : "✕"
                          )
                        );
                      })
                    ),

                    // Performer Assignment Controls for Cluster
                    React.createElement(
                      "div",
                      { className: "d-flex flex-column", style: { gap: "0.45rem" } },
                      React.createElement(
                        "div",
                        { className: "d-flex", style: { gap: "0.4rem" } },
                        React.createElement(
                          "select",
                          {
                            className: "sft-select flex-grow-1",
                            value:
                              selectedClusterPerformer[cluster.id] ||
                              (topSug && topSug.similarity >= 0.65 ? topSug.performerId : ""),
                            onChange: (e) =>
                              setSelectedClusterPerformer((prev) => ({
                                ...prev,
                                [cluster.id]: e.target.value,
                              })),
                          },
                          React.createElement(
                            "option",
                            { value: "" },
                            "-- Select Existing Performer --"
                          ),
                          performers.map((p) =>
                            React.createElement(
                              "option",
                              { key: p.id, value: p.id },
                              p.name
                            )
                          )
                        ),
                        React.createElement(
                          Button,
                          {
                            size: "sm",
                            variant: "primary",
                            disabled: activeFacesCount === 0,
                            onClick: () => handleApplyClusterPerformer(cluster),
                          },
                          "Label Group"
                        )
                      ),
                      React.createElement(
                        "div",
                        { className: "d-flex", style: { gap: "0.4rem" } },
                        React.createElement("input", {
                          type: "text",
                          className: "sft-input flex-grow-1",
                          placeholder: "Or create new performer name…",
                          value: newClusterPerfName[cluster.id] || "",
                          onChange: (e) =>
                            setNewClusterPerfName((prev) => ({
                              ...prev,
                              [cluster.id]: e.target.value,
                            })),
                        }),
                        React.createElement(
                          Button,
                          {
                            size: "sm",
                            variant: "outline-success",
                            disabled: !(newClusterPerfName[cluster.id] || "").trim(),
                            onClick: () => handleApplyClusterPerformer(cluster),
                          },
                          "+ Create & Label"
                        )
                      )
                    )
                  );
                })
              )
          : // TAB 2: SMART IMAGE SORTER & NAVIGATOR
            React.createElement(
              "div",
              null,
              // Batch action bar when images are selected
              Object.values(selectedImageIds).some(Boolean) &&
                React.createElement(
                  "div",
                  {
                    className:
                      "d-flex align-items-center justify-content-between bg-dark border border-info rounded p-2 mb-3",
                  },
                  React.createElement(
                    "span",
                    { className: "font-weight-bold text-info" },
                    `${
                      Object.values(selectedImageIds).filter(Boolean).length
                    } image(s) selected`
                  ),
                  React.createElement(
                    "div",
                    { className: "d-flex align-items-center", style: { gap: "0.5rem" } },
                    React.createElement(
                      "select",
                      {
                        className: "sft-select",
                        value: batchPerformerId,
                        onChange: (e) => setBatchPerformerId(e.target.value),
                      },
                      React.createElement("option", { value: "" }, "-- Assign Performer --"),
                      performers.map((p) =>
                        React.createElement("option", { key: p.id, value: p.id }, p.name)
                      )
                    ),
                    React.createElement(
                      Button,
                      {
                        size: "sm",
                        variant: "info",
                        disabled: !batchPerformerId,
                        onClick: handleBatchAssignSelectedImages,
                      },
                      "Apply to Selected"
                    ),
                    React.createElement(
                      Button,
                      {
                        size: "sm",
                        variant: "outline-secondary",
                        onClick: () => setSelectedImageIds({}),
                      },
                      "Clear Selection"
                    )
                  )
                ),

              React.createElement(
                "div",
                { className: "sft-image-grid" },
                sortedImages.map((imgObj, idx) => {
                  const thumbUrl =
                    (imgObj.paths &&
                      (imgObj.paths.thumbnail || imgObj.paths.preview || imgObj.paths.image)) ||
                    "";
                  const cachedFaces =
                    (faceMap[String(imgObj.id)] && faceMap[String(imgObj.id)].faces) || [];
                  const perfList = imgObj.performers || [];
                  const isSelected = Boolean(selectedImageIds[String(imgObj.id)]);
                  const clusterId = imageClusterMap[String(imgObj.id)];

                  return React.createElement(
                    "div",
                    {
                      key: imgObj.id,
                      className: `sft-image-card ${isSelected ? "selected" : ""}`,
                    },
                    React.createElement(
                      "div",
                      {
                        className: "sft-image-thumb-box",
                        onClick: () =>
                          setViewerState({ open: true, index: idx, focusBox: null }),
                      },
                      React.createElement("img", {
                        src: thumbUrl,
                        alt: imgObj.title || `Image ${imgObj.id}`,
                        loading: "lazy",
                      }),
                      React.createElement(
                        "div",
                        { className: "sft-card-badges" },
                        cachedFaces.length > 0 &&
                          React.createElement(
                            "span",
                            {
                              className: `sft-pill ${
                                perfList.length > 0 ? "assigned" : "unassigned"
                              }`,
                            },
                            `${cachedFaces.length} Face${cachedFaces.length > 1 ? "s" : ""}`
                          ),
                        clusterId &&
                          React.createElement("span", { className: "sft-pill" }, clusterId)
                      )
                    ),
                    React.createElement(
                      "div",
                      { className: "p-2 d-flex flex-column justify-content-between flex-grow-1" },
                      React.createElement(
                        "div",
                        { className: "d-flex align-items-center justify-content-between mb-1" },
                        React.createElement(
                          "label",
                          {
                            className: "mb-0 d-flex align-items-center small text-truncate",
                            style: { gap: "6px", cursor: "pointer" },
                          },
                          React.createElement("input", {
                            type: "checkbox",
                            checked: isSelected,
                            onChange: (e) => {
                              const checked = e.target.checked;
                              setSelectedImageIds((prev) => ({
                                ...prev,
                                [String(imgObj.id)]: checked,
                              }));
                            },
                          }),
                          React.createElement(
                            "span",
                            { className: "text-truncate", title: imgObj.title || "" },
                            imgObj.title || `Image #${imgObj.id}`
                          )
                        ),
                        React.createElement(
                          Button,
                          {
                            size: "sm",
                            variant: "link",
                            className: "p-0 text-info",
                            onClick: () =>
                              setViewerState({ open: true, index: idx, focusBox: null }),
                          },
                          "Zoom"
                        )
                      ),
                      React.createElement(
                        "div",
                        { className: "small text-muted text-truncate" },
                        perfList.length > 0
                          ? perfList.map((p) => p.name).join(", ")
                          : "No performers tagged"
                      )
                    )
                  );
                })
              )
            ),

        // Fullscreen Pan & Zoom Viewer Modal
        viewerState.open &&
          React.createElement(ZoomViewerModal, {
            images: sortedImages,
            initialIndex: viewerState.index,
            initialFocusBox: viewerState.focusBox,
            performers: performers,
            onClose: () => setViewerState({ open: false, index: 0, focusBox: null }),
            onImageUpdated: () => {
              setFaceMap({ ...loadFaceCache() });
              refetchImages();
            },
          })
      );
    };

    // =======================================================================
    // Component 3: Injection Wrapper for Individual Stash ImageDetailPanel
    // =======================================================================
    const ImageDetailFaceInspectorLauncher = function ({ image }) {
      const [openModal, setOpenModal] = useState(false);
      const { data: performersData } = GQL.useFindPerformersQuery({
        variables: {
          filter: { per_page: -1, sort: "name", direction: "ASC" },
        },
      });
      const performers =
        (performersData &&
          performersData.findPerformers &&
          performersData.findPerformers.performers) ||
        [];

      if (!image) return null;

      return React.createElement(
        "div",
        { className: "mt-3 pt-2 border-top border-secondary" },
        React.createElement(
          Button,
          {
            variant: "info",
            size: "sm",
            className: "w-100 mb-2",
            onClick: () => setOpenModal(true),
          },
          renderIcon(FA.faSearchPlus, "🔍"),
          " Open Zoom Navigator & Face Tagger"
        ),
        openModal &&
          React.createElement(ZoomViewerModal, {
            images: [image],
            initialIndex: 0,
            performers: performers,
            onClose: () => setOpenModal(false),
          })
      );
    };

    // =======================================================================
    // Register Routes & Patch Stash Core Components via PluginApi
    // =======================================================================

    // 1. Register "/plugins/face-tagger" route
    PluginApi.register.route("/plugins/face-tagger", FaceTaggerWorkspace);

    // 2. Patch MainNavBar.MenuItems to add a top-level "Face Tagger" nav button
    PluginApi.patch.before("MainNavBar.MenuItems", function (props) {
      return [
        {
          ...props,
          children: React.createElement(
            React.Fragment,
            null,
            props.children,
            React.createElement(
              "div",
              { className: "col-4 col-sm-3 col-md-2 col-lg-auto nav-link" },
              React.createElement(
                NavLink,
                {
                  className: "minimal p-4 p-xl-2 d-flex d-xl-inline-block flex-column justify-content-between align-items-center btn btn-primary",
                  to: "/plugins/face-tagger",
                  activeClassName: "active",
                },
                renderIcon(FA.faUserCheck || FA.faUser || FA.faImages, "👤"),
                React.createElement("span", { className: "ml-xl-2" }, "Face Tagger")
              )
            )
          ),
        },
      ];
    });

    // 3. Patch ImageDetailPanel so every single Image page has a one-click Zoom & Face Tagger launcher
    PluginApi.patch.after("ImageDetailPanel", function (props) {
      const originalResult = arguments[arguments.length - 1];
      return React.createElement(
        React.Fragment,
        null,
        originalResult,
        React.createElement(ImageDetailFaceInspectorLauncher, { image: props && props.image })
      );
    });
  }

  // Wait for window.PluginApi to be ready
  if (window.PluginApi) {
    initPluginWithApi(window.PluginApi);
  } else {
    let attempts = 0;
    const timer = setInterval(() => {
      attempts++;
      if (window.PluginApi) {
        clearInterval(timer);
        initPluginWithApi(window.PluginApi);
      } else if (attempts > 100) {
        clearInterval(timer);
      }
    }, 50);
  }
})();