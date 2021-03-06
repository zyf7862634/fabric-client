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
	fmt.Println("runDir=", runDir)

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

	getInfoList := make([]*GetInfo, 0)
	getInfo1 := &GetInfo{
		CardType: 1,
		CardID:   "身份证S3000",
	}
	// getInfo2 := &GetInfo{
	// 	CardType: 1,
	// 	CardID:   "身份证S2000",
	// }

	getInfoList = append(getInfoList, getInfo1)
	// getInfoList = append(getInfoList, getInfo2)

	arrayList, err := json.Marshal(getInfoList)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//printf
	{
		arrayList1, err := json.MarshalIndent(getInfoList, "", " ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(arrayList1))
	}

	payLoad, err := fabric.QueryChaincode(fabConfig.ChaincodeID, "User1", "get_info", [][]byte{arrayList})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(payLoad))

}
