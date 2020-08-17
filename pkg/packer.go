package kacker

import (
	"bufio"
	"io"
	"os"

	"github.com/mikefarah/yq/v3/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

// Packer collects the packer file and any YAML substitutions that
// need to be made.
type Packer struct {
	From  string      `yaml:"from"`
	Merge interface{} `yaml:"merge"`
}

var yq = yqlib.NewYqLib()

func (p *Packer) Resolve(out io.Writer) error {
	if !globalConf.VerboseLogging {
		logging.SetLevel(logging.ERROR, "yq")
	}

	fileNodes, err := getNodeContextFromFile(p.From)
	if err != nil {
		return err
	}
	fieldNodes, err := getNodeContextFromField(p.Merge)
	if err != nil {
		return err
	}
	mergeCommands := getMergeCommands(append(fileNodes, fieldNodes...))
	return yqlibUpdate(p.From, out, mergeCommands)
}

func (p *Packer) ResolveTempFile() (string, error) {
	return resolveToTempFile("./", "kacker-packer-*.yml", p.Resolve)
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
