// Copyright © 2018 David McCormick <davidemccormick@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	//	"fmt"
	convert "movconv/pkg/convert"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert a movie file",
	Long: `convert a movie file using sane defaults based on the
specific requirements of converting files for home use.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return convert.ValidateOptions(args, options)
	},
	Run: func(cmd *cobra.Command, args []string) {
		convert.Execute(args, options)
	},
}

var options convert.Options

func init() {
	options = convert.NewOptions()
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	convertCmd.Flags().StringVarP(&options.SourceFile, "from", "f", "", "the path to the source file to convert from")
	convertCmd.MarkFlagRequired("from")
	convertCmd.Flags().StringVarP(&options.OutFile, "out", "o", "", "the path of the resultant movie (default computed from source file name)")
	convertCmd.Flags().StringVarP(&options.Profile, "profile", "", "", "specify a profile of either tv or ps4 to inherit common parameters")
	convertCmd.Flags().StringVarP(&options.VideoFlags, "videoflags", "", "", "specify/override the ffmpeg video flags")
	convertCmd.Flags().StringVarP(&options.AudioFlags, "audioflags", "", "", "specify/override the ffmpeg audio flags")
	convertCmd.Flags().StringVarP(&options.AudioFormat, "audio", "", "eac3", "Choose audio format from eac3, aac or dts")
	convertCmd.Flags().StringVarP(&options.VideoFormat, "video", "", "x265", "Choose video format from x264 or x265")
	convertCmd.Flags().BoolVarP(&options.Subtitles, "subtitles", "s", true, "Copy subtitles")
	convertCmd.Flags().BoolVarP(&options.BurnPGCSubs, "pcg", "", false, "Burn PGC subtitles into movie (requires video recode)")
	convertCmd.Flags().BoolVarP(&options.ForceAudio, "forceaudio", "", false, "Force the recoding of the audio")
	convertCmd.Flags().BoolVarP(&options.NoiseReduction, "noisereduction", "n", false, "Apply hqdn3d noise filter on video")
	convertCmd.Flags().BoolVarP(&options.Sharpen, "sharpen", "", false, "Apply sharpen filter on video")
	convertCmd.Flags().BoolVarP(&options.Restore, "restore", "", false, "Apply hqdn3d then sharpen filter on video")
	convertCmd.Flags().IntVarP(&options.AudioTrack, "audiotrack", "a", 0, "audio track number (defaults to 0)")
	convertCmd.Flags().IntVarP(&options.VideoTrack, "videotrack", "v", 0, "video track number (defaults to 0)")
}
