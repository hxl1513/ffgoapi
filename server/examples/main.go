// Package examples demonstrates how to wire every framework component together.
//
// ExampleMain is the illustrative entry-point function.  It is NOT a real main
// function so that the package can be imported and tested; run it from a
// separate _main.go file in the same directory if you want an executable.
//
// Prerequisite: a server package (produced by earlier tasks) must expose:
//
//	type Server interface {
//	    DB() *model.DB
//	    RegisterRoutes(func(*gin.RouterGroup))
//	    Start() error
//	    Shutdown(context.Context) error
//	}
//
//	func NewServer(cfg *config.Config) (*Server, error)
package examples

import (
	"fmt"
	"log"

	"github.com/hxl1513/ffgoapi/config"
)

// ExampleMain demonstrates the wiring pattern for a production application.
//
// The server package referenced below is expected to be provided by earlier
// tasks in this project (Tasks 1-11).  The code is written so that when that
// package exists it can be uncommented and run as-is.
func ExampleMain() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// ── Bootstrap the HTTP server ─────────────────────────────────────────
	// srv, err := server.NewServer(cfg)
	// if err != nil {
	//     log.Fatalf("server: %v", err)
	// }
	// defer srv.Shutdown(context.Background())

	fmt.Println("go-web-api-framework example")
	fmt.Printf("  port : %d\n", cfg.Server.Port)
	fmt.Printf("  db   : %s (%s)\n", cfg.Database.Type, cfg.Database.DBName)

	// ── Register feature routes ────────────────────────────────────────────
	// Each feature module provides its own service and controller; main only
	// stitches them together.
	//
	// srv.RegisterRoutes(func(r *gin.RouterGroup) {
	//     userSvc  := users.NewUserService(srv.DB())
	//     userCtrl := users.NewUserController(userSvc)
	//     userCtrl.RegisterRoutes(r)
	// })
}

// _main.go (not in this package) would call ExampleMain() or replicate the
// pattern shown above under a real func main().
