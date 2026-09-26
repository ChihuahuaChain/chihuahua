#!/usr/bin/env bash
# Writes a software upgrade proposal for chihuahuad, ready for
#   chihuahuad tx gov submit-proposal proposal.json --from <key> ...
#
# The plan info is the upgrade-info.json that the release workflow attaches to
# the GitHub release, so cosmovisor (DAEMON_ALLOW_DOWNLOAD_BINARIES=true) can
# download and verify the new binary by itself. The binaries are downloaded and
# checked against it first.
#
# Usage:
#   scripts/upgrade-proposal.sh -n <plan name> -t <release tag> -s <summary.md> \
#     [-T <title>] [-H <height> | -a <UTC time, e.g. 2026-10-06T15:00:00Z>] \
#     [-d <deposit>] [-o proposal.json]
#
# With -a the height is estimated from the average block time of the last
# 10,000 blocks. The upgrade happens at the height, not at the time.
#
# Needs: bash, curl, jq, sha256sum, gh (logged in, to read draft releases).
set -euo pipefail

REPO=ChihuahuaChain/chihuahua
RPC=${RPC:-https://rpc.chihuahua.wtf}
GOV_AUTHORITY=chihuahua10d07y265gmmuvt4z0w9aw880jnsr700jeh7th3
DEPOSIT=5000000000000uhuahua
OUT=proposal.json
NAME= TAG= SUMMARY= TITLE= HEIGHT= AT=

usage() { sed -n '2,17p' "$0" | sed 's/^# \{0,1\}//'; exit 1; }
while getopts "n:t:s:T:H:a:d:o:h" opt; do
  case $opt in
    n) NAME=$OPTARG ;;
    t) TAG=$OPTARG ;;
    s) SUMMARY=$OPTARG ;;
    T) TITLE=$OPTARG ;;
    H) HEIGHT=$OPTARG ;;
    a) AT=$OPTARG ;;
    d) DEPOSIT=$OPTARG ;;
    o) OUT=$OPTARG ;;
    *) usage ;;
  esac
done
[ -n "$NAME" ] && [ -n "$TAG" ] && [ -f "$SUMMARY" ] || usage
[ -n "$HEIGHT" ] || [ -n "$AT" ] || usage
TITLE=${TITLE:-"Chihuahua $TAG software upgrade"}

block() { curl -fsS "$RPC/block${1:+?height=$1}" | jq -r '.result.block.header | "\(.height) \(.time)"'; }

if [ -z "$HEIGHT" ]; then
  read -r now now_time < <(block)
  read -r _ old_time < <(block $((now - 10000)))
  secs=$(( $(date -ud "$now_time" +%s) - $(date -ud "$old_time" +%s) ))
  target=$(date -ud "$AT" +%s)
  left=$(( target - $(date -ud "$now_time" +%s) ))
  [ "$left" -gt 0 ] || { echo "time $AT is in the past" >&2; exit 1; }
  HEIGHT=$(( now + left * 10000 / secs ))
  echo "block time $(echo "scale=3; $secs/10000" | bc) s, height now $now: upgrade height $HEIGHT for $AT" >&2
fi

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
gh release download "$TAG" -R "$REPO" -D "$work" \
  -p upgrade-info.json -p chihuahuad_linux_amd64 -p chihuahuad_linux_arm64
info=$(jq -c . "$work/upgrade-info.json")

for platform in linux/amd64 linux/arm64; do
  url=$(jq -r --arg p "$platform" '.binaries[$p]' <<<"$info")
  want=${url##*checksum=sha256:}
  file=chihuahuad_${platform/\//_}
  got=$(sha256sum "$work/$file" | cut -d' ' -f1)
  [ "$got" = "$want" ] || { echo "$file: checksum $got, upgrade-info says $want" >&2; exit 1; }
  case "$url" in "https://github.com/$REPO/releases/download/$TAG/$file?"*) ;; *) echo "unexpected url $url" >&2; exit 1 ;; esac
  echo "$platform ok: $want" >&2
done

jq -n \
  --arg authority "$GOV_AUTHORITY" --arg name "$NAME" --arg height "$HEIGHT" --arg info "$info" \
  --arg title "$TITLE" --rawfile summary "$SUMMARY" --arg deposit "$DEPOSIT" \
  '{
    messages: [{
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      authority: $authority,
      plan: {name: $name, height: $height, info: $info}
    }],
    metadata: "",
    deposit: $deposit,
    title: $title,
    summary: $summary,
    expedited: false
  }' > "$OUT"

echo "wrote $OUT: upgrade $NAME at height $HEIGHT with the $TAG binaries" >&2
