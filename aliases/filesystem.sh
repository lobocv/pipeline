
alias df="df -h | cCyan '[0-9]+\.?[0-9]?G' | cYellow '[0-9]+\.?[0-9]?M' | cGreen '[0-9]+\.?[0-9]?K' | cRed '9[0-9]\%' | cLightRed '8[0-9]\%'"
alias flw="tail -f -v"
alias fullpath="readlink -e"

function echoerr() {
	echo "$@" 1>&2
}

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
	DIR=${2:-.}
        find "${DIR}" -iname "*$1*"
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

# Writes a file containing random characters
# Example: randfile /tmp/randfile.txt 10  # Make 10KB file
# Example: randfile /tmp/randfile.txt 10 1048576 # Make 10MB file
# $1: Filepath
# $2: Number of blocks
# $3: Block size (Default: 1KB)
function randfile() {
	file=${1}
	COUNT=${2}	
	BS=${3:-1024}
	dd if=/dev/urandom of="${file}" bs=${BS} count=${COUNT}
}

# Rename file with .bak appended to it. If a file with .bak is already found, it will undo the process
# and remove the .bak
# Example: If a file named test.txt exists.
# The following renamed test.txt to test.txt.bak
# >> bak test.txt 
# And the following will renamed test.txt.bak back to test.txt
# >> bak test.txt.bak 
# OR 
# >> bak test.txt
function bak() {
	file="$1"
	REGEX='.*bak'
	if [[ ${file} =~ '.*bak' ]]; then
		# strip .bak from filename
		renamed="${file%".bak"}" 
	elif [[ -f ${file} ]]; then
		# add .bak to filename
		renamed="${file}.bak"
	elif  [[ ! -f ${file} ]] && [[ -f "${file}.bak" ]]; then
		renamed="${file}"
		file="${file}.bak"
	fi
	mv -v -i "${file}" "${renamed}" 1>&2
}

# List options for previous directories to change to
function fd() {
	PAST_DIRS=(`pushd`)
	if [[ ${#PAST_DIRS} -gt 1 ]]; then
		echo "Choose a path:"
		lc=1
		for x in $PAST_DIRS;
		do
			echo "$lc. $x"
			lc=$((lc+1))
		done
		read p
		if [[ ${p} -gt ${#PAST_DIRS} ]]; then
			echoerr "No such option $p exists"
			return 1
		fi
		cd $(echo "$PAST_DIRS[$(($p))]" | sed -r "s|~|$HOME|" )

	else
		echoerr "No directory history"
		return 1
	fi

}

# Compress a directory to .tar.gz
# $1 : Path to directory or file
# $2 : Output filename (Default: $1.tar.gz)
function tardir() {
	DIR="$1"
	OUT="${2:-$1.tar.gz}"
	tar -czf $OUT $DIR
}

# Untar / extract a .tar.gz file
# $1 : Path to .tar.gz
# $2 : Output path (Default: ./)
function untar() {
	TAR="$1"
	OUT="${2:-./}"
	tar -xzf $TAR -C $OUT
}

# List contents of tar file
# $1 : Path to .tar.gz
function tarview() {
	tar -tf "$1"
}
