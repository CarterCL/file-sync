package pipeline

import (
	"bufio"
	"file-sync/internal/config"
	"file-sync/internal/fetcher"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

type tmpFile struct {
	path string
	pair *config.FilePair
}

func Process(filePair *config.FilePair, wg *sync.WaitGroup) {
	defer wg.Done()

	tf, err := download(filePair)
	if err != nil {
		log.Printf("[%s] 下载失败: %v", filePair.Url, err)
		return
	}

	if err = preProcess(tf); err != nil {
		log.Printf("[%s] 预处理失败: %v", tf.path, err)
		return
	}

	if filePair.Convert {
		if err = convert(tf); err != nil {
			log.Printf("[%s] 转换失败: %v", tf.path, err)
			return
		}
	}

	if err = move(tf); err != nil {
		log.Printf("[%s] 移动失败: %v", tf.path, err)
	}
}

func download(pair *config.FilePair) (*tmpFile, error) {
	ua := pair.UA
	if ua == "" {
		ua = "file-sync/1.0"
	}

	data, err := fetcher.Fetch(pair.Url, ua)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("下载内容为空")
	}

	tmpPath := path.Join(os.TempDir(), path.Base(pair.Path))
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return nil, fmt.Errorf("写入临时文件失败: %w", err)
	}

	return &tmpFile{path: tmpPath, pair: pair}, nil
}

func preProcess(tf *tmpFile) error {
	lineFn := func(line string) string {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- '") {
			return strings.ReplaceAll(line, `- '+.`, "- '")
		}
		return line
	}

	preWriteFn := func(lines []string) []string {
		for _, ext := range tf.pair.Extensions {
			lines = append(lines, ext)
		}
		return lines
	}

	return editFile(tf, lineFn, preWriteFn)
}

func convert(tf *tmpFile) error {
	lineFn := func(line string) string {
		return strings.TrimSpace(strings.ReplaceAll(
			strings.TrimLeft(strings.TrimSpace(line), "-"),
			"'", "",
		))
	}
	return editFile(tf, lineFn, nil)
}

func move(tf *tmpFile) error {
	info, err := os.Stat(tf.path)
	if err != nil {
		return fmt.Errorf("无法访问文件: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("文件长度为 0，跳过移动")
	}

	return moveFile(tf.path, tf.pair.Path)
}

func editFile(tf *tmpFile, lineFn func(string) string, preWriteFn func([]string) []string) error {
	file, err := os.Open(tf.path)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "payload:") {
			continue
		}
		lines = append(lines, lineFn(line))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("关闭文件失败: %w", err)
	}

	if err := os.Remove(tf.path); err != nil {
		return fmt.Errorf("删除原文件失败: %w", err)
	}

	if preWriteFn != nil {
		lines = preWriteFn(lines)
	}

	return os.WriteFile(tf.path, []byte(strings.Join(lines, "\n")), 0644)
}

func moveFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	if err := in.Close(); err != nil {
		return fmt.Errorf("关闭源文件失败: %w", err)
	}

	if err := os.Remove(src); err != nil {
		return fmt.Errorf("删除源文件失败: %w", err)
	}

	return nil
}
