<div align="center">
  <img src=".github/nano.jpeg" alt="东云名乃" width = "256">
  <br>

  <h1>NanoBot-Plugin</h1>
  基于 <a href="https://github.com/fumiama/NanoBot">NanoBot</a> 的 QQ 机器人插件合集<br><br>

  <img src="https://counter.seku.su/cmoe?name=NanoBot&theme=r34" /><br>

  [![tencent-guild](https://img.shields.io/badge/%E9%A2%91%E9%81%93-Zer0BotPlugin-yellow?style=flat-square&logo=tencent-qq)](https://pd.qq.com/s/fjkx81mnr)

</div>

## 命令行参数
> `[]`代表是可选参数
```bash
nanobot [参数] ID1 ID2 ...
参数:
  -D    enable debug-level log output
  -T int
        api timeout (s) (default 60)
  -a string
        qq appid
  -c string
        load from config
  -h    print this help
  -public
        only listen to public intent
  -qq
        also listen QQ intent
  -s string
        qq secret
  -sandbox
        run in sandbox api
  -save string
        save bot config to filename (eg. config.yaml)
  -shardcount uint
        shard count
  -shardindex uint
        shard index
  -superallqq
        make all QQ users to be SuperUser
  -t string
        qq api token
```

其中公域配置参考如下，为一个数组，可自行增加更多 bot 实例。注意`Properties`不可为`[]`
```yaml
- AppID: "123456"
  Token: xxxxxxx
  Secret: ""
  SuperUsers:
    - "123456789"
  Timeout: 1m0s
  Intents: 1812730883
  ShardIndex: 0
  ShardCount: 0
  Properties: null
```

## 功能
> 在编译时，以下功能均可通过注释`main.go`中的相应`import`而物理禁用，减小插件体积。

<details>
  <summary>触发者撤回时也自动撤回(仅私域可用)</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/autowithdraw"`

  - [x] 撤回一条消息

</details>

<details>
  <summary>base16384加解密</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/b14"`

  - [x] 加密xxx

  - [x] 解密xxx

  - [x] 用yyy加密xxx

  - [x] 用yyy解密xxx

</details>

<details>
  <summary>base64卦加解密</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/base64gua"`

  - [x] 六十四卦加密xxx

  - [x] 六十四卦解密xxx

  - [x] 六十四卦用yyy加密xxx

  - [x] 六十四卦用yyy解密xxx

</details>

<details>
  <summary>base天城文加解密</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/baseamasiro"`

  - [x] 天城文加密xxx

  - [x] 天城文解密xxx

  - [x] 天城文用yyy加密xxx

  - [x] 天城文用yyy解密xxx

</details>

<details>
  <summary>英文字符翻转</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/chrev"`

  - [x] 翻转 I love you

</details>

<details>
  <summary>程序员做饭指南</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/dish"`

  - [x] 怎么做[xxx] | 烹饪[xxx]
  
  - [x] 随机菜谱 | 随便做点菜

</details>

<details>
  <summary>合成emoji</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/emojimix"`

  - [x] [emoji][emoji]

</details>

<details>
  <summary>每日运势</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/fortune"`

  - [x] 运势 | 抽签

  - [x] 设置底图[车万 DC4 爱因斯坦 星空列车 樱云之恋 富婆妹 李清歌 公主连结 原神 明日方舟 碧蓝航线 碧蓝幻想 战双 阴阳师 赛马娘 东方归言录 奇异恩典 夏日口袋 ASoul]

</details>

<details>
  <summary>原神抽卡</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/genshin"`

  - [x] 切换原神卡池

  - [x] 原神十连

</details>

<details>
  <summary>百人一首</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/hyaku"`

  - [x] 百人一首

  - [x] 百人一首之n

</details>

<details>
  <summary>bot管理相关</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/manager"`

  - [x] /exposeid @user1 @user2

</details>

<details>
  <summary>娶群友</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/qqwife"`

  - [x] 娶群友

  - [x] 群老婆列表
 
  - [x] (娶|嫁)@对方QQ
 
  - [x] 牛@对方QQ
 
  - [x] 闹离婚
 
  - [x] 买礼物给@对方QQ
 
  - [x] 做媒 @攻方QQ @受方QQ
 
  - [x] 查好感度@对方QQ
 
  - [x] 好感度列表
 
  - [x] [允许|禁止]自由恋爱
 
  - [x] [允许|禁止]牛头人
 
  - [x] 设置CD为xx小时    →(默认12小时)

</details>

<details>
  <summary>在线代码运行</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/runcode"`

  - [x] >runcode [language] help

  - [x] >runcode [language] [code block]

  - [x] >runcoderaw [language] [code block]

</details>

<details>
  <summary>签到</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/score"`

  - [x] 签到

  - [x] 获得签到背景

  - [x] 查看等级排名

</details>

<details>
  <summary>自检</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/status"`

  - [x] [检查身体 | 自检 | 启动自检 | 系统状态]

</details>

<details>
  <summary>塔罗牌</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/tarot"`

  - [x] 抽[塔罗牌|大阿卡纳|小阿卡纳]

  - [x] 解塔罗牌[牌名]

</details>

<details>
  <summary>抽老婆</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/wife"`

  - [x] 抽老婆

</details>

<details>
  <summary>猜单词</summary>

  `import _ "github.com/FloatTech/NanoBot-Plugin/plugin/wordle"`

  - [x] 个人猜单词

  - [x] 团队猜单词

  - [x] 团队六阶猜单词

  - [x] 团队七阶猜单词

</details>


## 特别感谢

- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
