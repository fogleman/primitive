import os

def run(in_folder, out_folder):
    for name in os.listdir(in_folder):
        if not name.endswith('.jpg'):
            continue
        name = name[:-4]
        for m in [1, 3, 5, 0]:
            print '<tr>'
            path = '%s.jpg' % name
            print '<td><img src="%s"></td>' % os.path.join(in_folder, path)
            for n in [50, 100, 200]:
                path = '%s.%d.128.4.%d.png' % (name, n, m)
                print '<td><img src="%s"></td>' % os.path.join(out_folder, path)
            print '</tr>'

def main():
    run('input', 'output')

if __name__ == '__main__':
    main()
