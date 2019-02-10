import re
import json
import hashlib
import gzip
from datetime import datetime

ROBOT_LIBRARY_SCOPE = 'TEST CASE'
salt = open('test/resources/salt','rb').read()


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


def should_be_valid_manifest_header(line, today=datetime.utcnow()):
    doc = json.loads(line)
    required_keys = ["ID", "date", "description"]
    type_expectations = [
        ("ID", str),
        ("date", str),
        ("description", str)
    ]

    for key in required_keys:
        if not key in doc:
            raise AssertionError("Header missing key: " + key)

    for key, expected_type in type_expectations:
        if not isinstance(doc[key], expected_type):
            raise AssertionError(
                "Unexpected type of key {}. Expected type {}, but got {}.".format(
                    key, expected_type, type(doc[key])
                )
            )

    print('Basing ID prefix expectation on date: {}'.format(today))
    expected_id_prefix = today.strftime('%Y/%m/%d')
    actual_id = doc['ID']
    if not actual_id.startswith(expected_id_prefix):
        raise AssertionError('Expected ID to start with {}, but got {}'.format(expected_id_prefix, actual_id))

    if not actual_id.endswith('.manifest'):
        raise AssertionError('Expected ID to end with ".manifest", but got {}'.format(actual_id))


def should_be_index_line_with_blocks(line, blocks):
    doc = json.loads(line)

    if not "blocks" in doc:
        raise AssertionError("Index line missing key: blocks")

    if len(blocks) != len(doc["blocks"]):
        raise AssertionError(
            "Expected index line to contain {} blocks, but found {} blocks.".format(
                len(blocks), len(doc["blocks"])
            )
        )

    for block in blocks:
        block["size"] = int(block["size"])
        block["offset"] = int(block["offset"])

        if not block in doc["blocks"]:
            raise AssertionError("Block not found in index line: " + json.dumps(block))


def should_be_index_line_with_block_count(line, num_blocks):
    num_blocks = int(num_blocks)
    doc = json.loads(line)

    if num_blocks == 0:
        if "blocks" in doc:
            raise AssertionError(
                'Index line without blocks should not have key "blocks"'
            )
        else:
            return

    if not "blocks" in doc:
        raise AssertionError("Index line missing key: blocks")

    num_actual_blocks = len(doc["blocks"])
    if num_blocks != num_actual_blocks:
        raise AssertionError(
            "Expected {} blocks, but found {} blocks".format(
                num_blocks, num_actual_blocks
            )
        )

    if not "blocks" in doc:
        raise AssertionError("Index line missing key: blocks")


def file_should_have_sha1_hash(path, expected_hash):
    if salt == '':
        raise RuntimeError('Salt has not been loaded')

    hasher = hashlib.sha1()
    hasher.update(salt)

    with open(path, "rb") as fp:
        for chunk in iter(lambda: fp.read(4096), b""):
            hasher.update(chunk)
    actual_hash = hasher.hexdigest().lower()

    if actual_hash != expected_hash:
        raise AssertionError(
            "Expected SHA1 hash {}, but got {}".format(expected_hash, actual_hash)
        )

def file_should_be_utf8_encoded(path):
    try:
        with open(path, 'rb') as fp:
            data = fp.read()
            data.decode('utf-8')
    except UnicodeError:
        raise AssertionError('File is not UTF-8 encoded')

def get_decompressed_file(path):
    with gzip.open(path, 'rb') as fp:
        return fp.read().decode('utf8')
