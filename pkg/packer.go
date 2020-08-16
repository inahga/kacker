package kacker

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/mikefarah/yq/v3/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

// PackerCustomization collects the packer file and any YAML substitutions that
// need to be made.
type PackerCustomization struct {
	From  string      `yaml:"from"`
	Merge interface{} `yaml:"merge"`
}

var yq = yqlib.NewYqLib()

func (p *PackerCustomization) readPackerFile(name string) (interface{}, error) {
	var ret interface{}

	f, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(f, &ret); err != nil {
		return nil, err
	}

	return nil, nil
}

type readDataFn func(dataBucket *yaml.Node) ([]*yqlib.NodeContext, error)
type yamlDecoderFn func(*yaml.Decoder) error
type updateDataFn func(dataBucket *yaml.Node, currentIndex int) error

func ReadYamlFile(filename string) ([]*yqlib.NodeContext, error) {
	var yqlibNodes []*yqlib.NodeContext
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err != nil {
		return nil, err
	}

	for {
		var yamlNode yaml.Node
		decodeErr := decoder.Decode(&yamlNode)
		if decodeErr == io.EOF {
			return yqlibNodes, nil
		} else if decodeErr != nil {
			return nil, err
		}

		yqlibNode, nodeErr := yq.Get(&yamlNode, "**", true)
		if nodeErr != nil {
			return nil, err
		}
		yqlibNodes = append(yqlibNodes, yqlibNode...)
	}
}

func ReadYamlField(field interface{}) ([]*yqlib.NodeContext, error) {
	var yamlNode yaml.Node
	if err := yamlNode.Encode(field); err != nil {
		return nil, err
	}
	return yq.Get(&yamlNode, "**", true)
}

func GetMergeCommands(nodes []*yqlib.NodeContext) []yqlib.UpdateCommand {
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

func ExecuteMergeCommands(filename string, out io.Writer, commands []yqlib.UpdateCommand) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	encoder := yaml.NewEncoder(out)
	defer encoder.Close()

	index := 0
	for {
		var yamlNode yaml.Node
		decodeErr := decoder.Decode(&yamlNode)
		if decodeErr == io.EOF {
			return nil
		} else if decodeErr != nil {
			return err
		}

		updateErr := yq.Update(&yamlNode, commands[index], false)
		if updateErr != nil {
			return nil
		}

		fmt.Printf("%s: %s\n", yamlNode.Tag, yamlNode.Value)
	}
}
