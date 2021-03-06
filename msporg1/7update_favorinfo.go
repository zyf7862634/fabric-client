package main

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"
	"strings"

	"github.com/devfabric/fabric-client/config"
	"github.com/devfabric/fabric-client/fabsdk"
)

func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

func main() {
	runDir, err := GetCurrentDirectory()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//设置环境变量，防止应用未设置
	workDirForFabSDK := os.Getenv("WORKDIR")
	if workDirForFabSDK == "" {
		os.Setenv("WORKDIR", runDir)
	}

	fabConfig, err := config.LoadFabircConfig(runDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// connectConfig, err := ioutil.ReadFile(fabConfig.ConfigFile)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	fabric := fabsdk.NewFabricClient(fabConfig.ConfigFile, fabConfig.ChannelID, fabConfig.UserName, fabConfig.UserOrg)
	err = fabric.Setup(runDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	putFavorList := make([]*UpdateFavor, 0)
	putItem1 := &PutFavor{
		CardType: 1,
		CardID:   "身份证S3000",
		AssetsInfo: &AssetsInfo{
			FieldKVMap: map[string]interface{}{
				"社保身份": 1,
			},
		},
	}

	putItem2 := &PutFavor{
		CardType: 1,
		CardID:   "身份证S4000",
		AssetsInfo: &AssetsInfo{
			FieldKVMap: map[string]interface{}{
				"社保身份": 1,
			},
		},
	}
	putFavorList = append(putFavorList, putItem1)
	putFavorList = append(putFavorList, putItem2)

	// putItem1 := &PutFavor{
	// 	CardType: 1,
	// 	CardID:   "身份证S1000",
	// 	BankCard: &BankCard{
	// 		FieldKVMap: map[string]interface{}{
	// 			"平安银行": 1000,
	// 		},
	// 	},
	// }
	// putFavorList = append(putFavorList, putItem1)

	arrayList, err := json.Marshal(putFavorList)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		arrayList1, err := json.MarshalIndent(putFavorList, "", " ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(arrayList1))
	}

	// // fmt.Println("success:", string(arrayList))
	// {
	// 	var puDataTest []*PutDataReq
	// 	err := json.Unmarshal(arrayList, &puDataTest)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		return
	// 	}
	// 	// fmt.Println("success1:", puDataTest)
	// 	for i := range puDataTest {
	// 		fmt.Println(puDataTest[i])
	// 	}
	// }

	payLoad, err := fabric.InvokeChaincodeWithEvent(fabConfig.ChaincodeID, "User1", "update_favorinfo", [][]byte{arrayList, []byte(EvUpdateFavorInfo)})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
