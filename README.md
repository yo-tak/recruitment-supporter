# 自動応援ツール
このツールは自社の採用を応援したいと思いつつなかなか毎日画面開いてボタンを押すオペレーションをこなすことができない人に、
快適に応援をしてもらうためのツールです。

## 利用方法
1. Goをインストール (http://golang.jp/install 参照)
2. webdriver, 必要なgoライブラリをインストール
webdriver(今はchromeのみサポートです）：
```
# for mac
brew install chromedriver
```
agoutiをインストール
```
go get github.com/sclevine/agouti
```
3. go ファイルを実行
//TODO いずれもっと簡略に！
```
git clone https://github.com/yohei-takeda/recruitment-supporter.git
cd ./recruiting-supporter
go run main.go <recruitment_page_url> <company_name> <userId> <password>
```
