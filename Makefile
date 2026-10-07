.PHONY: test smoke alfredworkflow clean

WORKFLOW_FILES := info.plist icon.png timein.js capitals.json workflow/screenshot.png

test:
	test/run.sh

smoke:
	test/run.sh --live

alfredworkflow: clean
	zip -j TimeIn.alfredworkflow $(WORKFLOW_FILES)

clean:
	rm -f TimeIn.alfredworkflow
