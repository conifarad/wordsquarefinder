package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	infoLog  = log.New(os.Stderr, "[info] ", 0)
	warnLog  = log.New(os.Stderr, "[warning] ", 0)
	errorLog = log.New(os.Stderr, "[error] ", 0)
)

func logf(format string, v ...any) {
	fmt.Fprintf(os.Stderr, format, v...)
}

type ProgState struct {
	square_size    int
	num_threads    int
	word_arr       []string
	word_tree      *WordTreeNode
	start_idx      int
	end_idx        int
	pretty_print   bool
	require_unique bool
	is_piped       bool
}

func all_unique_words(board []string) bool {
	squareSize := len(board)
	words := make([]string, squareSize, 2*squareSize)
	copy(words, board)

	for i, word := range words {
		if slices.Contains(words[:i], word) {
			return false
		}
	}

	vBytes := make([]byte, squareSize)
	for col := 0; col < squareSize; col++ {
		for row := 0; row < squareSize; row++ {
			vBytes[row] = board[row][col]
		}
		word := string(vBytes)
		if slices.Contains(words, word) {
			return false
		}
		words = append(words, word)
	}

	return true
}

func build_squares_recursive(
	squareSize int,
	requireUnique bool,
	wordTreeRoot *WordTreeNode,
	leftNode *WordTreeNode,
	aboveNodes []*WordTreeNode,
	currentBoard []byte,
	boardChan chan []string) {

	if len(currentBoard) == squareSize*squareSize {
		stringBoard := make([]string, 0, squareSize)
		for i := 0; i < squareSize; i++ {
			row := currentBoard[i*squareSize : (i+1)*squareSize]
			stringBoard = append(stringBoard, string(row))
		}
		if !requireUnique || all_unique_words(stringBoard) {
			boardChan <- stringBoard
		}
		return
	}

	aboveNode := aboveNodes[len(aboveNodes)-squareSize]

	for letter := byte('a'); letter <= byte('z'); letter++ {
		hNode := leftNode.get_child(letter)
		if hNode == nil {
			continue
		}
		vNode := aboveNode.get_child(letter)
		if vNode == nil {
			continue
		}

		nextAboveNodes := append(aboveNodes, vNode)
		nextBoard := append(currentBoard, letter)

		nextLeftNode := hNode
		if len(nextBoard)%squareSize == 0 {
			nextLeftNode = wordTreeRoot
		}

		build_squares_recursive(squareSize, requireUnique, wordTreeRoot, nextLeftNode,
			nextAboveNodes, nextBoard, boardChan)
	}
}

func build_squares_thread(
	state ProgState,
	boardChan chan []string,
	delegateChan chan int,
	doneChan chan int) {

	// Make a board with the appropriate capacity ahead of time to avoid memory allocations.
	numLetters := state.square_size * state.square_size
	board := make([]byte, state.square_size, numLetters)
	nodes := make([]*WordTreeNode, state.square_size, numLetters)

	for wordIdx := range delegateChan {
		word := state.word_arr[wordIdx]
		// Check if the word is valid
		success := true
		for i := 0; i < state.square_size; i++ {
			node := state.word_tree.get_child(word[i])
			if node == nil {
				success = false
				break
			}
			board[i] = word[i]
			nodes[i] = node
		}

		if success {
			// Search for squares
			build_squares_recursive(state.square_size, state.require_unique, state.word_tree,
				state.word_tree, nodes, board, boardChan)
		}

		// Notify main thread that we've finished
		doneChan <- wordIdx
	}
}

func build_squares(state ProgState) {
	nextIdx := state.start_idx
	activeIndices := make([]int, 0, state.num_threads)
	numSearched := 0
	numFound := 0
	needNewline := false

	log_status := func() {
		logf("\rSearching: %v... - %v... (%.2f%%, %v found)",
			state.word_arr[activeIndices[0]],
			state.word_arr[activeIndices[len(activeIndices)-1]],
			float32(numSearched)/float32(state.end_idx-state.start_idx+1)*100,
			numFound)
		needNewline = !state.is_piped
	}
	handle_square := func(board []string) {
		if needNewline {
			logf("\n")
			needNewline = false
		}
		if state.pretty_print {
			fmt.Print(strings.Join(board, "\n") + "\n\n")
		} else {
			fmt.Println(strings.Join(board, ","))
		}
		numFound++
	}
	handle_finished := func(idx int) {
		numSearched++
		i := slices.Index(activeIndices, idx)
		activeIndices = slices.Delete(activeIndices, i, i+1)
	}

	boardChan := make(chan []string)
	delegateChan := make(chan int)
	doneChan := make(chan int)

	for i := 0; i < state.num_threads; i++ {
		go build_squares_thread(state, boardChan, delegateChan, doneChan)
	}

	for nextIdx <= state.end_idx {
		select {
		case delegateChan <- nextIdx:
			// Store all words currently being processed
			activeIndices = append(activeIndices, nextIdx)
			log_status()
			nextIdx++
		case board := <-boardChan:
			handle_square(board)
		case doneIdx := <-doneChan:
			handle_finished(doneIdx)
		}
	}

	// Signal threads to stop
	close(delegateChan)

	// Print out remaining boards and wait for all threads to finish execution
	for len(activeIndices) > 0 {
		select {
		case board := <-boardChan:
			handle_square(board)
		case doneIdx := <-doneChan:
			handle_finished(doneIdx)
			if len(activeIndices) > 0 {
				log_status()
			}
		}
	}

	logf("\n")
	infoLog.Printf("Finished; found %v squares.", numFound)
}

func get_default_num_threads() int {
	num := runtime.NumCPU() / 2
	if num > 0 {
		return num
	}
	return 1 // e.g. if NumCPU() is 1
}

func read_words(path string, length int) []string {
	file, err := os.Open(path)
	if err != nil {
		errorLog.Fatalf("Failed to open input file: %v\n", err)
	}
	infoLog.Printf("Loading words from %q...\n", path)

	scanner := bufio.NewScanner(file)
	wordArr := []string{}
	numExcluded := 0

	is_invalid_character := func(r rune) bool {
		return r < 'a' || r > 'z'
	}

	for scanner.Scan() {
		word := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if len(word) == length {
			if strings.IndexFunc(word, is_invalid_character) == -1 {
				wordArr = append(wordArr, word)
			} else {
				numExcluded++
			}
		}
	}
	file.Close()

	if len(wordArr) == 0 {
		errorLog.Fatalln("No matching words found in word list.")
	}

	sort.Strings(wordArr)
	infoLog.Printf("Loaded %v %v-letter words, excluding %v invalid words.\n",
		len(wordArr), length, numExcluded)

	return wordArr
}

func word_arg_to_index(wordArr []string, arg string, dflt int) int {
	if arg == "" {
		return dflt
	}

	word := strings.ToLower(strings.TrimSpace(arg))
	index := sort.SearchStrings(wordArr, word)
	if index == len(wordArr) {
		index--
	}
	if wordArr[index] != word {
		warnLog.Printf("No word %q in word list, using %q.\n", word, wordArr[index])
	}
	return index
}

func stdout_is_piped() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeNamedPipe != 0
}

func main() {
	squareSize := flag.Int("n", 5, "Size of word squares")
	numThreads := flag.Int("t", get_default_num_threads(), "Number of threads to run")
	inputPath := flag.String("i", "english_filtered.txt", "Path to word list")
	startWord := flag.String("s", "", "Word to start search at")
	endWord := flag.String("e", "", "Word to end search at")
	prettyPrint := flag.Bool("p", false, "Pretty-print word squares (outputs one word per line).")
	requireUnique := flag.Bool("u", false,
		"Require unique words (no two rows/columns with the same word).")
	flag.Parse()

	if *squareSize < 2 {
		errorLog.Fatalf("Invalid square size: %v.\n", *squareSize)
	}
	if *numThreads < 1 {
		errorLog.Fatalf("Invalid number of threads: %v.\n", *numThreads)
	}

	isPiped := stdout_is_piped()
	infoLog.Printf("size = %v; # threads = %v; pretty-print = %v; unique = %v; piped = %v\n",
		*squareSize, *numThreads, *prettyPrint, *requireUnique, isPiped)

	wordArr := read_words(*inputPath, *squareSize)
	wordTree := NewWordTreeNode()
	for _, word := range wordArr {
		wordTree.insert(word)
	}

	startIdx := word_arg_to_index(wordArr, *startWord, 0)
	endIdx := word_arg_to_index(wordArr, *endWord, len(wordArr)-1)
	infoLog.Printf("Start word range: %v - %v\n", wordArr[startIdx], wordArr[endIdx])

	state := ProgState{
		square_size:    *squareSize,
		num_threads:    *numThreads,
		word_arr:       wordArr,
		word_tree:      &wordTree,
		start_idx:      startIdx,
		end_idx:        endIdx,
		pretty_print:   *prettyPrint,
		require_unique: *requireUnique,
		is_piped:       isPiped,
	}
	build_squares(state)
}
