function remove_pyc() {
	num_pyc=$(find . -name '*.pyc'| wc -l)
	find . -name '*.pyc' -delete
	echo "Deleted $num_pyc .pyc files."
}
