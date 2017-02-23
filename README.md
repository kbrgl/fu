# Find Unleashed
Fu is a more intuitive version of the Unix command find. It mainly searches
by matching on filenames, but some meta-information like file mode is
supported as well.

## Index
* [Installation](#installation)
* [Platform support](#platform-support)
* [Usage](#usage)
* [Examples](#examples)
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
  -s, --suffix     Use suffix-based search
  -p, --prefix     Use prefix-based search
  -d, --dir        Show only directories
  -m, --perm=PERM  Filter by Unix permissions
  -c, --parallel   Walk directories in parallel, may result in substantial
                   speedups for directories with many files

Args:
  <query>    Search query
  [<paths>]  Paths to search
```
Path is current working directory by default, and the program uses exact
filename matching by default.

### Examples
The general usage pattern of fu is:
```sh
fu [<matchers>...] <query> [<paths>...]
```
Where matcher may be any of the matchers specified in the usage.

The only thing that's actually necessary is the query. By default,
```
fu <query>
```
Will recursively search for files with the exact same name as <query> using
the current directory as the root.
#### Find all dotfiles in the home folder
```sh
fu -p . ~
```
### Find all Python scripts in current dir
```sh
fu -s .py
```
### Find files in current dir whose names start and end with 'a'
```
fu -r "^a.*a\$"
```
### Fuzzy matching
```
fu -f pkg
```
Will match any files that contain the letters p, k and g anywhere in the
filename provided they occur in that exact order, so names like 'package',
'parking' and 'pkg2.txt' will be matched but names like 'p13ag5uk' where the
letters do not appear in order won't.

Fuzzy search is the algorithm that Sublime Text and Atom use in their
Command Palettes.

## TODO
* Write tests for matchers
* Break up the main function into small, testable functions
* Beautify TTY output somehow (colorization or something)
* Find a way to estimate how much time a search will take, and print this
  estimate before performing the search. (I'm not sure if this is possible).

PRs for these would be appreciated.

## Contributing
Contributions are welcome. Just be sure to run `go fmt` and `go vet` on your
code.
