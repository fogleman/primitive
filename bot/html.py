import os
import sys

def run(in_folder, out_folder):
    seen = set()
    for name in os.listdir(out_folder):
        if not name.endswith('.png'):
            continue
        seen.add(name.split('.')[0])
    for name in os.listdir(in_folder):
        if not name.endswith('.jpg'):
            continue
        name = name[:-4]
        if name not in seen:
            continue
        for m in [1, 3, 5]:
            print '<tr>'
            path = '%s.jpg' % name
            print '<td><img src="%s"></td>' % os.path.join(in_folder, path)
            for n in [50, 100, 200]:
                path = '%s.%d.128.4.%d.png' % (name, n, m)
                print '<td><img src="%s"></td>' % os.path.join(out_folder, path)
            print '</tr>'

def main():
    args = sys.argv[1:]
    print HEADER
    run(args[0], args[1])
    print FOOTER

HEADER = '''
<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>PrimitivePic</title>
<style>
body {
    margin: 0;
    padding: 0;
    font-family: sans-serif;
}
table {
    border-collapse: collapse;
    margin: 4px;
}
img {
    width: 400px;
    display: block;
    margin: 4px;
}
td {
    padding: 0;
}
</style>
</head>
<body>
<table>

<tr>
<th>original</th>
<th>50 shapes</th>
<th>100 shapes</th>
<th>200 shapes</th>
</tr>
'''

FOOTER = '''
</table>
</body>
</html>
'''

if __name__ == '__main__':
    main()
