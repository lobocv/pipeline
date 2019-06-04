#!/bin/bash

# Run tests with coverage and create a colorized report
cover () {
    local t=$(mktemp -t coverXXXX)
    go test $COVERFLAGS -coverprofile=$t $@ \
        && go tool cover -func=$t | cRed '\s[0-9]\.?[0-9]?\%' | cLightRed '[1-5][0-9]\.[0-9]\%' | cLightGreen '[6-9][0-9]\.[0-9]\%'| cGreen '9[0-9]\.[0-9]\%' | cGreen '100\.?0?\%' 
    unlink $t   
}

