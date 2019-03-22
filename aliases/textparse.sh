# Formats stdin to pretty JSON
alias jsonp="python -m json.tool"

alias stripnl="tr -d '\r' | tr -d '\n'"

# Replace words with another
# $1: Word or regex to replace
# $2: Replacement word
function replace() {
        sed -r "s|$1|$2|g"
}


