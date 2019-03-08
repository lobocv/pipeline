alias containernames="docker ps --format "{{.Names}}""

function findcontainer() {
	if [[ -z "$1" ]]; then
		containers=( `docker ps --format "{{.Names}}"` )
	else
		containers=( `docker ps --format "{{.Names}}" | grep $1` )
	fi
	num=${#containers}
	if [[ $num -gt 1 ]]; then
		echo "More than one live container found for grep \"$1\". Choose a container from below:"
		lc=1
		for x in $containers; 
		do 
			echo $lc. $x
			lc=$((lc+1))
	       	done
		read c
		container=$containers[$(($c))]
	else
		container=$containers[1]
	fi
}

# SSH into a running docker container.
# $1: Container name to grep for
# $2: Shell to use (default: bash)
function dockerssh() {
	findcontainer "$1"
	docker exec -it "$container" ${2:-bash}
}

# Follow logs of a running docker container.
# $1: Container name to grep for
function dockerfollow() {
	findcontainer "$1"
	docker logs -f "$container"
}

# Dump logs of a running docker container.
# $1: Container name to grep for
# $2: Output file path
function dockerlog() {
	findcontainer "$1"
	docker logs "$container"
}

# Kill a running docker container.
# $1: Container name to grep for
function dockerkill() {
	findcontainer "$1"
	docker kill $2 "$container"
}

# Get the info of a running docker container.
# $1: Container name to grep for
function dockerinspect() {
	findcontainer "$1"
	docker inspect "$container"
}


# Get the IP of a running docker container.
# $1: Container name to grep for
function dockerip() {
	findcontainer "$1"
	docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{printf "\n"}}{{end}}' "$container"
}


# Get the host virtual network interface of the network interface of a running docker container.
# $1: Container name to grep for
# $2: container network interface (Default: eth0)
function dockeriface() {
	findcontainer "$1"
	iface=${2:-eth0}
	ifindex=$(docker exec -it $container sh -c "cat /sys/class/net/${iface}/iflink" | tr -d '\r')
	iface_path_host=$(grep -R "${ifindex}" /sys/class/net/*/ifindex)
	echo "$(basename $(dirname $iface_path_host))"
}

# List the IPs of all running docker container.
# $1: Container name to grep for
function dockeriplist() {
	docker ps -q | xargs -n 1 docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}} {{end}} {{ .Name }}' | sed 's/ \// /'
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
	
	docker run --rm -it --entrypoint ${2:-bash} $target
}

