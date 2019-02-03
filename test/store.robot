** Settings **
Library     OperatingSystem
Library     Process
Library     String
Library     Collections
Library     matchers.py
Resource    common.robot

Test Setup     Begin test
Test Teardown  End test

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
    File should exist   ${store dir}/${block dir}/${small file hash}.block

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
    Should be equal as integers     ${block size}  80

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

    ${block dir}=       Get substring  ${large file hash 1}  0  2
    File should exist   ${store dir}/${block dir}/${large file hash 1}.block

    ${block dir}=       Get substring  ${large file hash 2}  0  2
    File should exist   ${store dir}/${block dir}/${large file hash 2}.block

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
    Should be equal as integers     ${block size 1}     80
    Run keyword and expect error    *
    ...     File should be UTF8 encoded     ${block path 1}

    ${block dir 2}=                 Get substring  ${large file hash 2}  0  2
    Set test variable               ${block path 2}  ${store dir}/${block dir 2}/${large file hash 2}.block
    File should exist               ${block path 2}
    ${block size 2}=                Get file size       ${block path 2}
    Should be equal as integers     ${block size 2}     64
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

** Keywords **
Begin test
    Create directory        ${store dir}
    Copy file               test/resources/salt  ${store dir}/salt

End test
    Remove directory  ${store dir}  recursive=True
