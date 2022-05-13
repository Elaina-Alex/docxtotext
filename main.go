package main

import (
	"baliance.com/gooxml/document"
	"doc/config"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var Vars *config.ConfigClass

func getDocxInformation(fileName string) string {
	var content string
	//doc.Paragraphs() 得到包含文档所有的段落的切片
	if doc, err := document.Open(fileName); err == nil {
		for _, para := range doc.Paragraphs() {
			//run为每个段落相同格式的文字组成的片段
			for _, run := range para.Runs() {
				content += run.Text()
			}
			content += "\n　　"
		}
	} else {
		log.Fatalf("error opening document: %s", err)
	}
	return content
}

func init() {
	config.MkdirFile("./TextFile")
	Vars = config.InitConfig()
	if err := json.Unmarshal(Vars.FileInformation, &Vars.FileStruct); err != nil {
		log.Fatalf("error unmarshaling: %s", err)
	}
	if Vars.FileStruct.DelDocxFile {
		fmt.Println("[提醒]已经开启转换后删除旧docx文件选择")
	}
	if Vars.FileStruct.DelDocFile {
		fmt.Println("[提醒]已经开启转换后删除旧doc文件选择")
	}
	if Vars.FileStruct.DocToDocx {
		fmt.Println("[提醒]已经开启doc转换成docx选择")
	}
}

func delDocxFile() {
	for _, fileName := range Vars.DelFileList {
		if err := os.Remove(fileName); err != nil {
			log.Println(err)
		}
	}
	fmt.Println("已经删除所有docx文件")
	Vars.DelFileList = []string{} // 清空删除文件列表
}

func main() {
	if Vars.FileStruct.DocToDocx {
		config.CmdPythonSaveDocx() // 调用python脚本转换doc为docx
	}
	if NameList := Vars.FileNameList; NameList != nil || len(NameList) != 0 {
		for index, file := range Vars.FileNameList {
			fileName := filepath.Base(file)
			switch path.Ext(fileName) {
			case ".docx":
				if docxContent := getDocxInformation(fileName); docxContent != "" {
					config.SaveFile(strings.Replace(fileName, ".docx", ".txt", -1), docxContent)
					Vars.DelFileList = append(Vars.DelFileList, fileName)
				} else {
					fmt.Println("文件" + fileName + "处理失败")
				}
			default:
				fmt.Println("No:", index, fileName, "不是docx文件，跳过处理！")
			}
		}
		fmt.Println("文档转换处理完成！")
		if Vars.FileStruct.DelDocxFile && len(Vars.DelFileList) != 0 {
			delDocxFile()
		}
	}
	log.Println("文件列表获取失败或没有查找到doc docx 文档")
}
