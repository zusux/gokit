package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/zusux/gokit/yaml/yamlmerge"
)

func main() {
	// 命令行参数
	varsPath := flag.String("vars", "", "远程 vars.yaml URL")
	tmplPath := flag.String("tmpl", "./configs/config-tpl.yaml", "本地 template.yaml 路径")
	outPath := flag.String("out", "./configs/config.yaml", "输出的 config.yaml 路径")

	flag.Parse()

	if *varsPath == "" {
		fmt.Println("❌ 错误: 必须指定 --vars 参数")
		os.Exit(1)
	}
	var vars []byte
	var err error
	if strings.HasPrefix(*varsPath, "http://") || strings.HasPrefix(*varsPath, "https://") {
		// HTTP 下载
		resp, err := http.Get(*varsPath)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		vars, err = io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
	} else {
		// 本地文件
		vars, err = os.ReadFile(*varsPath)
		if err != nil {
			panic(err)
		}
	}

	// 读取本地 template.yaml
	tmpl, err := os.ReadFile(*tmplPath)
	if err != nil {
		panic(fmt.Errorf("读取本地模板 %s 失败: %w", *tmplPath, err))
	}

	// 合并
	out, err := yamlmerge.MergeWithComments(string(tmpl), vars)
	if err != nil {
		panic(fmt.Errorf("合并失败: %w", err))
	}

	// 写入 config.yaml
	if err := os.WriteFile(*outPath, []byte(out), 0644); err != nil {
		panic(fmt.Errorf("写入 %s 失败: %w", *outPath, err))
	}

	fmt.Printf("✅ 已生成 %s\n", *outPath)
}
