.PHONY: lint test live alfredworkflow clean

WORKFLOW_FILES := info.plist icon.png timein.js capitals.json workflow/screenshot.png
OUT ?= TimeIn.alfredworkflow

lint: # ES2020 is the syntax ceiling for older macOS JavaScriptCore
	npx --yes acorn@8 --ecma2020 --allow-hash-bang --silent timein.js

test:
	test/unit.js
	test/run.sh

live:
	test/run.sh live

alfredworkflow:
	rm -f "$(OUT)"
	zip -j "$(OUT)" $(WORKFLOW_FILES)

clean:
	rm -f *.alfredworkflow
