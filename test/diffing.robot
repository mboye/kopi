** Settings **
Library     OperatingSystem
Library     Process
Library     String
Library     Collections
Library     matchers.py

Test Teardown  End of test

** Variables **
${differ bin}   bin/kopi-diff
${indexer bin}  bin/kopi-index
${index a}      ${TEMPDIR}/index.a
${index b}      ${TEMPDIR}/index.b

** Test Cases **
Identical indices
    Create index from "test/resources/diff" and save it to "${index a}"

    ${lines}   Diff indices ${index a} and ${index a}
    Length should be    ${lines}    4

    ${line}=                    Get from list   ${lines}  0
    Should be valid index line  ${line}  path=test/resources/diff  size=0

    ${line}=                    Get from list   ${lines}  1
    Should be valid index line  ${line}  path=test/resources/diff/file-a.txt   size=10

    ${line}=                    Get from list   ${lines}  2
    Should be valid index line  ${line}  path=test/resources/diff/subdir  size=0

    ${line}=                    Get from list   ${lines}  3
    Should be valid index line  ${line}  path=test/resources/diff/subdir/file-b.txt  size=10

Timestamp changed
    Create index from "test/resources/diff" and save it to "${index a}"
    Sleep  2s
    Run and return RC  touch test/resources/diff/file-a.txt
    Create index from "test/resources/diff" and save it to "${index_b}"

    ${lines}   Diff indices ${index a} and ${index b}
    Length should be    ${lines}    4

    ${line}=                    Get from list   ${lines}  1
    Should be valid index line  ${line}  path=test/resources/diff/file-a.txt  size=10  modified=True

    ${line}=                    Get from list   ${lines}  0
    Should be valid index line  ${line}  path=test/resources/diff   size=0

Missing indices
    Run keyword and expect error  1 != 0
    ...  Diff indices ${index a} and /missing/index

    Run keyword and expect error  1 != 0
    ...  Diff indices /missing/index and ${index a}

** Keywords **
Diff indices ${path a} and ${path b}
    ${result}=  Run process  ${differ bin}  ${path a}  ${path b}
    Should be equal as integers  ${result.rc}  0
    ${lines}=   Split to lines  ${result.stdout}
    [Return]   ${lines}

Create index from "${path}" and save it to "${output path}"
    ${rc}=  Run and return RC  ${indexer bin} ${path} > ${output path} 2>/dev/null
    Should be equal as integers  ${rc}  0

End of test
    Remove file  ${index a}
    Remove file  ${index b}
