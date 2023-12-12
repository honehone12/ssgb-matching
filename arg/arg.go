package arg

import "flag"

type Args struct {
	SettingFile string
}

func ParseArgs() Args {
	settingFile := flag.String("s", "setting.json", "setting file name")
	flag.Parse()
	return Args{
		SettingFile: *settingFile,
	}
}
