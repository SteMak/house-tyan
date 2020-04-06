package messages

import (
	"bytes"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"

	"github.com/bwmarrin/discordgo"
)

var (
	tpls *template.Template
)

type Message struct {
	discordgo.MessageSend
	Reactions []string
}

func Init() {
	tpls = template.New("").Funcs(funcs)

	out.Infoln("Loading templates...")
	err := filepath.Walk(config.Bot.Templates, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name, err := filepath.Rel(config.Bot.Templates, path)
		if err != nil {
			return err
		}

		_, err = tpls.New(name).Parse(string(data))
		return err
	})

	if err != nil {
		out.Fatal(err)
	}

	for _, tpl := range tpls.Templates() {
		out.Infoln(tpl.Name())
	}
}

func Get(name string, data interface{}) (*Message, error) {
	buf := bytes.NewBufferString("")
	err := tpls.ExecuteTemplate(buf, name, data)
	if err != nil {
		return nil, err
	}

	var m shema
	s := bytes.NewBuffer(normalizeSpaces(buf.Bytes()))
	out.Debugln(s.String())
	if err := xml.NewDecoder(s).Decode(&m); err != nil {
		return nil, err
	}

	return buildMessage(&m)
}
