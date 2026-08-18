const $ = (id) => document.getElementById(id);

async function fetchJSON(url, opts) {
  const res = await fetch(url, opts);
  const data = await res.json();
  if (!res.ok) throw new Error("[" + data.code + "] " + data.message);
  return data;
}

function showError(msg) {
  $("error").textContent = msg;
  $("error").classList.remove("hidden");
}

function hideError() {
  $("error").classList.add("hidden");
}

async function loadMeta() {
  try {
    const meta = await fetchJSON("/api/meta", {});
    const sel = $("example");
    sel.innerHTML = "";
    for (const name of meta.examples) {
      const opt = document.createElement("option");
      opt.value = name;
      opt.textContent = name;
      sel.appendChild(opt);
    }
  } catch (e) {
    // 无 meta 时保留空下拉
  }
}

async function loadExample() {
  const name = $("example").value;
  const res = await fetch("/example/" + name);
  const spec = await res.json();
  $("kind").value = spec.kind;
  $("order").value = spec.order;
  $("cutoff").value = spec.cutoff;
  $("window").value = spec.window || "hamming";
  await design();
}

function drawCurve(freq, mag) {
  const canvas = $("curve");
  const ctx = canvas.getContext("2d");
  const dpr = window.devicePixelRatio || 1;
  const w = canvas.clientWidth;
  const h = 300;
  canvas.width = w * dpr;
  canvas.height = h * dpr;
  ctx.scale(dpr, dpr);
  ctx.clearRect(0, 0, w, h);

  let vmin = Infinity, vmax = -Infinity;
  for (const v of mag) {
    if (!isFinite(v)) continue;
    if (v < vmin) vmin = v;
    if (v > vmax) vmax = v;
  }
  if (vmin === vmax) { vmin -= 10; vmax += 10; }
  const pad = 14;
  const fmin = freq[0], fmax = freq[freq.length - 1];
  const fx = (x) => pad + ((x - fmin) / (fmax - fmin || 1)) * (w - 2 * pad);
  const fy = (y) => h - pad - ((y - vmin) / (vmax - vmin)) * (h - 2 * pad);

  ctx.strokeStyle = "#c9ced8";
  ctx.lineWidth = 1;
  ctx.strokeRect(pad, pad, w - 2 * pad, h - 2 * pad);
  ctx.font = "11px sans-serif";
  ctx.fillStyle = "#555";
  ctx.fillText(vmax.toFixed(1) + " dB", 2, pad + 4);
  ctx.fillText(vmin.toFixed(1) + " dB", 2, h - pad - 2);

  ctx.strokeStyle = "#2f6feb";
  ctx.lineWidth = 2;
  ctx.beginPath();
  for (let i = 0; i < freq.length; i++) {
    const x = fx(freq[i]);
    const y = fy(mag[i]);
    if (i === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
  }
  ctx.stroke();

  ctx.strokeStyle = "#d5482d";
  ctx.setLineDash([4, 4]);
  ctx.beginPath();
  ctx.moveTo(pad, fy(-3));
  ctx.lineTo(w - pad, fy(-3));
  ctx.stroke();
  ctx.setLineDash([]);
  ctx.fillStyle = "#d5482d";
  ctx.fillText("-3 dB", 4, fy(-3) - 2);
}

async function design() {
  hideError();
  const payload = {
    kind: $("kind").value,
    order: parseInt($("order").value, 10),
    cutoff: parseFloat($("cutoff").value),
    window: $("window").value,
  };
  try {
    const f = await fetchJSON("/api/design", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    const b = "b = [" + f.b.map((v) => v.toFixed(6)).join(", ") + "]";
    const a = "a = [" + f.a.map((v) => v.toFixed(6)).join(", ") + "]";
    $("coeff").textContent = "抽头数 " + f.b.length + "\n" + b + "\n" + a;

    const resp = await fetchJSON("/api/response", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ b: f.b, a: f.a, points: 256 }),
    });
    drawCurve(resp.freq, resp.mag_db);

    const zp = await fetchJSON("/api/zplane", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ b: f.b, a: f.a }),
    });
    $("zp").textContent = "零点 " + zp.zeros.length + " 个，极点 " + zp.poles.length +
      " 个，稳定: " + (zp.stable ? "是" : "否");
  } catch (e) {
    showError(e.message);
  }
}

$("design").addEventListener("click", design);
$("load").addEventListener("click", loadExample);
$("kind").addEventListener("change", () => {
  $("window").disabled = $("kind").value !== "fir";
  if ($("kind").value === "iir" && parseInt($("order").value, 10) > 16) {
    $("order").value = 4;
  }
});

(async function init() {
  await loadMeta();
  if ($("example").options.length === 0) {
    const opts = [
      ["fir_hamming.json", "fir_hamming.json"],
      ["fir_hann.json", "fir_hann.json"],
      ["iir_bw4.json", "iir_bw4.json"],
    ];
    for (const [v, label] of opts) {
      const opt = document.createElement("option");
      opt.value = v;
      opt.textContent = label;
      $("example").appendChild(opt);
    }
  }
  $("window").disabled = $("kind").value !== "fir";
  await loadExample();
})();
