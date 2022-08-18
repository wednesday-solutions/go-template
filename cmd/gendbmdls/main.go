package main

import (
	"fmt"
	"go-template/internal/config"
	"os"
	"strings"

	"github.com/spf13/viper"
	boilingcore "github.com/volatiletech/sqlboiler/v4/boilingcore"
	"github.com/volatiletech/sqlboiler/v4/drivers"
	importers "github.com/volatiletech/sqlboiler/v4/importers"
)

func allKeys(prefix string) []string {
	keys := make(map[string]bool)

	prefix += "."

	for _, e := range os.Environ() {
		splits := strings.SplitN(e, "=", 2)
		key := strings.ReplaceAll(strings.ToLower(splits[0]), "_", ".")

		if strings.HasPrefix(key, prefix) {
			keys[strings.ReplaceAll(key, prefix, "")] = true
		}
	}

	for _, key := range viper.AllKeys() {
		if strings.HasPrefix(key, prefix) {
			keys[strings.ReplaceAll(key, prefix, "")] = true
		}
	}

	keySlice := make([]string, 0, len(keys))
	for k := range keys {
		keySlice = append(keySlice, k)
	}
	return keySlice
}

func main() {
	err := config.LoadEnv()
	if err != nil {
		fmt.Println("failed while loading env")
	}
	driverName, _, err := drivers.RegisterBinaryFromCmdArg("psql")

	keys := allKeys(driverName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var cfg = drivers.Config{}
	for _, key := range keys {
		if key != "blacklist" && key != "whitelist" {
			prefixedKey := fmt.Sprintf("%s.%s", driverName, key)
			cfg[key] = viper.Get(prefixedKey)
		}
	}
	if err != nil {
		return
	}

	cmdState, err := boilingcore.New(&boilingcore.Config{
		DriverName:      driverName,
		DriverConfig:    cfg,
		NoHooks:         true,
		OutFolder:       "models",
		PkgName:         "models",
		StructTagCasing: "snake",
		RelationTag:     "-",
		Imports:         configureImports(),
	})
	if err != nil {
		panic(err)
	}
	err = cmdState.Run()
	if err != nil {
		panic(err)
	}
	err = cmdState.Cleanup()
	if err != nil {
		panic(err)
	}
}

func configureImports() importers.Collection {
	imports := importers.NewDefaultImports()

	mustMap := func(m importers.Map, err error) importers.Map {
		if err != nil {
			panic("failed to change viper interface into importers.Map: " + err.Error())
		}

		return m
	}

	if viper.IsSet("imports.all.standard") {
		imports.All.Standard = viper.GetStringSlice("imports.all.standard")
	}
	if viper.IsSet("imports.all.third_party") {
		imports.All.ThirdParty = viper.GetStringSlice("imports.all.third_party")
	}
	if viper.IsSet("imports.test.standard") {
		imports.Test.Standard = viper.GetStringSlice("imports.test.standard")
	}
	if viper.IsSet("imports.test.third_party") {
		imports.Test.ThirdParty = viper.GetStringSlice("imports.test.third_party")
	}
	if viper.IsSet("imports.singleton") {
		imports.Singleton = mustMap(importers.MapFromInterface(viper.Get("imports.singleton")))
	}
	if viper.IsSet("imports.test_singleton") {
		imports.TestSingleton = mustMap(importers.MapFromInterface(viper.Get("imports.test_singleton")))
	}
	if viper.IsSet("imports.based_on_type") {
		imports.BasedOnType = mustMap(importers.MapFromInterface(viper.Get("imports.based_on_type")))
	}

	return imports
}
