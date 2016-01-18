build:
	docker run -it --rm -v $(PWD):/go/src/github.com/goodguide/goodguide-git-hooks -w /go/src/github.com/goodguide/goodguide-git-hooks multiarch/goxc goxc -env=GO15VENDOREXPERIMENT=1

bump:
	goxc bump
	echo >> .goxc.json # put a newline at the end of the file because goxc fails to do so
	git commit -m 'Bump version [nostory]' -- .goxc.json
	git push
	goxc tag
	git push --tags origin
