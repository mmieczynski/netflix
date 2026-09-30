$ErrorActionPreference = 'Stop'
$prepGoCommand = Get-Command go -ErrorAction SilentlyContinue
if ($prepGoCommand) {
    $prepGoExe = $prepGoCommand.Source
} else {
    $prepGoExe = Join-Path $env:TEMP 'netflix-prep-go/go/bin/go.exe'
    if (-not (Test-Path -LiteralPath $prepGoExe)) {
        throw 'Go was not found. Install Go from https://go.dev/dl/ and rerun this script.'
    }
}
Push-Location $PSScriptRoot
try {
    & $prepGoExe test ./... -count=1 -cover
    if ($LASTEXITCODE -ne 0) { throw 'Go tests failed.' }
    & $prepGoExe vet ./...
    if ($LASTEXITCODE -ne 0) { throw 'Go vet failed.' }
} finally {
    Pop-Location
}
