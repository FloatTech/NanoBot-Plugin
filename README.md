<div align="center">
  <img src=".github/nano.jpeg" alt="东云名乃" width = "256">
  <br>

  <h1>NanoBot</h1>
  基于 NanoBot 的 QQ 机器人插件合集<br><br>

  <img src="https://counter.seku.su/cmoe?name=NanoBot&theme=r34" /><br>

  [![tencent-guild](https://img.shields.io/badge/%E9%A2%91%E9%81%93-Zer0BotPlugin-yellow?style=flat-square&logo=tencent-qq)](https://pd.qq.com/s/fjkx81mnr)

</div>

## 命令行参数
> `[]`代表是可选参数
```bash
nanobot [-Tadhst] ID1 ID2 ...

  -T int
        api timeout (s) (default 60)
  -a string
        qq appid
  -b    run in sandbox api
  -d    enable debug-level log output
  -h    print this help
  -s string
        qq secret
  -t string
        qq api token
```

## 功能
> 在编译时，以下功能均可通过注释`main.go`中的相应`import`而物理禁用，减小插件体积。

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


## 特别感谢

- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
