# Disk

Disk is a CLI tool that helps you identify files that you can remove.


## Motivation

I have a relatively small SSD as my main hard drive and I have to free up space more often than I'd like.
I haven't found any tools that are aggressive in suggesting files to remove, in particular when it comes to development
clutter like `node_modules`, `mod`, `.build`, `dist` and so on. I always have to fall back on manually looking for files.
Luckily it's a highly mechanical exercise, and therefore easy to automate.

## Usage

The primary command is `disk clean`, it will show a TUI with a table of files and folders that have been marked as candidates for deletion.

Examples of files and folders it will suggest:
- Clutter in the form caches, dependency folders, build folders, etc. above a given size and age.
- Steam games you haven't played in a while.
- Movies and TV shows that are easy to get a hold of even if you delete them.

To execute it run:
```
disk clean <path>
```
It has the following flags, though hopefully the defaults are good enough that you don't have to bother with them: 
```
Flags:
  -h, --help               help for clean
  -p, --max-playtime int   Maximum playtime of games to include in analysis results specified in hours. (default 20)
  -a, --min-age int        Minimum age of files to include in analysis results specified in days. (default 90)
  -s, --min-size int       Minimum size of files to include in analysis results specified in megabytes. (default 50)
```

The two other commands are there to help you find the appropriate `path` to run `disk clean` on.
```
disk usage
```
That outputs:
```
Disk    Size    Used    Available
C:\     464GB   410GB   54GB
D:\     931GB   387GB   544GB
Total:  1395GB  797GB   598GB
```

And
```
disk tree <path>
```
That outputs:
```
<path>: 11kB
├── file1: 1kB
├── file2: 2kB
└── file3: 8kB
```

The `disk tree` command has the following flags, where the `-d` flag is quite handy:
```
Flags:
  -d, --depth int   Depth of the tree structure. (default 1)
  -h, --help        help for tree
```

## Installation

Windows:
```
scoop bucket add disk https://github.com/sebastianappelberg/scoop-bucket 
scoop install disk
```

MacOS:
```
brew tap sebastianappelberg/homebrew-tap
brew install disk
```

Using Go:
```
go install github.com/sebastianappelberg/disk@latest
```

## Contributing

The easiest way to contribute is to take a look at the [clutter_folders.json](/pkg/config/clutter_folders.json) and
[unsafe_folders.json](/pkg/config/unsafe_folders.json) and see if there are folders that you think should be included 
for consideration when running `disk clean`. Other contributions are of course welcome. Please make sure your pull request
contains a description containing the motivation for the changes.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

## Acknowledgements

Shoutout to the following projects:
- https://github.com/Kei-K23/trashbox
- https://github.com/go-ole/go-ole
- https://github.com/razsteinmetz/go-ptn
- https://github.com/andygrunwald/vdf
- https://github.com/spf13/cobra
- https://github.com/charmbracelet

