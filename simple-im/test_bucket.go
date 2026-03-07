package main

import (
	// "fmt"
	// "hash/fnv"
)

// func main() {
// 	bucketCount := 32
// 	users := []string{
// 		"user1", "user2", "user3", "user4", "user5",
// 		"user123", "user456", "user789", "alice", "bob",
// 		"charlie", "david", "eve", "frank", "grace",
// 	}

// 	// 统计每个 bucket 的用户数
// 	bucketStats := make(map[uint32][]string)

// 	for _, userID := range users {
// 		hash := fnv.New32a()
// 		hash.Write([]byte(userID))
// 		hashValue := hash.Sum32()
// 		bucketIndex := hashValue % uint32(bucketCount)
// 		bucketStats[bucketIndex] = append(bucketStats[bucketIndex], userID)
// 	}

// 	fmt.Println("=== Bucket 分布测试 ===")
// 	fmt.Printf("总用户数: %d\n", len(users))
// 	fmt.Printf("Bucket 数量: %d\n", bucketCount)
// 	fmt.Printf("使用的 Bucket 数: %d\n\n", len(bucketStats))

// 	for bucketIndex, userList := range bucketStats {
// 		fmt.Printf("Bucket %d: %v\n", bucketIndex, userList)
// 	}
// }
