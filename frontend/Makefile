
include ./etc/help.mk

export UID=$(shell id -u)
export GID=$(shell id -g)

start: ##@docker Start the frontend
	touch dist/main.js
	chmod o+w dist/main.js
	mkdir -p elm-stuff
	docker-compose up --detach --build
.PHONY: start

stop: ##@docker Stop the frontend
	docker-compose stop --timeout 1
.PHONY: stop

clean: ##@docker Stop and remove the frontend
	docker-compose down --remove-orphans --volumes --timeout 1 
.PHONY: clean

logs: ##@docker Display the frontends logs
	docker-compose logs --follow 
.PHONY: logs

cli: ##@docker Access cli
	docker-compose exec frontend bash



run: ##@host start the frontend on port 8000
	elm-live src/Main.elm --dir=./dist --pushstate -- --output=./dist/main.js

open: ##@host start and open the frontend
	elm-live src/Main.elm --dir=./dist --pushstate --open -- --output=./dist/main.js

# https://github.com/rtfeldman/elm-spa-example
# $ elm make src/Main.elm --output elm.js --optimize
# $ uglifyjs elm.js --compress 'pure_funcs="F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9",pure_getters=true,keep_fargs=false,unsafe_comps=true,unsafe=true,passes=2' --output=elm.js && uglifyjs elm.js --mangle --output=elm.js


elmfmt: ##@qa format code and remove unused imports
	find src | grep \.elm\$$ | xargs -L 1 elm-impfix --replace
	find src | grep \.elm\$$ | xargs elm-format --yes

analyse: ##@qa run static code analyser
	elm-analyse . -s
