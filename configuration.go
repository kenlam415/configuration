package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-akka/configuration/hocon"
)

func ParseString(text string, includeCallback ...hocon.IncludeCallback) *Config {
	var callback hocon.IncludeCallback
	if len(includeCallback) > 0 {
		callback = includeCallback[0]
	} else {
		callback = defaultIncludeCallback
	}
	root := hocon.Parse(text, callback)
	return NewConfigFromRoot(root)
}

func checkFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// func directoryFallBack(filename string) string {
// 	splitDirAndFile := strings.Split(filename, "../")
// 	dirLevel := len(splitDirAndFile)
// 	newfile := func() string {
// 		switch {
// 			case dirLevel > 2:
// 				return strings.Join(splitDirAndFile[1:], "../")
// 			case dirLevel == 2:
// 				return strings.Join(splitDirAndFile, "./")
// 			default:
// 				return filename
// 			}
// 	}
// 	fileExists := checkFileExists(newfile())
// 	if (fileExists(newfile())) == true {
// 		return newfile()
// 	} else {
// 		return filename
// 	}

// }

func directoryFallBack(filename string) string {
	splitDirAndFile := strings.Split(filename, "../")
	dirLevel := len(splitDirAndFile)
	switch {
	case dirLevel > 2:
		return strings.Join(splitDirAndFile[1:], "../")
	case dirLevel == 2:
		return strings.Join(splitDirAndFile, "./")
	default:
		return filename
	}
}

func LoadConfig(filename string, upperDirCheck bool) *Config {
	if upperDirCheck == true {
		newFile := directoryFallBack(filename)
		if checkFileExists(newFile) == false {
			log.Println("Cannot Found the Configuration file - %s", newFile)
			log.Println("Fall back to orginal file directory")
		} else {
			filename = newFile
		}
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return ParseString(string(data), defaultIncludeCallback)
}

func FromObject(obj interface{}) *Config {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return ParseString(string(data), defaultIncludeCallback)
}

func defaultIncludeCallback(filename string) *hocon.HoconRoot {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ignore the include file: ", filename)
	}

	return hocon.Parse(string(data), defaultIncludeCallback)
}
