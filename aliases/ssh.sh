
# Inject bash aliases into a remote server before SSHing
# $1 : SSH host. Ex: root@192.168.1.120
function myssh() {
	HOST=$1
	ALIAS_DIR=${ALIAS_DIR:-$HOME/lobocv/mysetup/aliases}
	stty -echo
	printf "Password: "
	read PW
	stty echo
	printf "\n"
	sshpass -p "$PW" rsync -a --ignore-existing $ALIAS_DIR $HOST:/tmp
	sshpass -p "$PW" ssh $HOST <<zzz
if ! grep 'for f in /tmp/aliases/\*.sh; do source \$f; done' \$HOME/.bashrc; then
echo 'for f in /tmp/aliases/*.sh; do source \$f; done' >> \$HOME/.bashrc
fi
zzz

	sshpass -p "$PW" ssh $HOST
}
