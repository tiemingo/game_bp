package starter

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Liphium/magic/v2"
	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/scripting"
)

// BuildMagicConfig creates a magic config for the server.
func BuildMagicConfig() magic.Config {
	return magic.Config{
		AppName: "server",
		PlanDeployment: func(ctx *mconfig.Context) {

			// Get Port from config.. parse and use
			portValue := os.Getenv("PORT")
			portUint, err := strconv.ParseUint(portValue, 10, 32)
			if err != nil || portUint == 0 {
				portUint = 8000
			}
			port := ctx.ValuePort(uint(portUint))

			// Set the listeting address the server will use
			listen := mconfig.ValueWithBase([]mconfig.EnvironmentValue{port}, func(s []string) string {
				return fmt.Sprintf("127.0.0.1:%s", s[0])

			})

			ctx.WithEnvironment(mconfig.Environment{
				"LISTEN":    listen,
				"PORT":      port,
				"LOG_LEVEL": mconfig.ValueStatic("debug"),
				"LOG_ENV":   mconfig.ValueStatic("development"),
			})

			ctx.LoadSecretsToEnvironment(".env")
		},
		StartFunction: Start,
		Scripts:       []scripting.Script{},
	}
}
