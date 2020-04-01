package messages

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/SteMak/house-tyan/config"

	"github.com/bwmarrin/discordgo"
)

var (
	tpls = make(map[string]*template.Template)
)

func AddTpl(f string) error {
	file, err := os.Open(f)
	if err != err {
		return err
	}
	defer file.Close()

	var data shema
	err = xml.NewDecoder(file).Decode(&data)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	d = normalizeSpaces(d)

	name, err := filepath.Rel(config.Bot.Templates, f)
	if err != nil {
		return err
	}

	tpl, err := template.New(name).Funcs(funcs).Parse(string(d))
	if err != nil {
		return err
	}

	tpls[name] = tpl
	return nil
}

func Get(name string, data interface{}) (*discordgo.MessageSend, error) {
	tpl, ok := tpls[name]
	if !ok {
		return nil, fmt.Errorf("message '%s' no found", name)
	}

	buf := bytes.NewBufferString("")
	err := tpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		return nil, err
	}

	s := bytes.NewBufferString("")
	err = xml.EscapeText(s, buf.Bytes())
	if err != nil {
		return nil, err
	}

	var m shema

	err = xml.NewDecoder(buf).Decode(&m)
	if err != nil {
		return nil, err
	}

	return buildMessage(&m)
}
