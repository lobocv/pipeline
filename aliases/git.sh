alias gitcom="git commit"
alias gitam="git commit --amend"
alias gitch="git checkout"
alias gd="git diff"
alias gl="git log"
alias gs="git status"
alias newbranch="git checkout -b"

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

