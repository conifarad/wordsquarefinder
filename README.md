# wordsquarefinder

This is a command-line tool for generating word squares and double word squares,
[as described on Wikipedia](https://en.wikipedia.org/wiki/Word_square). It uses
a tree-search algorithm and supports multithreading for optimal efficiency.

### Examples:

```
    size 7              size 8:
(double square):
                    c r a b w i s e
 d e c o l o r      r a t l i n e s
 o p e r a t e      a t l a n t e s
 w i r e t a p      b l a s t e m a
 e t a g e r e      w i n t e r l y
 r o m a n i a      i n t e r t i e
 e m i n e n t      s e e m l i e r
 d e c o d e s      e s s a y e r s
```

Generated squares are currently limited to the English alphabet (26 letters and
no diacritics) only.

The inspiration for this project was the "weekly" mode in FurbyFubar's game
[Squardle](https://fubargames.se/squardle/).

## Compilation

Download and install [Go](https://go.dev/dl/) if you haven't already. Then
simply run `go build` in the project directory to build.

## Usage

`./wordsquarefinder [args]`

### Arguments:

* `-n`: Size of the word squares to generate. Default is 5.
* `-t`: Number of threads (goroutines) to run. If greater than one, squares are
  not guaranteed to be returned in alphabetical order! Defaults to half the
  number of logical CPU cores available. Experiment with different values to see
  which is most efficient for your system.
* `-i`: Path to the input word list. Expects one word per line; words that are
  too long, too short or contain unsupported characters are skipped. The
  included default list, `english_filtered.txt`, is based on `english3.txt` from
  [this website][1] with *some* offensive words (including some from
  [this list][2]) filtered out. `wordle_solutions.txt`, a sorted list of all
  5-letter Wordle solutions, is also included.
* `-s` and `-e`: Defines a range of starting words to search. For example, for
  n=5, `-s crate -e treks` will search the following range:
	```
	crate       treks
	aaron       zymes
	aaron  ...  zymes
	aaron       zymes
	aaron       zymes
	```
* `-p`: Pretty-print output, i.e. output words on separate lines. By default,
  squares are output in comma-separated form for easy CSV parsing.
* `-u`: Require unique words. This will eliminate squares like the second
  example above, where corresponding rows/columns have the same word. Note that
  the rows and columns of any square can be swapped to produce another valid
  square, so squares with unique words will appear in both configurations.

Each thread searches for squares starting with a particular word. When running,
the program displays the first and last starting words currently being searched.

Notice there is no option to specify an output file; instead, pipe output to a
file like so:

`./wordsquarefinder [args] > output.txt`

Log messages are written to `stderr` and are not included in the output file.

[1]: http://gwicks.net/justwords.htm
[2]: https://github.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words
