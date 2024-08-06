/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/spf13/cobra"
)

var FileName = "/.mind_assoc_ad.json"

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "通过一个词语去找到关联的其他词语，也可以查找到一个链条",
	Long:  `通过一个词语去找到关联的其他词语，也可以查找到一个链条`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//和assoc关联的功能
		//打开文件
		//获取当前用户的主目录路径
		usr, err := user.Current()
		if err != nil {
			fmt.Println("获取用户信息失败：", err)
			return
		}
		homeDir := usr.HomeDir

		file, err := os.OpenFile(homeDir+FileName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("打开文件失败：", err)
			return
		}
		defer file.Close()
		// 读取JSON数据
		data, err := ioutil.ReadAll(file)
		var jsonData map[string][]Assoc
		if err != nil {
			//fmt.Println("读取文件失败：", err)
			//return
		} else {
			// 将JSON数据转换为map
			err = json.Unmarshal(data, &jsonData)
			if err != nil {

				//fmt.Println("转换JSON失败：", err)
				//return
			}
		}
		first := args[0]
		flag := cmd.Flag("all")
		if flag.Value.String() == "true" {
			//查找过的存放在一个map里面，免得死循环
			chainMap = make(map[string][]Assoc)
			chainList = make([]string, 0)

			treeTheMap(jsonData, first)
			fmt.Println("关联的词链和图像：", chainMap)
			fmt.Println("关联的词链：", chainList)

		} else {
			assocs, ok := jsonData[first]
			if !ok {
				fmt.Println("没有找到关联词")
			} else {
				for _, assoc := range assocs {
					fmt.Printf("%v 和%v 关联：%v\n", assoc.Item1, assoc.Item2, assoc.Desc)
				}
			}
		}

	},
}

// 已经找过的词
var chainMap map[string][]Assoc
var chainList []string

func treeTheMap(m map[string][]Assoc, first string) {
	if first == "" {
		return
	}
	//fmt.Println(" 当前词语：", first)
	assocs, ok := m[first] //是个列表
	if !ok {
		fmt.Println("没有找到关联词")
		return
	} else {

		_, ok2 := chainMap[first]
		if !ok2 {
			//放入chainMap
			chainMap[first] = assocs
			chainList = append(chainList, first) //
		} else {

			return
		}
		for _, assoc := range assocs {
			if assoc.Item2 != first {
				treeTheMap(m, assoc.Item2)
			} else if assoc.Item1 != first {
				treeTheMap(m, assoc.Item1)
			} else {
				return
			}
		}

	}
}

func init() {
	rootCmd.AddCommand(searchCmd)

	//从一个词语关联到的下一个词语，那可能还有下下个，那么就可以一次性输出出来
	searchCmd.Flags().BoolP("all", "a", false, "输出一长篇记忆的词链条")
}
