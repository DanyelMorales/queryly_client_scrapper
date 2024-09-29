package news

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"news_scrapper/pkg/style"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
func WriteFile(filename string, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}
func CreateDirIfNotExist(dir string) error {
	err := os.MkdirAll(dir, 0755) // os.ModePerm sets permissions to 0777
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Generate MD5 hash from a string
func GenerateMD5Hash(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func FetchExistingNews(dirName string) []string {
	entries, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	var inMemoryExistingNews []string
	for _, entry := range entries {
		if !entry.IsDir() {
			inMemoryExistingNews = append(inMemoryExistingNews, strings.Replace(entry.Name(), ".", "", -1))
		}
	}
	return inMemoryExistingNews
}

func ArticleExists(inMemoryExistingNews *[]string, id string) bool {
	id = strings.Replace(id, ".", "", -1)
	for _, item := range *inMemoryExistingNews {
		if item == id {
			return true
		}
	}
	return false
}

func ArticleWatcher(q *Client) {
watchlist:
	for {
		select {
		case _, _ = <-q.ArticleStreamEnd:
			close(q.ArticleChan)
			close(q.ArticleStreamEnd)
			style.ErrorLogCF(q.FullLog, "[%s][ArticleStreamEnd] Channel closed", q.Context)
			break watchlist
		case value, ok := <-q.ArticleChan:
			if !ok {
				style.ErrorLogCF(q.FullLog, "[%s][ArticleChan] Channel closed", q.Context)
				break watchlist
			}
			if q.SaveData && q.ShouldScrapeLink != nil && q.ShouldScrapeLink(value.Id) {
				jsonData, err := json.Marshal(value)
				if err != nil {
					style.FailedActionF("[%s][ArticleChan] Error converting struct to JSON:%V", q.Context, err)
					continue
				} else {
					outputFile := fmt.Sprintf(q.OutputFile, value.Id)
					style.OkLogCF(q.FullLog, "[%V] Received Article, about to write the content into file: %V", q.Context, outputFile)
					err = WriteFile(outputFile, string(jsonData))
					if err != nil {
						style.FailedActionF("Error writing file:%V", err)
					}
				}
			} else {
				style.OkLogCF(q.FullLog, "[%V] Received Article, but skipping file writing to the system either saving is disabled or file already exists, %s", q.Context, value.Id)
			}
		}
	}
}

func GetDirNameAndOutputName(context, ext string) (string, string) {
	ext = strings.TrimLeft(ext, ".")
	dirName := filepath.Join(GetOrDefault(OutputPath, "./").(string), "data/", context, GetOrDefault(SubDirDefined, "").(string))
	_ = CreateDirIfNotExist(dirName)
	outputFile := dirName + "/%s." + ext
	return dirName, outputFile
}
func GetOrDefault(key string, defaultValue any) any {
	if !viper.IsSet(key) {
		return defaultValue
	}
	return viper.Get(key)
}
