package main

import (
	"context"
	"fmt"
	"os"
	"time"

	foundation "github.com/Pt-ll/OI-Desk/backend/app"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	foundation *foundation.Runtime
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	initialized, err := foundation.Bootstrap(ctx, "OI-Desk", os.Getenv("OIDESK_DATA_DIR"))
	if err != nil {
		fmt.Printf("OI-Desk foundation initialization failed: %v\n", err)
	} else {
		a.foundation = initialized
		a.foundation.Logger.Info("application foundation initialized", "database", a.foundation.Paths.Database)
	}
	if os.Getenv("OIDESK_PHASE0_AUTO_EXIT") == "1" {
		go func() {
			select {
			case <-time.After(15 * time.Second):
				fmt.Println("PHASE0_WEBVIEW_TIMEOUT")
				runtime.Quit(ctx)
			case <-ctx.Done():
			}
		}()
	}
}

func (a *App) shutdown(_ context.Context) {
	if a.foundation != nil {
		if err := a.foundation.Close(); err != nil {
			fmt.Printf("OI-Desk foundation shutdown failed: %v\n", err)
		}
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ReportFrontendProbe is a Phase 0 diagnostic, removed when Phase 1 begins.
func (a *App) ReportFrontendProbe(bindingOK bool, isolationOK bool) {
	fmt.Printf("PHASE0_WEBVIEW_RESULT binding=%t isolation=%t\n", bindingOK, isolationOK)
	if os.Getenv("OIDESK_PHASE0_AUTO_EXIT") == "1" {
		runtime.Quit(a.ctx)
	}
}
