package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	empty      = 0
	tree       = 1
	burnedTree = 2
)

type Coordinates struct {
	x, y int
}

type Forest struct {
	grid [][]int
}

func main() {
	forestSize := 100 //best size for the forest
	density := 0.4

	fmt.Println("Welcome to the Burning Forest simulation!")

	forest := CreateForest(forestSize, density)

	fmt.Println("What a peaceful day we have in here:")
	VisualizeForestInConsole(forest)
	fmt.Println("Image is being saved...")
	VisualizeForestOnImage(forest, 100, "data/forest-before.png")

	lightning_x := rand.Intn(forestSize)
	lightning_y := rand.Intn(forestSize)
	forest.SpreadFireWithWind(lightning_x, lightning_y, "W")
	fmt.Printf("Not for long! The lightning strikes at (%d, %d)!\n", lightning_x, lightning_y)

	fmt.Println("The forest after burning:")
	VisualizeForestInConsole(forest)
	fmt.Println("Image is being saved...")
	VisualizeForestOnImage(forest, 100, "data/forest-after.png")

	// percentage of burned trees
	burnedPercentage := calculatePercentageOfBurntTrees(forest, forestSize)
	fmt.Printf("%.2f%% of the forest got burned.\n", burnedPercentage)

	// optimal density
	fmt.Println()
	numberOfSamples := 1000
	// caution - it takes a while to calculate the optimal density
	optimalDensityForDefault := findOptimalDensity(forestSize, numberOfSamples)
	optimalDensityForNorth := findOptimalDensityForNorth(forestSize, numberOfSamples)
	optimalDensityForEast := findOptimalDensityForEast(forestSize, numberOfSamples)
	optimalDensityForWest := findOptimalDensityForWest(forestSize, numberOfSamples)
	optimalDensityForSouth := findOptimalDensityForSouth(forestSize, numberOfSamples)

	fmt.Printf("Optimal density for fire without any wind is %.2f%%.\n", optimalDensityForDefault*100)
	fmt.Printf("Optimal density for fire during nothern wind is %.2f%%.\n", optimalDensityForNorth*100)
	fmt.Printf("Optimal density for fire during eastern wind is %.2f%%.\n", optimalDensityForEast*100)
	fmt.Printf("Optimal density for fire during western wind is %.2f%%.\n", optimalDensityForWest*100)
	fmt.Printf("Optimal density for fire during southern wind is %.2f%%.\n", optimalDensityForSouth*100)

	optimalDensities := map[string]float64{
		"No Wind": optimalDensityForDefault * 100,
		"North":   optimalDensityForNorth * 100,
		"East":    optimalDensityForEast * 100,
		"West":    optimalDensityForWest * 100,
		"South":   optimalDensityForSouth * 100,
	}

	createBarChart(optimalDensities, "data/bestDens.png")
}

func createBarChart(optimalDensities map[string]float64, filename string) {
	p := plot.New()

	p.Title.Text = "Optimal Density for Different Wind Directions"
	p.Y.Label.Text = "Density (%)"

	var values plotter.Values
	var labels []string
	for direction, density := range optimalDensities {
		values = append(values, density)
		labels = append(labels, direction)
	}

	bars, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		panic(err)
	}

	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(1)

	p.Add(bars)
	p.NominalX(labels...)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
}

// creating square forect with given size and density
func CreateForest(size int, density float64) *Forest {
	forest := make([][]int, size)
	for i := 0; i < size; i++ {
		forest[i] = make([]int, size)
		for j := 0; j < size; j++ {
			if rand.Float64() < density {
				forest[i][j] = 1
			}
		}
	}
	return &Forest{
		grid: forest,
	}
}

// check if lightning struck specificly on the tree
func (f *Forest) StrikeLightning(x, y int) bool {
	return f.grid[x][y] == tree
}

// traditional spreading fire implemented with queue
func (f *Forest) SpreadFire(lightning_x, lightning_y int) {
	size := len(f.grid)
	treeQueue := []Coordinates{{lightning_x, lightning_y}}

	for len(treeQueue) > 0 {
		current := treeQueue[0]
		treeQueue = treeQueue[1:]

		x, y := current.x, current.y

		if x >= 0 && x < size && y >= 0 && y < size && f.grid[x][y] == tree {
			f.grid[x][y] = burnedTree

			neighbors := [][2]int{
				{x - 1, y}, {x + 1, y}, {x, y - 1}, {x, y + 1}, // up, down, left, right
				{x - 1, y - 1}, {x - 1, y + 1}, {x + 1, y - 1}, {x + 1, y + 1}, // diagonally
			}

			for _, neighbor := range neighbors {
				nx, ny := neighbor[0], neighbor[1]
				if nx >= 0 && nx < size && ny >= 0 && ny < size && f.grid[nx][ny] == tree {
					treeQueue = append(treeQueue, Coordinates{nx, ny})
				}
			}
		}
	}
}

func (f *Forest) getNeighbors(x, y int, windDirection string) []Coordinates {
	switch windDirection {
	case "N":
		return []Coordinates{{x + 1, y - 1}, {x + 1, y}, {x + 1, y + 1}} // Nothern wind - Boreas
	case "E":
		return []Coordinates{{x - 1, y - 1}, {x, y - 1}, {x + 1, y - 1}} // Eastern wind - Euros
	case "W":
		return []Coordinates{{x - 1, y + 1}, {x, y + 1}, {x + 1, y + 1}} // Western wind - Zephyros
	case "S":
		return []Coordinates{{x - 1, y - 1}, {x - 1, y}, {x - 1, y + 1}} // Southern wind - Notos
	default:
		return []Coordinates{{x - 1, y + 1}, {x, y + 1}, {x + 1, y + 1}} // default direction - Zephyros seemed the coolest one so he's in charge!
	}
}

func (f *Forest) SpreadFireWithWind(lightning_x, lightning_y int, windDirection string) {
	size := len(f.grid)
	treeQueue := []Coordinates{{lightning_x, lightning_y}}

	for len(treeQueue) > 0 {
		current := treeQueue[0]
		treeQueue = treeQueue[1:]

		x, y := current.x, current.y

		if x >= 0 && x < size && y >= 0 && y < size && f.grid[x][y] == tree {
			f.grid[x][y] = burnedTree

			neighbors := f.getNeighbors(x, y, windDirection)

			for _, neighbor := range neighbors {
				nx, ny := neighbor.x, neighbor.y
				if nx >= 0 && nx < size && ny >= 0 && ny < size && f.grid[nx][ny] == tree {
					treeQueue = append(treeQueue, Coordinates{nx, ny})
				}
			}
		}
	}
}

// the percentage of burned forest
func calculatePercentageOfBurntTrees(forest *Forest, size int) float64 {
	numOfAllTrees := 0.0
	numOfBurntTrees := 0.0
	for i := 0; i < size; i++ { // we ran through the whole forest
		for j := 0; j < size; j++ {
			if forest.grid[i][j] == tree { // and count the trees that survived
				numOfAllTrees++
			}
			if forest.grid[i][j] == burnedTree { // and the ones who burnt
				numOfAllTrees++
				numOfBurntTrees++
			}
		}
	}
	return (numOfBurntTrees / numOfAllTrees) * 100
}

// remaining trees
func counttreesLeft(forest Forest, size int) int {
	numOfTreesLeft := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if forest.grid[i][j] == tree {
				numOfTreesLeft++
			}
		}
	}
	return numOfTreesLeft
}

// find optimal density of trees (for fire without wind)
func findOptimalDensity(size, numberOfSamplesForDensity int) float64 {
	treesLeft := 1
	bestDensity := 1.0

	for density := 1.00; density > 0; density -= 0.01 {
		sumOfTreesLeft := 0

		for i := 0; i < numberOfSamplesForDensity; i++ {
			forest := CreateForest(size, density)
			lightning_x, lightning_y := rand.Intn(size), rand.Intn(size)
			if forest.StrikeLightning(lightning_x, lightning_y) {
				forest.SpreadFire(lightning_x, lightning_y)
				sumOfTreesLeft += counttreesLeft(*forest, size)
			}
		}

		averageResult := sumOfTreesLeft / numberOfSamplesForDensity
		if averageResult > treesLeft {
			treesLeft = averageResult
			bestDensity = density
		}
	}

	return bestDensity
}

func findOptimalDensityForNorth(size, numberOfSamplesForDensity int) float64 {
	treesLeft := 1
	bestDensity := 1.0

	for density := 1.00; density > 0; density -= 0.01 {
		sumOfTreesLeft := 0

		for i := 0; i < numberOfSamplesForDensity; i++ {
			forest := CreateForest(size, density)
			lightning_x, lightning_y := rand.Intn(size), rand.Intn(size)
			if forest.StrikeLightning(lightning_x, lightning_y) {
				forest.SpreadFireWithWind(lightning_x, lightning_y, "N")
				sumOfTreesLeft += counttreesLeft(*forest, size)
			}
		}

		averageResult := sumOfTreesLeft / numberOfSamplesForDensity
		if averageResult > treesLeft {
			treesLeft = averageResult
			bestDensity = density
		}
	}

	return bestDensity
}

func findOptimalDensityForEast(size, numberOfSamplesForDensity int) float64 {
	treesLeft := 1
	bestDensity := 1.0

	for density := 1.00; density > 0; density -= 0.01 {
		sumOfTreesLeft := 0

		for i := 0; i < numberOfSamplesForDensity; i++ {
			forest := CreateForest(size, density)
			lightning_x, lightning_y := rand.Intn(size), rand.Intn(size)
			if forest.StrikeLightning(lightning_x, lightning_y) {
				forest.SpreadFireWithWind(lightning_x, lightning_y, "E")
				sumOfTreesLeft += counttreesLeft(*forest, size)
			}
		}

		averageResult := sumOfTreesLeft / numberOfSamplesForDensity
		if averageResult > treesLeft {
			treesLeft = averageResult
			bestDensity = density
		}
	}

	return bestDensity
}

func findOptimalDensityForWest(size, numberOfSamplesForDensity int) float64 {
	treesLeft := 1
	bestDensity := 1.0

	for density := 1.00; density > 0; density -= 0.01 {
		sumOfTreesLeft := 0

		for i := 0; i < numberOfSamplesForDensity; i++ {
			forest := CreateForest(size, density)
			lightning_x, lightning_y := rand.Intn(size), rand.Intn(size)
			if forest.StrikeLightning(lightning_x, lightning_y) {
				forest.SpreadFireWithWind(lightning_x, lightning_y, "W")
				sumOfTreesLeft += counttreesLeft(*forest, size)
			}
		}

		averageResult := sumOfTreesLeft / numberOfSamplesForDensity
		if averageResult > treesLeft {
			treesLeft = averageResult
			bestDensity = density
		}
	}

	return bestDensity
}

func findOptimalDensityForSouth(size, numberOfSamplesForDensity int) float64 {
	treesLeft := 1
	bestDensity := 1.0

	for density := 1.00; density > 0; density -= 0.01 {
		sumOfTreesLeft := 0

		for i := 0; i < numberOfSamplesForDensity; i++ {
			forest := CreateForest(size, density)
			lightning_x, lightning_y := rand.Intn(size), rand.Intn(size)
			if forest.StrikeLightning(lightning_x, lightning_y) {
				forest.SpreadFireWithWind(lightning_x, lightning_y, "S")
				sumOfTreesLeft += counttreesLeft(*forest, size)
			}
		}

		averageResult := sumOfTreesLeft / numberOfSamplesForDensity
		if averageResult > treesLeft {
			treesLeft = averageResult
			bestDensity = density
		}
	}

	return bestDensity
}

// "drawing" forest in console
func VisualizeForestInConsole(forest *Forest) {
	for _, row := range forest.grid {
		for _, cell := range row {
			switch cell {
			case empty:
				fmt.Print(" . ")
			case tree:
				fmt.Print(" T ")
			case burnedTree:
				fmt.Print(" X ")
			}
		}
		fmt.Println()
	}
}

// creating image out of the forest grid
func VisualizeForestOnImage(forest *Forest, scale int, filename string) {
	size := len(forest.grid) * scale
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// fill the image with ground-like color
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			img.Set(x, y, color.RGBA{120, 100, 60, 255})
		}
	}

	for i, row := range forest.grid {
		for j, cell := range row {
			c := color.RGBA{}
			switch cell {
			case empty:
				c = color.RGBA{120, 100, 60, 255} // light brown - empty field
			case tree:
				c = color.RGBA{45, 145, 2, 255} // green - tree
			case burnedTree:
				c = color.RGBA{255, 20, 63, 255} // dark red - burned tree
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
