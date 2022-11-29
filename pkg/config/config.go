// Package config defines a manager to manage common configuration.
package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Manager is the configuration manager.
type Manager struct {
	name  string
	flags *pflag.FlagSet
	viper *viper.Viper
}

// New returns a new initialized manager with the given config.
func New(name string) *Manager {
	viper := viper.New()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AllowEmptyEnv(true)
	viper.SetEnvPrefix(name)
	viper.SetConfigName(name)
	viper.AutomaticEnv()

	return &Manager{
		name:  name,
		flags: pflag.NewFlagSet(name, pflag.ExitOnError),
		viper: viper,
	}
}

// AddConfigFlag adds the config flag to the persistent root flagset.
func (mgr *Manager) AddConfigFlag(persistentFlags *pflag.FlagSet) {
	persistentFlags.StringP("config", "c", "./config.yml", "config file (Default; config.yml in the current directory)")
}

// ReadFromFile reads the config from the file set with the `config` flag.
// Use `AddConfigFlag` to add the config flag to the root flagset.
func (mgr *Manager) ReadFromFile(fs *pflag.FlagSet) error {
	configFile, _ := fs.GetString("config")
	if configFile != "" {
		mgr.viper.SetConfigFile(configFile)
		return mgr.viper.MergeInConfig()
	}
	return nil
}

// InitFlags initializes the flagset with the provided config.
func (mgr *Manager) InitFlags(cfg any) error {
	rootStruct := reflect.TypeOf(cfg)

	if rootStruct.Kind() != reflect.Struct {
		panic("configuration is not a struct")
	}

	mgr.parseStructToFlags("", rootStruct)

	err := mgr.viper.BindPFlags(mgr.flags)
	if err != nil {
		panic(err)
	}

	return nil
}

// Flags returns pflag.FlagSet.
func (mgr *Manager) Flags() *pflag.FlagSet {
	return mgr.flags
}

// Unmarshal unmarshals the read config into the provided struct.
func (mgr *Manager) Unmarshal(config interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "name",
		Result:  config,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringSliceToStringMapHookFunc,
		),
	})
	if err != nil {
		return err
	}

	decoder.Decode(mgr.viper.AllSettings())
	return nil
}

// Viper returns viper.
func (mgr *Manager) Viper() *viper.Viper {
	return mgr.viper
}

// parseStructToFlags parses a struct and returns a flagset.
// It panics if there are parsing errors.
func (mgr *Manager) parseStructToFlags(prefix string, strT reflect.Type) {
	for i := 0; i < strT.NumField(); i++ {
		field := strT.Field(i)
		name := field.Tag.Get("name")
		kind := field.Type.Kind()
		if (name == "" || name == "-") && kind != reflect.Struct {
			continue
		}

		desc := field.Tag.Get("description")
		short := field.Tag.Get("short")

		if prefix != "" {
			name = prefix + "." + name
		}

		switch kind {
		case reflect.String:
			mgr.flags.StringP(name, short, "", desc)
		case reflect.Bool:
			mgr.flags.BoolP(name, short, false, desc)
		case reflect.Uint:
			mgr.flags.UintP(name, short, 0, desc)
		case reflect.Uint64:
			mgr.flags.Uint64P(name, short, 0, desc)
		case reflect.Int:
			mgr.flags.IntP(name, short, 0, desc)
		case reflect.Int64:
			mgr.flags.Int64P(name, short, 0, desc)
		case reflect.Float64:
			mgr.flags.Float64P(name, short, 0, desc)
		case reflect.Slice:
			mgr.flags.StringSliceP(name, short, nil, desc)
		case reflect.Struct:
			// This allows for recursion
			mgr.parseStructToFlags(name, field.Type)
		case reflect.Map:
			mgr.flags.StringSliceP(name, short, nil, desc)

		default:
			panic(fmt.Errorf("Unknown type in config: %v", kind))
		}
	}
}
