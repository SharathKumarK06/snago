package main

import (
    // "fmt"
    "log"
    "time"

    "github.com/gdamore/tcell/v2"
)

func run(screen tcell.Screen) {
    // Set default text style
    defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
    screen.SetStyle(defStyle)

    // Initialization of game
    xmax, ymax := screen.Size()
    snake := ':'
    pos_x := int(xmax / 2)
    pos_y := int(ymax / 2)

    // snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
    snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorRed)

    for {
        // Clear screen
        screen.Clear()

        screen.SetContent(pos_x, pos_y, snake, nil, snakeStyle)

        pos_x++

        // Update screen
        screen.Show()

        time.Sleep(80 * time.Millisecond)
    }

}

func main() {
    // initialize screen
    s, err := tcell.NewScreen()
    if err != nil {
        log.Fatalf("%+v", err)
    }
    if err := s.Init(); err != nil {
        log.Fatalf("%+v", err)

    }

    s.EnableMouse()
    s.EnablePaste()

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

    go run(s)

    // Event loop
    for {
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

