gp
---

> Holy Forking Shirtballs!

`gp` helps you keep track of directories you like.

## Installation

`go install github.com/packrat386/gp@main`

## Usage

Reference:

```
gp usage:
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

gp home page: https://github.com/packrat386/gp
```

Example:

```
$ pwd
/Users/packrat386/github/gp
$ gp ++
score for '/Users/packrat386/github/gp' is incremented to: 1
$ gp --
score for '/Users/packrat386/github/gp' is decremented to: 0
$ gp ++ 5
score for '/Users/packrat386/github/gp' is incremented to: 5
$ gp -- 2
score for '/Users/packrat386/github/gp' is decremented to: 3
$ gp query $PWD
/Users/packrat386/github/gp: 3
$ gp list
/Users/packrat386: 10
/Users/packrat386/github/gp: 3
/Users/packrat386/Downloads: 1
/Users/packrat386/aws_training: -10
```

## TODO

* Use advisory locking (i.e. the `flock` syscall) to ensure that concurrent access can't bork the storage
* Configurable storage backends
* More advanced querying and sorting
* Tests (lol)
