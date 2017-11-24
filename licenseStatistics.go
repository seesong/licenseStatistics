// Copyright 2017 sunny authors by cloudminds.com
// Licensed Copyright
// Source code and project home:
//
// https://github.com/sunnyregion/licenseStatistics
//
// Installation:
//
// go get  github.com/sunnyregion/licenseStatistics
//
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/sunnyregion/color"
	"github.com/sunnyregion/sunnyini"
	"github.com/sunnyregion/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	sumFlag          = flag.Bool("sum", false, "Use -sum 统计已经使用电表数量。")
	csvFlag          = flag.Bool("csv", false, "Use -csv 把每个电表最后一条记录输出到csv，不加file参数，默认文件名为当天日期。")
	csv_file *string = flag.String("file", util.SunnyTimeNow("day")+`.csv`, "Use -file <output file>")
)

const rootPEM = `-----BEGIN CERTIFICATE-----
MIIDmzCCAoOgAwIBAgIJAJRA9QxwGlXPMA0GCSqGSIb3DQEBCwUAMGQxCzAJBgNV
BAMMAkNOMQ8wDQYDVQQIDAZzaGFueGkxDTALBgNVBAcMBHhpYW4xDzANBgNVBAoM
Bmh1YXdlaTEMMAoGA1UECwwDUkRTMRYwFAYDVQQDDA1jYS5odWF3ZWkuY29tMB4X
DTE2MTIyMzA2NTQxMloXDTQ2MTIxNjA2NTQxMlowZDELMAkGA1UEAwwCQ04xDzAN
BgNVBAgMBnNoYW54aTENMAsGA1UEBwwEeGlhbjEPMA0GA1UECgwGaHVhd2VpMQww
CgYDVQQLDANSRFMxFjAUBgNVBAMMDWNhLmh1YXdlaS5jb20wggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQC5Ziz8T0DnoUYqqgLw5tDzw6+tcgpWrRQVFBy5
82do9jUrqy88nddheKpvJnkF4bJyPjBCAww1wplLDPCC9guwjOJZ7p1c8nZ/rXdL
5rXya3/fNl6h/JSpRW1laGUM5IjKPr/9bcjo9dlpr48cxl7P8sgWHfIlt7n0vuf0
qlQ+m4gTQrXAsGmcKyQPX1N04JP+4tkC4+lXtnChz9ncQKEMvTAq6EBrysZIDPDE
4PCkHSbTxUR0634BIxpxi3au12+P/AJJ/okSM/Aca2pDJuUuDJWkvnEfBqY3A2z+
Y/HaIA9xa8g+9yjh3EbvjYR84Fd2P3FFMMIWana2GxNtrnlNAgMBAAGjUDBOMB0G
A1UdDgQWBBRCUnkNIiaOaxRaQc+wgzNL9ZUD1DAfBgNVHSMEGDAWgBRCUnkNIiaO
axRaQc+wgzNL9ZUD1DAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAm
5YKFw25X2piJB2H8V2HdMvzUZo1aZkiQLEnGB0+VZKfKwJYJAAdAqfge2e/TmmDq
m6FMjdUXtviOQtdXYgiRxT3AVbF5coBiLVR1imUGxzc3kUtf0fjBHJ2Q9HKZwryS
ybCnRD6eTC3vD3wOzPb/bljDMo5e78Qsq1WD/Q/zaTBAdWyHCpQ/39ca9a1YJ4jD
hE0tecslhZlJUqup4SLYT8IJMea/JX418B8/jx4BQ+u02SOfSju1o0JgxOmuCPkJ
u2+HAwIjkDM5Hl02cO5aFLAyyl7N7cyz6gElv18ULKKYKllaMziOCP6iY6vkyrNU
RoOhR5k2fuJVS62O8Z/j
-----END CERTIFICATE-----`

//对象化
type LicenseStuct struct {
	err     error
	MongoDB *mgo.Database
}

//电表号的struct
type IDtype struct {
	DeviceId string `bson:"device_id"`
}

//数据库结构
type ElemeterMgo struct {
	ImageFile string    `bson:"image_file"`
	ID        IDtype    `bson:"_id"`
	Type      string    `bson:"type"`
	SubType   string    `bson:"sub_type"`
	Value     string    `bson:"value"`
	Flag      int       `bson:"flag"`
	AiImage   string    `bson:"ai_image"`
	IP        string    `bson:"ipaddress"`
	Pubtime   time.Time `bson:"pubtime"`
}

var (
	dbhostsip    = "127.0.0.1" //IP地址
	port         = "27017"
	dbusername   = ""         //用户名
	dbpassword   = ""         //密码
	database     = "elemeter" //表名
	collection   = "elemeter"
	mongodbStyle = "normal"
	crtfile      = "/root/ca.crt"
	url          = `mongodb://`
)

func NewLicenseStuct() (result *LicenseStuct) {

	f := sunnyini.NewIniFile()
	f.Readfile("config.ini")
	describ, v := f.GetValue("mongodb")
	if describ == "" { // 有数据
		dbhostsip = v[0]["dbhostsip"] //IP地址
		port = v[1]["port"]
		dbusername = v[2]["dbusername"] //用户名
		dbpassword = v[3]["dbpassword"] //密码
		database = v[4]["dbname"]       //表名
		collection = v[5]["collection"]
		mongodbStyle = v[6]["mongodbStyle"]
		crtfile = v[7]["crtfile"]
		url = url + v[8]["url"]
	} else {
		fmt.Println(describ)
	}

	if mongodbStyle == `ssl` {
		mogourl := url
		roots := x509.NewCertPool()
		if ca, err := ioutil.ReadFile(crtfile); err == nil {
			roots.AppendCertsFromPEM(ca)
		} else {
			ok := roots.AppendCertsFromPEM([]byte(rootPEM))
			if !ok {
				panic("failed to parse root certificate")
			}
		}
		tlsConfig := &tls.Config{
			RootCAs:            roots,
			InsecureSkipVerify: true,
		}

		dialInfo, err := mgo.ParseURL(mogourl)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			if err != nil {
				fmt.Println(err)
			}
			return conn, err
		}
		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			result := &LicenseStuct{err: err, MongoDB: nil}
			return result
		}
		result = &LicenseStuct{err: nil, MongoDB: session.DB(database)}
	} else {
		Host := []string{
			dbhostsip + ":" + port,
			// replica set addrs...
		}

		session, err := mgo.DialWithInfo(&mgo.DialInfo{
			Addrs:    Host,
			Timeout:  3 * time.Second, //10秒连接不到数据库
			Username: dbusername,
			Password: dbpassword,
			// Database: Database,
			// DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			// 	return tls.Dial("tcp", addr.String(), &tls.Config{})
			// },
		})
		if err != nil {
			result := &LicenseStuct{err: err, MongoDB: nil}
			return result
		}
		result = &LicenseStuct{err: nil, MongoDB: session.DB(database)}
	}

	return result
}

// 统计电表数量
func (this *LicenseStuct) SumLicense() {
	color := color.New(color.FgHiGreen).Add(color.BgBlack)
	c := this.MongoDB.C(collection)
	var results []interface{}
	c.Find(nil).Distinct("device_id", &results)
	color.Println("电表数量：", len(results))
	return
}

// 电表最后一条数据输出到csv
func (this *LicenseStuct) MeterDataToCsv() {
	color := color.New(color.FgHiGreen).Add(color.BgBlack)
	c := this.MongoDB.C(collection)
	pipe := c.Pipe([]bson.M{
		{
			"$group": bson.M{
				"_id":        bson.M{"device_id": "$device_id"},
				"pubtime":    bson.M{"$last": "$pubtime"},
				"image_file": bson.M{"$last": "$image_file"},
				"ai_image":   bson.M{"$last": "$ai_image"},
				`type`:       bson.M{"$last": "$type"},
				"sub_type":   bson.M{"$last": "$sub_type"},
				"value":      bson.M{"$last": "$value"},
				"flag":       bson.M{"$last": "$flag"},
				"ipaddress":  bson.M{"$last": "$ipaddress"},
			},
		},
		{"$sort": bson.M{"pubtime": 1}},
	})
	var result = []ElemeterMgo{}
	err := pipe.All(&result)
	if err != nil {
		panic(err)
	}
	//	color.Println(len(result))
	var data [][]string
	data = append(data, []string{"电表ID", "电表度数", "电表类型", "电表子类型", "web服务器保存的图片", "GPU服务器保存的图片", "访问IP", "访问时间"})
	for _, value := range result {
		//		if value.ID.DeviceId == `g121` {
		//			color.Println(value.ID.DeviceId, value)
		//			break
		//		}
		d := []string{value.ID.DeviceId, value.Value, value.Type, value.SubType, value.ImageFile, value.AiImage, value.IP, util.SunnyTimeToStr(value.Pubtime, "time")}
		data = append(data, d)
	}

	f, err := os.Create(*csv_file) //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)      //写入数据
	w.Flush()
	color.Println("csv over!")
}

//初始化
func init() {
	c := color.New(color.FgHiMagenta).Add(color.BgBlack)
	c.Println(`
 _      _                                 _____  _           _    _       _    _            
| |    (_)                               /  ___|| |         | |  (_)     | |  (_)           
| |     _   ___   ___  _ __   ___   ___  \ '--. | |_   __ _ | |_  _  ___ | |_  _   ___  ___ 
| |    | | / __| / _ \| '_ \ / __| / _ \  '--. \| __| / _' || __|| |/ __|| __|| | / __|/ __|
| |____| || (__ |  __/| | | |\__ \|  __/ /\__/ /| |_ | (_| || |_ | |\__ \| |_ | || (__ \__ \
\_____/|_| \___| \___||_| |_||___/ \___| \____/  \__| \__,_| \__||_||___/ \__||_| \___||___/ V0.1.2
                                                                                            
	`)
	flag.Parse()

}

//统计电报数量
func sumLicense() (result int) {
	return
}

//主程序
func main() {
	c := color.New(color.FgHiMagenta).Add(color.BgBlack)
	if *sumFlag || len(os.Args) == 1 {
		l := NewLicenseStuct()
		if l.err == nil {
			l.SumLicense()
		} else {
			c.Println(l.err.Error())
		}
	}

	if *csvFlag {
		l := NewLicenseStuct()
		if l.err == nil {
			l.MeterDataToCsv()
		} else {
			c.Println(l.err.Error())
		}
	}
}
