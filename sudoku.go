package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const size = 9

var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

type Board [size][size]int

func checkAlgorithm(board *Board, row, col, val int) bool {
	for i := 0; i < size; i++ {
		if board[row][i] == val || board[i][col] == val {
			return false
		}
	}
	startRow := (row / 3) * 3
	startCol := (col / 3) * 3
	for i := startRow; i < startRow+3; i++ {
		for j := startCol; j < startCol+3; j++ {
			if board[i][j] == val {
				return false
			}
		}
	}
	return true
}

func fillBoard(board *Board, row, col int) bool {
	if row == size {
		return true
	}
	if col == size {
		return fillBoard(board, row+1, 0)
	}
	if board[row][col] != 0 {
		return fillBoard(board, row, col+1)
	}

	candidates := make([]int, len(nums))
	copy(candidates, nums)
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	for v := 0; v < len(candidates); v++ {
		if checkAlgorithm(board, row, col, candidates[v]) {
			board[row][col] = candidates[v]
			if fillBoard(board, row, col+1) {
				return true
			}
			board[row][col] = 0
		}
	}
	return false
}

func hasUniqueSolution(puzzle Board) bool {
	var count int
	board := puzzle
	var solve func(row, col int) bool
	solve = func(row, col int) bool {
		if count > 1 {
			return true
		}
		if row == size {
			count++
			return count > 1
		}
		if col == size {
			return solve(row+1, 0)
		}
		if board[row][col] != 0 {
			return solve(row, col+1)
		}
		for checkValue := 0; checkValue < len(nums); checkValue++ {
			if checkAlgorithm(&board, row, col, nums[checkValue]) {
				board[row][col] = nums[checkValue]
				solve(row, col+1)
				board[row][col] = 0
				if count > 1 {
					return true
				}
			}
		}
		return false
	}
	count = 0
	solve(0, 0)
	return count == 1
}

func generateSudoku(amoutExist int) Board {
	for attempt := 0; attempt < 2000; attempt++ {
		var board Board
		if !fillBoard(&board, 0, 0) {
			continue
		}
		puzzle := board
		removeCount := size*size - amoutExist
		positions := make([][2]int, 0, size*size)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				positions = append(positions, [2]int{i, j})
			}
		}
		rand.Shuffle(len(positions), func(i, j int) {
			positions[i], positions[j] = positions[j], positions[i]
		})
		for i := 0; i < removeCount; i++ {
			randomRemoveROW, randomRemoveCOL := positions[i][0], positions[i][1]
			puzzle[randomRemoveROW][randomRemoveCOL] = 0
		}
		if hasUniqueSolution(puzzle) {
			return puzzle
		}
	}
	panic("Không thể tạo Sudoku hợp lệ sau nhiều lần thử!")
}

func writeJSON(w http.ResponseWriter, board Board) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func easyHandler(w http.ResponseWriter, r *http.Request) {
	board := generateSudoku(50)
	writeJSON(w, board)
}

func normalHandler(w http.ResponseWriter, r *http.Request) {
	board := generateSudoku(36)
	writeJSON(w, board)
}

func hardHandler(w http.ResponseWriter, r *http.Request) {
	board := generateSudoku(28)
	writeJSON(w, board)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/api/sudoku/easy", enableCORS(easyHandler))
	http.HandleFunc("/api/sudoku/normal", enableCORS(normalHandler))
	http.HandleFunc("/api/sudoku/hard", enableCORS(hardHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port", port)

	http.ListenAndServe(":"+port, nil)
}

/*
func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/api/sudoku/easy", easyHandler)
	http.HandleFunc("/api/sudoku/normal", normalHandler)
	http.HandleFunc("/api/sudoku/hard", hardHandler)
	fs := http.FileServer(http.Dir("D:/MyFrontend/Sudoku"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "D:/MyFrontend/Sudoku/index.html")
	})
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
*/
