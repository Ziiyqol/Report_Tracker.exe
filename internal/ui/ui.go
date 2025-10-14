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

// NewMainWindow создаёт содержимое главного окна приложения.
func NewMainWindow(w fyne.Window, s *service.StatsService, a fyne.App, cfg *config.Config) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("📞 Отчёт по звонкам", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	statsLabel := widget.NewLabel(s.GetStatsText())

	// Кнопки действий
	oldBtn := widget.NewButton("➕ Старый звонок", func() {
		s.AddOld()
		statsLabel.SetText(s.GetStatsText())
	})
	newBtn := widget.NewButton("🆕 Новый звонок", func() {
		s.AddNew()
		statsLabel.SetText(s.GetStatsText())
	})
	recorded := widget.NewButton("✅ Записан", func() {
		s.AddRecorded()
		statsLabel.SetText(s.GetStatsText())
	})
	thinking := widget.NewButton("💭 Думает", func() {
		s.AddThinking()
		statsLabel.SetText(s.GetStatsText())
	})
	reject := widget.NewButton("❌ Не подходит", func() {
		s.AddRejected()
		statsLabel.SetText(s.GetStatsText())
	})
	noAnswer := widget.NewButton("📵 Не дозвонился", func() {
		s.AddNoAnswer()
		statsLabel.SetText(s.GetStatsText())
	})

	undoBtn := widget.NewButton("↩️ Откатить действие", func() {
		s.UndoLast()
		statsLabel.SetText(s.GetStatsText())
	})

	resetBtn := widget.NewButton("🧹 Сбросить статистику", func() {
		dialog.NewConfirm("Подтверждение", "Вы уверены, что хотите сбросить статистику?", func(ok bool) {
			if ok {
				s.Reset()
				statsLabel.SetText(s.GetStatsText())
			}
		}, w).Show()
	})

	saveReportBtn := widget.NewButton("💾 Скачать отчёт", func() {
		if err := s.SaveReport(); err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Готово", "Отчёт сохранён в report.txt", w)
		}
	})

	// 🌗 Переключатель темы
	isDark := cfg.Theme != "light"
	themeSwitch := widget.NewButton("", nil)
	updateThemeButton := func() {
		if isDark {
			themeSwitch.SetText("🌙 Тёмная тема")
		} else {
			themeSwitch.SetText("☀️ Светлая тема")
		}
	}
	updateThemeButton()

	// ✨ Подпись внизу справа
	footerLabel := canvas.NewText("Сделано @Ziiyqol", getFooterColor(isDark))
	footerLabel.TextSize = 10
	footerLabel.Alignment = fyne.TextAlignTrailing

	// Плавная анимация смены темы — исправлена обработка альфы (присваиваем color.NRGBA)
	themeSwitch.OnTapped = func() {
		// Создаём полупрозрачный прямоугольник-оверлей
		overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		overlay.Resize(w.Canvas().Size()) // заполнить холст
		w.Canvas().Overlays().Add(overlay)

		// Затемнение/осветление (инкрементируем альфу, присваивая новый color.NRGBA)
		if isDark {
			// если сейчас тёмная — план: сначала "вспышка белого" (осветление), затем смена темы, потом убрать
			for aAlpha := 0; aAlpha <= 255; aAlpha += 25 {
				overlay.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: uint8(aAlpha)}
				overlay.Refresh()
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			// если сейчас светлая — затемнить чёрным
			for aAlpha := 0; aAlpha <= 255; aAlpha += 25 {
				overlay.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: uint8(aAlpha)}
				overlay.Refresh()
				time.Sleep(10 * time.Millisecond)
			}
		}

		// Переключаем тему
		if isDark {
			a.Settings().SetTheme(theme.LightTheme())
			cfg.Theme = "light"
		} else {
			a.Settings().SetTheme(theme.DarkTheme())
			cfg.Theme = "dark"
		}
		isDark = !isDark
		_ = config.Save(*cfg) // сохраняем выбор (игнорируем ошибку здесь)

		// Обновляем цвет подписи под текущую тему
		footerLabel.Color = getFooterColor(isDark)
		footerLabel.Refresh()

		// Анимация исчезновения оверлея (уменьшаем альфу)
		if isDark {
			// сейчас тёмная — убираем белый оверлей
			for aAlpha := 255; aAlpha >= 0; aAlpha -= 25 {
				overlay.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: uint8(aAlpha)}
				overlay.Refresh()
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			// сейчас светлая — убираем чёрный оверлей
			for aAlpha := 255; aAlpha >= 0; aAlpha -= 25 {
				overlay.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: uint8(aAlpha)}
				overlay.Refresh()
				time.Sleep(10 * time.Millisecond)
			}
		}

		w.Canvas().Overlays().Remove(overlay)
	}

	header := container.NewVBox(title, widget.NewSeparator())

	mainPanel := container.NewVBox(
		statsLabel,
		widget.NewSeparator(),
		container.NewHBox(oldBtn, newBtn),
		widget.NewSeparator(),
		container.NewGridWithColumns(2, recorded, thinking, reject, noAnswer),
		widget.NewSeparator(),
		container.NewHBox(undoBtn, resetBtn),
		widget.NewSeparator(),
		container.NewHBox(saveReportBtn, themeSwitch),
	)

	// Размещаем footer справа внизу с помощью Border
	content := container.NewBorder(nil, footerLabel, nil, nil,
		container.NewPadded(container.NewVBox(header, mainPanel)),
	)

	return content
}

// getFooterColor возвращает цвет подписи в зависимости от темы.
func getFooterColor(isDark bool) color.Color {
	if isDark {
		return color.NRGBA{R: 100, G: 200, B: 255, A: 255} // светло-голубой для тёмной темы
	}
	return color.NRGBA{R: 30, G: 100, B: 220, A: 255} // насыщенно-синий для светлой темы
}

// DefaultWindowSize — удобный дефолт
func DefaultWindowSize() fyne.Size {
	return fyne.NewSize(420, 500)
}
