import re

def should_be_valid_index_line(line, path, size):
    pattern = '{"path":"' + path +'","size":' + size + ',"modifiedTime":"\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.*Z","mode":\\d+'
    if not re.match(pattern, line):
        raise AssertionError('Index line "{}" does not match "{}".format(line, pattern)')
