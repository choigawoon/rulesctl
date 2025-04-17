#!/usr/bin/env node

const os = require("os");
const { spawn } = require("child_process");
const path = require("path");
const https = require("https");
const fs = require("fs");

const OWNER = "choigawoon";
const REPO = "rulesctl";
const VERSION = "v0.1.0";

async function downloadFile(url, dest) {
  console.log(`Downloading from: ${url}`);
  console.log(`Saving to: ${dest}`);
  
  return new Promise((resolve, reject) => {
    https.get(url, response => {
      console.log(`Response status: ${response.statusCode}`);
      console.log(`Response headers:`, response.headers);

      if (response.statusCode === 302) {
        console.log(`Following redirect to: ${response.headers.location}`);
        https.get(response.headers.location, redirectedResponse => {
          console.log(`Redirect response status: ${redirectedResponse.statusCode}`);
          if (redirectedResponse.statusCode !== 200) {
            return reject(new Error(`Failed to download: ${redirectedResponse.statusCode}`));
          }
          const file = fs.createWriteStream(dest);
          redirectedResponse.pipe(file);
          file.on("finish", () => {
            file.close();
            fs.chmodSync(dest, "755");
            resolve();
          });
        }).on("error", error => {
          console.error("Redirect request error:", error);
          reject(error);
        });
      } else if (response.statusCode === 200) {
        const file = fs.createWriteStream(dest);
        response.pipe(file);
        file.on("finish", () => {
          file.close();
          fs.chmodSync(dest, "755");
          resolve();
        });
      } else {
        reject(new Error(`Failed to download: ${response.statusCode}`));
      }
    }).on("error", error => {
      console.error("Initial request error:", error);
      reject(error);
    });
  });
}

async function install() {
  const platform = os.platform();
  const arch = os.arch();

  console.log(`Installing for platform: ${platform}, architecture: ${arch}`);

  const platformMap = {
    darwin: "Darwin",
    linux: "Linux",
    win32: "Windows"
  };

  const archMap = {
    x64: "x86_64",
    arm64: "arm64"
  };

  if (!platformMap[platform] || !archMap[arch]) {
    console.error(`Unsupported platform: ${platform} ${arch}`);
    process.exit(1);
  }

  const assetName = `rulesctl_${platformMap[platform]}_${archMap[arch]}${platform === "win32" ? ".zip" : ".tar.gz"}`;
  const downloadUrl = `https://github.com/${OWNER}/${REPO}/releases/download/${VERSION}/${assetName}`;

  const binDir = path.join(__dirname, "bin");
  if (!fs.existsSync(binDir)) {
    console.log(`Creating bin directory: ${binDir}`);
    fs.mkdirSync(binDir, { recursive: true });
  }

  const binPath = path.join(binDir, platform === "win32" ? "rulesctl-win.exe" : `rulesctl-${platform}`);

  console.log(`Downloading ${assetName}...`);
  try {
    await downloadFile(downloadUrl, binPath);
    console.log("rulesctl installed successfully!");
    console.log(`Binary path: ${binPath}`);
    // 파일 존재 여부와 권한 확인
    const stats = fs.statSync(binPath);
    console.log(`File exists: ${fs.existsSync(binPath)}`);
    console.log(`File permissions: ${stats.mode.toString(8)}`);
  } catch (error) {
    console.error("Failed to download rulesctl:", error);
    process.exit(1);
  }
}

if (process.argv[2] === "install") {
  install();
} else {
  const platform = os.platform();
  let bin = "";

  if (platform === "darwin") bin = "rulesctl-darwin";
  else if (platform === "linux") bin = "rulesctl-linux";
  else if (platform === "win32") bin = "rulesctl-win.exe";
  else {
    console.error(`Unsupported platform: ${platform}`);
    process.exit(1);
  }

  const binPath = path.join(__dirname, "bin", bin);
  
  // 실행 전 파일 존재 여부 확인
  if (!fs.existsSync(binPath)) {
    console.error(`Binary not found at: ${binPath}`);
    console.error("Please try reinstalling the package");
    process.exit(1);
  }

  const args = process.argv.slice(2);
  console.log(`Executing: ${binPath} ${args.join(" ")}`);
  const child = spawn(binPath, args, { stdio: "inherit" });

  child.on("error", (err) => {
    console.error("Failed to execute binary:", err);
    process.exit(1);
  });

  child.on("exit", code => process.exit(code));
}
