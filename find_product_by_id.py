from pymongo import MongoClient
from bson import ObjectId

client = MongoClient("mongodb://localhost:27017/")
db = client["app_db"]
collection = client["products"]

doc = collection.find_one({"_id": ObjectId(input())})

print(doc)
