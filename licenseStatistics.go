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

	"github.com/sunnyregion/color"
	"github.com/sunnyregion/sunnyini"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
)

var (
	sumFlag = flag.Bool("sum", false, "是否统计有多少块电表在统计。")
)

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

	if *sumFlag {
		fmt.Println("ok")
	}

}
