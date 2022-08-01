import requests
import json

result = requests.get("http://127.0.0.1:4040/api/tunnels")
myUrl = result.json()["tunnels"][0]["public_url"]

lineReq = json.dumps({'endpoint':myUrl+'/line/bot/callback'})

lineRes = requests.put("https://api.line.me/v2/bot/channel/webhook/endpoint", headers={"Authorization":"Bearer QMetcOLoLQzjvBfNBO0jxxHCAzrzhQu4lreovVEfPDUXxrol6m/a/PgAd9oXv9lHYn+r3nq9OgLWhNs6ZAtLhNN/7iRO1kro4hMJ/Drngv8Cb1iQciXM11vTFBPQVL4YyBlY+IilWyKD8FR5aqINqgdB04t89/1O/w1cDnyilFU=","Content-Type":"application/json"},data=lineReq)

print(lineRes,lineRes.content)

