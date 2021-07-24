package Logger

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

// 路径
var basePath string

// 日志目录
var logPath string

// error
var err error

var logger *logrus.Logger

// 初始化
func init() {
	// 获取日志目录
	GetLogsPath()
	// 初始化logs
	InitLogs()
}

/**
 * @Author yNsLuHan
 * @Description: 获取日志目录
 * @Time 2021-06-28 09:16:43
 * @return string
 */
func GetLogsPath() {
	basePath, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	// 日志目录
	logPath = path.Join(basePath, "logs")
}

/**
 * @Author yNsLuHan
 * @Description:
 * @Time 2021-06-28 09:20:20
 */
func InitLogs() {
	fmt.Println(logPath)
	// 日志文件
	fileName := path.Join(logPath, "bo.log")
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("err", err)
	}
	// 实例化
	logger = logrus.New()
	// 设置输出
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)
	// 设置普通日志输出模式  json模式注释
	logger.SetFormatter(
		&nested.Formatter{
			// 禁用键值对日志类型
			HideKeys: true,
			// 格式
			TimestampFormat: "2006-01-02 15:04:05",
			// 字段顺序  logger.WithFields(logrus.Fields{"client_ip": ip})  必须要使用该日志模式
			FieldsOrder: []string{"client_ip", "req_method", "req_uri", "status_code", "latency_time"},
			// 日志禁用颜色
			NoColors: true,
		},
	)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		path.Join(logPath, "bo-%Y-%m-%d_%H.log"),
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(30天)
		rotatelogs.WithMaxAge(30*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	// 普通写入设置日志切割  json模式注释
	logger.SetOutput(logWriter)
}

// 使用日志框架logrus日志记录到文件
func DefaultLogger() gin.HandlerFunc {
	// 使用logrus 日志系统
	return func(c *gin.Context) {
		// 开始时间
		//startTime := time.Now()
		start := time.Now()
		// 计算执行时间
		cost := time.Since(start)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 自定义日志
		logger.Info(clientIP+" ", reqMethod+" ", reqUri+" ", strconv.Itoa(statusCode)+" ", cost)
		c.Next()
	}
}
