package dgutils

import (
	"reflect"
	"testing"
)

func TestFindCommand(t *testing.T) {
	tests := make(map[string][]string)

	// Обычные команды
	tests["test"] = []string(nil)
	tests["test foo boo"] = []string{"foo", "boo"}

	// Обычные команды в группе
	tests["group test"] = []string(nil)
	tests["group test foo boo"] = []string{"foo", "boo"}

	// Обычные команды
	tests["raw_test"] = []string(nil)
	tests["raw_test foo boo"] = []string{"foo boo"}

	// Обычные команды в группе
	tests["group raw_test"] = []string(nil)
	tests["group raw_test foo boo"] = []string{"foo boo"}

	// Не существующая команда
	tests["usage foo boo"] = nil

	// вызов группы
	tests["group"] = nil

	commands := map[string]interface{}{
		"test":     &Command{},
		"raw_test": &Command{Raw: true},
		"group": &Group{
			Commands: map[string]interface{}{
				"test":     &Command{},
				"raw_test": &Command{Raw: true},
			},
		},
	}

	for key, value := range tests {
		_, args := findCommand(key, commands)
		if !reflect.DeepEqual(args, value) {
			t.Errorf("[\"%s\"] error\n\tgot: %#v\n\twant: %#v.", key, args, value)
		} else {
			t.Logf("[\"%s\"] ok", key)
		}
	}
}
