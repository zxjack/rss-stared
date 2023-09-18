import os
import datetime
from bs4 import BeautifulSoup

# 指定目录和扩展名
directory = '.'
extension = '.html'

# 遍历目录中的文件
files = [f for f in os.listdir(directory) if f.endswith(extension)]

# 创建index文件
with open('index.html', 'w') as f:
    # 写入HTML头部
    f.write('<html><body>')
    # 获取文件的修改时间
    for file in files:
        mtime = datetime.datetime.fromtimestamp(os.path.getmtime(file))
        
        # 按照修改时间对文件进行排序
        files.sort(key=lambda file: mtime)
    
    # 遍历html文件
    for file in files:
        # 记录文件名和链接
        name = file[:-len(extension)]
        link = file
        
        # 创建HTML链接
        soup = BeautifulSoup('<a></a>', 'html.parser')
        tag = soup.a
        tag.string = name
        tag['href'] = link
        
        # 将链接写入文件
        f.write(str(tag))
        f.write('<br>')
    
    # 写入HTML尾部
    f.write('</body></html>')
