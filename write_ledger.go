/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// write() - genric write variable into ledger
// 写入方法() - 写入账本
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
// ============================================================================================================================
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("starting write")
	// 除了方法名write 外还需要2个参数 key 和 value
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the ledger 写入账本
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}

// ============================================================================================================================
// delete_marble() - remove a marble from state and from marble index
// 删除小球 - 删除一个小球
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs - Array of strings
//      0      ,         1
//     id      ,  authed_by_company
// "m999999999", "united marbles"
// ============================================================================================================================
func delete_marble(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting delete_marble")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	authed_by_company := args[1]

	// get the marble
	marble, err := get_marble(stub, id)
	if err != nil {
		fmt.Println("Failed to find marble by id " + id)
		return shim.Error(err.Error())
	}

	// check authorizing company (see note in set_owner() about how this is quirky)
	if marble.Owner.Company != authed_by_company {
		return shim.Error("The company '" + authed_by_company + "' cannot authorize deletion for '" + marble.Owner.Company + "'.")
	}

	// remove the marble
	err = stub.DelState(id) //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete_marble")
	return shim.Success(nil)
}

// ============================================================================================================================
// Init Owner - create a new owner aka end user, store into chaincode state
// 这个用来创建新的User
// Shows off building key's value from GoLang Structure
//
// Inputs - Array of Strings
//           0     ,     1   ,   2 , 3,4,5
//      user id   , name , department, role, createTime, updateTime
// "u9999999999999","贾瞢", "技术部", "前段开发","2018-8-20 13:00:00",""
// ============================================================================================================================
func init_user(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("开始运行init_user，创建User")

	if (len(args) != 6) || (len(args) != 7) {
		return shim.Error("错误的参数个数. 期望 6 或 7 个参数")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 组建User,资产数量为0
	var user User
	user.ObjectType = "hk_user"
	user.Id = args[0]
	user.Name = args[1]
	user.Department = args[2]
	user.Role = args[3]
	user.Asset = 0
	user.CreateTime = args[4]
	user.UpdateTime = args[5]

	fmt.Println(user)

	//check if user already exists
	_, err = get_user(stub, user.Id)
	if err == nil {
		fmt.Println("This user already exists - " + user.Id)
		return shim.Error("This user already exists - " + user.Id)
	}

	//store user
	userAsBytes, _ := json.Marshal(user)      //convert to array of bytes
	err = stub.PutState(user.Id, userAsBytes) //store user by its Id
	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_user" + user.Name)
	return shim.Success(nil)
}

参数
0  ,  1   , 2   ,3     ,  4
A     B     X    交易内容   交易类型  
// 交易环科币的函数
func transaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	txId := stub.GetTxID()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	fmt.Println("交易：A → B ，币数量")

	if len(args) != 3 {
		return shim.Error("错误的参数个数. 期望 3 个参数")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	A = args[0]
	B = args[1]
	
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = json.Unmarshal(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = json.Unmarshal(Bvalbytes)

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	newAssatA = strconv.Atoi(Aval.Asset) - X
	Aval.Asset = newAssatA
	newAssatB = strconv.Atoi(Aval.Asset) + X
	Bval.Asset = newAssatB

	newAvalAsBytes, _ := json.Marshal(Aval)
	err = stub.PutState(A, newAvalAsBytes) //store user by its Id
	if err != nil {
		fmt.Println("Could not store" + A)
		return shim.Error(err.Error())
	}else{
		//添加交易记录
		var transactionLogA Transaction
		transactionLog.ObjectType = "hk_transactionlogs"
		transactionLog.TxId = txId
		transactionLog.Content = args[3]
		if args[4] == "赠送同事" {
			transactionLog.Type = "赠送同事"
		}else if args[4] == "公司兑换" {
			transactionLog.Type = "公司兑换"
		}
		transactionLog.Relation = B
		transactionLog.Volume = "-" + string(X)
		transactionLog.TimeStamp = args[5]

		// PutState 不确定能不能put进去
		var ATX User
		// 取出来A的value集
		ATXbytes, err := stub.GetState(A)
		if err != nil {
			return shim.Error("Failed to get state")
		}
		if Avalbytes == nil {
			return shim.Error("Entity not found")
		}
		//反序列化
		json.Unmarshal(ATXbytes,&ATX))
		//添加到List中
		ATX.List = append(ATX.List,transactionLogA)
		//序列化
		newAvalAsBytes := json.Marshal(ATX);
		//
		err = stub.PutState(A, newAvalAsBytes)
		if err != nil {
			fmt.Println("交易记录保存失败：" + ATX.Name + ATX.List[len(ATX.List)-1].Type + ATX.List[len(ATX.List)-1].Volume )
			return shim.Error(err.Error())
		}
	}
	newBvalAsBytes, _ := json.Marshal(Bval)
	err = stub.PutState(B, newBvalAsBytes) //store user by its Id
	if err != nil {
		fmt.Println("Could not store" + B)
		return shim.Error(err.Error())
	}else{
		//添加交易记录
		var transactionLogB Transaction
		transactionLog.ObjectType = "hk_transactionlogs"
		transactionLog.TxId = txId
		transactionLog.Content = args[3]
		if args[4] == "同事赠送" {
			transactionLog.Type = "同事赠送"
		}else if args[4] == "公司奖励" {
			transactionLog.Type = "公司奖励"
		}else if args[4] == "任务收益" {
			transactionLog.Type = "任务收益"
		}
		transactionLog.Relation = A
		transactionLog.Volume = "+" + string(X)
		transactionLog.TimeStamp = args[5]

		// PutState 不确定能不能put进去
		var BTX User
		// 取出来B的value集
		BTXbytes, err := stub.GetState(B)
		if err != nil {
			return shim.Error("Failed to get state")
		}
		if Bvalbytes == nil {
			return shim.Error("Entity not found")
		}
		//反序列化
		json.Unmarshal(BTXbytes,&BTX))
		//添加到List中
		BTX.List = append(BTX.List,transactionLogB)
		//序列化
		newBvalBsBytes := json.Marshal(BTX);
		//
		err = stub.PutState(B, newBvalBsBytes)
		if err != nil {
			fmt.Println("交易记录保存失败：" + BTX.Name + BTX.List[len(ATX.List)-1].Type + ATX.List[len(ATX.List)-1].Volume )
			return shim.Error(err.Error())
		}
	}

	fmt.Println("- end transaction from" + Aval.Name + "to" + Bval.Name + "for" + X +"HKCoin" )
	return shim.Success(nil)
}