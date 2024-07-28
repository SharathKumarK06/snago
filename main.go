package main

import (
    // "fmt"
    "log"

    "github.com/gdamore/tcell/v2"
)



func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
    x, y := x1, y1
    for _, c := range []rune(text) {
        s.SetContent(x, y, c, nil, style)
        x++

        if x > x2 {
            y++
            x = x1
        }
        if y > y2 {
            break
        }
    }
}

func fillBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
    if x1 > x2 {
        x1, x2 = x2, x1
    }
    if y1 > y2 {
        y1, y2 = y2, y1
    }

    for y := y1; y <= y2; y++ {
        for x := x1; x <= x2; x++ {
            s.SetContent(x, y, ' ', nil, style)
        }
    }
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
    if x1 > x2 {
        x1, x2 = x2, x1
    }
    if y1 > y2 {
        y1, y2 = y2, y1
    }

    // Draw borders
    for x := x1; x <= x2; x++ {
        s.SetContent(x, y1, ' ', nil, style)
        s.SetContent(x, y2, ' ', nil, style)
    }
    for y := y1 + 1; y < y2; y++ {
        s.SetContent(x1, y, ' ', nil, style)
        s.SetContent(x2, y, ' ', nil, style)
    }

    drawText(s, x1 + 1, y1 + 1, x2 - 1, y2 - 1, style, text)
}

func main() {
    defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
    textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

    // initialize screen
    s, err := tcell.NewScreen()
    if err != nil {
        log.Fatalf("%+v", err)
    }
    if err := s.Init(); err != nil {
        log.Fatalf("%+v", err)

    }

    // Set default text style
    s.SetStyle(defStyle)
    s.EnableMouse()
    s.EnablePaste()

    // Clear screen
    s.Clear()


    quit := func() {
        // You have to catch panics in a defer, clean up, and
        // re-raise them - otherwise your application can
        // die without leaving any diagnostic trace.
        maybePanic := recover()
        s.Fini()
        if maybePanic != nil {
            panic(maybePanic)
        }
    }
    defer quit()

    // Event loop
    for {
        // Update screen
        s.Show()

        // s.SetContent(5, 5, 'H', nil, textStyle)
        xmax, ymax := s.Size()

        // Draw borders
        drawBox(s, 0, 0, xmax - 1, ymax - 1, textStyle, "")


        // Poll event
        ev := s. PollEvent()

        switch ev := ev.(type) {
            case *tcell.EventResize:
                s.Sync()
            case *tcell.EventKey:
                switch {
                    case ev.Key() == tcell.KeyEscape:
                        return
                    case ev.Key() == tcell.KeyCtrlC:
                        return
                    case ev.Rune() == 'Q' || ev.Rune() == 'q':
                        return
                    case ev.Rune() == 'C' || ev.Rune() == 'c':
                        s.Clear()
                }
            case *tcell.EventMouse:
                // x, y := ev.Position()

                switch ev.Buttons() {
                    // case tcell.Button1, tcell.Button2:
                    // case tcell.ButtonNone:

                }

        }
    }
}

