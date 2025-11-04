// Package _mongo @Author:冯铁城 [17615007230@163.com] 2025-11-04 14:51:47
package _mongo

import (
	"log"
	"spiders/db/_mongo/client"

	"go.mongodb.org/mongo-driver/bson"
)

// DeleteAndSaveData 删除并保存数据
func DeleteAndSaveData[T any](data []T, collectionName string, dbName string) error {

	//1.获取客户端
	mongoClient := client.CreateMongoClient()
	defer client.CloseMongoClient(mongoClient)

	//2.声明数据库以及集合
	db := mongoClient.GetClient().Database(dbName)
	collection := db.Collection(collectionName)

	//3.将 []T 转换为 []interface{}
	var saveDataList []interface{}
	for _, item := range data {
		saveDataList = append(saveDataList, item)
	}

	//4.全量删除数据
	if _, err := collection.DeleteMany(mongoClient.GetCtx(), bson.M{}); err != nil {
		return err
	} else {
		log.Println("delete mongo data success")
	}

	//5.保存
	if saveRes, err := collection.InsertMany(mongoClient.GetCtx(), saveDataList); err != nil {
		return err
	} else {
		log.Printf("save to mongo success, data count:%v", len(saveRes.InsertedIDs))
		return nil
	}
}
