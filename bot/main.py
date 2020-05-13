import datetime
import os
import random
import requests
import subprocess
import time
import traceback
import twitter

RATE = 60 * 30
MENTION_RATE = 65

INPUT_FOLDER = ''
OUTPUT_FOLDER = ''

FLICKR_API_KEY = None
TWITTER_CONSUMER_KEY = None
TWITTER_CONSUMER_SECRET = None
TWITTER_ACCESS_TOKEN_KEY = None
TWITTER_ACCESS_TOKEN_SECRET = None

MODE_NAMES = [
    'primitives', # 0
    'triangles',  # 1
    'rectangles', # 2
    'ellipses',   # 3
    'circles',    # 4
    'rectangles', # 5
    'beziers',    # 6
    'ellipses',   # 7
    'polygons',   # 8
]

SINCE_ID = None
START_DATETIME = datetime.datetime.utcnow()
USER_DATETIME = {}

try:
    from config import *
except ImportError:
    print 'no config found!'

class AttrDict(dict):
    # This is a test for Doxygen
    def __init__(self, *args, **kwargs):
        super(AttrDict, self).__init__(*args, **kwargs)
        self.__dict__ = self

class Config(AttrDict):
    def randomize(self):
        self.m = random.choice([1, 5, 6, 7])
        self.n = random.randint(10, 50) * 10
        self.rep = 0
        self.a = 128
        self.r = 300
        self.s = 1200
    def parse(self, text):
        text = (text or '').lower()
        tokens = text.split()
        for i, name in enumerate(MODE_NAMES):
            if name in text:
                self.m = i
        for token in tokens:
            try:
                self.n = int(token)
            except Exception:
                pass
    def validate(self):
        self.m = clamp(self.m, 0, 8)
        if self.m == 6:
            self.a = 0
            self.rep = 19
            self.n = 100
        else:
            self.n = clamp(self.n, 1, 500)
    @property
    def description(self):
        total = self.n + self.n * self.rep
        return '%d %s' % (total, MODE_NAMES[self.m])

def clamp(x, lo, hi):
    if x < lo:
        x = lo
    if x > hi:
        x = hi
    return x

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

def primitive(**kwargs):
    args = []
    for k, v in kwargs.items():
        if v is None:
            continue
        args.append('-%s' % k)
        args.append(str(v))
    args = ' '.join(args)
    cmd = 'primitive %s' % args
    subprocess.call(cmd, shell=True)

def twitter_api():
    return twitter.Api(
        consumer_key=TWITTER_CONSUMER_KEY,
        consumer_secret=TWITTER_CONSUMER_SECRET,
        access_token_key=TWITTER_ACCESS_TOKEN_KEY,
        access_token_secret=TWITTER_ACCESS_TOKEN_SECRET)

def tweet(status, media, in_reply_to_status_id=None):
    api = twitter_api()
    api.PostUpdate(status, media, in_reply_to_status_id=in_reply_to_status_id)

def handle_mentions():
    global SINCE_ID
    print 'checking for mentions'
    api = twitter_api()
    statuses = api.GetMentions(200, SINCE_ID)
    for status in reversed(statuses):
        SINCE_ID = status.id
        print 'handling mention', status.id
        handle_mention(status)
    print 'done with mentions'

def handle_mention(status):
    mentions = status.user_mentions or []
    if len(mentions) != 1:
        print 'mention does not have exactly one mention'
        return
    media = status.media or []
    if len(media) != 1:
        print 'mention does not have exactly one media'
        return
    url = media[0].media_url or None
    if not url:
        print 'mention does not have a media_url'
        return
    created_at = datetime.datetime.strptime(
        status.created_at, '%a %b %d %H:%M:%S +0000 %Y')
    if created_at < START_DATETIME:
        print 'mention timestamp before bot started'
        return
    user_id = status.user.id
    now = datetime.datetime.utcnow()
    td = datetime.timedelta(minutes=5)
    if user_id in USER_DATETIME:
        if now - USER_DATETIME[user_id] < td:
            print 'user mentioned me too recently'
            return
    USER_DATETIME[user_id] = now
    in_path = os.path.join(INPUT_FOLDER, '%s.jpg' % status.id)
    out_path = os.path.join(OUTPUT_FOLDER, '%s.png' % status.id)
    print 'downloading', url
    download_photo(url, in_path)
    config = Config()
    config.randomize()
    config.parse(status.text)
    config.validate()
    status_text = '@%s %s.' % (status.user.screen_name, config.description)
    print status_text
    print 'running algorithm: %s' % config
    primitive(i=in_path, o=out_path, **config)
    if os.path.exists(out_path):
        print 'uploading to twitter'
        tweet(status_text, out_path, status.id)
        print 'done'
    else:
        print 'failed!'

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

def generate():
    date = random_date()
    print 'finding an interesting photo from', date
    photos = interesting(date)
    photo = random.choice(photos)
    print 'picked photo', photo['id']
    in_path = os.path.join(INPUT_FOLDER, '%s.jpg' % photo['id'])
    out_path = os.path.join(OUTPUT_FOLDER, '%s.png' % photo['id'])
    url = photo_url(photo, 'z')
    print 'downloading', url
    download_photo(url, in_path)
    config = Config()
    config.randomize()
    config.validate()
    status_text = '%s. %s' % (config.description, flickr_url(photo['id']))
    print status_text
    print 'running algorithm: %s' % config
    primitive(i=in_path, o=out_path, **config)
    if os.path.exists(out_path):
        print 'uploading to twitter'
        tweet(status_text, out_path)
        print 'done'
    else:
        print 'failed!'

def main():
    previous = 0
    mention_previous = 0
    while True:
        now = time.time()
        if now - previous > RATE:
            previous = now
            try:
                generate()
            except Exception:
                traceback.print_exc()
        if now - mention_previous > MENTION_RATE:
            mention_previous = now
            try:
                handle_mentions()
            except Exception:
                traceback.print_exc()
        time.sleep(5)

def download_photos(folder, date=None):
    try:
        os.makedirs(folder)
    except Exception:
        pass
    date = date or random_date()
    photos = interesting(date)
    for photo in photos:
        url = photo_url(photo, 'z')
        path = '%s.jpg' % photo['id']
        path = os.path.join(folder, path)
        download_photo(url, path)

if __name__ == '__main__':
    main()
