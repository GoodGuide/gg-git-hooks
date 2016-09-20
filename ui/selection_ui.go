package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

const (
	defaultBg = termbox.ColorBlack
	defaultFg = termbox.ColorWhite
	headerMsg = "Select one or more Pivotal Tracker stories to tag this commit:"
	footerMsg = "Move cursor with standard directionals; Select = SPACE; Confirm = ENTER; Reload Stories = R"
)

type SelectionUI struct {
	OptionsFunc func(forceReload bool) ([]string, error)
	Selections  []bool // 1:1 list indicating which were selected

	options      []string // all the available options
	cursorIdx    int
	errorMessage string
}

func (s *SelectionUI) loadOptions(forceReload bool) {
	if s.OptionsFunc != nil {
		s.errorMessage = "updating..."
		s.printAll()
		data, err := s.OptionsFunc(forceReload)
		if err == nil {
			s.options = data
			s.setCursor(0)
			s.Selections = make([]bool, len(s.options))
			s.errorMessage = ""
		} else {
			s.errorMessage = fmt.Sprintf("Error while updating: %s", err)
		}
	}
}

func (s *SelectionUI) Run() error {
	if err := s.init(); err != nil {
		return err
	}
	s.printAll()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			s.errorMessage = ""
			if ev.Ch == 0 {
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC:
					s.resetSelections()
					break mainloop

				case termbox.KeyEnter:
					if isAnyTrue(s.Selections) {
						break mainloop
					} else {
						s.errorMessage = "Make a selection before continuing!"
					}

				case termbox.KeyArrowUp, termbox.KeyCtrlP:
					s.moveCursorUp()

				case termbox.KeyArrowDown, termbox.KeyCtrlN:
					s.moveCursorDown()

				case termbox.KeySpace:
					s.toggleSelectedUnderCursor()
				}

			} else {
				switch ev.Ch {
				case 'k', 'p':
					s.moveCursorUp()

				case 'j', 'n':
					s.moveCursorDown()

				case 'g':
					s.moveCursorToStart()

				case 'G':
					s.moveCursorToEnd()

				case 'r':
					s.loadOptions(true)

				case 'q':
					s.resetSelections()
					break mainloop
				}
			}

		case termbox.EventError:
			s.deinit()
			return ev.Err
		}

		s.printAll()
	}

	s.deinit()

	return nil
}

func (s *SelectionUI) init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc)

	s.loadOptions(false)
	return nil
}

func (s *SelectionUI) deinit() {
	termbox.Close()
}

func (s *SelectionUI) setCursor(rowIdx int) {
	s.cursorIdx = rowIdx
}

func (s *SelectionUI) moveCursorToStart() {
	s.setCursor(0)
}

func (s *SelectionUI) moveCursorToEnd() {
	s.setCursor(len(s.options) - 1)
}

func (s *SelectionUI) moveCursorUp() {
	if s.cursorIdx > 0 {
		s.setCursor(s.cursorIdx - 1)
	}
}

func (s *SelectionUI) moveCursorDown() {
	if s.cursorIdx < len(s.options)-1 {
		s.setCursor(s.cursorIdx + 1)
	}
}

func (s *SelectionUI) toggleSelectedUnderCursor() {
	s.toggleSelected(s.cursorIdx)
}

func (s *SelectionUI) toggleSelected(row int) {
	s.Selections[row] = !s.Selections[row]
}

func (s *SelectionUI) resetSelections() {
	for i := 0; i < len(s.Selections); i++ {
		s.Selections[i] = false
	}
}

func (s *SelectionUI) printHeader(x int, y int) (newY int) {
	newY = printText(x, y, headerMsg, termbox.ColorCyan, defaultBg)
	return
}

func (s *SelectionUI) printFooter(x, y int) (newY int) {
	newY = printText(x, y, footerMsg, termbox.ColorCyan, defaultBg)
	return
}

func (s *SelectionUI) printAll() {
	termbox.Clear(defaultFg, defaultBg)
	termbox.HideCursor()

	var y int
	y = s.printHeader(2, 1)
	y += 1
	y = s.printOptions(2, y)
	y += 1
	y = s.printFooter(2, y)
	y += 1

	if s.errorMessage != "" {
		y = printError(2, y, s.errorMessage)
	}

	termbox.Flush()
}

func (s *SelectionUI) printOptions(x, originY int) (newY int) {
	var bgColor, fgColor termbox.Attribute
	var y int = originY

	for i, option := range s.options {
		if s.Selections[i] {
			bgColor = defaultBg
			fgColor = termbox.ColorGreen | termbox.AttrBold
		} else {
			bgColor = defaultBg
			fgColor = defaultFg
		}

		if i == s.cursorIdx {
			termbox.SetCell(x, y, '➜', defaultFg, defaultBg)
		}

		y = printText(x+2, y, option, fgColor, bgColor)
	}

	return y
}

func printText(originX, originY int, text string, fg, bg termbox.Attribute) int {
	var x, y int = originX, originY
	var width, _ = termbox.Size()
	wrapWidth := width - originX - 3

	for i, r := range text {
		y = originY + i/wrapWidth
		x = originX + i%wrapWidth
		if y > originY {
			x += 2
		}
		termbox.SetCell(x, y, r, fg, bg)
	}

	return y + 1
}

func printError(x, y int, msg string) int {
	return printText(x, y, msg, termbox.ColorRed, defaultBg)
}

func isAnyTrue(a []bool) bool {
	for _, b := range a {
		if b {
			return true
		}
	}
	return false
}
