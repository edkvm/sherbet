import csv

from os import path

basepath = path.dirname(__file__)
filepath = path.abspath(path.join(basepath, "./", "sample/media_list.csv"))
savepath = path.abspath(path.join(basepath, "./", "sample/media.html"))

doc_tmpl = '''
<!DOCTYPE html>
<html>
<head>

</head>
<body>
{}
</body>
</html>
'''

img_tmlp = '''
<div class="container"> 
  <div class="item">
  <div class="item">
    <img src="http://localhost:6061{}" />
  </div>
</div>'''

imgs = []
with open(filepath, 'r') as f:
  r = csv.reader(f)
  for row in r:
    imgs.append(img_tmlp.format(row[0], row[0]))


with open(savepath, 'w') as f:
  f.write(doc_tmpl.format(' '.join(imgs)))

