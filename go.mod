module github.com/FloatTech/NanoBot-Plugin

go 1.20

require (
	github.com/FloatTech/AnimeAPI v1.6.1-0.20230207081411-573533b18194
	github.com/FloatTech/floatbox v0.0.0-20230827160415-f0865337a824
	github.com/FloatTech/gg v1.1.2
	github.com/FloatTech/imgfactory v0.2.1
	github.com/FloatTech/zbpctrl v1.5.3-0.20230130095145-714ad318cd52
	github.com/fumiama/NanoBot v0.0.0-20231015152604-ce34c996ef31
	github.com/fumiama/go-base16384 v1.7.0
	github.com/fumiama/unibase2n v0.0.0-20221020155353-02876e777430
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/sirupsen/logrus v1.9.3
	github.com/wdvxdr1123/ZeroBot v1.7.4
)

require (
	github.com/FloatTech/sqlite v1.5.7 // indirect
	github.com/FloatTech/ttl v0.0.0-20220715042055-15612be72f5b // indirect
	github.com/RomiChan/syncx v0.0.0-20221202055724-5f842c53020e // indirect
	github.com/RomiChan/websocket v1.4.3-0.20220227141055-9b2c6168c9c5 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/ericpauley/go-quantize v0.0.0-20200331213906-ae555eb2afa4 // indirect
	github.com/fumiama/cron v1.3.0 // indirect
	github.com/fumiama/go-registry v0.2.6 // indirect
	github.com/fumiama/go-simple-protobuf v0.1.0 // indirect
	github.com/fumiama/gofastTEA v0.0.10 // indirect
	github.com/fumiama/imgsz v0.0.2 // indirect
	github.com/fumiama/jieba v0.0.0-20221203025406-36c17a10b565 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	golang.org/x/image v0.3.0 // indirect
	golang.org/x/sys v0.1.1-0.20221102194838-fc697a31fa06 // indirect
	golang.org/x/text v0.6.0 // indirect
	modernc.org/libc v1.21.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.4.0 // indirect
	modernc.org/sqlite v1.20.0 // indirect
)

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.20.0-with-win386

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b
