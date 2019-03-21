alias gitcom="git commit"
alias gitam="git commit --amend"
alias gitch="git checkout"
alias gd="git diff"
alias gl="git log"
alias gs="git status"

# Rebase the current branch onto the most up to date specified branch
# $1: Branch to rebase onto (Default: master)
function gitrebase() {
	set -e
	base=${1:-master}
	echo git checkout "$base"
	echo git pull
	echo git checkout -
	echo git rebase "$base"
}
