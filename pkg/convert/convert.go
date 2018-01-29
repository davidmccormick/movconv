package convert

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//Options models the different conversion settings and ways of selecting them.
type Options struct {
	Profile        string
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

var validProfiles = []string{"tv", "ps4"}
var validVideoFormats = []string{"x264", "x265"}
var validAudioFormats = []string{"aac", "eac3", "dts"}
var ffmpegMappings = map[string]string{
	"x264": "libx264 -profile:v high -level 4.2 -preset slow -crf 12 -pix_fmt yuv420p -movflags faststart",
	"x265": "libx265 -pix_fmt yuv420p10le -preset fast -x265-params level=5.2:vbv-bufsize=60000:vbv-maxrate=60000:crf=20",
	"aac":  "aac -b:a 512k -ac 2 -clev 1.414 -slev .5 -strict -2",
	"eac3": "eac3 -ab 1536k -strict -2",
	"dts":  "dca -ab 1536k -strict -2",
}

//NewOptions - constructor for an Options object with default settings.
func NewOptions() Options {
	return Options{
		Profile:        "",
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

func valueInList(value string, validList []string) bool {
	for _, s := range validList {
		if value == s {
			return true
		}
	}
	return false
}

//ValidateOptions - called to check that the options are sane.
func ValidateOptions(args []string, flags Options) error {
	fmt.Printf("Validating options\n")
	fmt.Printf("Options: %+v\n", flags)

	if flags.Profile != "" && !valueInList(flags.Profile, validProfiles) {
		return fmt.Errorf("invalid profile specified: %s", flags.Profile)
	}
	if flags.VideoFormat != "" && !valueInList(flags.VideoFormat, validVideoFormats) {
		return fmt.Errorf("invalid video format specified: %s", flags.VideoFormat)
	}
	if flags.AudioFormat != "" && !valueInList(flags.AudioFormat, validAudioFormats) {
		return fmt.Errorf("invalid audio format specified: %s", flags.AudioFormat)
	}
	return nil
}

//Execute - do to the requested conversion.
func Execute(args []string, flags Options) error {
	fmt.Printf("Running convert command from within convert package...\n")
	fmt.Printf("Options: %+v\n", flags)

	if _, err := os.Stat(flags.SourceFile); os.IsNotExist(err) {
		return err
	}

	if flags.OutFile == "" {
		// setoutfile to input with extra .converted before final file type.
		r := regexp.MustCompile(`^(.+)\.([^.]+)$`)
		flags.OutFile = r.ReplaceAllString(flags.SourceFile, "$1.converted.$2")
		fmt.Printf("OutFile set to %s\n", flags.OutFile)
	}

	if flags.Profile != "" {
		switch flags.Profile {
		case "ps4":
			flags.VideoFormat = "x265"
			flags.AudioFormat = "aac"
		default:
			flags.VideoFormat = "x265"
			flags.AudioFormat = "eac3"
		}
	}

	codecs, err := findCodecs(flags.SourceFile)
	if err != nil {
		return err
	}

	fmt.Printf("Found these codecs in the source: %v\n", codecs)
	return nil
}

// return the codecs found in a movie file by name
func findCodecs(filepath string) ([]string, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return []string{}, fmt.Errorf("File %s does not exist", filepath)
	}

	c1 := exec.Command("ffprobe", "'filepath'", "-show_entries", "stream=codec_name")
	c2 := exec.Command("grep", "codec_name")
	c3 := exec.Command("cut", "-d=", "-f2")

	//set up the pipeline
	c2.Stdin, _ = c1.StdoutPipe()
	c3.Stdin, _ = c2.StdoutPipe()

	_ = c2.Start()
	_ = c3.Start()
	_ = c1.Run()
	_ = c2.Wait()
	_ = c3.Wait()

	out, _ := c3.Output()
	return strings.Fields(string(out)), nil
}
