** Settings **
Library     OperatingSystem
Library     Process
Library     String
Library     Collections
Library     matchers.py
Resource    common.robot

** Variables **
${relative path}    test/resources/index
${absolute path}    ${CURDIR}/resources/index

** Test Cases **
Create index from relative path
    ${lines}=   Create index from "${relative path}" and return lines
    Length should be    ${lines}    4

    ${line}=                    Get from list   ${lines}  0
    Should be valid index line  ${line}  path=test/resources/index  size=0

    ${line}=                    Get from list   ${lines}  1
    Should be valid index line  ${line}  path=test/resources/index/file-a.txt   size=10

    ${line}=                    Get from list   ${lines}  2
    Should be valid index line  ${line}  path=test/resources/index/subdir  size=0

    ${line}=                    Get from list   ${lines}  3
    Should be valid index line  ${line}  path=test/resources/index/subdir/file-b.txt  size=10

Create index from absolute path
    ${lines}=   Create index from "${absolute path}" and return lines
    Length should be    ${lines}    4

    ${line}=                    Get from list   ${lines}  0
    Should be valid index line  ${line}  path=${CURDIR}/resources/index  size=0

    ${line}=                    Get from list   ${lines}  1
    Should be valid index line  ${line}  path=${CURDIR}/resources/index/file-a.txt  size=10

    ${line}=                    Get from list   ${lines}  2
    Should be valid index line  ${line}  path=${CURDIR}/resources/index/subdir  size=0

    ${line}=                    Get from list   ${lines}  3
    Should be valid index line  ${line}  path=${CURDIR}/resources/index/subdir/file-b.txt  size=10

Create index without recursing
    ${result}=  Run process  ${indexer bin}  --recursive\=false  ${relative path}
    Should be equal as integers  ${result.rc}  0
    ${lines}=   Split to lines  ${result.stdout}
    Length should be    ${lines}    2

    ${line}=                    Get from list   ${lines}  0
    Should be valid index line  ${line}  path=test/resources/index  size=0

    ${line}=                    Get from list   ${lines}  1
    Should be valid index line  ${line}  path=test/resources/index/file-a.txt  size=10

Create index from missing path
    ${result}=  Run process  ${indexer bin}  /tmp/missing/path
    Should be equal as integers  ${result.rc}  1
    Should contain  ${result.stderr}  Failed to walk path: /tmp/missing/path
    Should contain  ${result.stderr}  no such file or directory
