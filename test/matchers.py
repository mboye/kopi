import re
import json


def should_be_valid_index_line(line, path, size, modified=None):
    doc = json.loads(line)
    required_keys = ["path", "modifiedTime", "mode"]
    type_expectations = [
        ("path", str),
        ("modifiedTime", str),
        ("mode", int),
        ("size", int),
    ]

    if modified:
        required_keys.append("modified")
        type_expectations.append(("modified", bool))

    for key in required_keys:
        if not key in doc:
            raise AssertionError("Index line is missing key: " + key)

    for key, expected_type in type_expectations:
        if not isinstance(doc[key], expected_type):
            raise AssertionError(
                "Unexpected type of key {}. Expected type {}, but got {}.".format(
                    key, expected_type, type(doc[key])
                )
            )

    if not re.match(
        "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.*Z", doc["modifiedTime"]
    ):
        raise AssertionError("Format of modifiedTime value is invalid")

    if doc["path"] != path:
        raise AssertionError("Unexpected path")

    if doc["size"] != int(size):
        raise AssertionError(
            "Unexpected size. Expected {}, but got {}".format(size, doc["size"])
        )

    if modified and doc["modified"] != True:
        raise AssertionError(
            'Value key "modified" must always be true if key is present'
        )


class KopiFile(object):
    def __init__(self, file_json):
        self.path = file_json["path"]
        # self.diff = file_json['diff']


def parse_file_line(lines, index):
    index = int(index)
    data = json.loads(lines[index])
    return KopiFile(data)
