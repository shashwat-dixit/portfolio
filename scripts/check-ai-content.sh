#!/usr/bin/env bash
# Verify that humans and AI crawlers can fetch markdown in the expected format.
#
# Usage:
#   ./scripts/check-ai-content.sh
#   ./scripts/check-ai-content.sh https://shashwatdixit.com https://api.shashwatdixit.com
#   ./scripts/check-ai-content.sh http://localhost:4321 http://localhost:8080
set -euo pipefail

SITE_URL="${1:-https://shashwatdixit.com}"
API_URL="${2:-https://api.shashwatdixit.com}"
SITE_URL="${SITE_URL%/}"
API_URL="${API_URL%/}"

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

fetch() {
  local url="$1"
  shift
  curl -fsS "$@" "$url"
}

echo "Checking AI markdown endpoints on ${SITE_URL} (API ${API_URL})"

content_type() {
  curl -fsS -D - -o /dev/null "$@" | tr -d '\r' | awk -F': ' 'tolower($1)=="content-type"{print $2; exit}'
}

llms_type="$(content_type "${SITE_URL}/llms.txt")"
[[ "$llms_type" == text/plain* ]] || fail "/llms.txt Content-Type should be text/plain, got '${llms_type}'"
llms="$(fetch "${SITE_URL}/llms.txt")"
[[ "$llms" == \#* ]] || fail "/llms.txt should start with a markdown heading"
[[ "$llms" == *"index.md"* ]] || fail "/llms.txt should link to index.md"
[[ "$llms" == *".md"* ]] || fail "/llms.txt should list .md post URLs"
echo "  /llms.txt OK"

robots="$(fetch "${SITE_URL}/robots.txt")"
[[ "$robots" == *"GPTBot"* ]] || fail "/robots.txt should allow GPTBot"
[[ "$robots" == *"/llms.txt"* ]] || fail "/robots.txt should point at llms.txt"
echo "  /robots.txt OK"

home_md="$(fetch "${SITE_URL}/index.md")"
[[ "$home_md" == \#* ]] || fail "/index.md should start with a markdown heading"
echo "  /index.md OK"

home_accept="$(curl -fsS -H 'Accept: text/markdown' "${SITE_URL}/blog")"
[[ "$home_accept" == \#* ]] || fail "Accept: text/markdown on /blog should return markdown"
echo "  Accept: text/markdown on /blog OK"

home_ua="$(curl -fsS -A 'ChatGPT-User/1.0' "${SITE_URL}/blog")"
[[ "$home_ua" == \#* ]] || fail "ChatGPT-User on /blog should return markdown, not HTML"
echo "  ChatGPT-User on /blog OK"

browser_home="$(curl -fsS -H 'Accept: text/html' -A 'Mozilla/5.0' "${SITE_URL}/")"
[[ "$browser_home" == *"<html"* ]] || fail "Browsers should still get HTML on /"
echo "  Browser HTML on / OK"

slug="$(printf '%s\n' "$llms" | sed -n 's/.*\/blog\/\([^)]*\)\.md.*/\1/p' | head -n 1 || true)"
if [[ -z "${slug}" ]]; then
  echo "  No blog posts listed in llms.txt yet; skipping post checks."
else
  post_md="$(fetch "${SITE_URL}/blog/${slug}.md")"
  [[ "$post_md" == ---* ]] || fail "/blog/${slug}.md should start with YAML frontmatter"
  [[ "$post_md" == *"title:"* ]] || fail "/blog/${slug}.md frontmatter should include title"
  [[ "$post_md" == *"slug:"* ]] || fail "/blog/${slug}.md frontmatter should include slug"
  echo "  /blog/${slug}.md OK"

  api_md="$(fetch "${API_URL}/api/posts/${slug}?format=md")"
  [[ "$api_md" == ---* ]] || fail "API ?format=md should return markdown with frontmatter"
  echo "  API ?format=md OK"

  api_accept_type="$(curl -fsS -D - -o /tmp/ai-md-body -H 'Accept: text/markdown' "${API_URL}/api/posts/${slug}" | tr -d '\r' | awk -F': ' 'tolower($1)=="content-type"{print $2; exit}')"
  [[ "$api_accept_type" == text/markdown* ]] || fail "API Accept: text/markdown Content-Type was '${api_accept_type}'"
  [[ "$(cat /tmp/ai-md-body)" == ---* ]] || fail "API Accept: text/markdown body should be markdown"
  echo "  API Accept: text/markdown OK"
fi

echo
echo "All AI markdown checks passed."
echo "In ChatGPT (web search on), ask:"
echo "  Read ${SITE_URL}/llms.txt and summarize this site. Then open one .md blog URL and quote the first heading."
