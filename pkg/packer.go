package kacker

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/mikefarah/yq/v3/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

// Packer collects the packer file and any YAML substitutions that
// need to be made.
type Packer struct {
	From  string      `yaml:"from"`
	Force bool        `yaml:"force"`
	Merge interface{} `yaml:"merge"`
}

// Customization represents the customization file: modifiers for the kickstart
// file and the packer file.
type Customization struct {
	Kickstart Kickstart `yaml:"kickstart"`
	Packer    Packer    `yaml:"packer"`
}

func (c *Customization) Run(packerFlags []string) error {
	kickstart, err := c.Kickstart.ResolveTempFile()
	if err != nil {
		return err
	}
	os.Setenv("KICKSTART", kickstart)
	if !globalConf.KeepTempFiles {
		defer os.Remove(kickstart)
	}

	packer, err := c.Packer.ResolveTempFile()
	if err != nil {
		return err
	}
	if !globalConf.KeepTempFiles {
		defer os.Remove(packer)
	}

	var commands []*exec.Cmd
	commands = append(commands, exec.Command("packer",
		append(append([]string{"validate"}, packerFlags...), packer)...))
	if c.Packer.Force {
		packerFlags = append(packerFlags, "-force")
	}
	commands = append(commands, exec.Command("packer",
		append(append([]string{"build"}, packerFlags...), packer)...))

	if err = runCommands(commands); err != nil {
		return err
	}
	return nil
}

func Run(filename string, packerFlags []string) error {
	if !HasPacker() {
		return fmt.Errorf("Packer is not installed or is not in your PATH")
	}
	var cust Customization
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(f, &cust); err != nil {
		return err
	}
	return cust.Run(packerFlags)
}

var yq = yqlib.NewYqLib()

func (p *Packer) Resolve(out io.Writer) error {
	logging.SetLevel(logging.ERROR, "yq")

	from := path.Join(globalConf.RelativeDir, p.From)
	fileNodes, err := getNodeContextFromFile(from)
	if err != nil {
		return err
	}
	fieldNodes, err := getNodeContextFromField(p.Merge)
	if err != nil {
		return err
	}
	mergeCommands := getMergeCommands(append(fileNodes, fieldNodes...))
	return yqlibUpdate(from, out, mergeCommands)
}

func (p *Packer) ResolveTempFile() (string, error) {
	return resolveToTempFile("./", ".kacker-packer-*.yml", p.Resolve)
}

func getNodeContextFromFile(filename string) ([]*yqlib.NodeContext, error) {
	var yqlibNodes []*yqlib.NodeContext
	return yqlibNodes, decodeAndExecute(filename, func(node *yaml.Node) error {
		yqlibNode, nodeErr := yq.Get(node, "**", true)
		if nodeErr != nil {
			return nodeErr
		}
		yqlibNodes = append(yqlibNodes, yqlibNode...)
		return nil
	})
}

func yqlibUpdate(filename string, out io.Writer, commands []yqlib.UpdateCommand) error {
	writer := bufio.NewWriter(out)
	defer writer.Flush()
	encoder := yqlib.NewJsonEncoder(writer, true, 2)
	return decodeAndExecute(filename, func(node *yaml.Node) error {
		for _, update := range commands {
			updateErr := yq.Update(node, update, true)
			if updateErr != nil {
				return updateErr
			}
		}
		encodeErr := encoder.Encode(node)
		if encodeErr != nil {
			return encodeErr
		}
		return nil
	})
}

func decodeAndExecute(filename string, fn func(*yaml.Node) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	if err != nil {
		return err
	}

	for {
		var yamlNode yaml.Node
		decodeErr := decoder.Decode(&yamlNode)
		if decodeErr == io.EOF {
			return nil
		} else if decodeErr != nil {
			return decodeErr
		}
		if fnErr := fn(&yamlNode); fnErr != nil {
			return fnErr
		}
	}
}

func getNodeContextFromField(field interface{}) ([]*yqlib.NodeContext, error) {
	var yamlNode yaml.Node
	if err := yamlNode.Encode(field); err != nil {
		return nil, err
	}
	return yq.Get(&yamlNode, "**", true)
}

func getMergeCommands(nodes []*yqlib.NodeContext) []yqlib.UpdateCommand {
	var mergeCommands []yqlib.UpdateCommand = []yqlib.UpdateCommand{}
	for _, node := range nodes {
		mergePath := yq.MergePathStackToString(node.PathStack, false)
		mergeCommands = append(mergeCommands, yqlib.UpdateCommand{
			Command:   "update",
			Path:      mergePath,
			Value:     node.Node,
			Overwrite: true,
		})
	}
	return mergeCommands
}
