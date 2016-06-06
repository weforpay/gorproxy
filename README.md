# gorproxy
简单的http反向代码,基于go语言的ReverseProxy
配置文件,在当前目录下
{
	"CertFilePath":"xxx.crt",
	"KeyFilePath":"xxx.key",
	"Router":[
		{
			"Src":"www.weforpay.com",
			"Dst":"http://host:port"
		}
	]
}