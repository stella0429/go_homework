package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"
	"time"
)
import "github.com/olivere/elastic/v7"

const (
	//todo ip修改
	address = "http://localhost:9200"
	esIndex = "blogs"
)

var (
	esClient *elastic.Client
	ctx      = context.Background()
)

type Article struct {
	Title   string    // 文章标题
	Content string    // 文章内容
	Author  string    // 作者
	Created time.Time // 发布时间
}

func init() {
	// 创建ES client用于后续操作ES
	// 创建client连接ES
	client, err := elastic.NewClient(
		// elasticsearch 服务地址，多个服务地址使用逗号分隔
		elastic.SetURL(address),
		// 基于http base auth验证机制的账号和密码
		elastic.SetBasicAuth("", ""),
		elastic.SetSniff(false),
		// 启用gzip压缩
		elastic.SetGzip(true),
		// 设置监控检查时间间隔
		elastic.SetHealthcheckInterval(10*time.Second),
		// 设置请求失败最大重试次数
		elastic.SetMaxRetries(5),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	esClient = client
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
		panic(err)
	} else {
		fmt.Println("连接成功")
	}
}
func main() {
	//insertDemo("6")
	//queryDemo("4")
	//queryMultiDemo("1", "2", "3")
	//updateDemo("6", map[string]interface{}{"Author": "lili", "Title": "新的文章标题"})
	//updateByQueryDemo("Author", "jack", "ctx._source['Title']='1111111'")
	//queryDemo("1")
	//deleteDemo("4")
	//deleteByQueryDemo("Author", "jack")
	//searchDemo("term", map[string]interface{}{"termname": "Author", "termvalue": "jack"})
	//searchDemo("match", map[string]interface{}{"matchname": "Title", "matchvalue": "新的文章标题"})
	searchDemo("must", map[string]interface{}{
		"termname": "Author", "termvalue": "jack",
		"matchname": "Title", "matchvalue": "新的文章标题"})
}

//添加文档
func insertDemo(id string) {
	// 定义一篇博客
	blog := Article{
		Title:   CreateRandomString(50),
		Content: CreateRandomString(100),
		Author:  CreateRandomString(10),
		Created: time.Now()}

	// 使用client创建一个新的文档
	insertRes, err := esClient.Index().
		Index(esIndex). // 设置索引名称
		Id(id).         // 设置文档id
		BodyJson(blog). // 指定前面声明struct对象
		Do(ctx)         // 执行请求，需要传入一个上下文对象
	if err != nil {
		fmt.Println("更新失败：", err)
		panic(err)
	}
	fmt.Printf("文档Id %s, 索引名 %s\n", insertRes.Id, insertRes.Index)
}

//查询文档
func queryDemo(id string) {
	res, err := esClient.Get().
		Index(esIndex). // 指定索引名
		Id(id).         // 设置文档id
		Do(ctx)         // 执行请求
	if err != nil {
		panic(err)
	}
	if res.Found {
		fmt.Printf("文档id=%s 版本号=%d 索引名=%s\n", res.Id, res.Version, res.Index)
	}

	//手动将文档内容转换成go struct对象
	article := Article{}
	// 提取文档内容，原始类型是json数据
	data, _ := res.Source.MarshalJSON()
	// 将json转成struct结果
	json.Unmarshal(data, &article)
	// 打印结果
	fmt.Println("title:", article.Title)
	fmt.Println("content:", article.Content)
	fmt.Println("author:", article.Author)
}

//批量查询
func queryMultiDemo(id1, id2, id3 string) {
	// 查询id等于1,2,3的博客内容
	result, err := esClient.MultiGet().
		Add(elastic.NewMultiGetItem().Index(esIndex).Id(id1)).
		Add(elastic.NewMultiGetItem().Index(esIndex).Id(id2)).
		Add(elastic.NewMultiGetItem().Index(esIndex).Id(id3)).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	// 遍历文档
	for _, doc := range result.Docs {
		// 转换成struct对象
		var content Article
		tmp, _ := doc.Source.MarshalJSON()
		err := json.Unmarshal(tmp, &content)
		if err != nil {
			panic(err)
		}

		fmt.Println("title:", content.Title)
		fmt.Println("content:", content.Content)
		fmt.Println("author:", content.Author)
	}
}

//更新文档
func updateDemo(id string, data map[string]interface{}) {
	res, err := esClient.Update().
		Index(esIndex). // 设置索引名
		Id(id).         // 文档id
		Doc(data).      // 更新字段，支持传入键值结构
		Do(ctx)         // 执行ES查询
	if err != nil {
		panic(err)
	}
	fmt.Println("更新成功:", res)
}

//根据条件更新文档
func updateByQueryDemo(name, value, script string) {
	res, err := esClient.UpdateByQuery(esIndex).
		// 设置查询条件，这里设置Author=tizi
		Query(elastic.NewTermQuery(name, value)).
		// 通过脚本更新内容
		Script(elastic.NewScript(script)).
		// 如果文档版本冲突继续执行
		ProceedOnVersionConflict().
		Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("更新成功:", res)
}

//删除文档
func deleteDemo(id string) {
	// 根据id删除一条数据
	res, err := esClient.Delete().
		Index(esIndex).
		Id(id). // 文档id
		Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("删除成功:", res)
}

//根据条件删除文档
func deleteByQueryDemo(name, value string) {
	res, err := esClient.DeleteByQuery(esIndex). // 设置索引名
		// 设置查询条件为: Author = tizi
		Query(elastic.NewTermQuery(name, value)).
		// 文档冲突也继续删除
		ProceedOnVersionConflict().
		Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("删除成功:", res)
}

//查询
func searchDemo(t string, data map[string]interface{}) {
	var query interface{}
	switch t {
	case "term":
		query = elastic.NewTermQuery(data["termname"].(string), data["termvalue"])
	case "match":
		query = elastic.NewMatchQuery(data["matchname"].(string), data["matchvalue"])
	case "must":
		// 创建bool查询
		tempQuery := elastic.NewBoolQuery().Must()
		// 创建term查询
		termQuery := elastic.NewTermQuery(data["termname"].(string), data["termvalue"])
		matchQuery := elastic.NewMatchQuery(data["matchname"].(string), data["matchvalue"])
		// 设置bool查询的must条件, 组合了两个子查询
		tempQuery.Must(termQuery, matchQuery)
		query = termQuery
	}

	searchResult, err := esClient.Search().
		Index(esIndex).               // 设置索引名
		Query(query.(elastic.Query)). // 设置查询条件
		Sort("Created", true).        // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(0).                      // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(10).                     // 设置分页参数 - 每页大小
		Pretty(true).                 // 查询结果返回可读性较好的JSON格式
		Do(ctx)                       // 执行请求

	if err != nil {
		panic(err)
	}

	fmt.Printf("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())

	if searchResult.TotalHits() > 0 {
		// 查询结果不为空，则遍历结果
		var b1 Article
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(b1)) {
			// 转换成Article对象
			if t, ok := item.(Article); ok {
				fmt.Println("title:", t.Title)
				fmt.Println("content:", t.Content)
				fmt.Println("author:", t.Author)
			}
		}
	}
}

//随机指定长度字符串
func CreateRandomString(len int) string {
	var (
		res string
		str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	)
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		res += string(str[randomInt.Int64()])
	}
	return res
}
