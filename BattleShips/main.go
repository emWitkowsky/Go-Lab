package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

// Ships game

// thats the part where they will be changing things
func main() {
	ShipsSize := 100
	// density := 0.6

	// Ships := CreateShips(ShipsSize, density)
	Ships := CreateShips(ShipsSize)

	fmt.Println("What a peaceful day we have in here:")
	VisualizeShipsInConsole(Ships)
	VisualizeShipsInImage(Ships, 100, "data/ships-before.png")

	shot_x, shot_y := rand.Intn(ShipsSize), rand.Intn(ShipsSize)
	Ships.SpreadFireWithWind(shot_x, shot_y, "W") // ????
	fmt.Printf("Not for long! The shot has been fired at (%d, %d)!\n", shot_x, shot_y)

	fmt.Println()
	fmt.Println("The Ships after target shooting:")
	VisualizeShipsInConsole(Ships)
	VisualizeShipsInImage(Ships, 100, "data/ships-after.png")

	// percentage of burned trees //????
	// fmt.Println()
	// burnedPercentage := calculatePercentageOfshipHits(Ships, ShipsSize)
	// fmt.Printf("%.2f%% of the Ships got burned.\n", burnedPercentage)

	// optimal density
	// fmt.Println()
	// numberOfSamples := 1000
	// optimalDensity := findOptimalDensity(ShipsSize, numberOfSamples)
	// fmt.Printf("Optimal density for 'traditional' fire is %.2f%%.\n", optimalDensity*100)
}

const (
	sea        = 0
	ship       = 1
	shipHit = 2
)

type Coordinates struct {
	x, y int
}

type Ships struct {
	grid [][]int
}

// creating square map with given size and 3 3x1 ships
func CreateShips(size int) *Ships {
	sea := make([][]int, size)
	for i := 0; i < size; i++ {
		sea[i] = make([]int, size)
	}

	// Place 3 3x1 ships
	for i := 0; i < 3; i++ {
		shipStartX := rand.Intn(size - 2) // Ensure ship fits in the map
		shipStartY := rand.Intn(size)

		// Place ship
		for j := 0; j < 3; j++ {
			sea[shipStartX+j][shipStartY] = 1
		}
	}

	return &Ships{
		grid: sea,
	}
}

// check if lightning struck on the tree
func (f *Ships) StrikeShip(x, y int) bool {
	return f.grid[x][y] == ship
}

// traditional spreading fire implemented with queue
func (f *Ships) SpreadFire(lightning_x, lightning_y int) {
	size := len(f.grid)
	treeQueue := []Coordinates{{lightning_x, lightning_y}}

	for len(treeQueue) > 0 {
		current := treeQueue[0]
		treeQueue = treeQueue[1:]

		x, y := current.x, current.y

		if x >= 0 && x < size && y >= 0 && y < size && f.grid[x][y] == ship {
			f.grid[x][y] = shipHit

			neighbors := [][2]int{
				{x - 1, y}, {x + 1, y}, {x, y - 1}, {x, y + 1}, // up, down, left, right
				{x - 1, y - 1}, {x - 1, y + 1}, {x + 1, y - 1}, {x + 1, y + 1}, // diagonally
			}

			for _, neighbor := range neighbors {
				nx, ny := neighbor[0], neighbor[1]
				if nx >= 0 && nx < size && ny >= 0 && ny < size && f.grid[nx][ny] == ship {
					treeQueue = append(treeQueue, Coordinates{nx, ny})
				}
			}
		}
	}
}

// spreading fire with wind, implemented with queue
func (f *Ships) SpreadFireWithWind(lightning_x, lightning_y int, windDirection string) {
	size := len(f.grid)
	treeQueue := []Coordinates{{lightning_x, lightning_y}}

	for len(treeQueue) > 0 {
		current := treeQueue[0]
		treeQueue = treeQueue[1:]

		x, y := current.x, current.y

		if x >= 0 && x < size && y >= 0 && y < size && f.grid[x][y] == ship {
			f.grid[x][y] = shipHit

			var neighbors []Coordinates

			switch windDirection {
			case "N":
				neighbors = []Coordinates{{x + 1, y - 1}, {x + 1, y}, {x + 1, y + 1}} // Nothern wind - Boreas
			case "E":
				neighbors = []Coordinates{{x - 1, y - 1}, {x, y - 1}, {x + 1, y - 1}} // Eastern wind - Euros
			case "W":
				neighbors = []Coordinates{{x - 1, y + 1}, {x, y + 1}, {x + 1, y + 1}} // Western wind - Zephyros
			case "S":
				neighbors = []Coordinates{{x - 1, y - 1}, {x - 1, y}, {x - 1, y + 1}} // Southern wind - Notos
			default:
				windDirection = "W" // default direction - Zephyros seemed the coolest one so he's in charge!
			}

			for _, neighbor := range neighbors {
				nx, ny := neighbor.x, neighbor.y
				if nx >= 0 && nx < size && ny >= 0 && ny < size && f.grid[nx][ny] == ship {
					treeQueue = append(treeQueue, Coordinates{nx, ny})
				}
			}
		}
	}
}

// the percentage of burned Ships
func calculatePercentageOfshipHits(Ships *Ships, size int) float64 {
	numberOfAllShips := 0.0
	numberOfshipHits := 0.0
	for i := 0; i < size; i++ { // we ran through the whole Ships
		for j := 0; j < size; j++ {
			if Ships.grid[i][j] == ship { // and count the trees that survived
				numberOfAllShips++
			}
			if Ships.grid[i][j] == shipHit { // and the ones who burnt
				numberOfAllShips++
				numberOfshipHits++
			}
		}
	}
	return (numberOfshipHits / numberOfAllShips) * 100
}

// remaining trees
func countRemainingShips(Ships Ships, size int) int {
	numberOfHealthyShips := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if Ships.grid[i][j] == ship {
				numberOfHealthyShips++
			}
		}
	}
	return numberOfHealthyShips
}

// find optimal density of trees (with traditional fire spreading)
func findOptimalDensity(size, numberOfSamplesForDensity int) float64 {
	remainingTrees := 1
	bestDensity := 1.0

	for density := 1.00; density > 0; density -= 0.01 {
		// fmt.Printf("Checking for %.2f... \n", density)
		sumOfRemainingTrees := 0

		for i := 0; i < numberOfSamplesForDensity; i++ {
			// Ships := CreateShips(size, density)
			Ships := CreateShips(size)
			lightning_x, lightning_y := rand.Intn(size), rand.Intn(size)
			if Ships.StrikeShip(lightning_x, lightning_y) {
				Ships.SpreadFire(lightning_x, lightning_y)
				sumOfRemainingTrees += countRemainingShips(*Ships, size)
			}
		}

		averageResult := sumOfRemainingTrees / numberOfSamplesForDensity
		if averageResult > remainingTrees {
			remainingTrees = averageResult
			bestDensity = density
		}
	}

	return bestDensity
}

// Simple illustration of the Ships grid
func VisualizeShipsInConsole(Ships *Ships) {
	for _, row := range Ships.grid {
		for _, cell := range row {
			switch cell {
			case sea:
				fmt.Print(" ~ ")
			case ship:
				fmt.Print(" S ")
			case shipHit:
				fmt.Print(" X ")
			}
		}
		fmt.Println()
	}
}

// Visualising very same
func VisualizeShipsInImage(Ships *Ships, scale int, filename string) {
	size := len(Ships.grid) * scale
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// fill the image with sea color blue
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	for i, row := range Ships.grid {
		for j, cell := range row {
			c := color.RGBA{}
			switch cell {
			case sea:
				c = color.RGBA{0, 0, 225, 255} // blue - sea
			case ship:
				c = color.RGBA{211, 211, 211, 255} // grey - ship
			case shipHit:
				c = color.RGBA{255, 20, 63, 255} // dark red - ship hit
			}

			for x := i * scale; x < (i+1)*scale; x++ {
				for y := j * scale; y < (j+1)*scale; y++ {
					img.Set(y, x, c)
				}
			}
		}
	}

	f, _ := os.Create(filename)
	png.Encode(f, img)
	f.Close()
}
