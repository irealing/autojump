# 微信跳一跳辅助工具

自动跳一跳，算法参考项目[wechat_jump_game](https://github.com/wangshub/wechat_jump_game)。使用golang开发，免除配置和搭建环境的麻烦。

## 安装和使用

```shell
$ go get github.com/irealing/autojump
$ go install github.com/irealing/autojump
$ export PATH=$PATH:$GOPATH\bin
```

*启动程序前须启动调试模式，并允许USB调试。启动程序时会提示是否允许调试，允许即可。*

*须将手机界面调整至微信跳一跳界面开始游戏方可自动跳一跳.*

## 其它

`img`目录为截图目录，运行结束后可删除其中的内容。