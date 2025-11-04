Write-Host "Building version: $env:VERSION"
$env:CGO_CFLAGS="-I$PWD/vcpkg_installed/x64-windows/include"
$env:CGO_LDFLAGS="-L$PWD/vcpkg_installed/x64-windows/lib -lSDL2 -X 'main.VERSION=$env:VERSION'"
wails build