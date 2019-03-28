
alias df="df -h | cCyan '[0-9]+\.?[0-9]?G' | cYellow '[0-9]+\.?[0-9]?M' | cGreen '[0-9]+\.?[0-9]?K' | cRed '9[0-9]\%' | cLightRed '8[0-9]\%'"

alias fullpath="readlink -e"

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

# Recursively search for text within files in a folder
# $1: Search text
# $2: Search directory (Default: $PWD)
function textsearch() {
	DIR=${2:-.}
	grep -R "$1" ${DIR}
}

# Capture a filepath into a buffer that can be used with other commands (drop, etc) 
function hold() {
	if [ -f "$1" ]; then
		export HOLD_BUFFER="$(readlink -e $1)"
		echo "Holding $HOLD_BUFFER" 1>&2 
	else
		echo "Hold failed: path $1 does not exist" 1>&2 
		return "1"
	fi
}

function whatsheld() {
	echo "$HOLD_BUFFER"
}
function droppath() {
	export HOLD_BUFFER=""
}

# Executes a command providing the hold buffer as the first argument
# Example1: release cp myfile.txt
# Example2: release chmod +x
# $1: Command to run
# $@: Additional arguments or flags to the command called after the hold buffer
function release() {
	cmd="$1"
	shift
	FLAGS=()
	ARGS=()
	while [ ! -z "$1" ];
	do
		if [[ "$1" == [-+]* ]]; then
			FLAGS+=$1
		else
			ARGS+=$1
		fi
		shift
	done
	
	if [ ! -z "$HOLD_BUFFER" ]; then
		$cmd $FLAGS "$HOLD_BUFFER" $ARGS
	else
		echo "Hold buffer is empty" 1>&2
		return 1
	fi
}

