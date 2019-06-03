
# Pretty print mount
function mount() {
	/usr/bin/mount $@ | column -t
}

function ports() {
	netstat -tulpe | cGreen $USER | cLightRed root | cCyan LISTEN | cYellow 'tcp6?' | cPink 'udp6?' | cBold "Proto.*name"
}
