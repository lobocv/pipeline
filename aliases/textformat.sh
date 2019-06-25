# Text formatting tools for terminal output

DEFAULT_HL_COLOR=93
# Highlight text
function _highlight() {
	CODE=$(printf "\033[%qm" "$1")
	RESET=$(printf "\033[%qm" "$2")
	MATCH="$3"

	while read TEXT;
	do
		# & acts as the match group in sed
		echo "$TEXT" | sed -r "s|"$MATCH"|${CODE}&${RESET}|g"
	done
}

# Reads from stdin and removes a line matching the regex
# $1 : Regex
function removeline() {
	MATCH="$1"
	while read TEXT;
	do
		# & acts as the match group in sed
		fmtline=$(echo "$TEXT" | sed -r "s|"$MATCH"||g")
		if [[ ! -z "$fmtline" ]]; then
			echo "$fmtline"
		fi

	done

}
alias hl="_highlight $DEFAULT_HL_COLOR 39"

# Colors
alias cBlack="_highlight 30 39"
alias cLightRed="_highlight 91 39"
alias cRed="_highlight 31 39"
alias cLightBlue="_highlight 94 39"
alias cBlue="_highlight 34 39"
alias cLightGreen="_highlight 92 39"
alias cGreen="_highlight 32 39"
alias cPink="_highlight 95 39"
alias cBrown="_highlight 33 39"
alias cMag="_highlight 35 39"
alias cWhite="_highlight 0 39"
alias cCyan="_highlight 36 39"
alias cGrey="_highlight 90 39"
alias cYellow="_highlight 93 39"

# Styles
alias cBold="_highlight 1 21"
alias cItalic="_highlight 3 24"
alias cUnderline="_highlight 4 24"
alias cBlink="_highlight 5 25"
alias cHighlight="_highlight 7 27"
alias cHide="_highlight 8 28"
alias cStrike="_highlight 9 0"

# Right justify stdin 
rjust() {
    read TEXT
   
    padlimit=$(tput cols)
    pad=$(printf '%*s' "$padlimit")
    pad=${pad// / }
    padlength=$(tput cols)
    
    printf '%*.*s' 0 $((padlength - ${#TEXT})) "$pad"
    printf '%s' "$TEXT"
}

