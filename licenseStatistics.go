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
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sunnyregion/color"
	"github.com/sunnyregion/sunnyini"
	"github.com/sunnyregion/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	sumFlag = flag.Bool("sum", false, "是否统计有多少块电表在统计。")
	cvsFlag = flag.Bool("cvs", false, "把最后一块电表的数据输出到cvs。")
)

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
	dbhostsip  = "127.0.0.1" //IP地址
	port       = "27017"
	dbusername = ""         //用户名
	dbpassword = ""         //密码
	database   = "elemeter" //表名
	collection = "elemeter"
)

func NewLicenseStuct() *LicenseStuct {

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
	} else {
		fmt.Println(describ)
	}
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
	result := &LicenseStuct{err: nil, MongoDB: session.DB(database)}
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

// 电表最后一条数据输出到cvs
func (this *LicenseStuct) MeterDataToCvs() {
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
	for _, value := range result {
		//		if value.ID.DeviceId == `g121` {
		//			color.Println(value.ID.DeviceId, value)
		//			break
		//		}
		d := []string{value.ID.DeviceId, value.Value, value.Type, value.SubType, value.ImageFile, value.AiImage, value.IP, util.SunnyTimeToStr(value.Pubtime, "time")}
		data = append(data, d)
	}

	f, err := os.Create("test.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)      //写入数据
	w.Flush()
	color.Println("cvs over!")
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
\_____/|_| \___| \___||_| |_||___/ \___| \____/  \__| \__,_| \__||_||___/ \__||_| \___||___/
                                                                                            
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
	if *sumFlag {
		l := NewLicenseStuct()
		if l.err == nil {
			l.SumLicense()
		} else {
			c.Println(l.err.Error())
		}
	}

	if *cvsFlag {
		l := NewLicenseStuct()
		if l.err == nil {
			l.MeterDataToCvs()
		} else {
			c.Println(l.err.Error())
		}
	}

}
