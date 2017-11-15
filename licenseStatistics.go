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
	"flag"
	"fmt"
	"time"

	"github.com/sunnyregion/color"
	"github.com/sunnyregion/sunnyini"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
)

var (
	sumFlag = flag.Bool("sum", false, "是否统计有多少块电表在统计。")
)

//对象化
type LicenseStuct struct {
	err     error
	MongoDB *mgo.Database
}

//数据库结构
type ElemeterMgo struct {
	ImageFile string    `bson:"image_file"`
	DeviceId  string    `bson:"device_id"`
	Type      string    `bson:"type"`
	SubType   string    `bson:"sub_type"`
	Value     string    `bson:"value"`
	Flag      int       `bson:"flag"`
	IP        string    `bson:"ipaddress"`
	Pubtime   time.Time `bson:"pubtime"`
}

var (
	dbhostsip  = "127.0.0.1" //IP地址
	port       = "27017"
	dbusername = ""         //用户名
	dbpassword = ""         //密码
	database   = "elemeter" //表名
	collection = ""
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
	//	_ = c
	//	o1 := bson.M{"_id": bson.M{"device_id": "$device_id"}}

	//	o2 := bson.M{"pubtime": bson.M{"$last": "$pubtime"}}
	//	o3 := bson.M{"$group": []bson.M{o1, o2}}
	//	pipe := c.Pipe([]bson.M{o3})
	//	var result []bson.M{}
	//	_ = pipe.All(&result)
	//	fmt.Println(result)
	var results []interface{}
	//c.Find(nil).Select(bson.M{"device_id": 1, "pubtime": 1}).Distinct("device_id", &results)
	c.Find(nil).Distinct("device_id", &results)
	color.Println("电表数量：", len(results))
	return
}
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

}
