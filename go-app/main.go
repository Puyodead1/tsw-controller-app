package main

import (
	"embed"
	"fmt"
	"os"
	"path"
	"tsw_controller_app/logger"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var VERSION = "0.0.0-dev"

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	fmt.Printf("Version %s\n", VERSION)

	config_dir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("could not find user config directory %e", err))
	}

	exec_file, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("could not find executable %e", err))
	}

	global_config_dir := path.Join(config_dir, "tswcontrollerapp/config")
	local_config_dir := path.Join(path.Dir(exec_file), "config")
	os.MkdirAll(global_config_dir, 0o755)
	os.MkdirAll(local_config_dir, 0o755)

	app := NewApp(AppConfig{
		GlobalConfigDir: global_config_dir,
		LocalConfigDir:  local_config_dir,
	})

	err = wails.Run(&options.App{
		Title:  "TSW Controller Utility",
		Width:  600,
		Height: 400,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewGpuIsDisabled: false,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
		},
	})

	if err != nil {
		logger.Logger.Error("[main] error", "error", err)
	}
}
