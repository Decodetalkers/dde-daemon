package launcher

import (
	"pkg.linuxdeepin.com/lib/glib-2.0"
	"fmt"
	"os"
	"strings"
)

const (
	FavorConfigFile  string = "launcher/favor.ini"
	FavorConfigGroup string = "FavorConfig"
	FavorConfigKey   string = "Ids"
	IndexKey         string = "Index"
	FixedKey         string = "Fixed"
)

// TODO: add a signal for content changed.

type FavorItem struct {
	Id    string
	Index int64
	Fixed bool
}

type FavorItemList []FavorItem

func defaultConfigPath() string {
	languange := os.Getenv("LANGUAGE")
	var defaultPath string
	if strings.HasPrefix(languange, "zh") {
		defaultPath = "/usr/share/dde/data/config/launcher/zh/favor.ini"
	} else {
		defaultPath = "/usr/share/dde/data/config/launcher/en/favor.ini"
	}
	return defaultPath
}

func getFavorIdList(file *glib.KeyFile) []string {
	_, list, err := file.GetStringList(FavorConfigGroup, FavorConfigKey)
	if err != nil {
		logger.Info(fmt.Errorf("getFavorIdList: %s", err))
		return make([]string, 0)
	}
	return uniqueStringList(list)
}

func getFavors() FavorItemList {
	favors := make(FavorItemList, 0)
	file, err := configFile(FavorConfigFile, defaultConfigPath())
	defer file.Free()
	if err != nil {
		logger.Info(fmt.Errorf("getFavors: %s", err))
		return favors
	}

	ids := getFavorIdList(file)
	for _, id := range ids {
		fixed, err := file.GetBoolean(id, FixedKey)
		if err != nil {
			continue
		}
		index, err := file.GetInt64(id, IndexKey)
		if err != nil {
			continue
		}

		favors = append(favors, FavorItem{id, index, fixed})
	}

	return favors
}

func saveFavors(items FavorItemList) bool {
	file, err := configFile(FavorConfigFile, defaultConfigPath())
	defer file.Free()
	if err != nil {
		logger.Info(fmt.Errorf("saveFavors: %s", err))
		return false
	}

	previousIds := getFavorIdList(file)
	previousIdMap := make(map[string]bool, 0)
	for _, id := range previousIds {
		previousIdMap[id] = true
	}

	ids := make([]string, 0)
	itemMap := make(map[string]FavorItem, 0)
	for _, item := range items {
		itemMap[item.Id] = item
		ids = append(ids, item.Id)
	}

	ids = uniqueStringList(ids)
	file.SetStringList(FavorConfigGroup, FavorConfigKey, ids)
	for id, item := range itemMap {
		file.SetBoolean(id, FixedKey, item.Fixed)
		file.SetInt64(id, IndexKey, item.Index)
	}

	for id, _ := range previousIdMap {
		if _, ok := itemMap[id]; !ok {
			file.RemoveGroup(id)
		}
	}

	err = saveKeyFile(file, configFilePath(FavorConfigFile))
	if err != nil {
		logger.Info(err)
		return false
	}
	return true
}
