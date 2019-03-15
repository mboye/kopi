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
${index a}          ${TEMPDIR}/index.a
${index b}          ${TEMPDIR}/index.b
${diff}             ${TEMPDIR}/index.diff
${stored index}     ${TEMPDIR}/index.stored

** Test Cases **
Store small file
    Create index from "${small file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}" to "${store dir}" and return lines
    Length should be        ${lines}  1
    ${line}                 Get from list  ${lines}  0

    ${expected block}=      Create dictionary  hash=${small file hash}  offset=0  size=64
    ${expected blocks}=     Create list  ${expected block}

    Should be index line with blocks        ${line}  ${expected blocks}
    Should be index line with block count   ${line}  1

    ${block dir}=       Get substring  ${small file hash}  0  2
    Set test variable   ${block path}  ${store dir}/${block dir}/${small file hash}.block
    File should exist   ${block path}

    ${block size}=                  Get file size   ${block path}
    Should be equal as integers     ${block size}   64

Store small file with encryption
    Create index from "${small file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}" with encryption to "${store dir}" and return lines
    Length should be        ${lines}  1
    ${line}                 Get from list  ${lines}  0

    ${expected block}=      Create dictionary  hash=${small file hash}  offset=0  size=64
    ${expected blocks}=     Create list  ${expected block}

    Should be index line with blocks        ${line}  ${expected blocks}
    Should be index line with block count   ${line}  1

    ${block dir}=       Get substring  ${small file hash}  0  2
    Set test variable   ${block path}  ${store dir}/${block dir}/${small file hash}.block
    File should exist   ${block path}

    ${block size}=                  Get file size   ${block path}
    Should be equal as integers     ${block size}  92

    Run Keyword and expect error    *
    ...     File should be UTF8 encoded  ${block path}

Store large file
    Create index from "${large file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}" to "${store dir}" and return lines
    Length should be        ${lines}  1
    ${line}                 Get from list  ${lines}  0

    ${expected block 1}=      Create dictionary  hash=${large file hash 1}  offset=0  size=64
    ${expected block 2}=      Create dictionary  hash=${large file hash 2}  offset=64  size=48
    ${expected blocks}=       Create list  ${expected block 1}  ${expected block 2}

    Should be index line with blocks        ${line}  ${expected blocks}
    Should be index line with block count   ${line}  2

    ${block dir 1}=                 Get substring  ${small file hash}  0  2
    Set test variable               ${block path 1}  ${store dir}/${block dir 1}/${large file hash 1}.block
    File should exist               ${block path 1}
    ${block size 1}=                Get file size       ${block path 1}
    Should be equal as integers     ${block size 1}     64

    ${block dir 2}=                 Get substring  ${large file hash 2}  0  2
    Set test variable               ${block path 2}  ${store dir}/${block dir 2}/${large file hash 2}.block
    File should exist               ${block path 2}
    ${block size 2}=                Get file size       ${block path 2}
    Should be equal as integers     ${block size 2}     48

Store large file with encryption
    Create index from "${large file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}" with encryption to "${store dir}" and return lines
    Length should be        ${lines}  1
    ${line}                 Get from list  ${lines}  0

    ${expected block 1}=      Create dictionary  hash=${large file hash 1}  offset=0  size=64
    ${expected block 2}=      Create dictionary  hash=${large file hash 2}  offset=64  size=48
    ${expected blocks}=       Create list  ${expected block 1}  ${expected block 2}

    Should be index line with blocks        ${line}  ${expected blocks}
    Should be index line with block count   ${line}  2

    ${block dir 1}=                 Get substring  ${small file hash}  0  2
    Set test variable               ${block path 1}  ${store dir}/${block dir 1}/${large file hash 1}.block
    File should exist               ${block path 1}
    ${block size 1}=                Get file size       ${block path 1}
    Should be equal as integers     ${block size 1}     92
    Run keyword and expect error    *
    ...     File should be UTF8 encoded     ${block path 1}

    ${block dir 2}=                 Get substring  ${large file hash 2}  0  2
    Set test variable               ${block path 2}  ${store dir}/${block dir 2}/${large file hash 2}.block
    File should exist               ${block path 2}
    ${block size 2}=                Get file size       ${block path 2}
    Should be equal as integers     ${block size 2}     76
    Run keyword and expect error    *
    ...     File should be UTF8 encoded     ${block path 1}

Store multiple files
    Create index from "${small file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}" to "${store dir}" and return lines
    Length should be        ${lines}  1
    ${line}                 Get from list  ${lines}  0

    ${expected block}=      Create dictionary  hash=${small file hash}  offset=0  size=64
    ${expected blocks}=       Create list  ${expected block}
    Should be index line with blocks        ${line}  ${expected blocks}
    Should be index line with block count   ${line}  1

    ${block dir}=       Get substring  ${small file hash}  0  2
    File should exist   ${store dir}/${block dir}/${small file hash}.block

    # Store large file and expect first block to be reused from small file
    Create index from "${large file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}" to "${store dir}" and return lines
    Length should be        ${lines}  1
    ${line}                 Get from list  ${lines}  0

    ${expected block 1}=      Create dictionary  hash=${large file hash 1}  offset=0  size=64
    ${expected block 2}=      Create dictionary  hash=${large file hash 2}  offset=64  size=48
    ${expected blocks}=       Create list  ${expected block 1}  ${expected block 2}

    Should be index line with blocks        ${line}  ${expected blocks}
    Should be index line with block count   ${line}  2

    ${block dir}=       Get substring  ${large file hash 1}  0  2
    File should exist   ${store dir}/${block dir}/${large file hash 1}.block

    ${block dir}=       Get substring  ${large file hash 2}  0  2
    File should exist   ${store dir}/${block dir}/${large file hash 2}.block

Store multiple files and print progress
    Create index from "${backup source dir}" and save it to "${index}"

    ${result}=  Run process  ${store bin} --progress 1 ${store dir} < ${index}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    Should contain  ${result.stderr}  progress
    Should contain  ${result.stderr}  elapsed_time
    Should contain  ${result.stderr}  remaining_time
    Should contain  ${result.stderr}  byte_progress
    Should contain  ${result.stderr}  file_progress
    Should contain  ${result.stderr}  \= 0.00%
    Should contain  ${result.stderr}  \= 100.00%

Store missing file
    Copy file           ${small file}  ${small file}-copy
    Create index from "${backup source dir}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  5

    Remove file         ${small file}-copy

    ${result}=  Run process  ${store bin} ${store dir} < ${index}  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}

    ${stored lines}         Split to lines  ${result.stdout}
    Length should be        ${stored lines}  4

    Should contain          ${result.stderr}  "File not found"
    Should contain          ${result.stderr}  ${small file}-copy

Store empty file
    Create index from "${empty file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}" to "${store dir}" and return lines
    Length should be        ${lines}  1
    ${line}                 Get from list  ${lines}  0

    Should be index line with block count   ${line}  0

    ${block dir}=       Get substring  ${empty file hash}  0  2
    Set test variable   ${block path}  ${store dir}/${block dir}/${empty file hash}.block
    Run keyword and expect error  *
    ...     File should exist   ${block path}

Store diff
    Create index from "test/resources/diff" and save it to "${index a}"
    Sleep   2s
    Touch   test/resources/diff/file-a.txt
    Create index from "test/resources/diff" and save it to "${index b}"

    ${lines}   Diff indices ${index a} and ${index b}, and save result to ${diff}
    Length should be    ${lines}    4

    Store index "${diff}" to "${store dir}" and save output to "${stored index}"

    ${stored index data}    Get file                        ${stored index}

    # Stored index should contain 4 index lines
    ${stored index lines}   Split to lines                  ${stored index data}
    Length should be        ${lines}  4

    # Stored index should contain only 1 modified index line
    ${matches}=             Get lines containing string     ${stored index data}    "modified":true
    ${modified lines}       Split to lines                  ${matches}
    Length should be        ${modified lines}   1

** Keywords **
Begin test
    Create directory        ${store dir}
    Copy file               test/resources/salt  ${store dir}/salt

End test
    Remove directory    ${store dir}  recursive=True
    Remove file         ${small file}-copy
    Remove file         ${index a}
    Remove file         ${index b}
    Remove file         ${diff}
    Remove file         ${stored index}

Diff indices ${path a} and ${path b}, and save result to ${diff output}
    ${result}=  Run process  ${differ bin} ${path a} ${path b} | tee "${diff output}"  shell=True
    Log many    ${result.stdout}
    Log many    ${result.stderr}
    Should be equal as integers  ${result.rc}  0  ${result.stderr}
    ${lines}=   Split to lines  ${result.stdout}
    [Return]   ${lines}
