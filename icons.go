package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type iconMap map[string]string

func parseIcons() iconMap {
	im := make(iconMap)

		defaultIcons := []string{
    "fi=📃",
    "di=📁",
    "tw=🤝",
    "ow=📂",
    "ln=⛓" ,
    "or=❌",
    "ex=🎯",
    "pi=|",
    "so=ﯲ ",
    "db= ",
    "cd=c",
    "su= ",
    "sg= ",
    "tw= ",
    "st= ",
    "*.sh= ",
    "*.c++= ",
    "*.c= ",
    "*.cp= ",
    "*.cpio= ",
    "*.cpp= ",
    "*.cs= ",
    "*.css= ",
    "*.cvs= ",
    "*.cxx= ",
    "*.cmd= ",
    "*.db= ",
    "*.exe=🎯",
    "*.conf=📃",
    "*.desktop=🖥️",
    "*.lnk=↩️",
    "*.txt=✍️ ",
    "*.mom=✍️ ",
    "*.me=✍️ ",
    "*.ms=✍️ ",
    "*.png=📷",
    "*.yuf=📷",
    "*.ico=📷",
    "*.jpg=📸",
    "*.jpe=📸",
    "*.jpeg=📸",
    "*.gif=🖼️",
    "*.webp=🖼️",
    "*.avif=🖼️",
    "*.svg=🗺️",
    "*.tif=📷",
    "*.tiff=📷",
    "*.xcf=🖌️",
    "*.html=🌎",
    "*.xml=📰",
    "*.sqlite=📰",
    "*.yml=📰",
    "*.json=📰",
    "*.gpg=🔒",
    "*.pub=🔐",
    "*.css=🎨",
    "*.pdf=📚",
    "*.djvu=📚",
    "*.epub=📚",
    "*.csv=📓",
    "*.xlsx=📓",
    "*.tex=📜",
    "*.md=📘",
    "*.r=📊",
    "*.R=📊",
    "*.rmd=📊",
    "*.Rmd=📊",
    "*.m=📊",
    "*.mp3=💿",
    "*.ogg=🎤",
    "*.opus=🎵",
    "*.cue=🎵",
    "*.m4a=🎵",
    "*.flac=🎼",
    "*.wav=🎼",
    "*.mkv=🎥",
    "*.mp4=🎬",
    "*.webm=🎥",
    "*.mpeg=🎥",
    "*.avi=🎥",
    "*.mov=🎥",
    "*.mpg=🎥",
    "*.wmv=🎥",
    "*.m4b=🎥",
    "*.flv=🎥",
    "*.webm=🎥",
    "*.rar=📦",
    "*.z=📦",
    "*.7z=📦",
    "*.zip=📦",
    "*.zoo=📦",
    "*.tar=📦",
    "*.gz=📦",
    "*.zip=📦",
    "*.z64=🎮",
    "*.v64=🎮",
    "*.n64=🎮",
    "*.gba=🎮",
    "*.nes=🎮",
    "*.gdi=🎮",
    "*.WAD=🕹️ ",
    "*.1=ℹ️ ",
    "*.nfo=ℹ️ ",
    "*.info=ℹ️ ",
    "*.log=📙",
    "*.iso=📀",
    "*.img=📀",
    "*.bib=🎓",
    "*.ged=👪",
    "*.part=💔",
    "*.torrent=🔽",
    "*.jar=🫙",
    "*.java=♨️ ",
    "*.fish=🐟",
	}

	im.parseEnv(strings.Join(defaultIcons, ":"))

	if env := os.Getenv("LF_ICONS"); env != "" {
		im.parseEnv(env)
	}

	for _, path := range gIconsPaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			im.parseFile(path)
		}
	}

	return im
}

func (im iconMap) parseFile(path string) {
	log.Printf("reading file: %s", path)

	f, err := os.Open(path)
	if err != nil {
		log.Printf("opening icons file: %s", err)
		return
	}
	defer f.Close()

	pairs, err := readPairs(f)
	if err != nil {
		log.Printf("reading icons file: %s", err)
		return
	}

	for _, pair := range pairs {
		key, val := pair[0], pair[1]

		key = replaceTilde(key)

		if filepath.IsAbs(key) {
			key = filepath.Clean(key)
		}

		im[key] = val
	}
}

func (im iconMap) parseEnv(env string) {
	for _, entry := range strings.Split(env, ":") {
		if entry == "" {
			continue
		}

		pair := strings.Split(entry, "=")

		if len(pair) != 2 {
			log.Printf("invalid $LF_ICONS entry: %s", entry)
			return
		}

		key, val := pair[0], pair[1]

		key = replaceTilde(key)

		if filepath.IsAbs(key) {
			key = filepath.Clean(key)
		}

		im[key] = val
	}
}

func (im iconMap) get(f *file) string {
	if val, ok := im[f.path]; ok {
		return val
	}

	if f.IsDir() {
		if val, ok := im[f.Name()+"/"]; ok {
			return val
		}
	}

	var key string

	switch {
	case f.linkState == working:
		key = "ln"
	case f.linkState == broken:
		key = "or"
	case f.IsDir() && f.Mode()&os.ModeSticky != 0 && f.Mode()&0002 != 0:
		key = "tw"
	case f.IsDir() && f.Mode()&0002 != 0:
		key = "ow"
	case f.IsDir() && f.Mode()&os.ModeSticky != 0:
		key = "st"
	case f.IsDir():
		key = "di"
	case f.Mode()&os.ModeNamedPipe != 0:
		key = "pi"
	case f.Mode()&os.ModeSocket != 0:
		key = "so"
	case f.Mode()&os.ModeDevice != 0:
		key = "bd"
	case f.Mode()&os.ModeCharDevice != 0:
		key = "cd"
	case f.Mode()&os.ModeSetuid != 0:
		key = "su"
	case f.Mode()&os.ModeSetgid != 0:
		key = "sg"
	case f.Mode()&0111 != 0:
		key = "ex"
	}

	if val, ok := im[key]; ok {
		return val
	}

	if val, ok := im[f.Name()+"*"]; ok {
		return val
	}

	if val, ok := im["*"+f.Name()]; ok {
		return val
	}

	if val, ok := im[filepath.Base(f.Name())+".*"]; ok {
		return val
	}

	if val, ok := im["*"+f.ext]; ok {
		return val
	}

	if val, ok := im["fi"]; ok {
		return val
	}

	return " "
}
