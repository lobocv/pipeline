# This is a collection of aliases that I find helpful

########### NAVIGATION ###########
# Go to previous dir
alias cdc="cd -"

# Go up one dir
alias cx="cd .."

# Go up two dir
alias cxx="cd ../.."

alias dl="cd ~/Downloads"

########### RC FILE #############
alias zshrc="vi ~/.zshrc"
alias zshrcl="source ~/.zshrc"

########### TEXT FORMATTING ###############

# Formats stdin to pretty JSON
alias jsonp="python -m json.tool"

alias stripnl="tr -d '\r' | tr -d '\n'"


########### TOOLS ###############

# Echo the return code of the last command
alias rc="echo \$?"

# Swap two  paths
# $1 : path1
# $2 : path2
function swap() {
	mv "$2" "$2.swaptemp"
	mv "$1" "$2"
	mv "$2.swaptemp" "$1"
}

# Retry a command repeatedly until it exits with status code 0
function retry() {
	ec="1"
	count=0
	cmd=${@:1}
	out=$(mktemp /tmp/retry.XXXXXX)
	while : ; do
		eval $cmd 2> "$out"
		ec="$?"
		result=$(cat "$out")
		# Only echo output if it is different
		if [[ "$result" != "$prev_result" ]]; then
			echo $result
		fi
		if [[ "$ec" = "0" ]]; then
			break
		fi
		echo -n "."
		sleep 5
		prev_result=$result
	done
}
