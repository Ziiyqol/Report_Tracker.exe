package app

import (
	"os"
	"report/internal/config"
	"report/internal/services"
	"report/internal/storage"
	"report/internal/ui"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func Run() {
	a := fyneApp.NewWithID("report_app")
	cfg := config.Load()

	if cfg.Theme == "light" {
		a.Settings().SetTheme(ui.NewForcedTheme(theme.VariantLight))
	} else {
		a.Settings().SetTheme(ui.NewForcedTheme(theme.VariantDark))
	}

	w := a.NewWindow("Report Tracker v1.1")

	// Загрузка иконки (без изменений)
	iconData, err := os.ReadFile("icon.png")
	if err == nil {
		w.SetIcon(fyne.NewStaticResource("icon.png", iconData))
	}

	store := storage.NewFileStorage()
	statsService := services.NewStatsService(store)

	content := ui.NewMainWindow(w, statsService, a, &cfg)
	w.SetContent(content)

	w.Resize(ui.DefaultWindowSize())
	w.CenterOnScreen()
	w.ShowAndRun()
}
