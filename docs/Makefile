all: build/openapi.yaml
	redoc-cli bundle build/openapi.yaml

define appendFilesWithTabs
	for file in $(1); do \
		while IFS= read -r line || [ -n "$$line" ]; do \
			echo "$(2)$$line"; \
		done < $$file; \
	done >> $3
endef

build/models.yaml: ./models/*.yaml
	echo "components:" > $@
	echo "  schemas:" >> $@
	$(call appendFilesWithTabs, $^,""    "", $@)

build/paths.yaml: ./paths/*.yaml
	echo "paths:" > $@
	$(call appendFilesWithTabs, $^,""  "", $@)

build/openapi.yaml: index.yaml build/paths.yaml build/models.yaml
	echo "#This file has been generated automatically. DO NOT MODIFY!" > $@
	for file in $^; do \
		while IFS= read -r line; do \
			echo "$$line"; \
		done < $$file; \
	done >> $@

clean:
	rm -rf build/*
