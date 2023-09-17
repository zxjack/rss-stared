# rss-stared

这是一个使用 github 的工作流自动备份 tinytinyrss 的脚本。实现的功能如下：

1. 自动抓取 tinytinyrss 中 star feeds 链接，自动下载相关的文章
2. 以每一个文章的 title 作为文件名作为链接，自动建立相关的 index.html
3. 自动生成 github pages .

以上是因为 ttrss 作为自建服务，有一定风险丢失数据，对于重要的，感兴趣的文章有备份的必要，使用这个脚本可以完成基本需求。

（这个仓库的所有脚本都有 chatgpt 生成，根据对话提交需求，修改 bug,最终完成的。）
