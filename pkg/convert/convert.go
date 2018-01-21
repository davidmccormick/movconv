package convert

import (
	"fmt"
	"regexp"
)

type Options struct {
	SourceFile     string
	OutFile        string
	Subtitles      bool
	BurnPGCSubs    bool
	ForceAudio     bool
	NoiseReduction bool
	Sharpen        bool
	Restore        bool
	AudioTrack     int
	VideoFormat    string
	VideoTrack     int
	VideoFlags     string
	AudioFlags     string
	AudioFormat    string
}

var audioFlags map[string]string = map[string]string{
	"aac":  "aac -b:a 512k -ac 2 -clev 1.414 -slev .5 -strict -2",
	"eac3": "eac3 -ab 1536k -strict -2",
	"dts":  "dca -ab 1536k -strict -2",
}

var videoFlags map[string]string = map[string]string{
	"x264": "libx264 -profile:v high -level 4.2 -preset slow -crf 12 -pix_fmt yuv420p -movflags faststart",
	"x265": "libx265 -pix_fmt yuv420p10le -preset fast -x265-params level=5.2:vbv-bufsize=60000:vbv-maxrate=60000:crf=20",
}

func NewOptions() Options {
	return Options{
		SourceFile:     "",
		OutFile:        "",
		Subtitles:      true,
		BurnPGCSubs:    false,
		ForceAudio:     false,
		NoiseReduction: false,
		Sharpen:        false,
		Restore:        false,
		AudioTrack:     0,
		VideoTrack:     0,
		VideoFlags:     "",
		AudioFlags:     "",
		AudioFormat:    "",
	}
}

func Execute(args []string, flags Options) {
	fmt.Printf("Running convert command from within convert package...\n")
	fmt.Printf("Options: %+v\n", flags)

	if flags.OutFile == "" {
		// setoutfile to input with extra .converted before final file type.
		r := regexp.MustCompile(`^(.+)\.([^.]+)$`)
		flags.OutFile = r.ReplaceAllString(flags.SourceFile, "$1.converted.$2")
		fmt.Printf("OutFile set to %s\n", flags.OutFile)
	}

	switch flags.Profile {
		case "ps4":
			flags.VideoFormat = "x265"
			flags.AudioFormat = "aac"
		default:
			flags.VideoFormat = "x265"
			flags.AudioFormat = "eac3"
		}

}

}
