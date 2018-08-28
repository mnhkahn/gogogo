include .env

d := $(shell date)
branch := $(shell git branch | grep \* | cut -d ' ' -f2)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in ${BINARY}:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## pull: Revert & pull code from remote by current branch.
pull:
	git reset --hard
	git pull origin ${branch}

## push: Auto commit & push code to remote by current branch.
push: version
	git commit -am "chore: date is $d."
	git push origin ${branch}

## cover: Run all test code & print cover profile.
cover:
	go test -coverprofile=/tmp/${BINARY}-makecover.out ${TEST_PACKAGES} && go tool cover -html=/tmp/${BINARY}-makecover.out

## version: Show current code version.
version:
	@git remote -v
	@echo branch: $(branch)

## chglog: Run git-chglog & save to CHANGELOG
chglog:
	@git-chglog -o CHANGELOG.md
	git commit -am "add chglog: date is $d."
	git push origin ${branch}

