package kitdemo

import (
	"context"
	"sync"

	"github.com/MrZhangjicheng/kitdemo/transport"
)

// 程序入口以及整个生命周期管理
type App struct {
	// 控制生命周期
	ctx     context.Context
	cancel  func()
	servers []transport.Server
}

func New() *App {
	a := &App{}
	a.ctx, a.cancel = context.WithCancel(context.Background())
	return a

}

// 该函数负责将各个服务启动，并进行注册，---> 实现服务的优雅退出
func (a *App) Run() error {
	// 先启动服务
	wg := sync.WaitGroup{}
	for _, srv := range a.servers {
		sve := srv
		wg.Add(1)
		go func() error {
			return sve.Start(context.Background())
		}()
	}

	return nil

}
