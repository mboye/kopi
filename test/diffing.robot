** Settings **
Library     OperatingSystem
Library     Process
Library     String
Library     Collections
Library     matchers.py
Resource    common.robot

Test Setup     Begin test
Test Teardown  End test

** Variables **
${differ bin}       bin/kopi-diff
${index a}          ${TEMPDIR}/index.a
${index b}          ${TEMPDIR}/index.b
${stored index a}   ${TEMPDIR}/index.a.stored

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
    Sleep   2s
    Touch   test/resources/diff/file-a.txt
    Create index from "test/resources/diff" and save it to "${index b}"

    ${lines}   Diff indices ${index a} and ${index b}
    Length should be    ${lines}    4

    ${line}=                    Get from list   ${lines}  1
    Should be valid index line  ${line}  path=test/resources/diff/file-a.txt  size=10  modified=True

    ${line}=                    Get from list   ${lines}  0
    Should be valid index line  ${line}  path=test/resources/diff   size=0

Preserve blocks of unmodified files
    Create index from "test/resources/diff" and save it to "${index a}"
    Store index "${index a}" to "${store dir}" and save output to "${stored index a}"
    Sleep   2s
    Touch   test/resources/diff/file-a.txt
    Create index from "test/resources/diff" and save it to "${index b}"

    ${lines}            Diff indices ${stored index a} and ${index b}
    Length should be    ${lines}    4

    # Blocks of modified files should be discarded
    ${line}=                    Get from list   ${lines}  1
    Should be valid index line  ${line}  path=test/resources/diff/file-a.txt  size=10  modified=True
    Should be index line with block count  ${line}  0

    # Blocks of unmodified files should be preserved
    ${line}=                    Get from list   ${lines}  3
    Should be valid index line  ${line}  path=test/resources/diff/subdir/file-b.txt  size=10  modified=False
    Should be index line with block count  ${line}  1

Missing indices
    Run keyword and expect error  *failed to open index*
    ...  Diff indices ${index a} and /missing/index

    Run keyword and expect error  *failed to open index*
    ...  Diff indices /missing/index and ${index a}

** Keywords **
Diff indices ${path a} and ${path b}
    ${result}=  Run process  ${differ bin} ${path a} ${path b}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}
    ${lines}=   Split to lines  ${result.stdout}
    [Return]   ${lines}

Begin test
    Create directory        ${store dir}
    Copy file               test/resources/salt  ${store dir}/salt

End test
    Remove file         ${index a}
    Remove file         ${index b}
    Remove directory    ${store dir}  recursive=True
