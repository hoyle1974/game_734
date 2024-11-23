package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// faint := color.New(color.Faint).SprintFunc()
	// buf := NewBuffer(10, 10)
	// buf.DrawBox(0, 0, 9, 9)
	// buf.WriteString(1, 1, "🌍")
	// buf.WriteString(2, 2, "🌍 Hi🌍")
	// buf.WriteString(3, 3, faint("Hi")+"Hi🌍")
	// fmt.Println(buf.String())
	// os.Exit(-1)

	g := NewGame(160, 40)

	g.logger.Log("System starting . . .")

	go func() {
		ticker := time.NewTicker(1 * time.Second) // Create a ticker that ticks every second
		defer ticker.Stop()                       // Ensure the ticker is stopped when done

		for range ticker.C {
			g.logger.Warn(fmt.Sprintf("Fired at: %v", time.Now().Format("15:04:05")))
		}
	}()

	if _, err := g.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
