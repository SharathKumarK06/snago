package main

import (
    "fmt"
    "log"
    "time"

    "github.com/gdamore/tcell/v2"
)

type direction int

const (
    up direction = iota
    right
    down
    left
)

type position struct {
    x, y int
}

type size struct {
    w, h int
}

type snakePart struct {
    pos position
    dir direction
    char string
    tail *snakePart
}

type snake struct {
    head *snakePart
    length int
    style tcell.Style
}

type boardState struct {
    screen tcell.Screen
    size size
    snakeUpdated int
}

// func updateDirection(snake []SnakePart, dir Direction) {
//     if dir + 2 == snake[0].dir || dir - 2 == snake[0].dir {
//         return
//     }

//     snake[0].dir = dir
// }

// func Run(screen tcell.Screen, snake []SnakePart) {
//     if screen == nil {
//         panic("Error: null screen")
//     }

//     if snake == nil {
//         panic("Error: null snake")
//     }

//     for {
//         // Clear screen
//         screen.Clear()

//         // Render snake
//         for i := 0; i < len(snake); i++ {
//             screen.SetContent(snake[i].pos.x, snake[i].pos.y, snake[i].char, nil, snakeStyle)
//         }

//         // Update snake position
//         for i := 0; i < len(snake); i++ {
//             switch (snake[i].dir) {
//                 case Up:
//                     snake[i].pos.y--
//                 case Down:
//                     snake[i].pos.y++
//                 case Right:
//                     snake[i].pos.x++
//                 case Left:
//                     snake[i].pos.x--
//                 default:
//                     panic("Error: unknown direction")
//             }
//         }

//         // Update screen
//         screen.Show()

//         time.Sleep(90 * time.Millisecond)
//     }

// }

func (board *boardState) printStats(snake *snake) {
    pos := position{2, board.size.h + 2}
    statusStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorReset).Attributes(tcell.AttrBold)
    text := fmt.Sprintf("Position(x, y): (%v, %v)  Direction: (%v)\n", snake.head.pos.x, snake.head.pos.y, snake.head.dir)

    var char rune
    for i := 0; i < len(text); i++ {
        char = rune(text[i])
        if char == rune('\n') {
            pos.y++
            pos.x = 0
        }
        board.screen.SetContent(pos.x + i, pos.y, char, nil, statusStyle)
    }
}

func (board *boardState) gameLoop(snake *snake) {
    for {
        // clearing the screen
        board.screen.Clear()

        borderStyle := tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorReset)
        board.drawBorder(borderStyle)

        // render snake
        n := snake.head
        for n != nil {
            for x := len(n.char) - 1; x >= 0; x-- {
                board.screen.SetContent(n.pos.x + x, n.pos.y, rune(n.char[x]), nil, snake.style)
            }
            n = n.tail
        }

        snake.head.pos.x++

        board.printStats(snake)

        // screen update interval
        time.Sleep(180 * time.Millisecond)

        // update screen
        board.screen.Show()
    }
}

func (board *boardState) eventLoop() {
    for {
        // Poll event
        ev := board.screen.PollEvent()

        switch ev := ev.(type) {
            case *tcell.EventResize:
                board.screen.Sync()
            case *tcell.EventKey:
                switch {
                    // quit
                    case ev.Key() == tcell.KeyEscape:
                        return
                    case ev.Key() == tcell.KeyCtrlC:
                        return
                    case ev.Rune() == 'Q' || ev.Rune() == 'q':
                        return

                    // clear screen
                    case ev.Rune() == 'C' || ev.Rune() == 'c':
                        board.screen.Clear()

                    // control keys
                    case ev.Rune() == 'W' || ev.Rune() == 'w':
                        // updateDirection(snake, Up)
                    case ev.Rune() == 'A' || ev.Rune() == 'a':
                        // updateDirection(snake, Left)
                    case ev.Rune() == 'S' || ev.Rune() == 's':
                        // updateDirection(snake, Down)
                    case ev.Rune() == 'D' || ev.Rune() == 'd':
                        // updateDirection(snake, Right)
                }
        }
    }
}

func createSnake(pos position, style tcell.Style) *snake {
    return &snake{
        head: &snakePart{
            pos: pos,
            dir: right,
            char: " :",
            tail: nil,
        },
        length: 1,
        style: style,
    }
}

func (board *boardState) drawBorder(style tcell.Style) {
    for y := 0; y <= board.size.h + 1; y++ {
        board.screen.SetContent(0, y, ' ', nil, style)
        board.screen.SetContent(board.size.w + 1, y, ' ', nil, style)
    }

    for x := 1; x <= board.size.w; x++ {
        board.screen.SetContent(x, 0, ' ', nil, style)
        board.screen.SetContent(x, board.size.h + 1, ' ', nil, style)
    }
}

func createBoard(screen tcell.Screen, size size) *boardState {
    board := boardState{
        screen: screen,
        size: size,
        snakeUpdated: -1,
    }

    return &board
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

    // creating board
    board := createBoard(s, size{80, 24})

    // board style
    boardStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
    board.screen.SetStyle(boardStyle)

    // hiding cursor
    board.screen.HideCursor()

    // creating snake
    snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorRed)

    snakePosition := position{
        x: int(board.size.w / 2),
        y: int(board.size.h / 2),
    }
    snake := createSnake(snakePosition, snakeStyle)

    // cleanup
    quit := func() {
        maybePanic := recover()

        board.screen.Fini()

        if maybePanic != nil {
            panic(maybePanic)
        }
    }

    defer quit()

    // game loop
    go board.gameLoop(snake)

    // event loop
    board.eventLoop()

    // snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
    // snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorRed)
}

