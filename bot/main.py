import datetime
import os
import random
import requests
import time
import twitter

RATE = 60 * 15

FLICKR_API_KEY = None
TWITTER_CONSUMER_KEY = None
TWITTER_CONSUMER_SECRET = None
TWITTER_ACCESS_TOKEN_KEY = None
TWITTER_ACCESS_TOKEN_SECRET = None

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

def primitive(i, o, n, a=128, s=1, m=1):
    args = (i, o, n, a, s, m)
    os.system('primitive -i %s -o %s -n %d -a %d -s %d -m %d' % args)

def tweet(status, media):
    api = twitter.Api(
        consumer_key=TWITTER_CONSUMER_KEY,
        consumer_secret=TWITTER_CONSUMER_SECRET,
        access_token_key=TWITTER_ACCESS_TOKEN_KEY,
        access_token_secret=TWITTER_ACCESS_TOKEN_SECRET)
    api.PostUpdate(status, media)

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
    n = random.choice([50, 100, 150])
    a = 128
    s = 4
    m = random.randint(1, 4)
    if random.random() < 0.5:
        a /= 2
        n *= 2
    print 'running algorithm, n=%d, a=%d, s=%d, m=%d' % (n, a, s, m)
    primitive(in_path, out_path, n=n, a=a, s=s, m=m)
    if os.path.exists(out_path):
        print 'uploading to twitter'
        tweet('', out_path)
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
        run()

if __name__ == '__main__':
    main()
