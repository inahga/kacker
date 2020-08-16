package kacker

// PackerCustomization collects the packer file and any YAML substitutions that
// need to be made.
type PackerCustomization struct {
	From  string      `yaml:"from"`
	Merge interface{} `yaml:"merge"`
}
