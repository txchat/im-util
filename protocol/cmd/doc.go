/*
Copyright © 2022 oofpgDLD <oofpgdld@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var (
	format string
	dir    string
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "生成Protocol工具参考文档",
	Long:  ``,
	RunE:  docRunE,
}

func init() {
	docCmd.DisableAutoGenTag = true
	rootCmd.AddCommand(docCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	docCmd.Flags().StringVarP(&format, "format", "f", "", "文档格式")
	docCmd.Flags().StringVarP(&dir, "dir", "d", "./", "文档输出目录")
}

func docRunE(cmd *cobra.Command, args []string) error {
	switch format {
	case "man":
		return genMan(dir + "man")
	case "markdown":
		return genMD(dir + "markdown")
	}
	return nil
}

func checkAndCreateDir(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		//新建文件夹
		err = os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func genMan(filepath string) error {
	err := checkAndCreateDir(filepath)
	if err != nil {
		return err
	}
	header := &doc.GenManHeader{
		Title:   "IM Protocol Tools",
		Section: "1",
	}
	err = doc.GenManTree(rootCmd, header, filepath)
	if err != nil {
		return err
	}
	return nil
}

func genMD(filepath string) error {
	err := checkAndCreateDir(filepath)
	if err != nil {
		return err
	}
	err = doc.GenMarkdownTree(rootCmd, filepath)
	if err != nil {
		return err
	}
	return nil
}
