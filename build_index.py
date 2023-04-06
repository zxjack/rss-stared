import os
import requests
from bs4 import BeautifulSoup

# 获取需要的参数
repo_name = os.getenv('GITHUB_REPOSITORY').split("/")[1]
branch_name = os.getenv('GITHUB_REF').split("/")[2]
repo_url = f"https://github.com/{os.getenv('GITHUB_REPOSITORY')}"
raw_files_url = f"{repo_url}/tree/{branch_name}"
file_extension = ".html"

# 从GitHub仓库下载Raw文件列表
response = requests.get(raw_files_url)
soup = BeautifulSoup(response.text, 'html.parser')

# 筛选出HTML文件
html_files = []
for a in soup.find_all('a', href=True):
    if a['href'].endswith(file_extension):
        html_files.append(a['href'])

# 生成HTML链接列表
link_list = ""
for f in html_files:
    link_list += f'<a href="{repo_url}/blob/{branch_name}/{f}">{f}</a><br>'

# 创建index.html文件
with open('index.html', 'w') as index_file:
    index_file.write(f'<html><body>{link_list}</body></html>')