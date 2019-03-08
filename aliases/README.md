Add this line to your rc file to load all the aliases:

```bash
for f in <PATH_TO_ALIASES_DIR>/*.sh; do
        source $f
done
```

or add individual ones if you prefer.