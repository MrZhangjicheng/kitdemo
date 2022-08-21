package kitdemo

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MrZhangjicheng/kitdemo/log"
	"github.com/MrZhangjicheng/kitdemo/registry"
	"github.com/MrZhangjicheng/kitdemo/transport"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// 1. 程序的控制入口，需要用户传入 服务名称以及对应的endpoint,然后开启服务，并进行注册
// 2. 程序启动以及优雅退出

// 该结构对应的是需要用户传入参数,给用户使用
type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Endpoints() []string
}

// 该结构对应对服务启动以及优雅退出的参数
// 该模块包含两个部分，用户传入的参数以及生命周期管理 采用选项模式进行参数配置
type App struct {
	// 用户可配置参数
	opts options
	// 控制生命周期
	ctx    context.Context
	cancel func()
	// 注册到注册中心的实例，该参数 由 opts 构建 serviceInstance 进行真实服务注册
	mu       sync.Mutex
	instance *registry.ServiceIntance
}

func New(opts ...Option) *App {
	o := options{
		ctx:     context.Background(),
		singles: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	if id, err := uuid.NewUUID(); err != nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.logger != nil {
		log.SetLogger(o.logger)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}

}

// 该函数负责将各个服务启动，并进行注册，---> 实现服务的优雅退出
func (a *App) Run() error {
	// 构建对应注册服务的实例
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	// 赋值给a时需要考虑并发吗？
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()

	// 管理多个协程中可能出现的错误
	ewg, ctx := errgroup.WithContext(a.ctx)

	// 开启启动各个服务，使用多个协程
	wg := sync.WaitGroup{}

	for _, srv := range a.opts.Servers {
		srv := srv
		// 使用 errorGroup 开启协程 监听错误
		ewg.Go(func() error {
			<-ctx.Done() //监听通道，用来判断是否需要退出
			// 需确认
			return srv.Stop(context.Background())

		})
		wg.Add(1)
		// ewg 开启的协程中的执行函数只要有错误，会让通道有值
		ewg.Go(func() error {
			wg.Done()
			return srv.Start(context.Background())
		})
	}
	wg.Wait() // 等待所有服务启动好，其实最多就两个

	// 进行注册
	if a.opts.registar != nil {
		if err := a.opts.registar.Register(context.Background(), a.instance); err != nil {
			return err
		}
	}
	// 当所以组件启动没问题后，进行拦截用户输入的信号量
	s := make(chan os.Signal, 1)
	signal.Notify(s, a.opts.singles...)

	ewg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-s:
			return a.Stop()

		}
	})
	if err := ewg.Wait(); err != nil && errors.Is(err, context.Canceled) {
		return err
	}
	return nil

}

// 正常关闭，将服务从注册中心取消  是不是可以将服务也进行退出
func (a *App) Stop() error {
	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()
	if a.opts.registar != nil && instance != nil {
		if err := a.opts.registar.Deregister(context.Background(), instance); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) ID() string { return a.opts.id }

func (a *App) Name() string { return a.opts.name }

func (a *App) Version() string { return a.opts.version }

func (a *App) Endpoint() []string {
	if a.instance != nil {
		return a.instance.Endpoints
	}
	return nil
}

func (a *App) buildInstance() (*registry.ServiceIntance, error) {
	endpoints := make([]string, 0, len(a.opts.endpoints))

	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range a.opts.Servers {
			if r, ok := srv.(transport.Endpointer); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}

				endpoints = append(endpoints, e.String())
			}

		}
	}
	return &registry.ServiceIntance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Endpoints: endpoints,
	}, nil
}
