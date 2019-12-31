/*
Package i18n provides a simple utility to support I18N.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r1
*/
package i18n

import (
	"fmt"
	hocon "github.com/go-akka/configuration"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultLocale = "en-us"

func NewI18n(dir string) *I18n {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	log.Printf("Loading i18n files from directory [%s]", dir)

	textFromAllFiles := make([]string, 0)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".i18n" {
			return nil
		}
		log.Printf("\tLoading i18n file [%s]...", path)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		textFromAllFiles = append(textFromAllFiles, string(data))
		return nil
	})
	if err != nil {
		panic(err)
	}
	i18n := &I18n{config: hocon.ParseString(strings.Join(textFromAllFiles, "\n"))}
	i18n.locale = defaultLocale
	return i18n
}

// TODO: change locale?
type I18n struct {
	config *hocon.Config
	locale string
}

func (i18n *I18n) Text(path string, params ...interface{}) string {
	if i18n.locale != "" {
		path = i18n.locale + ".text." + path
	}
	format := i18n.config.GetString(path, path)
	return fmt.Sprintf(format, params...)
}
