package main

import (
	"testing"

	"golang.org/x/exp/slices"
)

func assert_eq[C comparable](t testing.TB, a C, b C) {
	if a != b {
		t.Fatalf("expected equal, got %v != %v", a, b)
	}
}

func assert_neq[C comparable](t testing.TB, a C, b C) {
	if a == b {
		t.Fatalf("expected not equal, got %v == %v", a, b)
	}
}

func TestAllUniqueWords(t *testing.T) {
	hDuplicates := []string{
		"abcd",
		"efgh",
		"ijkl",
		"abcd",
	}
	vDuplicates := []string{
		"abca",
		"bcdb",
		"cdec",
		"defd",
	}
	crossDuplicates := []string{
		"eafg",
		"abcd",
		"hcij",
		"kdlm",
	}
	noDuplicates := []string{
		"abcd",
		"efgh",
		"ijkl",
		"mnop",
	}
	assert_eq(t, all_unique_words(hDuplicates), false)
	assert_eq(t, all_unique_words(vDuplicates), false)
	assert_eq(t, all_unique_words(crossDuplicates), false)
	assert_eq(t, all_unique_words(noDuplicates), true)
}

func TestWordTree(t *testing.T) {
	root := NewWordTreeNode()
	includedWords := []string{"test", "tests", "zebra"}
	for _, word := range includedWords {
		root.insert(word)
	}

	assert_eq(t, root.get_child('a'), nil)
	assert_neq(t, root.get_child('z'), nil)

	node := root.get_child('t').get_child('e').get_child('s').get_child('t')
	assert_neq(t, node.get_child('s'), nil)
	assert_eq(t, node.get_child('s').get_child('s'), nil)
	assert_eq(t, node.get_child('y'), nil)
}

var testBoards = [][]string{
	{"palays", "abadan", "lexeme", "alisma", "cisted", "eaters"},
	{"palays", "abadan", "lexeme", "alisma", "tested", "esters"},
	{"palays", "aboral", "karaka", "eyelet", "halite", "assays"},
	{"palays", "aboral", "torero", "trices", "encash", "reasty"},
	{"palays", "aboral", "torero", "trines", "encash", "reasty"},
	{"palays", "adonai", "rogers", "crimes", "ancile", "essays"},
	{"palays", "adonai", "tracks", "iodous", "niente", "strass"},
	{"palays", "adonai", "trucks", "iodous", "niente", "strass"},
	{"palays", "adonai", "trunks", "iodous", "niente", "strass"},
	{"palays", "agorae", "rubati", "giants", "esteem", "shears"},
	{"palays", "agorot", "unlade", "palled", "etoile", "repass"},
	{"palays", "amazon", "reside", "insole", "atoner", "lasers"},
	{"palays", "amebae", "namers", "azalea", "monism", "ansate"},
	{"palays", "amenta", "reests", "enters", "nellie", "tremas"},
	{"palays", "amoret", "rafale", "attila", "neesed", "adders"},
	{"palays", "amoret", "regime", "irised", "accend", "nessie"},
	{"palays", "amoret", "tendre", "engobe", "neural", "tressy"},
	{"palays", "amulet", "lascar", "attune", "tories", "elands"},
	{"palays", "amulet", "lessor", "intima", "nerkas", "greens"},
	{"palays", "amulet", "lessor", "intime", "nerkas", "greens"},
	{"palays", "amulet", "retama", "abeles", "narine", "assais"},
	{"palays", "amulet", "ruddle", "aslope", "meused", "ormers"},
	{"palays", "amulet", "woggle", "argala", "weiter", "steeds"},
	{"palays", "amulet", "woggle", "argala", "weiter", "steers"},
	{"palays", "anemia", "ligers", "elands", "settee", "treads"},
	{"palays", "anemia", "ligers", "elands", "settee", "troads"},
	{"palays", "anicut", "topeka", "riatas", "ensate", "stelas"},
	{"palays", "animal", "litera", "amened", "turtle", "essays"},
	{"palays", "animal", "rivera", "amened", "dartle", "essays"},
	{"palays", "animal", "rivera", "amened", "gentle", "essays"},
	{"palays", "aranea", "potass", "aments", "yankee", "asters"},
	{"palays", "aranea", "sanely", "trompe", "easier", "steads"},
	{"palays", "aranea", "sanely", "trompe", "easier", "stears"},
	{"palays", "ararat", "caribe", "enrobe", "reused", "sapors"},
	{"palays", "ararat", "lanate", "incite", "neesed", "gaters"},
	{"palays", "ararat", "lorate", "amrita", "mauser", "aspers"},
	{"palays", "ararat", "retake", "aneled", "danite", "essays"},
	{"palays", "ararat", "retake", "aniler", "danite", "essays"},
	{"palays", "arouet", "coggle", "illipe", "flotel", "yagers"},
	{"palays", "arouet", "measle", "pattle", "achier", "shends"},
	{"palays", "atabal", "nagari", "abused", "denise", "agents"},
	{"palays", "atabek", "coarse", "ingate", "needer", "orrery"},
	{"palays", "atabek", "thyrse", "remote", "enamel", "senary"},
	{"palays", "ativan", "polite", "uniate", "letter", "ashery"},
	{"palays", "avidin", "rebore", "ananke", "mitier", "oreads"},
	{"palays", "avocet", "rattle", "attila", "naiver", "greeds"},
	{"palays", "avowal", "litera", "latten", "otiose", "reests"},
	{"palays", "egeria", "alecky", "retake", "calder", "emeers"},
	{"palays", "elanet", "settle", "exhale", "tiered", "anears"},
	{"palays", "eleven", "liaise", "inmate", "teeter", "eddery"},
	{"palays", "eleven", "liaise", "innate", "teeter", "eddery"},
	{"palays", "eleven", "lupine", "oleate", "tartar", "assess"},
	{"palays", "enerve", "teston", "attend", "lierne", "scryer"},
	{"palays", "evolve", "waggon", "tigons", "elaine", "rended"},
	{"palays", "evovae", "radars", "ungirt", "stelae", "eident"},
	{"palays", "ibadan", "lexeme", "alisma", "fitted", "sayers"},
	{"palays", "ibadan", "lexeme", "alisma", "tested", "esters"},
	{"palays", "ibadan", "lexeme", "alisma", "witted", "sayers"},
	{"palays", "igorot", "lovage", "uranin", "latent", "assais"},
	{"palays", "imaret", "neagle", "engild", "reeved", "orrery"},
	{"palays", "imaret", "nubile", "aslope", "teased", "arbors"},
	{"palays", "imaret", "nubile", "aslope", "teasel", "arbors"},
	{"palays", "imaret", "nubile", "aslope", "teaser", "arbors"},
	{"palays", "imaret", "nuncle", "escape", "reeded", "orders"},
	{"palays", "imaret", "nuncle", "escape", "reeded", "orrery"},
	{"palays", "imaret", "nuncle", "escape", "reeden", "orders"},
	{"palays", "imaret", "nuncle", "escape", "reeder", "orders"},
	{"palays", "imaret", "tingle", "ancile", "reeved", "asters"},
	{"palays", "isobel", "tadema", "aneled", "ranine", "assais"},
	{"palays", "loiret", "unable", "fibula", "falter", "sneers"},
	{"palays", "oberon", "ureide", "sapele", "sieger", "eddery"},
	{"palays", "olivet", "debase", "sprite", "opaled", "loners"},
	{"palays", "olivet", "oppose", "depute", "leered", "enders"},
	{"palays", "olivet", "pinole", "ungula", "leered", "idlers"},
	{"palays", "origan", "monera", "amener", "candle", "essays"},
	{"palays", "origan", "monera", "amener", "dandle", "essays"},
	{"palays", "origan", "soneri", "amened", "dandle", "assays"},
	{"palays", "orison", "tinkle", "ashake", "starer", "sayids"},
	{"palays", "swerve", "orator", "racing", "achene", "skyres"},
}

func BenchmarkFinderFuncs(b *testing.B) {
	const squareSize = 6
	const numLetters = squareSize * squareSize
	const firstWord = "palays"

	wordArr := read_words("english_filtered.txt", 6)

	wordTree := NewWordTreeNode()
	for _, word := range wordArr {
		wordTree.insert(word)
	}

	b.Run("recursive", func(b2 *testing.B) {
		board := make([]byte, 0, numLetters)
		aboveNodes := make([]*WordTreeNode, 0, numLetters)
		boardChan := make(chan []string, 100)

		for i := 0; i < squareSize; i++ {
			board = append(board, firstWord[i])
			aboveNodes = append(aboveNodes, wordTree.get_child(firstWord[i]))
		}

		b2.ResetTimer()
		build_squares_recursive(squareSize, true, &wordTree, &wordTree, aboveNodes, board,
			boardChan)
		b2.StopTimer()
		close(boardChan)

		for i, testBoard := range testBoards {
			gotBoard, ok := <-boardChan
			if !(ok && slices.Equal(testBoard, gotBoard)) {
				b2.Fatalf("Failed on board %v.\n", i)
			}
		}

		extraBoard, ok := <-boardChan
		if ok {
			b2.Fatalf("Returned extra board: %v.\n", extraBoard)
		}
	})
}
