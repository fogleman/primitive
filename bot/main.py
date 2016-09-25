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

SINCE_ID = None
START_DATETIME = datetime.datetime.utcnow()
USER_DATETIME = {}

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
    in_path = '%s.jpg' % status.id
    out_path = '%s.png' % status.id
    print 'downloading', url
    download_photo(url, in_path)
    n = random.randint(10, 40) * 10
    a = 128
    m = random.choice([1, 3, 5, 1, 3, 5, 1, 3, 4])
    text = (status.text or '').lower()
    text = ''.join(x for x in text if x.isalnum() or x.isspace())
    tokens = text.split()
    for mode in tokens:
        if mode in MODE_NAMES:
            m = MODE_NAMES.index(mode)
    for count, mode in zip(tokens, tokens[1:]):
        if count.isdigit() and mode in MODE_NAMES:
            n = int(count)
            if n < 10:
                n = 10
            if n > 400:
                n = 400
            m = MODE_NAMES.index(mode)
            break
    status_text = '@%s %d %s.' % (status.user.screen_name, n, MODE_NAMES[m])
    print status_text
    print 'running algorithm, n=%d, a=%d, m=%d' % (n, a, m)
    primitive(in_path, out_path, n=n, a=a, m=m)
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
    in_path = '%s.jpg' % photo['id']
    out_path = '%s.png' % photo['id']
    url = photo_url(photo, 'm')
    print 'downloading', url
    download_photo(url, in_path)
    n = random.randint(10, 40) * 10
    a = 128
    m = random.choice([1, 3, 5, 1, 3, 5, 1, 3, 4])
    status_text = '%d %s. %s' % (n, MODE_NAMES[m], flickr_url(photo['id']))
    print status_text
    print 'running algorithm, n=%d, a=%d, m=%d' % (n, a, m)
    primitive(in_path, out_path, n=n, a=a, m=m)
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
        url = photo_url(photo, 'm')
        path = '%s.jpg' % photo['id']
        path = os.path.join(folder, path)
        download_photo(url, path)

if __name__ == '__main__':
    main()
