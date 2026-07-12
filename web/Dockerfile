FROM node:22-alpine AS base
RUN corepack enable && corepack prepare pnpm@11.1.2 --activate
WORKDIR /app

FROM base AS deps
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc* ./
RUN pnpm install --frozen-lockfile --ignore-scripts && \
    pnpm rebuild esbuild

FROM base AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ARG PUBLIC_API_URL
ARG PUBLIC_POSTHOG_KEY
ARG PUBLIC_POSTHOG_HOST
ENV PUBLIC_API_URL=${PUBLIC_API_URL}
ENV PUBLIC_POSTHOG_KEY=${PUBLIC_POSTHOG_KEY}
ENV PUBLIC_POSTHOG_HOST=${PUBLIC_POSTHOG_HOST}
RUN pnpm build

FROM node:22-alpine AS runtime
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
ENV HOST=0.0.0.0
ENV PORT=4321
EXPOSE 4321
CMD ["node", "dist/server/entry.mjs"]
