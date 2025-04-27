import re

ACTIONS = ["bật", "mở", "tắt", "đóng"]
DEVICES = ["đèn", "quạt", "máy lạnh", "máy điều hòa", "máy quạt", "tất cả"]
LOCATIONS = ["phòng ngủ", "phòng khách", "phòng bếp", "tầng 1", "tầng 2", "ban công"]

def extract_keywords(text):
  text = text.lower()
  all_keywords = ACTIONS + DEVICES + LOCATIONS
  keywords = re.findall(r'\b(?:' + '|'.join(all_keywords) + r')\b', text)
  print("keyword", keywords)
  return keywords

def split_keywords_by_action(keywords):
  result = []
  temp = []
  pre = "ACTION"
  for keyword in keywords:
    if keyword in LOCATIONS:
      temp.append(keyword)
      pre = "LOCATION"
    elif keyword in ACTIONS:
      if pre == "LOCATION":
        if temp:
          result.append(temp)
        temp = [keyword]
      else:
        temp.append(keyword)
      pre = "ACTION"
    else:
      temp.append(keyword)
      pre = "DEVICE"
  if temp:
    result.append(temp)
  print("tách câu",result)
  return result

def split_keywords_by_all(keywords: list):
  dist = {}
  pre = "ACTION"
  ok_word = False
  for keyword in keywords:
    if keyword in ACTIONS:
      if pre != "DEVICE":
        dist["ACTIONS"] = [{
          "ACTION": keyword,
        }]
      else:
        dist["ACTIONS"].append({
          "ACTION": keyword,
        })
        ok_word = True
      pre = "ACTION"
    if keyword in DEVICES:
      if pre == "ACTION":
        dist["ACTIONS"][-1]["DEVICES"] = [{
          "DEVICE": [keyword],
          "LOCATION": []
        }]
        pre = "DEVICE"
      elif pre == "DEVICE":
        dist["ACTIONS"][-1]["DEVICES"][-1]["DEVICE"].append(keyword)
        pre = "DEVICE"
      elif pre == "LOCATION":
        print(keyword)
        dist["ACTIONS"][-1]["DEVICES"].append({
          "DEVICE": [keyword],
          "LOCATION": []
        })
        pre = "DEVICE"
    if keyword in LOCATIONS:
      if pre == "DEVICE":
        dist["ACTIONS"][-1]["DEVICES"][-1]["LOCATION"].append(keyword)
        pre = "LOCATION"
      elif pre == "LOCATION":
        dist["ACTIONS"][-1]["DEVICES"][-1]["LOCATION"].append(keyword)
        pre = "LOCATION"
      if ok_word:
        dist["ACTIONS"][-2]["DEVICES"][-1]["LOCATION"].append(keyword)
  print("phân tích câu", dist)

# Test
if __name__ == "__main__":
  test_texts = [
    "Bật ở phòng ngủ và phòng khách, tắt máy lạnh ở phòng bếp, mở máy điều hòa ở tầng 2, đóng máy quạt ở ban công.",
    "Tắt đèn và quạt trong phòng bếp, mở máy lạnh ở tầng 1, đóng máy điều hòa ở phòng khách, bật đèn ở ban công",
    "Mở máy quạt và đèn ở phòng khách, tắt máy điều hòa ở phòng ngủ, bật máy lạnh ở tầng 2, đóng đèn ở phòng bếp",
    "Bật quạt ở phòng ngủ, tắt đèn ở phòng khách, mở máy điều hòa ở ban công, đóng máy lạnh ở phòng bếp.",
    "Tắt máy lạnh và đèn ở tầng 1, mở quạt và máy điều hòa ở tầng 2, đóng đèn và quạt ở phòng khách",
    "Bật đèn và quạt ở phòng khách rồi tắt máy lạnh và mở máy điều hòa ở phòng ngủ, sau đó đóng máy quạt ở tầng 1 và bật đèn ở ban công, mở quạt ở phòng bếp và tắt máy lạnh trong phòng khách, và cuối cùng đóng máy điều hòa ở tầng 2",
  ]
  for test_text in test_texts:
    print(f"Text: {test_text}")
    keywords = split_keywords_by_action(extract_keywords(test_text))
    finals = []
    for keyword in keywords:
      final = split_keywords_by_all(keyword)