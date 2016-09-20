build:
	docker run -it --rm -v $(PWD):/go/src/github.com/goodguide/goodguide-git-hooks -w /go/src/github.com/goodguide/goodguide-git-hooks multiarch/goxc bash -c 'go get -v . && goxc -env=GO15VENDOREXPERIMENT=1'

bump:
	goxc bump
	echo >> .goxc.json # put a newline at the end of the file because goxc fails to do so
	git commit -m 'Bump version [nostory]' -- .goxc.json
	git push
	goxc tag
	git push --tags origin

sign:
	bash -x -c 'find dist/$$(jq -r .PackageVersion < .goxc.json) -name "*.tar.gz" | xargs -n1 gpg -a -b'
