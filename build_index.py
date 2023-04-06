import os
import requests
from bs4 import BeautifulSoup
from git import Repo

# 获取需要的参数
repo_name = os.getenv('GITHUB_REPOSITORY').split("/")[1]
branch_name = os.getenv('GITHUB_REF').split("/")[2]
repo_url = f"https://github.com/{os.getenv('GITHUB_REPOSITORY')}"
raw_files_url = f"{repo_url}/tree/{branch_name}"
file_extension = ".html"
publish_branch = "gh-pages"

# 克隆要推送的分支并切换到其中
if os.path.exists(publish_branch):
    repo = Repo(publish_branch)
    repo.git.checkout(publish_branch)
else:
    repo = Repo.clone_from(f"{repo_url}.git", publish_branch, branch=publish_branch)
    repo.git.checkout(publish_branch)

# 从GitHub仓库下载Raw文件列表
response = requests.get(raw_files_url)
soup = BeautifulSoup(response.text, 'html.parser')

# 筛选出HTML文件
html_files = []
for a in soup.find_all('a', href=True):
    if a['href'].endswith(file_extension):
        html_files.append(a['href'])

# 构建HTML链接列表
link_list = ""
for f in html_files:
    file_name = os.path.basename(f)  # 获取文件名称
    link_list += f'<a href="{f}">{file_name}</a><br>'

# 创建index.html文件
with open('index.html', 'w') as index_file:
    index_file.write(f'<html><body>{link_list}</body></html>')
    
# 读取环境变量
author_name = os.getenv('GIT_AUTHOR_NAME')
author_email = os.getenv('GIT_AUTHOR_EMAIL')
committer_name = os.getenv('GIT_COMMITTER_NAME')
committer_email = os.getenv('GIT_COMMITTER_EMAIL')    

# 设置用户名和电子邮件
Repo().config_writer().set_value('user', 'name', author_name).release()
Repo().config_writer().set_value('user', 'email', author_email).release()

# 提交并推送到代码库
repo.git.add("index.html")
repo.git.commit("-m", "Update index.html")
repo.git.push("origin", publish_branch)