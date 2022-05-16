具体架构请参考文档：[Design](design.md)（[设计](design_cn.md)）

**ATTENTION: This project is a Work-in-Progress.**

**注意，目前尚在开发，暂时无法运行，仅供代码参考。**

## KubeOps structure
```
.
├── api                     // API&Error Proto files & Generated codes
│   └── v1
│       └── errors
├── app                     // microservices projects
│   ├── ansible
│   ├── cache
│   ├── client
│   ├── config
│   ├── constant
│   ├── inventory
│   └── server
├── cmd
│   ├── client
│   │   ├── adhoc
│   │   ├── playbook
│   │   ├── project
│   │   ├── root
│   │   └── task
│   ├── inventory
│   └── server
├── conf
├── deploy
│   ├── build
│   ├── docker-compose
│   └── kubernetes
├── docs
├── example
├── pkg                       // common used packages
│   └── util
├── plugin
└── third_party
    ├── errors
    ├── google
    │   └── api
    └── validate


```
