package app

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"report/internal/config"
	"report/internal/services"
	"report/internal/storage"
	"report/internal/ui"
)

func Run() {
	a := app.NewWithID("report_app")

	cfg := config.Load()

	if cfg.Theme == "light" {
		a.Settings().SetTheme(theme.LightTheme())
	} else {
		a.Settings().SetTheme(theme.DarkTheme())
	}

	w := a.NewWindow("📞 Отчёт по звонкам")

	store := storage.NewFileStorage()
	service := service.NewStatsService(store)

	content := ui.NewMainWindow(w, service, a, &cfg)
	w.SetContent(content)
	w.Resize(ui.DefaultWindowSize())
	w.CenterOnScreen()

	w.ShowAndRun()
}
