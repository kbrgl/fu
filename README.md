```
______       ____
\ \ \ \     / __/_  __
 \ \ \ \   / /_/ / / /
 / / / /  / __/ /_/ /
/_/_/_/  /_/  \__,_/

```

# Find Unleashed
Fu is an intuitive alternative to the Unix command find. It mainly searches
filenames, but some meta-information like file mode is supported as well.
It can walk directories in parallel, resulting in substantial speedups for
directories with thousands of files.

## Index
* [Installation](#installation)
    + [Platform support](#platform-support)
* [Usage](#usage)
    + [Examples](#examples)
* [TODO](#todo)
* [Contributing](#contributing)

## Installation
If you have Go installed, run
```sh
go get -u github.com/kbrgl/fu
```
If you don't, download the latest release and put it somewhere on your $PATH.

### Platform support
Unix-based systems (macOS, Linux distros, \*BSDs, etc.) are supported.
I'm guessing Windows will work as well, although I haven't tried it out on a
Windows machine yet.

## Usage
```sh
usage: fu [<flags>] <query> [<paths>...]

Flags:
  -h, --help       Show context-sensitive help (also try --help-long and
                   --help-man).
  -f, --fuzzy      Use fuzzy search
  -r, --regexp     Use regexp-based search
  -a, --suffix     Use suffix-based search (short flag 'a' is short for for
                   'after')
  -b, --prefix     Use prefix-based search (short flag 'b' is short for
                   'before')
  -s, --substring  Use substring-based search allowing the query to be at any
                   position in the filename
  -d, --dir        Show only directories
  -m, --perm=PERM  Filter by Unix permissions
  -c, --parallel   Walk directories in parallel, may result in substantial
                   speedups for directories with many files
  -v, --version    Show application version.

Args:
  <query>    Search query
  [<paths>]  Paths to search
```
Path is current working directory by default, and the program uses exact
filename matching by default.

### Examples
The general usage pattern of fu is:
```sh
fu [<search modifiers>...] <query> [<paths>...]
```
Where the search modifiers may be any of the options specified in the usage.

The only thing that's actually necessary is the query. By default,
```
fu <query>
```
Will search under the current directory for files with the name <query>.

#### Find all dotfiles in the home folder
```sh
fu -b . ~
```
#### Find all Python scripts in current dir
```sh
fu -a .py
```
#### Find files in current dir whose names start and end with 'a'
```
fu -r "^a.*a\$"
```
#### Fuzzy matching
```
fu -f pkg
```

Fuzzy search is the algorithm that Sublime Text and Atom use in their
Command Palettes. For an explanation, check out [this Wikipedia article]
(https://en.wikipedia.org/wiki/Approximate_string_matching).

## TODO
* Write tests for matchers
* Break up the main function into small, testable functions
* Beautify TTY output somehow (colorization or something)
* Find a way to estimate how much time a search will take, and print this
  estimate before performing the search. (I'm not sure if this is possible).
* Remove -c flag and automatically switch between powerwalk and filepath Walk
  functions

PRs for these would be appreciated.

## Contributing
I've made an effort to keep the Fu source code simple, modular and easy to browse.
Contributions are welcome. Just be sure to run `go fmt` and `go vet` on your
code.
