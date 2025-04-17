#!/usr/bin/env node

const os = require("os");
const { spawn } = require("child_process");
const path = require("path");
const https = require("https");
const fs = require("fs");
const { execSync } = require("child_process");

const OWNER = "choigawoon";
const REPO = "rulesctl";
const VERSION = "v0.1.3";

async function downloadFile(url, dest) {
  console.log(`Downloading from: ${url}`);
  console.log(`Saving to: ${dest}`);
  
  return new Promise((resolve, reject) => {
    const options = {
      headers: {
        'Cache-Control': 'no-cache',
        'Pragma': 'no-cache',
        'User-Agent': 'rulesctl-installer'
      }
    };

    https.get(url, options, response => {
      console.log(`Response status: ${response.statusCode}`);
      
      if (response.statusCode === 302 || response.statusCode === 301) {
        console.log(`Following redirect to: ${response.headers.location}`);
        https.get(response.headers.location, options, redirectedResponse => {
          if (redirectedResponse.statusCode !== 200) {
            return reject(new Error(`Failed to download: ${redirectedResponse.statusCode}`));
          }
          const file = fs.createWriteStream(dest);
          redirectedResponse.pipe(file);
          file.on("finish", () => {
            file.close();
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

async function extractBinary(archivePath, targetDir, platform) {
  console.log(`Extracting binary from ${archivePath} to ${targetDir}`);
  
  const isWindows = platform === "win32";
  const binName = isWindows ? "rulesctl.exe" : "rulesctl";
  const finalName = isWindows ? "rulesctl-win.exe" : `rulesctl-${platform}`;
  
  try {
    // 임시 디렉토리 생성
    const tempDir = path.join(os.tmpdir(), `rulesctl-${Date.now()}`);
    fs.mkdirSync(tempDir, { recursive: true });
    
    // 압축 해제
    if (isWindows) {
      // Windows: unzip 사용
      execSync(`powershell -command "Expand-Archive -Path '${archivePath}' -DestinationPath '${tempDir}'"`, { stdio: 'inherit' });
    } else {
      // macOS/Linux: tar 사용
      execSync(`tar -xzf "${archivePath}" -C "${tempDir}"`, { stdio: 'inherit' });
    }
    
    // 바이너리 찾기
    const extractedBinary = path.join(tempDir, binName);
    const targetBinary = path.join(targetDir, finalName);
    
    // 바이너리 이동
    fs.copyFileSync(extractedBinary, targetBinary);
    fs.chmodSync(targetBinary, "755");
    
    // 임시 파일 정리
    fs.rmSync(tempDir, { recursive: true, force: true });
    fs.rmSync(archivePath);
    
    console.log(`Binary extracted and installed at: ${targetBinary}`);
    return targetBinary;
  } catch (error) {
    console.error("Failed to extract binary:", error);
    throw error;
  }
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
  
  // Try direct download URL first
  const directUrl = `https://github.com/${OWNER}/${REPO}/releases/download/${VERSION}/${assetName}`;
  // Fallback to API URL with cache bypass
  const apiUrl = `https://api.github.com/repos/${OWNER}/${REPO}/releases/latest/assets/${assetName}?t=${Date.now()}`;

  const binDir = path.join(__dirname, "bin");
  if (!fs.existsSync(binDir)) {
    console.log(`Creating bin directory: ${binDir}`);
    fs.mkdirSync(binDir, { recursive: true });
  }

  const archivePath = path.join(os.tmpdir(), assetName);

  try {
    // Try direct download first
    console.log(`Downloading ${assetName} using direct URL...`);
    try {
      await downloadFile(directUrl, archivePath);
    } catch (directError) {
      console.log("Direct download failed, trying API URL...");
      await downloadFile(apiUrl, archivePath);
    }
    
    // Extract and install
    await extractBinary(archivePath, binDir, platform);
    
    console.log("rulesctl installed successfully!");
  } catch (error) {
    console.error("Failed to install rulesctl:", error);
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
