#!/usr/bin/env bash
set -euo pipefail

REPO="sergiocarracedo/on-a-meet"
INSTALL_DIR="/usr/local/bin"

main() {
  local os arch version

  os="$(detect_os)"
  arch="$(detect_arch)"
  echo "==> Detected: ${os}/${arch}"

  version="$(fetch_latest_version)"
  echo "==> Latest release: ${version}"

  local asset="on-a-meet_${version}_${os}_${arch}.tar.gz"
  local url="https://github.com/${REPO}/releases/download/${version}/${asset}"

  local tmpdir
  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT

  echo "==> Downloading ${asset}..."
  curl -fsSL -o "${tmpdir}/on-a-meet.tar.gz" "$url"

  echo "==> Extracting..."
  tar xzf "${tmpdir}/on-a-meet.tar.gz" -C "$tmpdir" on-a-meet
  chmod +x "${tmpdir}/on-a-meet"

  echo "==> Installing to ${INSTALL_DIR}/on-a-meet..."
  if [ "$(id -u)" -eq 0 ]; then
    mv "${tmpdir}/on-a-meet" "${INSTALL_DIR}/on-a-meet"
  else
    sudo mv "${tmpdir}/on-a-meet" "${INSTALL_DIR}/on-a-meet"
  fi

  echo "==> Done! Run 'on-a-meet --help' to get started."
}

detect_os() {
  case "$(uname -s)" in
    Linux)  echo "linux" ;;
    Darwin) echo "darwin" ;;
    *)      echo "fatal: unsupported OS: $(uname -s)" >&2; exit 1 ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *)            echo "fatal: unsupported architecture: $(uname -m)" >&2; exit 1 ;;
  esac
}

fetch_latest_version() {
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | cut -d '"' -f 4
}

main "$@"
