#!/usr/bin/env bash
# Create a deploy-only SSH key and install it on the EC2 host you already SSH into.
# Prints the values to paste into GitHub. Does not upload anything to GitHub.
#
# Usage:
#   ./scripts/setup-github-deploy.sh ubuntu@YOUR_EC2_HOST
set -euo pipefail

TARGET="${1:-}"
if [[ -z "$TARGET" || "$TARGET" != *@* ]]; then
  echo "Usage: $0 <user>@<host>" >&2
  echo "  Example: $0 ubuntu@ec2-xx-xx-xx-xx.compute.amazonaws.com" >&2
  echo "  Use the same user@host you type when you SSH in by hand." >&2
  exit 1
fi

DEPLOY_USER="${TARGET%@*}"
DEPLOY_HOST="${TARGET#*@}"
KEY_PATH="${HOME}/.ssh/portfolio-github-deploy"

mkdir -p "${HOME}/.ssh"
chmod 700 "${HOME}/.ssh"

if [[ ! -f "$KEY_PATH" ]]; then
  ssh-keygen -t ed25519 -C "github-actions-portfolio-deploy" -f "$KEY_PATH" -N ""
  echo "Created ${KEY_PATH}"
else
  echo "Reusing existing ${KEY_PATH}"
fi

echo "Installing the public key on ${TARGET} (you may be asked for your usual SSH password or key)..."
# Prefer ssh-copy-id; fall back to appending authorized_keys.
if command -v ssh-copy-id >/dev/null 2>&1; then
  ssh-copy-id -i "${KEY_PATH}.pub" "$TARGET"
else
  ssh "$TARGET" "mkdir -p ~/.ssh && chmod 700 ~/.ssh && cat >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys" < "${KEY_PATH}.pub"
fi

echo "Checking that the new key can log in..."
ssh -i "$KEY_PATH" -o IdentitiesOnly=yes -o BatchMode=yes "$TARGET" 'echo ok' >/dev/null

SYNC_HINT=""
if ssh -i "$KEY_PATH" -o IdentitiesOnly=yes -o BatchMode=yes "$TARGET" 'test -f ~/portfolio/.env'; then
  SYNC_VALUE="$(ssh -i "$KEY_PATH" -o IdentitiesOnly=yes -o BatchMode=yes "$TARGET" "grep -E '^SYNC_API_KEY=' ~/portfolio/.env | tail -1 | cut -d= -f2- | tr -d '\"' || true")"
  if [[ -n "$SYNC_VALUE" ]]; then
    SYNC_HINT="$SYNC_VALUE"
  fi
fi

cat <<EOF

Key works. Add these in GitHub → Settings → Secrets and variables → Actions.

Repository variables (Variables tab) — not secret:
  DEPLOY_USER = ${DEPLOY_USER}
  DEPLOY_HOST = ${DEPLOY_HOST}

Repository secret (Secrets tab):
  SSH_PRIVATE_KEY = the full contents of ${KEY_PATH}
  (run: cat ${KEY_PATH})

Optional secret:
  SYNC_API_KEY = value of SYNC_API_KEY in ~/portfolio/.env on the server
EOF

if [[ -n "$SYNC_HINT" ]]; then
  echo "  Found on the server. Paste this as SYNC_API_KEY:"
  echo "  ${SYNC_HINT}"
else
  echo "  Not found. Deploy still works without it; blog sync will be skipped."
fi

echo
echo "Do not commit the private key or paste it into chat. After saving in GitHub, merge the deploy PR and push to main (or run the workflow from the Actions tab)."
