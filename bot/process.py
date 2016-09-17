from Queue import Queue
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
        result.append((in_path, out_path, n, a, s, m))
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
    for i in range(nworkers):
        t = threading.Thread(target=worker, args=(jobs, done))
        t.setDaemon(True)
        t.start()
    count = 0
    for n, a, s, m in itertools.product(nlist, alist, slist, mlist):
        for job in create_jobs(in_folder, out_folder, n, a, s, m):
            jobs.put(job)
            count += 1
    for i in range(count):
        done.get()

log_lock = threading.Lock()

def log(x):
    with log_lock:
        print x

if __name__ == '__main__':
    nlist = [50, 100, 150, 200]
    alist = [128]
    slist = [4]
    mlist = [0, 1, 3, 5, 2, 4]
    nworkers = 4
    process('input', 'output', nlist, alist, slist, mlist, nworkers)
