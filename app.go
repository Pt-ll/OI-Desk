package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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
