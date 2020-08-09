package kacker

// Customization represents the customization file: modifiers for the kickstart
// file and the packer file.
type Customization struct {
	Kickstart KickstartCustomization `yaml:"kickstart"`
	Packer    PackerCustomization    `yaml:"packer"`
}
