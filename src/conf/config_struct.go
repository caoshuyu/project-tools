package conf

type config struct {
	Http       httpConf
	Log        logConf
	KitVersion string `toml:"kit_version"`
}

type httpConf struct {
	Port int
}

type logConf struct {
	SavePath string `toml:"save_path"`
}
