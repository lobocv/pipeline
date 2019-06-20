
# View PEM encoded certificate
# $1 : Path to PEM encoded certificate
function readcert() {
	openssl x509 -text -noout -in "$1" | less
}
