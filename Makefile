build:
	goxc

bump:
	goxc bump
	echo >> .goxc.json # put a newline at the end of the file because goxc fails to do so
	git commit -m 'Bump version [nostory]' -- .goxc.json
	git push
	goxc tag
	git push --tags origin
