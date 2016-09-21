from Queue import Queue
import itertools
import os
import subprocess
import sys
import threading

def makedirs(x):
    try:
        os.makedirs(x)
    except Exception:
        pass

def primitive(i, o, n, a, m):
    makedirs(os.path.split(o)[0])
    args = (i, o, n, a, m)
    cmd = 'primitive -r 128 -s 512 -i %s -o %s -n %d -a %d -m %d' % args
    subprocess.call(cmd, shell=True)

def create_jobs(in_folder, out_folder, n, a, m):
    result = []
    for name in os.listdir(in_folder):
        base, ext = os.path.splitext(name)
        if ext.lower() not in ['.jpg', '.jpeg', '.png']:
            continue
        out_name = '%d.%%d.png' % (m)
        in_path = os.path.join(in_folder, name)
        out_path = os.path.join(out_folder, base, out_name)
        if os.path.exists(out_path):
            continue
        key = (base, n, m)
        args = (in_path, out_path, n, a, m)
        result.append((key, args))
    return result

def worker(jobs, done):
    while True:
        job = jobs.get()
        log(job)
        primitive(*job)
        done.put(True)

def process(in_folder, out_folder, nlist, alist, mlist, nworkers):
    jobs = Queue()
    done = Queue()
    for i in xrange(nworkers):
        t = threading.Thread(target=worker, args=(jobs, done))
        t.setDaemon(True)
        t.start()
    count = 0
    items = []
    for n, a, m in itertools.product(nlist, alist, mlist):
        for item in create_jobs(in_folder, out_folder, n, a, m):
            items.append(item)
    items.sort()
    for _, job in items:
        jobs.put(job)
        count += 1
    for i in xrange(count):
        done.get()

log_lock = threading.Lock()

def log(x):
    with log_lock:
        print x

if __name__ == '__main__':
    args = sys.argv[1:]
    nlist = [500]
    alist = [128]
    mlist = [0, 1, 3, 5]
    nworkers = 4
    process(args[0], args[1], nlist, alist, mlist, nworkers)
