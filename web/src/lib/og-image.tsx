import { readFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import satori from "satori";
import { Resvg } from "@resvg/resvg-js";
import { DATA } from "@/data/resume";
import { CONFIG } from "@/data/config";

export type OgImageInput = {
  title: string;
  description: string;
  eyebrow?: string;
  author?: string;
};

const OG_WIDTH = 1200;
const OG_HEIGHT = 630;
const OG_SCALE = 2;

function resolvePublicFile(rel: string): string | null {
  const cwd = process.cwd();
  const candidates = [
    join(cwd, "public", rel),
    join(cwd, "dist", "client", rel),
    join(cwd, "client", rel),
  ];
  return candidates.find((path) => existsSync(path)) ?? null;
}

let fonts:
  | { regular: Buffer; semibold: Buffer }
  | null
  | undefined;
let avatarDataUri: string | null | undefined;

function loadFonts() {
  if (fonts) return fonts;
  if (fonts === null) {
    throw new Error("OG fonts missing");
  }

  const regularPath = resolvePublicFile("fonts/Outfit-Regular.ttf");
  const semiboldPath = resolvePublicFile("fonts/Outfit-SemiBold.ttf");
  if (!regularPath || !semiboldPath) {
    fonts = null;
    throw new Error("OG fonts missing");
  }

  fonts = {
    regular: readFileSync(regularPath),
    semibold: readFileSync(semiboldPath),
  };
  return fonts;
}

function loadAvatarDataUri(): string | null {
  if (avatarDataUri !== undefined) return avatarDataUri;

  const avatarPath = resolvePublicFile("avatar.jpg");
  if (!avatarPath) {
    avatarDataUri = null;
    return null;
  }

  avatarDataUri = `data:image/jpeg;base64,${readFileSync(avatarPath).toString("base64")}`;
  return avatarDataUri;
}

function siteHost(): string {
  return CONFIG.site.url.replace(/^https?:\/\//, "").replace(/\/$/, "");
}

function OgCard({
  title,
  description,
  eyebrow = "Blog",
  author = DATA.name,
  avatar,
}: OgImageInput & { avatar: string | null }) {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        width: "100%",
        height: "100%",
        backgroundColor: "#171717",
        color: "#fafafa",
        padding: "56px 72px 48px",
        fontFamily: "Outfit",
        position: "relative",
      }}
    >
      <div
        style={{
          display: "flex",
          position: "absolute",
          top: 0,
          left: 0,
          right: 0,
          height: 160,
          backgroundImage:
            "linear-gradient(to bottom, rgba(255,255,255,0.08), transparent)",
        }}
      />

      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: 28,
          position: "relative",
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
          }}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              border: "1px solid rgba(255,255,255,0.12)",
              backgroundColor: "#e8e8e8",
              borderRadius: 14,
              padding: "6px 16px",
            }}
          >
            <div
              style={{
                display: "flex",
                color: "#171717",
                fontSize: 15,
                fontWeight: 500,
                letterSpacing: "0.02em",
              }}
            >
              {eyebrow}
            </div>
          </div>
        </div>

        <div
          style={{
            display: "flex",
            flexDirection: "column",
            gap: 16,
            maxWidth: 980,
          }}
        >
          <div
            style={{
              display: "flex",
              fontSize: 58,
              fontWeight: 600,
              letterSpacing: "-0.045em",
              lineHeight: 1.08,
              color: "#fafafa",
              maxHeight: 188,
              overflow: "hidden",
            }}
          >
            {title}
          </div>
          {description ? (
            <div
              style={{
                display: "flex",
                fontSize: 24,
                fontWeight: 400,
                lineHeight: 1.45,
                color: "#a3a3a3",
                maxHeight: 104,
                overflow: "hidden",
              }}
            >
              {description}
            </div>
          ) : null}
        </div>
      </div>

      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          position: "relative",
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 16,
          }}
        >
          {avatar ? (
            <img
              src={avatar}
              width={56}
              height={56}
              style={{
                width: 56,
                height: 56,
                borderRadius: 9999,
                objectFit: "cover",
                border: "1px solid rgba(255,255,255,0.12)",
              }}
            />
          ) : null}
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: 4,
            }}
          >
            <div
              style={{
                display: "flex",
                fontSize: 20,
                fontWeight: 600,
                color: "#fafafa",
              }}
            >
              {author}
            </div>
            <div
              style={{
                display: "flex",
                fontSize: 16,
                color: "#a3a3a3",
              }}
            >
              {siteHost()}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export async function renderOgImage(input: OgImageInput): Promise<Uint8Array> {
  const loadedFonts = loadFonts();
  const svg = await satori(
    <OgCard {...input} avatar={loadAvatarDataUri()} />,
    {
      width: OG_WIDTH,
      height: OG_HEIGHT,
      fonts: [
        {
          name: "Outfit",
          data: loadedFonts.regular,
          weight: 400,
          style: "normal",
        },
        {
          name: "Outfit",
          data: loadedFonts.semibold,
          weight: 500,
          style: "normal",
        },
        {
          name: "Outfit",
          data: loadedFonts.semibold,
          weight: 600,
          style: "normal",
        },
      ],
    }
  );

  const resvg = new Resvg(svg, {
    fitTo: { mode: "width", value: OG_WIDTH * OG_SCALE },
    font: { loadSystemFonts: false },
  });

  return resvg.render().asPng();
}

export function readFallbackOgImage(): Uint8Array | null {
  const path = resolvePublicFile("og.png");
  if (!path) return null;
  return readFileSync(path);
}

export function pngResponse(png: Uint8Array, cache = true): Response {
  return new Response(Buffer.from(png), {
    headers: {
      "Content-Type": "image/png",
      "Cache-Control": cache
        ? "public, max-age=86400, stale-while-revalidate=604800"
        : "no-store",
    },
  });
}

export function fallbackOgResponse(): Response {
  const fallback = readFallbackOgImage();
  if (!fallback) {
    return new Response("OG image unavailable", { status: 500 });
  }
  return pngResponse(fallback);
}
