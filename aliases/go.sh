#!/bin/bash

# Run tests with coverage and create a colorized report
gocover () {
    local t=$(mktemp -t coverXXXX)
    go test $COVERFLAGS -coverprofile=$t $@ \
        && go tool cover -func=$t | cRed '\s[0-9]\.?[0-9]?\%' | cLightRed '[1-5][0-9]\.[0-9]\%' | cLightGreen '[6-9][0-9]\.[0-9]\%'| cGreen '9[0-9]\.[0-9]\%' | cGreen '100\.?0?\%' 
    unlink $t   
}

gotest function() {
	SUITE="$1"
	shift;
	local NOCACHE=""
	for i in "$@"; do
		case "$i" in
		-n|--no-cache)
			NOCACHE="-count=1"
			shift
			;;
		*)
			
		esac
	done
	go test -run Test$SUITE $NOCACHE $@ | cGreen "--- PASS:" | cRed "--- FAIL:" | cBold "--- PASS:" | cBold "--- FAIL:" | removeline "=== (RUN\|CONT\|PAUSE).*"
	if [[ "$?" == "0" ]]; then
		echo "All tests passed" | cGreen ".*" | cBold ".*"
	else
		echo "One or more tests have failed" | cRed ".*" | cBold ".*"
	fi

}

goint function() {
	gotest Integration $@
}

gounit function() {
	gotest Unit $@
}
