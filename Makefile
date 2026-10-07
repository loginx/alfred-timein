.PHONY: test smoke alfredworkflow clean

WORKFLOW_FILES := info.plist icon.png timein.js capitals.json workflow/screenshot.png
OUT ?= TimeIn.alfredworkflow

test:
	test/run.sh

smoke:
	test/run.sh --live

alfredworkflow: clean
	zip -j "$(OUT)" $(WORKFLOW_FILES)

clean:
	rm -f *.alfredworkflow
