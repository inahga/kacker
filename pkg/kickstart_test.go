package kacker

import "testing"

func TestVariableIsValid(t *testing.T) {
	validVariables := map[string]Variable{
		"Valid": Variable{
			Name:  "test",
			Value: "test",
		},
		"Valid array": Variable{
			Name: "test",
			Values: []string{
				"test",
				"test",
			},
		},
		"Valid blanks": Variable{
			Name:      "test",
			Value:     "test",
			Values:    []string{},
			URL:       "",
			Fragment:  "",
			Fragments: []string{},
		},
	}

	invalidVariables := map[string]Variable{
		"Missing name": Variable{
			Name: "test",
		},
		"Too many parameters": Variable{
			Name:  "test",
			Value: "test",
			URL:   "test",
		},
		"Too many parameters, one of them is slice": Variable{
			Name:  "test",
			Value: "test",
			Values: []string{
				"test",
				"test",
			},
		},
	}

	for key, elem := range validVariables {
		t.Run(key, func(t *testing.T) {
			_, err := elem.isValid()
			if err != nil {
				t.Errorf("%+v should be valid, got err %s", elem, err)
			}
		})
	}

	for key, elem := range invalidVariables {
		t.Run(key, func(t *testing.T) {
			_, err := elem.isValid()
			if err == nil {
				t.Errorf("%+v should be invalid", elem)
			}
		})
	}
}

func TestResolveVariables(t *testing.T) {
	variables := []Variable{
		Variable{
			Name:  "test",
			Value: "test",
		},
		Variable{
			Name:  "test2",
			Value: "test2",
		},
		Variable{
			Name: "test3",
			Values: []string{
				"test",
				"test",
			},
		},
		Variable{
			Name: "test4",
			Values: []string{
				"test",
				"test",
			},
		},
	}

	redefinedVariables := []Variable{
		Variable{
			Name:  "test",
			Value: "test",
		},
		Variable{
			Name:  "test",
			Value: "test2",
		},
	}

	_, err := resolveVariables(variables)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	_, err = resolveVariables(redefinedVariables)
	if err == nil {
		t.Errorf("Expected redefined variable error")
	}
}
