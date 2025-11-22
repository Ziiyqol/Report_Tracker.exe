package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"report/internal/config"
	"report/internal/services"
)

// NewMainWindow создает контент окна
func NewMainWindow(w fyne.Window, s *services.StatsService, a fyne.App, cfg *config.Config) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("📞 Отчёт по звонкам", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	statsLabel := widget.NewLabel(s.GetStatsText())

	// Хелпер для обновления текста статистики
	refreshStats := func() {
		statsLabel.SetText(s.GetStatsText())
	}

	// -- Кнопки верхнего уровня --
	// Используем Grid с 2 колонками, чтобы кнопки делили ширину пополам
	oldBtn := widget.NewButton("➕ Старый звонок", func() { s.AddOld(); refreshStats() })
	newBtn := widget.NewButton("🆕 Новый звонок", func() { s.AddNew(); refreshStats() })
	topButtons := container.NewGridWithColumns(2, oldBtn, newBtn)

	// -- Кнопки результатов --
	recordedBtn := widget.NewButton("✅ Записан", func() { s.AddRecorded(); refreshStats() })
	thinkingBtn := widget.NewButton("💭 Думает", func() { s.AddThinking(); refreshStats() })
	rejectBtn := widget.NewButton("❌ Не подходит", func() { s.AddRejected(); refreshStats() })
	noAnswerBtn := widget.NewButton("📵 Не дозвонился", func() { s.AddNoAnswer(); refreshStats() })

	// Сетка 2x2 для основных действий
	resultGrid := container.NewGridWithColumns(2, recordedBtn, thinkingBtn, rejectBtn, noAnswerBtn)

	// -- Кнопка Резерв --
	// Она идет отдельной строкой на всю ширину под сеткой результатов
	reservedBtn := widget.NewButton("🗂 Записан в резерв", func() { s.AddReserved(); refreshStats() })

	// -- Управление --
	undoBtn := widget.NewButton("↩️ Откатить", func() { s.UndoLast(); refreshStats() })
	resetBtn := widget.NewButton("🧹 Сброс", func() {
		dialog.NewConfirm("Подтверждение", "Сбросить всю статистику?", func(ok bool) {
			if ok {
				s.Reset()
				refreshStats()
			}
		}, w).Show()
	})
	controlButtons := container.NewGridWithColumns(2, undoBtn, resetBtn)

	saveReportBtn := widget.NewButton("💾 Скачать отчёт", func() {
		if err := s.SaveReportToFile(); err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Готово", "Отчёт сохранён в report.txt", w)
		}
	})

	// -- Тема --
	isDark := cfg.Theme != "light"
	themeSwitch := widget.NewButton("", nil)

	// Функция обновления текста кнопки темы
	updateThemeText := func() {
		if isDark {
			themeSwitch.SetText("☀️ Светлая тема") // Если сейчас темная, предлагаем светлую
		} else {
			themeSwitch.SetText("🌙 Тёмная тема") // Если сейчас светлая, предлагаем темную
		}
	}
	updateThemeText()

	// -- Footer --
	footerLabel := canvas.NewText("Telegram Автора: @Ziiyqol", getFooterColor(isDark))
	footerLabel.TextSize = 10
	footerLabel.Alignment = fyne.TextAlignTrailing

	// Логика смены темы
	themeSwitch.OnTapped = func() {
		overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		overlay.Resize(w.Canvas().Size())
		w.Canvas().Overlays().Add(overlay)

		// Анимация затемнения
		targetR, targetG, targetB := uint8(0), uint8(0), uint8(0)
		if isDark {
			targetR, targetG, targetB = 255, 255, 255
		} // Вспышка белого при переходе на светлую

		for a := 0; a <= 200; a += 25 {
			overlay.FillColor = color.NRGBA{R: targetR, G: targetG, B: targetB, A: uint8(a)}
			overlay.Refresh()
			time.Sleep(5 * time.Millisecond)
		}

		// Смена темы
		if isDark {
			// Включаем светлую
			a.Settings().SetTheme(NewForcedTheme(theme.VariantLight))
			cfg.Theme = "light"
		} else {
			// Включаем темную
			a.Settings().SetTheme(NewForcedTheme(theme.VariantDark))
			cfg.Theme = "dark"
		}
		isDark = !isDark
		_ = config.Save(*cfg)

		updateThemeText() // Обновляем текст кнопки
		footerLabel.Color = getFooterColor(isDark)
		footerLabel.Refresh()

		// Анимация проявления
		for a := 200; a >= 0; a -= 25 {
			overlay.FillColor = color.NRGBA{R: targetR, G: targetG, B: targetB, A: uint8(a)}
			overlay.Refresh()
			time.Sleep(5 * time.Millisecond)
		}
		w.Canvas().Overlays().Remove(overlay)
	}

	// -- Компоновка --
	header := container.NewVBox(title, widget.NewSeparator())

	// Основная панель
	mainPanel := container.NewVBox(
		statsLabel,
		widget.NewSeparator(),
		topButtons,
		widget.NewSeparator(),
		resultGrid,  // Сетка 2x2
		reservedBtn, // Резерв (на всю ширину)
		widget.NewSeparator(),
		controlButtons,
		widget.NewSeparator(),
		container.NewGridWithColumns(2, saveReportBtn, themeSwitch),
	)

	return container.NewBorder(nil, footerLabel, nil, nil,
		container.NewPadded(container.NewVBox(header, mainPanel)),
	)
}

func getFooterColor(isDark bool) color.Color {
	if isDark {
		return color.NRGBA{R: 100, G: 200, B: 255, A: 255}
	}
	return color.NRGBA{R: 30, G: 100, B: 220, A: 255}
}

func DefaultWindowSize() fyne.Size {
	return fyne.NewSize(420, 550)
}
