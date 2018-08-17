# 自動応援ツール
このツールは自社の採用を応援したいと思いつつなかなか毎日画面開いてボタンを押すオペレーションをこなすことができない人に、
快適に応援をしてもらうためのツールです。

## 利用方法
1. Goをインストール
2. webdriverをインストール
chromeなら以下実行：
```
# for mac
brew install chromedriver
```
3. go ファイルを実行
//TODO いずれもっと簡略に！
```
git clone <this_repository>
cd ./recruiting-supporter //このプロジェクトのmain.goが置いてある場所
go run main.go <recruitment_page_url> <company_name> <userId> <password>
```
