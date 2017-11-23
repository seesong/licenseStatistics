# 电表License统计
	 _      _                                 _____  _           _    _       _    _            
	| |    (_)                               /  ___|| |         | |  (_)     | |  (_)           
	| |     _   ___   ___  _ __   ___   ___  \ '--. | |_   __ _ | |_  _  ___ | |_  _   ___  ___ 
	| |    | | / __| / _ \| '_ \ / __| / _ \  '--. \| __| / _' || __|| |/ __|| __|| | / __|/ __|
	| |____| || (__ |  __/| | | |\__ \|  __/ /\__/ /| |_ | (_| || |_ | |\__ \| |_ | || (__ \__ \
	\_____/|_| \___| \___||_| |_||___/ \___| \____/  \__| \__,_| \__||_||___/ \__||_| \___||___/
                               

## 功能一：电表数量统计

主要是学习学习MongoDB的distinct。

	db.elemeter.distinct("device_id")

## 功能二：电表详情cvs表格输出          

https://stackoverflow.com/questions/26062658/mongodb-aggregation-in-golang

	db.elemeter.aggregate({"$group":
                        {"_id":{"device_id":"$device_id"},
                         "pubtime":{"$last": "$pubtime"} 
                         }});
这个方法有些问题，还不知道怎么解决。                       

	http://www.01happy.com/golang-mongodb-find-demo/
	https://stackoverflow.com/questions/26062658/mongodb-aggregation-in-golang
	https://godoc.org/labix.org/v2/mgo#Collection.Pipe               


## 参数表

|参数|含义|
|:--|:--|
|-h --help|显示帮助文档。|
|-sum --sum|统计已经使用电表数量。|  
|-csv --csv|把每个电表最后一条记录输出到csv，不加file参数，默认文件名为当天日期。|
|-file --file|指定输出的csv文件名。|       
|无参数|默认是与-sum功能相同。|      


## config.ini的设定

	mongodbStyle = ssl 如果是华为云，就设定为ssl，不是就设定为normal