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

########### TOOLS ###############

# Formats stdin to pretty JSON
alias jsonp="python -m json.tool"

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
