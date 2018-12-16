## 视频播放

该项目主要解决使用手机或者平板看PC上的视频资源，另外还支持一些收藏功能。<br>
因为我从tumblr上下载了大量资源后我发现看的时候还不是特别方便，于是写了这个小工具。<br>

# 编译

```
git clone https://github.com/abelQJ/video-web-player.git
cd video-web-player && go build
```

# 启动

```
./video-web-player -confAccessKey key -port num -https true
```
请记住这个key，在后续配置系统时会用到。不指定confAccessKey，默认key=pi。端口不设定默认值是8080，https默认不开启

# 配置

网站起来后，进入配置页配置视频资源。可以指定本机上供浏览的视频资源路径，该路径下所有视频（递归）会呈现在一个列表中。配置格式如下，配置时需要指定confAccessKey,否则会提示没权限修改。

```
{
	"dirs": [{
			"id": 1,
			"path": "/Users/abel/tmp/mp4",
			"desc": "测试"
		},
		{
			"id": 2,
			"path": "/Volumes/Abel/tumblr/video/v1",
			"desc": "关注用户"
		},
		{
			"id": 3,
			"path": "/Volumes/Abel/tumblr/video/v2",
			"desc": "视频1"
		},
		{
			"id": 4,
			"path": "/Volumes/Abel/tumblr/video/v3",
			"desc": "视频2"
		}
	],
	"confAccessKey": "pi"
}
```
./video-web-player -confAccessKey key
