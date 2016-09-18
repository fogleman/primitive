from Queue import Queue
import itertools
import os
import subprocess
import threading

def primitive(i, o, n, a, s, m):
    args = (i, o, n, a, s, m)
    cmd = 'primitive -i %s -o %s -n %d -a %d -s %d -m %d' % args
    subprocess.call(cmd, shell=True)

def create_jobs(in_folder, out_folder, n, a, s, m):
    result = []
    try:
        os.makedirs(out_folder)
    except Exception:
        pass
    for name in os.listdir(in_folder):
        if not name.endswith('.jpg'):
            continue
        out_name = '%s.%d.%d.%d.%d.png' % (name[:-4], n, a, s, m)
        in_path = os.path.join(in_folder, name)
        out_path = os.path.join(out_folder, out_name)
        if os.path.exists(out_path):
            continue
        key = (name[:-4], n, m)
        args = (in_path, out_path, n, a, s, m)
        result.append((key, args))
    return result

def worker(jobs, done):
    while True:
        job = jobs.get()
        log(job)
        primitive(*job)
        done.put(True)

def process(in_folder, out_folder, nlist, alist, slist, mlist, nworkers):
    jobs = Queue()
    done = Queue()
    for i in xrange(nworkers):
        t = threading.Thread(target=worker, args=(jobs, done))
        t.setDaemon(True)
        t.start()
    count = 0
    items = []
    for n, a, s, m in itertools.product(nlist, alist, slist, mlist):
        for item in create_jobs(in_folder, out_folder, n, a, s, m):
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
    nlist = [50, 100, 200]
    alist = [128]
    slist = [4]
    mlist = [1, 3, 5]
    nworkers = 4
    process('input1', 'output1', nlist, alist, slist, mlist, nworkers)
