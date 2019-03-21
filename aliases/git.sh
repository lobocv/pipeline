alias gitcom="git commit"
alias gitam="git commit --amend"
alias gitch="git checkout"
alias gd="git diff"
alias gl="git log"
alias gs="git status"

# Rebase the current branch onto the most up to date specified branch
# $1: Branch to rebase onto (Default: master)
function gitrebase() {
	base=${1:-master}
	git checkout "$base"
	git pull
	git checkout -
	git rebase "$base"
}
