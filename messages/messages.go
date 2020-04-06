package messages

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/SteMak/house-tyan/config"

	"github.com/bwmarrin/discordgo"
)

var (
	tpls = make(map[string]*template.Template)
)

type Message struct {
	discordgo.MessageSend
	Reactions []string
}

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

func Get(name string, data interface{}) (*Message, error) {
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

func Random(pattern string) (string, error) {
	matcher, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	var finded []string
	var wg sync.WaitGroup

	for s := range tpls {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			if matcher.MatchString(s) {
				finded = append(finded, s)
			}
		}(s)
	}

	wg.Wait()

	if len(finded) == 0 {
		return "", errors.New("not found")
	}

	return finded[rand.Intn(len(finded))], nil
}
