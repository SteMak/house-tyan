package mafia

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type module struct {
	session *discordgo.Session
	cmds    *dgutils.Discord

	running bool

	stopHandlers []func()

	config *mafiaConfig
	game   *Game

	roles sync.Map
}

func (*module) ID() string {
	return "mafia"
}

func (bot *module) IsRunning() bool {
	return bot.running
}

func (bot *module) Init(prefix string) error {
	out.Infoln("loading images...")
	bot.loadImages()
	out.Infoln("images loaded")

	bot.cmds = &dgutils.Discord{
		Prefix:   prefix,
		Commands: commands,
	}

	return nil
}

func (bot *module) Start(session *discordgo.Session) {
	bot.session = session

	bot.stopHandlers = []func(){
		bot.session.AddHandler(bot.mafiaJoinHandler),
		bot.session.AddHandler(bot.mafiaLeaveHandler),
		bot.session.AddHandler(bot.mafiaVoteHandler),
	}

	bot.cmds.Start(session)
	bot.running = true
}

func (bot *module) LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.WithStack(err)
	}

	err = yaml.Unmarshal(data, &bot.config)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (bot *module) SetLogger(logger *logrus.Logger) {
	log = logger
}

func (bot *module) Stop() {
	bot.cmds.Stop()

	for _, stopHandler := range bot.stopHandlers {
		stopHandler()
	}
}

func (bot *module) loadImages() {
	files, err := filepath.Glob(filepath.Join(bot.config.ImagesPath, "*"))
	if err != nil {
		out.Err(true, err)
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		card, err := os.Stat(file)
		if err != nil {
			out.Err(true, err)
		}

		if card.IsDir() {
			go bot.loadDir(file, &wg)
			continue
		}

		go bot.loadFile(file, &wg)
	}

	wg.Wait()
}

func (bot *module) loadFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	name := filepath.Base(path)
	ext := filepath.Ext(name)

	imgFile, err := os.Open(path)
	if err != nil {
		out.Err(true, err)
	}
	defer imgFile.Close()

	buf, err := ioutil.ReadAll(imgFile)
	if err != nil {
		out.Err(true, err)
	}

	bot.roles.Store(strings.TrimSuffix(name, ext), []*bytes.Buffer{bytes.NewBuffer(buf)})
}

func (bot *module) loadDir(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	name := filepath.Base(path)
	backgroundFile, err := os.Open(filepath.Join(path, "background.png"))
	if err != nil {
		out.Err(true, err)
	}
	defer backgroundFile.Close()

	background, err := png.Decode(backgroundFile)
	if err != nil {
		out.Err(true, err)
	}

	layerPaths, err := filepath.Glob(filepath.Join(path, "layers/*.png"))
	if err != nil {
		out.Err(true, err)
	}

	wg.Add(len(layerPaths))
	for _, layerPath := range layerPaths {
		go func(path string) {
			defer wg.Done()

			layerFile, err := os.Open(path)
			if err != nil {
				out.Err(true, err)
			}
			defer layerFile.Close()

			layer, err := png.Decode(layerFile)
			if err != nil {
				out.Err(true, err)
			}

			img := image.NewRGBA(background.Bounds())
			draw.Draw(img, background.Bounds(), background, image.ZP, draw.Src)
			draw.Draw(img, layer.Bounds(), layer, image.ZP, draw.Over)

			buf := new(bytes.Buffer)
			err = png.Encode(buf, img)
			if err != nil {
				out.Err(true, err)
			}

			if role, ok := bot.roles.Load(name); ok {
				bot.roles.Store(name, append(role.([]*bytes.Buffer), buf))
				return
			}
			bot.roles.Store(name, []*bytes.Buffer{buf})
		}(layerPath)
	}
}
