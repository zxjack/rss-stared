from bs4 import BeautifulSoup
from pathlib import Path

# Load all HTML files in the repository
html_files = []
for file_path in Path('.').rglob('*.html'):
    with open(file_path, 'r') as f:
        html_files.append(f.read())

# Merge HTML files into a single index.html file
merged_html = '\n'.join(html_files)
soup = BeautifulSoup(merged_html, 'html.parser')
with open('index.html', 'w') as f:
    f.write(str(soup))
