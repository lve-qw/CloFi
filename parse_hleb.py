import requests
from bs4 import BeautifulSoup
import sys
import json
from pymongo import MongoClient
from pymongo.errors import ConnectionFailure

try:
    client = MongoClient('mongodb://localhost:27017')
    # Test connection
    client.admin.command('ping')
    print("Connected to MongoDB successfully!")
except ConnectionFailure:
    print("Failed to connect to MongoDB")

# Alternatively, with connection string
# client = MongoClient('mongodb://username:password@localhost:27017/')

# Get database
db = client['app_db']
collection = db['products']

animation = ["-", "\\", "|", "/"]

def download_page(num):
    url = f'https://hlebspace.ru/search/?page={num}&query='
    return download_url_page(url)

def download_url_page(url):
    res = requests.get(url)
    return res.text

def find_items(contents):
    soup = BeautifulSoup(contents, "html.parser")
    items = soup.find_all(class_='product-item__image')
    item_url_list = []
    for item in items:
        item_url_list.append('https://hlebspace.ru'+item['href']) # type: ignore
    with open('prod_urls.txt', 'a') as file:
        file.write('\n'.join(item_url_list)+'\n')

def find_page_qt():
    contents = download_page(1)
    find_items(contents)
    soup = BeautifulSoup(contents, "html.parser")
    items = soup.find(class_='paging__list').find_all('a') # type: ignore
    page_qt = int(items[-2].text)
    return page_qt

def find_all_items_url():
    page_qt = find_page_qt()
    for i in range(2,page_qt+1):
        sys.stdout.write(f"\rLoading {animation[i % len(animation)]} \t processed {round(i/page_qt*100,2)}%")
        sys.stdout.flush()
        find_items(download_page(i))

def read_item(url):
    contents = download_url_page(url)
    soup = BeautifulSoup(contents, "html.parser")
    
    name = soup.find(itemprop = 'name').text # type: ignore
    
    price = soup.find(class_ = 'product-price').get('data-price')# type: ignore
    
    photos_urls = []
    try:
        photos_urls = ['https://hlebspace.ru'+i['href'] for i in soup.find(id='s-photos-list').find_all('a')] # type: ignore
    except:
        try:
            photos_urls = ['https://hlebspace.ru'+soup.find(class_='product-photos__image')['href']] # type: ignore
        except:
            print(url, "no photo")

    # for i in photos_urls:
    # TODO: write a picture download
    #     response = requests.get(i)
    #     photo = response.content
    
    if soup.find(class_ = 'stock-none'): availability = False
    else: availability = True
    
    brand = "-"
    try: 
        features = soup.find('div',class_='features-list').find_all(class_='features-list__item') # type: ignore
        for feature in features: 
            if feature.find(class_='features-list__name').text == 'Бренд': # type: ignore
                brand = feature.find(class_='features-list__value').text   # type: ignore
    except:
        brand = "-"

    # brand = soup.find('div',class_='features-list').find_all(class_='features-list__item')[1].find(class_='features-list__value').text # type: ignore
    
    description_html = soup.find('div', class_='product-overview__description col-12')
    if description_html:
        br_tags = description_html.find_all('br') # type: ignore
        for br_tag in br_tags: br_tag.replace_with('\n') # type: ignore
        description = description_html.get_text() # type: ignore
    else: 
        description=''
        
    data = {
        "name": name.strip(),
        "url": url,
        "price": int(float(str(price))), 
        "brand": brand,
        "photos_urls": photos_urls,
        "availability": availability,
        "description": description,
    }
    
    return data
    # print(json_string)
def write_all():
    with open('prod_urls.txt', 'r') as file:
        urls = file.readlines()
        
    n = len(urls)

    all_documents = collection.find().sort("_id", -1).limit(1)
    for doc in all_documents: last_url = doc["url"]
    start = 0
    if last_url: 
        start = urls.index(last_url)+1
        print(f"Last url: {last_url}")

    for i in range(start, len(urls)):
        document = read_item(urls[i])
        result = collection.insert_one(document)
        
        sys.stdout.write(f"\rLoading {animation[i % len(animation)]} \t Inserted document with ID: {result.inserted_id}\t processed {round((i+1)/n*100,2)}%")
        sys.stdout.flush()
    
def main():
    c = "0"
    while not c in ["1", "2", "3"]:
        print("1. Get urls")
        print("2. Save to a MongoDB")
        c = input("-> ")
    if c == "1": find_all_items_url()
    if c == "2": write_all()
    if c == "3":
        all_documents = collection.find().sort("_id", -1).limit(1)
        for doc in all_documents: print(doc["url"])
        #if input("are you sure? y/n ") == "n": return
        #result = collection.delete_many({})
        #print(f"Deleted {result.deleted_count} document")

if __name__ == '__main__':
    main()
    
    
