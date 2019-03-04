
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

###########  GIT  ###############
alias gitcom="git commit"
alias gitam="git commit --amend"
alias gitch="git checkout"
alias gd="git diff"
alias gl="git log"
alias gs="git status"

########### DOCKER ###############
function dockerssh() {
	container=( `docker ps --format "{{.Names}}" | grep $1` )
	num=${#container}
	if [[ $num -gt 1 ]]; then
		echo "More than one live container found for grep \"$1\". Choose a container from below:"
		lc=1
		for x in $container; 
		do 
			echo $lc. $x
			lc=$((lc+1))
	       	done
		read c
		target=$container[$(($c))]
	else
		target=$container[1]
	fi
	
	docker exec -it $target ${2:-bash}
}

function dockerrun() {
	images=( `docker images --format "{{.Repository}}:{{.Tag}}" | grep $1` )
	if [[ ${#images} -gt 1 ]]; then
		echo "More than one image:tag found for grep \"$1\". Choose an image below:"
		lc=1
		for x in $images; 
		do 
			echo $lc. $x
			lc=$((lc+1))
	       	done
		read c
		target=$images[$(($c))]
	else
		target=$images[1]
	fi
	
	docker run --rm -it $target ${2:-bash}
}


########### TOOLS ###############

alias jsonp="python -m json.tool"
