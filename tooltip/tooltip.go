package tooltip

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Style определяет внешний вид тултипа.
type Style struct {
	Background  color.Color
	BorderColor color.Color
	BorderWidth float32
	FontBold    bool
}

// DefaultStyle возвращает стиль по умолчанию.
func DefaultStyle() Style {
	return Style{
		Background:  color.NRGBA{R: 22, G: 22, B: 22, A: 220},
		BorderColor: color.NRGBA{R: 255, G: 192, B: 26, A: 180},
		BorderWidth: 1.5,
		FontBold:    false,
	}
}

type Manager struct {
	current fyne.Widget
	popUp   *widget.PopUp
	timer   *time.Timer
	canvas  fyne.Canvas
	style   Style
}

func NewManager(canvas fyne.Canvas) *Manager {
	return &Manager{
		canvas: canvas,
		style:  DefaultStyle(),
	}
}

// SetStyle устанавливает новый стиль для всех последующих тултипов.
func (m *Manager) SetStyle(s Style) {
	m.style = s
}

// Show отображает тултип с текущим стилем.
func (m *Manager) Show(w fyne.Widget, text string, pos fyne.Position) {
	if m.current == w && m.popUp != nil && m.popUp.Visible() {
		m.popUp.Move(pos.Add(fyne.NewPos(10, 10)))
		return
	}
	m.Hide()
	m.current = w
	m.timer = time.AfterFunc(300*time.Millisecond, func() {
		fyne.Do(func() {
			m.showPopUp(text, pos)
		})
	})
}

func (m *Manager) showPopUp(text string, pos fyne.Position) {
	const maxWidth = 333
	const padding = float32(1)      // 8
	const outerPadding = float32(1) // 6

	// Создаём Label с переносом
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.TextStyle = fyne.TextStyle{Bold: m.style.FontBold}

	// Устанавливаем ширину Label, чтобы он мог рассчитать высоту
	label.Resize(fyne.NewSize(maxWidth, label.MinSize().Height))
	labelSize := label.MinSize()

	// Вычисляем размеры контента
	width := maxWidth + 2*padding
	height := labelSize.Height + 2*padding

	// Фон
	bg := canvas.NewRectangle(m.style.Background)
	bg.CornerRadius = 4

	// Рамка
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = m.style.BorderWidth
	border.StrokeColor = m.style.BorderColor
	border.CornerRadius = 4

	// Контейнер без автоматического layout
	content := container.NewWithoutLayout()
	content.Add(bg)
	content.Add(border)
	content.Add(label)

	// Позиционируем элементы вручную
	bg.Move(fyne.NewPos(0, 0))
	bg.Resize(fyne.NewSize(width, height))

	border.Move(fyne.NewPos(0, 0))
	border.Resize(fyne.NewSize(width, height))

	label.Move(fyne.NewPos(padding, padding))
	label.Resize(fyne.NewSize(maxWidth, labelSize.Height))

	// Задаём размер контента
	content.Resize(fyne.NewSize(width, height))

	// Создаём PopUp с небольшими внешними отступами
	popUp := widget.NewPopUp(content, m.canvas)
	popUp.Modal = false
	popUpSize := fyne.NewSize(width+outerPadding*2, height+outerPadding*2)
	popUp.Resize(popUpSize)

	// Позиционируем с учётом границ экрана
	winSize := m.canvas.Size()
	x := pos.X // + 10
	y := pos.Y // + 10
	if x+popUpSize.Width > winSize.Width {
		x = winSize.Width - popUpSize.Width - 5
	}
	if y+popUpSize.Height > winSize.Height {
		y = winSize.Height - popUpSize.Height - 5
	}
	popUp.Move(fyne.NewPos(x, y))
	popUp.Show()
	m.popUp = popUp
}

func (m *Manager) Hide() {
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
	if m.popUp != nil {
		m.popUp.Hide()
		m.popUp = nil
	}
	m.current = nil
}
