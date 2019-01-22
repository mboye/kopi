** Settings **
Library     OperatingSystem
Library     Process
Library     String
Library     Collections
Library     matchers.py

Test Setup     Begin test
Test Teardown  End test

** Variables **
${differ bin}           bin/kopi-diff
${indexer bin}          bin/kopi-index
${store bin}            bin/kopi-store
${small file}           test/resources/store/small-file.txt
${small file hash}      144062aa1d1186d6ef1c122d645b567a
${large file}           test/resources/store/large-file.txt
${large file hash 1}    ${small file hash}
${large file hash 2}    074e8e431cc1335d6a44f366adf0eb11
${store dir}            ${TEMPDIR}/simple_store_data
${index}                ${TEMPDIR}/index
${max block size}       64

** Test Cases **
File of same size as block size
    Create index from "${small file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}"
    Length should be        ${lines}  1
    ${line}                 Get from list  ${lines}  0

    ${expected block}=      Create dictionary  hash=${small file hash}  offset=0  size=64
    ${expected blocks}=     Create list  ${expected block}

    Should be index line with blocks        ${line}  ${expected blocks}
    Should be index line with block count   ${line}  1

    ${block dir}=       Get substring  ${small file hash}  0  2
    File should exist   ${store dir}/${block dir}/${small file hash}.block

File larger than block size
    Create index from "${large file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}"
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


Files reuse existing blocks
    Create index from "${small file}" and save it to "${index}"
    ${index data}       Get file        ${index}
    ${index lines}      Split to lines  ${index data}
    Length should be    ${index lines}  1

    ${lines}=               Store index "${index}"
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

    ${lines}=               Store index "${index}"
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
Create index from "${path}" and save it to "${output path}"
    ${rc}=  Run and return RC  ${indexer bin} --init=true ${path} > ${output path} 2>/dev/null
    Should be equal as integers  ${rc}  0

Store index "${index}"
    ${rc}  ${stdout}=  Run and return RC and output  ${store bin} --maxBlockSize ${max block size} ${store dir} < ${index} 2>/dev/null
    Should be equal as integers  ${rc}  0
    Log many  ${stdout}

    ${lines}=   Split to lines  ${stdout}
    [Return]   ${lines}

Begin test
    Create directory        ${store dir}

End test
    Remove directory  ${store dir}  recursive=True
