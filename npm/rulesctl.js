#!/usr/bin/env node

const os = require("os");
const { spawn } = require("child_process");
const path = require("path");

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
const args = process.argv.slice(2);
const child = spawn(binPath, args, { stdio: "inherit" });

child.on("exit", code => process.exit(code));
