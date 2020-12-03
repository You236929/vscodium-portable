//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"path"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

type config struct {
	Cleanup bool `yaml:"cleanup" mapstructure:"cleanup"`
}

var (
	app *portapps.App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Cleanup: false,
	}

	// Init app
	if app, err = portapps.NewWithCfg("vscodium-portable", "VSCodium", cfg); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "VSCodium.exe")
	app.Args = []string{
		"--log debug",
	}

	// Cleanup on exit
	if cfg.Cleanup {
		defer func() {
			utl.Cleanup([]string{
				path.Join(os.Getenv("APPDATA"), "VSCodium"),
			})
		}()
	}

	os.Setenv("VSCODE_APPDATA", utl.PathJoin(app.DataPath, "appdata"))
	os.Setenv("VSCODE_LOGS", utl.PathJoin(app.DataPath, "logs"))
	os.Setenv("VSCODE_EXTENSIONS", utl.PathJoin(app.DataPath, "extensions"))

	defer app.Close()
	app.Launch(os.Args[1:])
}
