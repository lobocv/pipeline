
alias df="df -h"

# Show files over a certain file size
# $1: Human-readable file size (ex. 50M, 100K, 2G)
function filesover() {
        USAGE=$(du --all --human-readable --threshold="$1")
        echo "$USAGE" | while read line;
        do
                fp=$(echo "$line" | cut -f 2)
                if [[ ! -d "$fp" ]]; then
                        echo $line
                fi

        done
}

# Show files over a certain file size
# $1: Human-readable file size (ex. 50M, 100K, 2G)
function dirsover() {
        du --human-readable --threshold="$1"
}

# Swap two  paths
# $1 : path1
# $2 : path2
function swap() {
        mv "$2" "$2.swaptemp"
        mv "$1" "$2"
        mv "$2.swaptemp" "$1"
}

# Search for a file recursively from the current directory
# $1: File name (case insensitive)
function search() {
        find . -iname "*$1*"
}


