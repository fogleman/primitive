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
    run(args[0], args[1])

if __name__ == '__main__':
    main()
