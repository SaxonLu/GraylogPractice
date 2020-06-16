package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zbindenren/logrus_mail"
	grayloghook "gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

var (
	logFilePath  = "./logdoc"
	logFileName  string
	createFile   bool   = true
	graylogHost  string = "127.0.0.1:12201"
	LogTimeLimit int    = 1
	LogSliceTime int    = 1
)

// 為根目錄systemlog做分類用
const (
	E          = 0 //Exception
	R          = 1 //Request
	InfoLevel  = logrus.InfoLevel
	FatalLevel = logrus.FatalLevel
	DebugLevel = logrus.DebugLevel
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
	PanicLevel = logrus.PanicLevel
)

func init() {
	//獲取項目的執行路徑
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config := viper.New()
	//設置讀取的文件路徑
	config.AddConfigPath(path)
	//設置讀取的文件名
	config.SetConfigName("config")
	//設置文件的類型
	config.SetConfigType("yaml")
	//嘗試進行配置讀取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	createFile = config.Get("set.createFile").(bool)
	graylogHost = config.Get("set.graylogHost").(string)
	LogTimeLimit = config.Get("set.LogTimeLimit").(int)
	LogSliceTime = config.Get("set.LogSliceTime").(int)
}

// log實作區塊
func Logger(status int) *logrus.Logger {

	logger := logrus.New()

	//Log 記錄等級
	logger.SetLevel(DebugLevel)

	//  ----  system.log 輸出至根目錄  ---- Start
	if createFile == true {

		if status == E {
			logFileName = "Exception.log"
		} else if status == R {
			logFileName = "System.log"
		}

		fileName := path.Join(logFilePath, logFileName)
		if _, err := os.Stat(fileName); err != nil {
			if _, err := os.Create(fileName); err != nil {
				fmt.Println(err.Error())
			}
		}

		// 打開文件
		src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("err", err)
		}

		//設置輸出
		logger.Out = src

		// 設置 rotatelogs
		logWriter, err := rotatelogs.New(
			// 分割後的文件名稱
			fileName+".%Y_%m_%d.log",

			// 生成软链，指向最新日志文件
			rotatelogs.WithLinkName(fileName),

			// Log最大保存時間(7天)

			rotatelogs.WithMaxAge(time.Duration(LogTimeLimit)*time.Hour),

			// 設置Log切割時間(1天)
			rotatelogs.WithRotationTime(time.Duration(LogSliceTime)*time.Hour),
		)
		writeMap := lfshook.WriterMap{
			InfoLevel:  logWriter,
			FatalLevel: logWriter,
			DebugLevel: logWriter,
			WarnLevel:  logWriter,
			ErrorLevel: logWriter,
			PanicLevel: logWriter,
		}
		// 輸出改成JSON 格式 並Formate 時間格式
		logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		}))
	}
	//  ----  system.log 輸出至根目錄  ---- Done

	//  ----  log Broadcast至Graylog ----

	grayHook := grayloghook.NewGraylogHook(graylogHost, map[string]interface{}{"server": "apiName"})

	logger.AddHook(grayHook)

	return logger

}

// 藉由gin監聽，收集對應Request/Respone 用
func logerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		//呼叫處理
		c.Next()

		endTime := time.Now()
		//Request 處理時間
		latencyTime := endTime.Sub(startTime)
		//請求方式
		reqMethod := c.Request.Method
		//請求路由
		reqUrl := c.Request.RequestURI
		//狀態碼
		statusCode := c.Writer.Status()
		//Request Ip
		clientIP := c.ClientIP()

		// log 格式
		Logger(R).WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info("LogLevel", "ThisInfoLevel 提示等級")
	}
}

func main() {
	app := gin.Default()
	//啟用監聽
	app.Use(logerMiddleware())

	//以下範例
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	app.GET("/zz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Message": "TT",
		})
	})
	app.GET("/TES", LoggerToES())
	app.Run()
}

type para struct {
	Name string
	ID   int
	Drop bool
}

// Sample
func LoggerToES() gin.HandlerFunc {
	logPara := para{
		Name: "Ada",
		ID:   234567,
		Drop: true,
	}

	//Info
	Logger(E).WithFields(logrus.Fields{
		"LogName":     "Iris",
		"Status":      "AccountNotExist",
		"SendToAdmin": true,
		"ParseCode":   1234,
		"PingList":    []string{"222", "#333"},
		"Para":        logPara,
	}).Info("LogLevel", "ThisInfoLevel 提示等級")
	//Error
	// Logger(E).WithFields(logrus.Fields{
	// 	"name": "Endless",
	// }).Error("LogLevel", "Error")
	//Warn
	// Logger(E).WithFields(logrus.Fields{
	// 	"name": "Warcraft",
	// }).Warn("LogLevel", "Warn")
	//Debug
	// Logger(E).WithFields(logrus.Fields{
	// 	"name": "Dell",
	// }).Debug("LogLevel", "Debug")
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Message": "efluthku",
		})
	}
}

// google會擋 未測試完成
func Email() {
	logger := logrus.New()
	//parameter"APPLICATION_NAME", "HOST", PORT, "FROM", "TO"
	hook, err := logrus_mail.NewMailAuthHook("testapp", "smtp.gmail.com", 8080, "From MailAdress", "To MailAdress", "smtp_account", "smtp_password")
	if err == nil {
		logger.Hooks.Add(hook)
	}
	//生成*Entry
	var filename = "123.txt"
	contextLogger := logger.WithFields(logrus.Fields{
		"file":    filename,
		"content": "GG",
	})
	//設置時間戳和message
	contextLogger.Time = time.Now()
	contextLogger.Message = "這是一個藉由hook發來的Ｍail"
	//只能發送Error,Fatal,Panic等級的log
	contextLogger.Level = logrus.FatalLevel

	//使用Fire發送,包含時間戳,message
	hook.Fire(contextLogger)
}
