const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const https = require("https");

const pkg = require("./package.json");
const version = pkg.version;
const name = "on-a-meet";

function platform() {
  const os = process.platform;
  if (os === "darwin") return "darwin";
  if (os === "linux") return "linux";
  if (os === "win32") return "windows";
  throw new Error("Unsupported OS: " + os);
}

function arch() {
  const a = process.arch;
  if (a === "x64") return "amd64";
  if (a === "arm64") return "arm64";
  throw new Error("Unsupported arch: " + a);
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, (res) => {
      if (res.statusCode !== 200) {
        reject(new Error("Download failed: " + res.statusCode));
        return;
      }
      res.pipe(file);
      file.on("finish", () => {
        file.close();
        resolve();
      });
    }).on("error", reject);
  });
}

async function main() {
  const vendorDir = path.join(__dirname, "..", "vendor");
  fs.mkdirSync(vendorDir, { recursive: true });

  const os = platform();
  const a = arch();
  const archiveName = `${name}_${version}_${os}_${a}.tar.gz`;
  const url = `https://github.com/sergiocarracedo/${name}/releases/download/v${version}/${archiveName}`;
  const archivePath = path.join(vendorDir, archiveName);

  console.log(`Downloading ${name} v${version} (${os}/${a})...`);
  await download(url, archivePath);

  execSync(`tar xzf "${archivePath}" -C "${vendorDir}"`, { stdio: "inherit" });
  fs.unlinkSync(archivePath);

  const binaryPath = path.join(vendorDir, name);
  fs.chmodSync(binaryPath, 0o755);
  console.log(`Installed ${name} to ${binaryPath}`);
}

main().catch((err) => {
  console.error(err.message);
  process.exit(1);
});
