/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "通过一个词语去找到关联的其他词语，也可以查找到一个链条",
	Long:  `通过一个词语去找到关联的其他词语，也可以查找到一个链条`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//和assoc关联的功能
		// 打开文件
		// 获取当前用户的主目录路径
		usr, err := user.Current()
		if err != nil {
			fmt.Println("获取用户信息失败：", err)
			return
		}
		homeDir := usr.HomeDir
		file, err := os.OpenFile(homeDir+"/.mind_assoc.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("打开文件失败：", err)
			return
		}
		defer file.Close()
		// 读取JSON数据
		data, err := ioutil.ReadAll(file)
		var old map[string]Assoc
		if err != nil {
			fmt.Println("读取文件失败：", err)
			//return
		} else {
			// 将JSON数据转换为map
			err = json.Unmarshal(data, &old)
			if err != nil {
				fmt.Println("转换JSON失败：", err)
				//return
			}
		}

		first := args[0]
		flag := cmd.Flag("all")
		if flag.Value.String() == "true" {
			fmt.Println("all")
			//遍历查找，直到结束
			next := treeTheMap(old, first)
			log.Println("next=", next)
			for next != "" {
				next = treeTheMap(old, next)
				log.Println("next=", next)
			}
		} else {
			assoc, ok := old[first]
			if !ok {
				fmt.Println("没有找到关联词")
			} else {
				fmt.Printf("%v 和%v 关联：%v\n", assoc.Item1, assoc.Item2, assoc.Desc)

			}
		}

	},
}

func treeTheMap(m map[string]Assoc, first string) (next string) {
	if first == "" {
		return ""
	}
	log.Println(first)
	assoc, ok := m[first]
	if !ok {
		fmt.Println("没有找到关联词")
		return ""
	} else {
		fmt.Printf("%v 和%v 关联：%v\n", assoc.Item1, assoc.Item2, assoc.Desc)
		//这里从Item1开始
		fmt.Println(" assoc.Item2 != first?", assoc.Item2 != first)
		fmt.Println(" assoc.Item1 != first?", assoc.Item1 != first)

		if assoc.Item2 != first {
			return assoc.Item2
		} else if assoc.Item1 != first {
			return assoc.Item1
		} else {
			return ""
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(searchCmd)

	//从一个词语关联到的下一个词语，那可能还有下下个，那么就可以一次性输出出来
	//todo 还有死循环
	//searchCmd.Flags().BoolP("all", "a", false, "输出一长篇记忆的词链条")
}
