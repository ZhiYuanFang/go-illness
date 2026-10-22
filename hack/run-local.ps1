# Local run: load env from manifest/docker/env then gf run main.go
# Default env: .env.prod (override with -EnvFile)
param(
    [string]$EnvFile = "",
    [switch]$SkipOssEndpoint
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

if (-not $EnvFile) {
    $EnvFile = Join-Path $Root "manifest\docker\env\.env.prod"
} elseif (-not [System.IO.Path]::IsPathRooted($EnvFile)) {
    $EnvFile = Join-Path $Root $EnvFile
}

if (-not (Test-Path -LiteralPath $EnvFile)) {
    Write-Error "Env file not found: $EnvFile. Copy manifest\docker\.env.example to manifest\docker\env\.env.prod and fill values."
    exit 1
}

Write-Host "Loading env: $EnvFile"
Get-Content -LiteralPath $EnvFile -Encoding UTF8 | ForEach-Object {
    $line = $_.Trim()
    if ($line -eq "" -or $line.StartsWith("#")) { return }
    $i = $line.IndexOf("=")
    if ($i -lt 1) { return }
    $name = $line.Substring(0, $i).Trim()
    $value = $line.Substring($i + 1).Trim()
    Set-Item -Path "Env:$name" -Value $value
}

if (-not $SkipOssEndpoint -and [string]::IsNullOrWhiteSpace($env:ILLNESS_OSS_ENDPOINT)) {
    $env:ILLNESS_OSS_ENDPOINT = "oss-cn-hangzhou.aliyuncs.com"
    Write-Host "Set ILLNESS_OSS_ENDPOINT=$($env:ILLNESS_OSS_ENDPOINT) (local public endpoint; use -SkipOssEndpoint to skip)"
}

$gf = Get-Command gf -ErrorAction SilentlyContinue
if (-not $gf) {
    Write-Error "gf CLI not found. Install GoFrame CLI first (e.g. make cli), then retry."
    exit 1
}

Write-Host "Starting: gf run main.go (root: $Root)"
& gf run main.go
exit $LASTEXITCODE
