package main

import (
	"file-sync/internal/config"
	"file-sync/internal/pipeline"
	"flag"
	"log"
	"os"
	"sync"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	configFile := flag.String("c", "./config.yaml", "config file path")
	flag.Parse()

	syncConfig := config.InitConfig(*configFile)
	if len(syncConfig.SyncTasks) == 0 {
		return
	}

	wg := &sync.WaitGroup{}
	for _, task := range syncConfig.SyncTasks {
		if len(task.FilePairs) == 0 {
			continue
		}
		wg.Add(1)
		go runTask(task, wg)
	}
	wg.Wait()
}

func runTask(task *config.SyncTask, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("任务 [%s] 开始执行", task.Tag)

	wwg := &sync.WaitGroup{}
	for _, fp := range task.FilePairs {
		wwg.Add(1)
		go pipeline.Process(fp, wwg)
	}
	wwg.Wait()

	log.Printf("任务 [%s] 执行完毕", task.Tag)
}
