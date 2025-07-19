package main

type Point struct {
	I, J int // I: строка, J: столбец
}

func getAdjacent(i, j int) []Point {
	var adj []Point
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			if di == 0 && dj == 0 {
				continue
			}
			ni, nj := i+di, j+dj
			if ni >= 0 && ni < Rows && nj >= 0 && nj < Cols {
				adj = append(adj, Point{ni, nj})
			}
		}
	}
	return adj
}

func getClose(point Point) int8 {
	var quantity int8
	adj := getAdjacent(point.I, point.J)
	for i := 0; i < len(adj); i++ {
		if field[adj[i].I][adj[i].J] == -1 {
			quantity++
		}
	}
	return quantity
}

func getCloseBomb(point Point) int8 {
	var quantity int8
	adj := getAdjacent(point.I, point.J)
	for i := 0; i < len(adj); i++ {
		if field[adj[i].I][adj[i].J] == 10 {
			quantity++
		}
	}
	return quantity
}

func findSafeCells() []Point {
	safeMap := make(map[Point]bool)
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			if field[i][j] > 0 && field[i][j] < 7 { // N от 1 до 6
				N := int(field[i][j])
				adj := getAdjacent(i, j)
				F := 0
				var U []Point
				for _, a := range adj {
					if field[a.I][a.J] == mine { // mine определено как 10
						F++
					} else if field[a.I][a.J] == -1 {
						U = append(U, a)
					}
				}
				if F == N {
					for _, u := range U {
						safeMap[u] = true
					}
				}
			}
		}
	}
	var safeList []Point
	for p := range safeMap {
		safeList = append(safeList, p)
	}
	return safeList
}

func findBomb() {
	for i := 0; i < Rows; i++ {
		for j := 0; j < Cols; j++ {
			cls := getClose(Point{i, j})
			if cls != 0 && cls == field[i][j]-getCloseBomb(Point{i, j}) {
				for _, point := range getAdjacent(i, j) {
					if field[point.I][point.J] == -1 {
						field[point.I][point.J] = mine
					}
				}
			}
		}
	}
}
