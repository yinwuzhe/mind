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

type Assoc struct {
	Item1, Item2, Desc string
}

// assocCmd represents the assoc command
var assocCmd = &cobra.Command{
	Use:   "assoc",
	Short: "将一个词语和另一个词关联记忆",
	Long: `联想词是一种有效的记忆技巧，通过建立两个词语之间的联系，帮助我们更容易地记住它们。
联想词的基本原理是通过建立一副图像，将一个词与另一个词联系起来。
这个图像可以是任何东西，例如一个场景、一个物品、一个人物或者一个动作。
通过将两个词联系起来，我们可以更容易地记住它们，因为我们的大脑更容易记住图像和故事。`,
	Example: "assoc 耳朵 笔 笔刺穿了耳朵",
	Args:    cobra.ExactArgs(3),
	//ValidArgs: []string{"词语1", "词语2", "用一句话描述关联的图像"},
	Run: func(cmd *cobra.Command, args []string) {
		assoc := Assoc{args[0], args[1], args[2]}
		m := make(map[string]Assoc)
		m[assoc.Item1] = assoc
		m[assoc.Item2] = assoc
		log.Println(m)
		//将这个map和已有的map合并，然后序列化写入文件

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
				//fmt.Println("转换JSON失败：", err)
				//return
			}
		}
		log.Println(old)
		m3 := make(map[string]Assoc)
		for k, v := range m {
			m3[k] = v
		}
		for k, v := range old {
			m3[k] = v
		}

		// 写入JSON数据
		marshal, err := json.Marshal(m3)
		if err != nil {
			return
		}
		_, err = file.WriteAt(marshal, 0)
		if err != nil {
			fmt.Println("写入文件失败：", err)
			return
		}

		fmt.Println("关联成功")
	},
}

func init() {
	rootCmd.AddCommand(assocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//assocCmd.PersistentFlags().String("list", "", "special use may use ai help")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//ai 功能，暂时没法实现
	//assocCmd.Flags().BoolP("ai", "a", false, "special use may use ai help")
	//写入进去就是为了查找这个词的另外一个关联词，这个数据结构是个问题，如果是用数据库，那还是最方便的，每个列都能查找
	//assocCmd.Flags().BoolP("search", "s", false, "search the other word and a picture")
}
