local ALIAS_DIR=$(dirname $0)/aliases

function load_aliases() {
	local f
	
	for f in $(ls $ALIAS_DIR/*.sh); do
		if [[ $(basename $f)  == _* ]]; then
			continue
		fi
	        source $f
	done
}

load_aliases

