package main
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)
type TodoUnit struct{
	Title	string `json:"title"`
	Desc	string `json:"desc"`
	Done	bool `json:"done"`
}
func initDataFile()(string, error){
	configDir, err := os.UserConfigDir()
	if(err != nil){
		return "", err
	}
	appDir := filepath.Join(configDir, "todo-go-cui")
	filePath := filepath.Join(appDir, "todos.json")
	if err := os.MkdirAll(appDir, 0755); err != nil{
		return "", err
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err){
		emptyJson := []byte("[]\n")
		if err := os.WriteFile(filePath, emptyJson, 0644); err != nil {
			return "", err
		}
	}
	return filePath, nil
}
func main(){
	if(len(os.Args) < 2){
		fmt.Println("コマンドを指定してください。")
		fmt.Println("コマンド一覧:")
		fmt.Println("	add")
		fmt.Println("	list")
		fmt.Println("	done")
		os.Exit(1)
	}
	dataPath, err := initDataFile()
	if(err != nil){
		fmt.Printf("データファイルの初期化に失敗しました:%v\n", err)
		os.Exit(1)
	}
	switch(os.Args[1]){
		case "add":
			if(len(os.Args) < 4){
				fmt.Println("使用法:todo add [タイトル] [内容]")
				os.Exit(1)
			}
			fileData, err := os.ReadFile(dataPath)
			if(err != nil){
				fmt.Printf("ファイルの読み込みに失敗しました:%v\n", err)
				os.Exit(1)
			}
			var todos []TodoUnit
			if err := json.Unmarshal(fileData, &todos); err != nil{
				fmt.Printf("データの解析に失敗しました:%v\n", err)
				os.Exit(1)
			}
			for _, item := range todos {
				if(item.Title == os.Args[2] && !item.Done){
					fmt.Printf("未完了のタスク %s はすでに存在します。\n", os.Args[2])
					os.Exit(1)
				}
			}
			newTodo := TodoUnit{
				Title: os.Args[2],
				Desc: os.Args[3],
				Done: false,
			}
			todos = append(todos, newTodo)
			updatedData, err := json.MarshalIndent(todos, "", "  ")
			if(err != nil){
				fmt.Printf("データの変換に失敗しました:%v\n", err)
				os.Exit(1)
			}
			if err := os.WriteFile(dataPath, updatedData, 0644); err != nil{
				fmt.Printf("データの保存に失敗しました:%v\n", err)
				os.Exit(1)
			}
			fmt.Println("操作は正常に完了しました。")
		case "list":
			fileData, err := os.ReadFile(dataPath)
			if(err != nil){
				fmt.Printf("ファイルの読み込みに失敗しました:%v\n", err)
				os.Exit(1)
			}
			var todos []TodoUnit
			if err := json.Unmarshal(fileData, &todos); err != nil{
				fmt.Printf("データの解析に失敗しました:%v\n", err)
				os.Exit(1)
			}
			showAll := false
			if(len(os.Args) > 2 && os.Args[2] == "all"){
				showAll = true
			}
			for _, item := range todos{
				if(!showAll && item.Done){
					continue
				}
				status := "[TODO]"
				if(item.Done){
					status = "[DONE]"
				}
				fmt.Printf("%s %s (%s)\n", status, item.Title, item.Desc)
			}
		case "done":
			if(len(os.Args) < 3){
				fmt.Println("使用法:todo done [タイトル]")
				os.Exit(1)
			}
			fileData, err := os.ReadFile(dataPath)
			if(err != nil){
				fmt.Printf("ファイルの読み込みに失敗しました:%v\n", err)
				os.Exit(1)
			}
			var todos []TodoUnit
			if err := json.Unmarshal(fileData, &todos); err != nil{
				fmt.Printf("データの解析に失敗しました:%v\n", err)
				os.Exit(1)
			}
			found := false
			for i := 0; i < len(todos); i++ {
				if(todos[i].Title == os.Args[2] && todos[i].Done == false){
					todos[i].Done = true
					found = true
					break
				}
			}
			if(!found){
				fmt.Printf("タスク %s は見つかりませんでした。\n", os.Args[2])
				os.Exit(1)
			}
			updatedData, err := json.MarshalIndent(todos,"", "  ")
			if(err != nil){
				fmt.Printf("データの変換に失敗しました:%v\n", err)
				os.Exit(1)
			}
			if err := os.WriteFile(dataPath, updatedData, 0644); err != nil {
				fmt.Printf("データの保存に失敗しました:%v\n", err)
				os.Exit(1)
			}
			fmt.Println("操作は正常に完了しました。")
		case "erase":
			if(len(os.Args) < 3){
				fmt.Println("使用法:todo erase [タイトル]")
				os.Exit(1)
			}
			fileData, err := os.ReadFile(dataPath)
			if(err != nil){
				fmt.Printf("ファイルの読み込みに失敗しました:%v\n", err)
				os.Exit(1)
			}
			var todos []TodoUnit
			if err := json.Unmarshal(fileData, &todos); err != nil {
				fmt.Printf("データの解析に失敗しました:%v\n", err)
				os.Exit(1)
			}
			var keptTodos []TodoUnit
			found := false
			for _, item := range todos {
				if(item.Title == os.Args[2]){
					found = true
					continue
				}
				keptTodos = append(keptTodos, item)
			}
			if(!found){
				fmt.Printf("タスク %s は見つかりませんでした。\n", os.Args[2])
				os.Exit(1)
			}
			updatedData, err := json.MarshalIndent(keptTodos, "", "  ")
			if(err != nil){
				fmt.Printf("データの変換に失敗しました:%v\n", err)
				os.Exit(1)
			}
			if err := os.WriteFile(dataPath, updatedData, 0644); err != nil {
				fmt.Printf("データの保存に失敗しました:%v\n", err)
				os.Exit(1)
			}
			fmt.Println("操作は正常に完了しました。")
		default:
			fmt.Printf("不明なコマンド:%s\n", os.Args[1])
			os.Exit(1)
	}
}
