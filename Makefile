CONVERT=./convert.py
PACKER=packer
ARGS=-var-file=./secrets.json

echo:
	cat $(TARGET) | $(CONVERT)

validate:
	cat $(TARGET) | $(CONVERT) | $(PACKER) validate $(ARGS) $(OVERRIDE) -

build:
	cat $(TARGET) | $(CONVERT) | $(PACKER) build $(ARGS) $(OVERRIDE) -

debug:
	cat $(TARGET) | $(CONVERT) | $(PACKER) build -debug $(ARGS) $(OVERRIDE) -
