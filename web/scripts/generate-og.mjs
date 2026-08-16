#!/usr/bin/env node
/**
 * Renders scripts/og-template.html to public/og.png (1200x630).
 *
 * Usage (from web/): node scripts/generate-og.mjs
 */
import { spawn } from "node:child_process";
import fs from "node:fs";
import http from "node:http";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const webDir = path.join(__dirname, "..");
const templatePath = path.join(__dirname, "og-template.html");
const avatarPath = path.join(webDir, "public", "avatar.jpg");
const outPath = path.join(webDir, "public", "og.png");

const MIME = {
  ".html": "text/html; charset=utf-8",
  ".jpg": "image/jpeg",
  ".jpeg": "image/jpeg",
  ".png": "image/png",
  ".woff2": "font/woff2",
};

function chromePath() {
  const candidates = [
    process.env.CHROME_PATH,
    "/usr/bin/google-chrome-stable",
    "/usr/local/bin/google-chrome",
    "/usr/bin/google-chrome",
    "/usr/bin/chromium-browser",
    "/usr/bin/chromium",
  ].filter(Boolean);
  return candidates.find((p) => fs.existsSync(p));
}

function startServer() {
  return new Promise((resolve, reject) => {
    const server = http.createServer((req, res) => {
      const url = new URL(req.url, "http://127.0.0.1");
      let file = url.pathname === "/" ? templatePath : null;
      if (url.pathname === "/avatar.jpg" || url.pathname === "./avatar.jpg") {
        file = avatarPath;
      }
      if (!file || !fs.existsSync(file)) {
        res.writeHead(404);
        res.end("not found");
        return;
      }
      const ext = path.extname(file);
      res.writeHead(200, { "Content-Type": MIME[ext] || "application/octet-stream" });
      fs.createReadStream(file).pipe(res);
    });
    server.listen(0, "127.0.0.1", () => {
      const { port } = server.address();
      resolve({ server, url: `http://127.0.0.1:${port}/` });
    });
    server.on("error", reject);
  });
}

function run(cmd, args) {
  return new Promise((resolve, reject) => {
    const child = spawn(cmd, args, { stdio: "inherit" });
    child.on("error", reject);
    child.on("exit", (code) => {
      if (code === 0) resolve();
      else reject(new Error(`${cmd} exited with code ${code}`));
    });
  });
}

const chrome = chromePath();
if (!chrome) {
  console.error("No Chrome/Chromium found. Set CHROME_PATH or install Chrome.");
  process.exit(1);
}

if (!fs.existsSync(avatarPath)) {
  console.error(`Missing avatar at ${avatarPath}`);
  process.exit(1);
}

const { server, url } = await startServer();
try {
  await run(chrome, [
    "--headless=new",
    "--no-sandbox",
    "--disable-gpu",
    "--hide-scrollbars",
    "--force-device-scale-factor=2",
    "--window-size=1200,630",
    "--virtual-time-budget=8000",
    `--screenshot=${outPath}`,
    url,
  ]);
} finally {
  server.close();
}

const stat = fs.statSync(outPath);
console.log(`Wrote ${outPath} (${stat.size} bytes)`);
