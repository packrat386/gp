package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

func main() {
	if len(os.Args) == 1 {
		help()
		fail(errors.New("no command supplied"))
	}

	cmd := os.Args[1]
	rest := os.Args[2:]

	switch cmd {
	case "increment", "++":
		increment(rest)
	case "decrement", "--":
		decrement(rest)
	case "query":
		query(rest)
	case "list":
		list(rest)
	case "help":
		help()
	default:
		help()
		fmt.Println("---")
		fail(errors.New("unrecognized command"))
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	os.Exit(1)
}

func failIf(err error) {
	if err != nil {
		fail(err)
	}
}

func must[T any](result T, err error) T {
	if err != nil {
		fail(err)
	}

	return result
}

func storageLocation() string {
	if fname := os.Getenv("GP_STORAGE"); fname != "" {
		return fname
	} else {
		return filepath.Join(must(os.UserHomeDir()), ".gp_storage")
	}
}

func scoreArg(args []string) int64 {
	if len(args) > 1 {
		fail(fmt.Errorf("unexpected arguments: %v", args[1:]))
	}

	if len(args) == 0 {
		return 1
	}

	inc, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fail(fmt.Errorf("could not parse number"))
	}

	return inc
}

func increment(args []string) {
	score := scoreArg(args)
	wd := must(os.Getwd())

	s := must(openFileStorage(storageLocation()))
	e := must(s.read(wd))

	e.score += score

	failIf(s.write(e))
	failIf(s.persistAndClose())

	fmt.Printf("score for '%s' is incremented to: %d\n", e.name, e.score)
}

func decrement(args []string) {
	score := scoreArg(args)
	wd := must(os.Getwd())

	s := must(openFileStorage(storageLocation()))
	e := must(s.read(wd))

	e.score -= score

	failIf(s.write(e))
	failIf(s.persistAndClose())

	fmt.Printf("score for '%s' is decremented to: %d\n", e.name, e.score)
}

func query(args []string) {
	if len(args) == 0 {
		args = append(args, must(os.Getwd()))
	}

	s := must(openFileStorage(storageLocation()))

	for _, dname := range args {
		e := must(s.read(dname))

		fmt.Printf("%s: %d\n", e.name, e.score)
	}
}

func list(args []string) {
	if len(args) != 0 {
		fail(fmt.Errorf("unexpected arguments: %v", args))
	}

	s := must(openFileStorage(storageLocation()))
	ee := must(s.all())

	slices.SortFunc(ee, func(l, r entry) int {
		return int(r.score - l.score)
	})

	for _, e := range ee {
		fmt.Printf("%s: %d\n", e.name, e.score)
	}
}

func help() {
	fmt.Println(`gp usage:
gp (increment|++) <score>? => increment the score for the current working
                              directory by <score>, default 1
gp (decrement|++) <score>? => decrement the score for the current working
                              directory by <score>, default 1
gp query <dirname>...      => get the score for <dirname>(s), defaults to
                              current working directory
gp list                    => list the scores for all directories, sorted
                              from highest to lowest
gp help                    => display this help message

the environment variable GP_STORAGE points to the storage file, which will
be created if it does not exist.

gp home page: https://github.com/packrat386/gp`)
}
