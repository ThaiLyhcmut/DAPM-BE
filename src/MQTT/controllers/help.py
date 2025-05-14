import re

class CommandParser:
    # Class constants
    ACTIONS = ["bật", "mở", "tắt", "đóng"]
    DEVICES = ["đèn", "quạt", "máy lạnh", "máy điều hòa", "máy quạt", "tất cả"]
    LOCATIONS = ["phòng ngủ", "phòng khách", "phòng bếp", "tầng 1", "tầng 2", "ban công"]
    
    def __init__(self):
        # Initialize any instance variables if needed
        pass
    
    def extract_keywords(self, text):
        text = text.lower()
        all_keywords = self.ACTIONS + self.DEVICES + self.LOCATIONS
        keywords = re.findall(r'\b(?:' + '|'.join(all_keywords) + r')\b', text)
        print("keyword", keywords)
        return keywords
    
    def split_keywords_by_action(self, keywords):
        result = []
        temp = []
        pre = "ACTION"
        for keyword in keywords:
            if keyword in self.LOCATIONS:
                temp.append(keyword)
                pre = "LOCATION"
            elif keyword in self.ACTIONS:
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
        print("tách câu", result)
        return result
    
    def split_keywords_by_all(self, keywords):
        dist = {}
        pre = "ACTION"
        ok_word = False
        for keyword in keywords:
            if keyword in self.ACTIONS:
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
            if keyword in self.DEVICES:
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
            if keyword in self.LOCATIONS:
                if pre == "DEVICE":
                    dist["ACTIONS"][-1]["DEVICES"][-1]["LOCATION"].append(keyword)
                    pre = "LOCATION"
                elif pre == "LOCATION":
                    dist["ACTIONS"][-1]["DEVICES"][-1]["LOCATION"].append(keyword)
                    pre = "LOCATION"
                if ok_word:
                    try:
                        dist["ACTIONS"][-2]["DEVICES"][-1]["LOCATION"].append(keyword)
                    except IndexError:
                        pass  # Handle potential index error
        print("phân tích câu", dist)
        return dist
    
    def parse_command(self, text):
        """Main method to parse a command text into structured format"""
        keywords = self.extract_keywords(text)
        keyword_groups = self.split_keywords_by_action(keywords)
        results = []
        for keyword_group in keyword_groups:
            result = self.split_keywords_by_all(keyword_group)
            results.append(result)
        return results