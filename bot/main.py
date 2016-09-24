import datetime
import os
import random
import requests
import subprocess
import time
import traceback
import twitter

RATE = 60 * 30

FLICKR_API_KEY = None
TWITTER_CONSUMER_KEY = None
TWITTER_CONSUMER_SECRET = None
TWITTER_ACCESS_TOKEN_KEY = None
TWITTER_ACCESS_TOKEN_SECRET = None

MODE_NAMES = [
    'primitives',
    'triangles',
    'rectangles',
    'ellipses',
    'circles',
    'rectangles',
]

try:
    from config import *
except ImportError:
    print 'no config found!'

def random_date(max_days_ago=1000):
    today = datetime.date.today()
    days = random.randint(1, max_days_ago)
    d = today - datetime.timedelta(days=days)
    return d.strftime('%Y-%m-%d')

def interesting(date=None):
    url = 'https://api.flickr.com/services/rest/'
    params = dict(
        api_key=FLICKR_API_KEY,
        format='json',
        nojsoncallback=1,
        method='flickr.interestingness.getList',
    )
    if date:
        params['date'] = date
    r = requests.get(url, params=params)
    return r.json()['photos']['photo']

def photo_url(p, size=None):
    # See: https://www.flickr.com/services/api/misc.urls.html
    if size:
        url = 'https://farm%s.staticflickr.com/%s/%s_%s_%s.jpg'
        return url % (p['farm'], p['server'], p['id'], p['secret'], size)
    else:
        url = 'https://farm%s.staticflickr.com/%s/%s_%s.jpg'
        return url % (p['farm'], p['server'], p['id'], p['secret'])

def download_photo(url, path):
    r = requests.get(url)
    with open(path, 'wb') as fp:
        fp.write(r.content)

def primitive(i, o, n, a=128, m=1):
    args = (i, o, n, a, m)
    cmd = 'primitive -i %s -o %s -n %d -a %d -m %d' % args
    subprocess.call(cmd, shell=True)

def tweet(status, media):
    api = twitter.Api(
        consumer_key=TWITTER_CONSUMER_KEY,
        consumer_secret=TWITTER_CONSUMER_SECRET,
        access_token_key=TWITTER_ACCESS_TOKEN_KEY,
        access_token_secret=TWITTER_ACCESS_TOKEN_SECRET)
    api.PostUpdate(status, media)

def flickr_url(photo_id):
    alphabet = '123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ'
    return 'https://flic.kr/p/%s' % base_encode(alphabet, int(photo_id))

def base_encode(alphabet, number, suffix=''):
    base = len(alphabet)
    if number >= base:
        div, mod = divmod(number, base)
        return base_encode(alphabet, div, alphabet[mod] + suffix)
    else:
        return alphabet[number] + suffix

def run():
    date = random_date()
    print 'finding an interesting photo from', date
    photos = interesting(date)
    photo = random.choice(photos)
    print 'picked photo', photo['id']
    in_path = '%s.jpg' % photo['id']
    out_path = '%s.png' % photo['id']
    url = photo_url(photo, 'm')
    print 'downloading', url
    download_photo(url, in_path)
    n = random.randint(10, 40) * 10
    a = 128
    m = random.choice([1, 3, 5, 1, 3, 5, 1, 3, 4])
    status = '%d %s. %s' % (n, MODE_NAMES[m], flickr_url(photo['id']))
    print status
    print 'running algorithm, n=%d, a=%d, m=%d' % (n, a, m)
    primitive(in_path, out_path, n=n, a=a, m=m)
    if os.path.exists(out_path):
        print 'uploading to twitter'
        tweet(status, out_path)
        print 'done'
    else:
        print 'failed!'

def main():
    previous = 0
    while True:
        while True:
            now = time.time()
            elapsed = now - previous
            if elapsed > RATE:
                previous = now
                break
            time.sleep(5)
        try:
            run()
        except Exception:
            traceback.print_exc()

def download_photos(folder, date=None):
    try:
        os.makedirs(folder)
    except Exception:
        pass
    date = date or random_date()
    photos = interesting(date)
    for photo in photos:
        url = photo_url(photo, 'm')
        path = '%s.jpg' % photo['id']
        path = os.path.join(folder, path)
        download_photo(url, path)

if __name__ == '__main__':
    main()
