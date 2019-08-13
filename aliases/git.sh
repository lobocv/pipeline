alias gitcom="git commit"
alias gitam="git commit --amend"
alias gd="git diff"
alias gl="git log"
alias gs="git status"
alias gitch="git checkout"
alias newbranch="git checkout -b"

# Checkout a local git branch
# $1: Name of branch or filter
function gc() {
	local BRANCHES BRANCH="$1"
	BRANCHES=(`git branch --list | grep $BRANCH`)
	if [[ ${#BRANCHES} -gt 1 ]]; then
		echo "More than one branch is found for grep \"$BRANCH\". Choose a branch:"
		lc=1
		for x in $BRANCHES;
		do
			if [[ "$x" = "*" ]]; then
				continue;
			fi
			echo $lc. $x
                        lc=$((lc+1))
                done
                read c
                B=$BRANCHES[$(($c))]
	else
		B=$BRANCHES[1]
	fi
	git checkout $B
}

# Rebase the current branch onto the most up to date specified branch
# $1: Branch to rebase onto (Default: master)
function gitrebase() {
	base=${1:-master}
	git checkout "$base"
	git pull
	git checkout -
	git rebase "$base"
}

# Get the last hash of the given branch
# $1: Branch name (Default: HEAD)
function githash() {
	branch=${1:-HEAD}
	git rev-parse $branch
}

function gitlastdiff() {	
	git diff $(git rev-parse HEAD~1) $(git rev-parse HEAD) 
}

# Delete a branch
# $1: Branch name (Default: current branch)
function gitdel() {
	CURRENT=$(git rev-parse --abbrev-ref HEAD)
	BRANCH=${1:-$CURRENT}
	echo "Are you sure you want to delete branch $BRANCH? [y|N]" 
	read cont
	if ! _confirm_yesno "$cont"; then
	        return 0
	fi
	
	if [[ "$CURRENT" = "$BRANCH" ]]; then
	        git checkout - > /dev/null
	fi
	
	git branch -d "$BRANCH"
}

