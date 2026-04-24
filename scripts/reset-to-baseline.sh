#!/usr/bin/env bash
# Force-restore main to the pristine baseline. Wipes everything after
# baseline-v1. Use whenever the test app has drifted from pristine
# (PRs merged, manual edits, whatever).
set -euo pipefail

TAG="${1:-baseline-v1}"

git fetch origin --tags
git checkout main
git reset --hard "$TAG"
git push --force origin main

echo "main force-reset to $TAG. All bugs restored."
echo "Any commits/PRs after $TAG are gone from main (still in reflog locally)."
