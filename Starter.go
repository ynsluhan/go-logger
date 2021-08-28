package Logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"time"
)

// 日志实体类
var logger *logrus.Logger

// 获取日志路径
func GetLogPath() string {
	getwd, err := os.Getwd()
	//
	if err != nil {
		log.Fatal("ERROR 获取项目路径失败:", err)
	}
	//
	return path.Join(getwd, "logs")
}

// 获取logger
func GetLogger() *logrus.Logger {
	return logger
}


// 初始化日志
func InitLogger(logType, logPath *string) {
	//
	if logType != nil {

		if *logType == "json" {
			InitJsonLogger(logPath)
		}

		if *logType == "text" {
			InitTextLogger(logPath)
		}
	} else {

		InitTextLogger(logPath)

	}
}

/**
 * @Author yNsLuHan
 * @Description:
 * @Time 2021-06-28 09:20:20
 */
func InitTextLogger(logPath *string) {
	var lp string
	if logPath == nil {
		lp = GetLogPath()
	} else {
		lp = *logPath
	}
	// 日志文件
	fileName := path.Join(lp, "bo.log")
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("ERROR 日志文件创建失败:", err)
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
		path.Join(lp, "bo-%Y-%m-%d_%H.log"),
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

/**
 * @Author yNsLuHan
 * @Description:
 * @Time 2021-06-28 09:20:20
 */
func InitJsonLogger(logPath *string) {
	//
	var lp string
	if logPath == nil {
		lp = GetLogPath()
	} else {
		lp = *logPath
	}
	// 日志文件
	fileName := path.Join(lp, "bo.log")
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("ERROR 日志文件创建失败:", err)
	}
	// 实例化
	logger = logrus.New()
	// 设置输出
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)
	// 设置普通文本日志输出模式  json模式注释
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
		path.Join(lp, "bo-%Y-%m-%d_%H.log"),
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(30天)
		rotatelogs.WithMaxAge(30*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	// json格式  普通文本日志模式注释
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	// 钩子函数，格式化日期 json格式 普通模式注释
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 新增 Hook json格式 普通模式注释
	logger.AddHook(lfHook)
}