package shared

import "fmt"

type PlayerVarHandlerFunc func(g *Game, newVal uint16)

func HandleShowTutorial(g *Game, val uint16) {
	fmt.Printf("handling tut show %d\n", val)
	if val == 0 {
		g.RenderExclamOnGuide = true
	} else if val == 1 {
		g.RenderExclamOnGuide = false
	}
}
