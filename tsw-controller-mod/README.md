# Building the go-socket lib
- Install TDM-GCC and set `CGO_ENABLED=1` and `CC=gcc`
- Run `go build -buildmode=c-shared -o tsw_controller_app_socket_lib.dll`
- Run `./generate_lib.ps1` to generate the .lib file for linking on Windows

# Building the mod
- Add the mod to the CppMods directory in UE4SS
- Update the proxy path before running `xmake config` to `user32.dll` for vendoring
- Update the ue4ss default path (` const fs::path ue4ssPath = currentPath`) in `UE4ss/proxy_generator/main.cpp` for vendoring

