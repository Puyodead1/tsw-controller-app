$env:CGO_CFLAGS="-I$PWD/vcpkg_installed/x64-windows/include"
$env:CGO_LDFLAGS="-L$PWD/vcpkg_installed/x64-windows/lib -lSDL2"
wails build ldflags="-X 'main.VERSION=$VERSION'"