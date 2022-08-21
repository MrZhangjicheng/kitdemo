module github.com/MrZhangjicheng/kitdemo/contrib/registry/consul

go 1.16

require (
    github.com/hashicorp/consul/api v1.14.0   
   
)

require  github.com/MrZhangjicheng/kitdemo v0.0.0 

replace github.com/MrZhangjicheng/kitdemo => ../../../
