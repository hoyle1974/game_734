package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {

	g := NewGame(160, 40)

	fmt.Printf("Solid block: %s\n", "\u2588")

	NewPlasma(g.display, g)

	// for y := 0; y < len(image); y++ {
	// 	g.display.buffer.WriteString(0, y, image[y])
	// }

	g.logger.Log("System starting . . ." + "\u2588")

	// go func() {
	// 	ticker := time.NewTicker(1 * time.Second) // Create a ticker that ticks every second
	// 	defer ticker.Stop()                       // Ensure the ticker is stopped when done

	// 	for range ticker.C {
	// 		g.logger.Warn(fmt.Sprintf("Fired at: %v", time.Now().Format("15:04:05")))
	// 	}
	// }()

	if _, err := g.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
