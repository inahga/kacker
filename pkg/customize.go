package kacker

// Customization represents the customization file: modifiers for the kickstart
// file and the packer file.
type Customization struct {
	Kickstart KickstartCustomization `yaml:"kickstart"`
	Packer    PackerCustomization    `yaml:"packer"`
}

// Configuration sets the behavior of this module.
type Configuration struct {
	VerboseLogging bool
	UseInsecureSSL bool
}

var globalConf Configuration

// UseConfiguration sets the behavior of this module.
func UseConfiguration(conf *Configuration) {
	globalConf = *conf
}
